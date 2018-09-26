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

	logger.Print("Starting CMP Banner Click Server:")
	//commands
	cmd.Init()
	//db
	config.InitializeRedisSlave(2000)
	config.InitializeRedisTranxn(200)
	config.InitializeRedisMaster(200)
	config.InitialiseMasterES()
	//subscribers
	//helpers.InitialiseCMPBrokers()
	//api
	server.Init(
		server.Name("go.micro.api.v1.click"),
	)
	//server
	server.Handle(
		server.NewHandler(
			new(listner.Click),
		),
	)
	//start
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
