package templates

import (
	_ "embed"
	"html/template"
	"log"
)

//go:embed home.html
var home string
var Home *template.Template

func init() {
	var err error
	Home, err = template.New("home").Parse(home)
	if err != nil {
		log.Fatal(err)
	}
}
