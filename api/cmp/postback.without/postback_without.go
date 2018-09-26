package main

import (
	"log"

	logger "github.com/Sirupsen/logrus"
	"github.com/adcamie/adserver/db/config"
	helpers "github.com/adcamie/adserver/helpers/v1"
	listner "github.com/adcamie/adserver/listners/cmp"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
)

func main() {
	logger.Print("Starting CMP Postback-Without Transaction Server")
	//commands
	cmd.Init()
	//subscribers
	helpers.InitialiseCMPBrokers()
	//dbs
	config.InitializeRedisMaster(1000)
	//api
	server.Init(
		server.Name("go.micro.api.v1.postback"),
	)
	//server
	server.Handle(
		server.NewHandler(
			new(listner.Postback),
		),
	)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
