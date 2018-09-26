package main

import (
	"log"

	logger "github.com/Sirupsen/logrus"
	listner "github.com/adcamie/adserver/listners/v1/tracker"
	"github.com/micro/go-micro/server"
)

func main() {
	logger.Print("Starting Tracker Report Server")
	server.Init(
		server.Name("go.micro.api.v1.report"),
	)
	server.Handle(
		server.NewHandler(
			new(listner.Report),
		),
	)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
