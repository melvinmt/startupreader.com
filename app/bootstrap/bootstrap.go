package bootstrap

import (
	"app/controller"
	"net/http"
)

func init() {
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/fetch/startups", controller.FetchStartups)
}
