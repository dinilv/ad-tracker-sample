package main

import (
	"log"

	logger "github.com/Sirupsen/logrus"
	listner "github.com/adcamie/adserver/listners/tracker"
	"github.com/micro/go-micro/server"
)

func main() {
	logger.Print("Starting Tracker Setting Server")
	server.Init(
		server.Name("go.micro.api.v1.tracker"),
	)
	server.Handle(
		server.NewHandler(
			new(listner.Tracker),
		),
	)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
