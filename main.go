package main

import (
	"github.com/joho/godotenv"
	"github.com/zuma206/socktopus/cli"
	"github.com/zuma206/socktopus/controllers"
	"github.com/zuma206/socktopus/models"
	"github.com/zuma206/socktopus/web"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	cli.Run()
	router := web.New()

	router.Route("/", func(w web.ResponseWriter, r web.Request) error {
		return w.SendString(200, "<h1>Socktopus</h1>")
	})

	router.Route("/count", func(w web.ResponseWriter, r web.Request) error {
		return w.SendJson(200, models.DefaultSocketManager.Count())
	})

	router.Route("/recieve", controllers.HandleRecieve)
	router.Route("/send", controllers.HandleSend)

	router.Listen()
}
