package url

import (
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"net/http"
)

/*
	Retrieve client from Appengine Context
*/
func getClient(r *http.Request) *http.Client {
	// retrieve appengine context
	c := appengine.NewContext(r)

	// create urlfetch client 
	return urlfetch.Client(c)
}

/* 
	Fetch body contents from url
*/
func Fetch(url string, r *http.Request) ([]byte, error) {

	// fetch url w/ client
	client := getClient(r)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	// make sure to finally close Body
	defer resp.Body.Close()

	// read response to []byte body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
