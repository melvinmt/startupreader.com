package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"regexp"
)

type Startup struct {
	Id            bson.ObjectId "_id"
	Name          string
	Blog_Url      string
	Blog_Feed_Url string
	Homepage_Url  string
}

func main() {
	fmt.Println("startupreader!")

	// connect with db
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	// clean up connection
	defer session.Close()

	// retrieve collection
	c := session.DB("startupreader").C("startups")

	// query collection
	startups := []Startup{}
	err = c.Find(
		bson.M{"$and": []bson.M{
			bson.M{"tc_posts": bson.M{"$gt": 1}},
			bson.M{"blog_feed_url": bson.M{"$ne": ""}},
			bson.M{"blog_url": bson.M{"$ne": ""}},
		}}).Sort("-tc_posts").Limit(10).All(&startups)

	if err != nil {
		panic(err)
	}

	// initialize goroutine channel
	ch := make(chan string)
	it := 0

	// loop through results
	for _, startup := range startups {

		fmt.Printf("_Id: %s, Name: %s, BlogURL: %s, BlogFeedUrl: %s\n", startup.Id, startup.Name, startup.Blog_Url, startup.Blog_Feed_Url)

		// validate blog feed url
		var urlValidator = regexp.MustCompile("^http")

		if !urlValidator.MatchString(startup.Blog_Feed_Url) {
			fmt.Printf("not a valid url")
			continue
		}

		// fire off a goroutine to fetch url
		go func(blogFeedUrl string, ch chan string) {

			// build Google Feed API request
			loadFeedUrl := "https://ajax.googleapis.com/ajax/services/feed/load"

			v := url.Values{}
			v.Set("v", "1.0")
			v.Add("q", blogFeedUrl)

			apiRequest := loadFeedUrl + "?" + v.Encode()

			// fetch contents from url
			resp, err := http.Get(apiRequest)
			if err != nil {
				ch <- err
				return
			}

			// read entire contents into []byte
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				ch <- err
				return
			}

			// write bytes into buffer
			buf := bytes.NewBuffer(bodyBytes)

			// convert buffer to string
			bodyStr := buf.String()

			// return to channel
			ch <- bodyStr

		}(startup.Blog_Feed_Url, ch)

		it++
	}

	// fetch contents from channel
	for i := 0; i < it; i++ {
		fmt.Printf("%s", <-ch)
	}
}
