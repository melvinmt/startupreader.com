/*
	@copyright 2012 Melvin Tercan
	@license http://creativecommons.org/licenses/by-nc/3.0/
	@repository https://github.com/melvinmt/startupreader
*/
package bootstrap

import (
	"app/controller"
	"net/http"
)

func init() {
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/fetch/startups", controller.FetchStartups)
}
