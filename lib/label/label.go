package label

import (
	htmltemplate "html/template"
)

type Label struct {
	IconURI htmltemplate.URL
	Text    string
	URI     htmltemplate.URL
}
