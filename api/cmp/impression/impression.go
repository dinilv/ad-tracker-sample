package main

import (
	"log"

	logger "github.com/Sirupsen/logrus"
	"github.com/adcamie/adserver/db/config"
	listner "github.com/adcamie/adserver/listners/cmp"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/googlepubsub"
)

func main() {
	logger.Print("Starting CMP Impression Server:")
	//commands
	cmd.Init()
	//db
	config.InitializeRedisSlave(2000)
	//subscribers
	//helpers.InitialiseCMPBrokers()
	//api
	server.Init(
		server.Name("go.micro.api.v1.impression"),
	)
	//server
	server.Handle(
		server.NewHandler(
			new(listner.Impression),
		),
	)
	//start
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
