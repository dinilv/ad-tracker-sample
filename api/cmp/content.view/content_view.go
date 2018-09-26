package main

import (
	"log"

	logger "github.com/Sirupsen/logrus"
	helpers "github.com/adcamie/adserver/helpers"
	listner "github.com/adcamie/adserver/listners/cmp"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/googlepubsub"
)

func main() {
	logger.Print("Starting CMP Content View Server:")
	//commands
	cmd.Init()
	//dbz
	//subscribers
	helpers.InitialiseCMPBrokers()
	//api
	server.Init(
		server.Name("go.micro.api.v1.cv"),
	)
	//server
	server.Handle(
		server.NewHandler(
			new(listner.Cv),
		),
	)
	//start
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
