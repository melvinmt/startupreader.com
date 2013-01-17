package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Startup struct {
	_Id           bson.ObjectIdHex
	Name          string
	Blog_Url      string
	Blog_Feed_Url string
	Homepage_Url  string
}

func main() {
	fmt.Println("startupreader!")

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("startupreader").C("startups")

	startups := []Startup{}
	err = c.Find(bson.M{"tc_posts": bson.M{"$gt": 1}}).Sort("tc_posts").All(&startups)
	if err != nil {
		panic(err)
	}

	for _, startup := range startups {
		fmt.Printf("_Id: %s, Name: %s, BlogURL: %s, BlogFeedUrl: %s\n", startup._Id, startup.Name, startup.Blog_Url, startup.Blog_Feed_Url)
	}

}
