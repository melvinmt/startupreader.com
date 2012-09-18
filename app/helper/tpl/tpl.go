/*
	@copyright 2012 Melvin Tercan
	@license http://creativecommons.org/licenses/by-nc/3.0/
	@repository https://github.com/melvinmt/startupreader
*/
package tpl

import (
	"github.com/hoisie/mustache"
)

/*
	Render mustache template from given path
*/
func Render(path string) string {
	// look for template path in folder tpl/ and with .mustache extension
	return mustache.RenderFile("tpl/" + path + ".mustache")
}
