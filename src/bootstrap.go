package bootstrap

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	"io/ioutil"
	"net/http"
)

/*
//	Set up current page environment
type Page struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func (p *Page) init(w http.ResponseWriter, r *http.Request) {
	p.Writer = w
	p.Request = r
}

var page = new(Page)*/

/** ROUTES	 **/

func init() {
	http.HandleFunc("/", controllerIndex)
	http.HandleFunc("/fetch/startups", controllerFetchStartups)
}

/** MODELS **/

type Startup struct {
	Name      string
	Permalink string
}

/** CONTROLLERS **/

func controllerIndex(w http.ResponseWriter, r *http.Request) {
	page.init(w, r)
	fmt.Fprintf(w, renderTpl("index"))
}

func controllerFetchStartups(w http.ResponseWriter, r *http.Request) {

	// url := "http://api.crunchbase.com/v/1/companies.js"
	url := "http://f.cl.ly/items/1L3V1b0v453C133v1n1L/companies.js"

	// fetch contents from url
	body, err := fetchUrl(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// decode json and map to []Startup
	var startups []Startup
	json.Unmarshal(body, &startups)

	// loop through startups
	for _, startup := range startups {
		fmt.Fprintf(w, "Name: %s\nPermalink: %s\n\n", startup.Name, startup.Permalink)
	}

}

/** HELPERS **/

/* 
	Fetch body contents from url
*/
func fetchUrl(url string, r *http.Request) ([]byte, error) {

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
	Render mustache template from given path
*/
func renderTpl(path string) string {
	// look for template path in folder tpl/ and with .mustache extension
	return mustache.RenderFile("tpl/" + path + ".mustache")
}
