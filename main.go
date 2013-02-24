package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
)

type Startup struct {
	Id            bson.ObjectId "_id"
	Name          string
	Blog_Url      string
	Blog_Feed_Url string
	Homepage_Url  string
	Feed          []byte
}

type Post struct {
	// Id        bson.ObjectId "_id"
	StartupId bson.ObjectId
	Title     string
	Link      string
	Date      string
}

func urlGetContents(url string) ([]byte, error) {
	// fetch contents from url
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Could not fetch url")
	}

	// read entire contents into []byte
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Could not read contents")
	}

	return body, nil
}

func main() {
	fmt.Println("startupreader!")

	runtime.GOMAXPROCS(runtime.NumCPU())

	// connect with db
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	// clean up connection
	defer session.Close()

	// retrieve collection
	c := session.DB("startupreader").C("startups")
	p := session.DB("startupreader").C("posts")

	// query collection
	startups := []Startup{}
	err = c.Find(
		bson.M{"$and": []bson.M{
			bson.M{"tc_posts": bson.M{"$gt": 1}},
			bson.M{"blog_feed_url": bson.M{"$ne": ""}},
			bson.M{"blog_url": bson.M{"$ne": ""}},
		}}).Sort("-tc_posts").All(&startups)

	if err != nil {
		panic(err)
	}

	// initialize goroutine channel
	ch := make(chan Startup)
	it := 0

	// loop through results
	for _, startup := range startups {

		// fmt.Printf("_Id: %s, Name: %s, BlogURL: %s, BlogFeedUrl: %s\n", startup.Id, startup.Name, startup.Blog_Url, startup.Blog_Feed_Url)

		// validate blog feed url
		var urlValidator = regexp.MustCompile("^http")

		if !urlValidator.MatchString(startup.Blog_Feed_Url) {
			// fmt.Printf("not a valid url")
			continue
		}

		// fire off a goroutine to fetch url
		go func(s Startup, c chan Startup) {

			// build Google Feed API request url
			loadFeedUrl := "https://ajax.googleapis.com/ajax/services/feed/load"

			v := url.Values{}
			v.Set("v", "1.0")
			v.Add("q", s.Blog_Feed_Url)

			apiRequest := loadFeedUrl + "?" + v.Encode()

			// retrieve contents
			body, err := urlGetContents(apiRequest)
			if err != nil {
				c <- s
				return
			}

			// return to channel
			s.Feed = body
			c <- s

		}(startup, ch)

		it++
	}

	// fetch contents from channel
	for i := 0; i < it; i++ {

		// retrieve Startup from channel
		startup := <-ch

		// raw JSON content
		blob := startup.Feed

		if startup.Feed == nil {
			continue
		}

		// ref: https://ajax.googleapis.com/ajax/services/feed/load?v=1.0&q=http%3A%2F%2Fgoogleblog.blogspot.com%2Ffeeds%2Fposts%2Fdefault%3Falt%3Drss

		type Entry struct {
			Title         string
			Link          string
			PublishedDate string
		}

		type Response struct {
			ResponseData struct {
				Feed struct {
					FeedUrl string
					Title   string
					Link    string
					Entries []Entry
				}
			}
		}

		var r Response

		// decode json contents
		err := json.Unmarshal(blob, &r)
		if err != nil {
			// fmt.Println("error:", err)
			continue
		}

		// retrieve data
		feed := r.ResponseData.Feed
		entries := feed.Entries

		if entries == nil || len(entries) == 0 {
			continue
		}

		fmt.Println()
		fmt.Printf("[%s]\n", feed.Title)
		fmt.Println()

		for _, entry := range entries {
			fmt.Printf("- %s\n", entry.Title)
			fmt.Println(entry.Link)
			fmt.Println(entry.PublishedDate)
			fmt.Println()

			// save into database
			go func(s Startup, e Entry) {

				// check if post already exists
				posts := []Post{}
				err = p.Find(
					bson.M{"$and": []bson.M{
						bson.M{"title": e.Title},
						bson.M{"link": e.Link},
						bson.M{"date": e.PublishedDate},
					}}).All(&posts)

				if len(posts) > 0 {
					return
				}

				// insert post into database
				post := Post{
					StartupId: s.Id,
					Title:     e.Title,
					Link:      e.Link,
					Date:      e.PublishedDate,
				}
				err := p.Insert(post)
			}(startup, entry)
		}

		// make sure to clean up feed
		startup.Feed = nil

	}
}
