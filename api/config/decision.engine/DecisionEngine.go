package main

import (
		"log"

		logger "github.com/Sirupsen/logrus"
		listner "github.com/adcamie/adserver/listners/v1/tracker"
		"github.com/micro/go-micro/server"
	)

	func main() {
		logger.Print("Starting Tracker Decision Engine Server")
		server.Init(
			server.Name("go.micro.api.v1.engine"),
		)
		server.Handle(
			server.NewHandler(

				new(listner.Engine),

			),
		)
		if err := server.Run(); err != nil {
			log.Fatal(err)
		}

	}
