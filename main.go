package main

import (
	"github.com/joho/godotenv"
	"github.com/zuma206/socktopus/cli"
	"github.com/zuma206/socktopus/controllers"
	"github.com/zuma206/socktopus/web"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	cli.Run()
	router := web.New()

	router.Route("/", controllers.HandleHome)
	router.Route("/recieve", controllers.HandleRecieve)
	router.Route("/send", controllers.HandleSend)
	router.Route("/kick", controllers.HandleKick)

	router.Listen()
}
