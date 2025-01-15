package verboten

import (
	_ "embed"
	htmltemplate "html/template"
)

//go:embed webpage.html
var webpage string

var template *htmltemplate.Template = htmltemplate.Must(htmltemplate.New("webpage").Parse(webpage))
