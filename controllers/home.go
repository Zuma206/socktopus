package controllers

import (
	"github.com/zuma206/socktopus/models"
	"github.com/zuma206/socktopus/templates"
	"github.com/zuma206/socktopus/utils"
	"github.com/zuma206/socktopus/web"
)

type HomeData struct {
	Count    int
	Services []string
}

func HandleHome(w web.ResponseWriter, r web.Request) error {
	return templates.Home.Execute(w, HomeData{
		Count:    models.DefaultSocketManager.Count(),
		Services: utils.GetSecretNames(),
	})
}
