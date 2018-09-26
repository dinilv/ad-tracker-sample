package main

import (
	"log"

	logger "github.com/Sirupsen/logrus"
	listner "github.com/adcamie/adserver/listners/tracker"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/googlepubsub"
)

func main() {
	logger.Print("Starting Tracker Impression Server")
	cmd.Init()
	server.Init(
		server.Name("go.micro.api.v1.impression"),
	)

	server.Handle(
		server.NewHandler(
			new(listner.Impression),
		),
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
