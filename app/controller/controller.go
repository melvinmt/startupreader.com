/*
	@copyright 2012 Melvin Tercan
	@license http://creativecommons.org/licenses/by-nc/3.0/
	@repository https://github.com/melvinmt/startupreader
*/
package controller

import (
	"app/helper/tpl"
	"app/helper/url"
	"app/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, tpl.Render("index"))
}

func FetchStartups(w http.ResponseWriter, r *http.Request) {

	// href := "http://api.crunchbase.com/v/1/companies.js"
	href := "http://f.cl.ly/items/1L3V1b0v453C133v1n1L/companies.js"

	// fetch contents from url
	body, err := url.Fetch(href, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// decode json and map to []model.Startup
	var startups []model.Startup
	json.Unmarshal(body, &startups)

	// loop through startups
	for _, startup := range startups {
		fmt.Fprintf(w, "Name: %s\nPermalink: %s\n\n", startup.Name, startup.Permalink)
	}

}
