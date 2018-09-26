package main

import (
	"log"

	listner "github.com/adcamie/adserver/listners/v1/adserver"
	subscriber "github.com/adcamie/adserver/subscribers/v1"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
)

func main() {

	cmd.Init()

	// Initialise Server
	server.Init(
		server.Name("go.micro.api.v1.adserve"),
		//Checking authentication before serving service
		//server.WrapHandler(listner.NewAuthHandler()),
	)

	// Register Handlers
	server.Handle(
		server.NewHandler(
			new(listner.Adserve),
		),
	)

	// Register Subscribers
	if err := server.Subscribe(
		server.NewSubscriber(
			"go.micro.sub.moadserve", new(subscriber.Moadserve),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
