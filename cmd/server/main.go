package main

import (
	"log"

	"github.com/Kalinin-Andrey/redditclone/pkg/config"

	"github.com/Kalinin-Andrey/redditclone/internal/app/api"
)


func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalln("Can not load the config")
	}
	app := api.New(*cfg)

	if err := app.Run(); err != nil {
		log.Fatalf("Error while application is running: %s", err.Error())
	}
}
