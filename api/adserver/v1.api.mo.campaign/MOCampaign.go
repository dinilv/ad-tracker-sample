package main

import (
	"log"

	listner "github.com/adcamie/adserver/listners/v1/adserver"

	"github.com/micro/go-micro/server"
)

func main() {

	// Initialise Server
	server.Init(
		server.Name("go.micro.api.v1.mocampaign"),
	)

	// Register Handlers
	server.Handle(
		server.NewHandler(
			new(listner.Mocampaign),
		),
	)

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
