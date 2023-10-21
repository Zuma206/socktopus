package controllers

import (
	"github.com/zuma206/socktopus/templates"
	"github.com/zuma206/socktopus/web"
)

func HandleHome(w web.ResponseWriter, r web.Request) error {
	return templates.Home.Execute(w, nil)
}
