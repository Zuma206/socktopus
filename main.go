package main

import (
	"github.com/joho/godotenv"
	"github.com/zuma206/socktopus/cli"
	"github.com/zuma206/socktopus/web"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	cli.Run()
	router := web.New()

	router.Route("/", func(w web.ResponseWriter, r web.Request) error {
		return w.SendString(200, "Hello, World!")
	})

	router.Listen()
}
