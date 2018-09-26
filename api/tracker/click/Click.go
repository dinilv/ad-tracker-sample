package main

import (
	"log"

	logger "github.com/Sirupsen/logrus"
	"github.com/adcamie/adserver/db/config"
	listner "github.com/adcamie/adserver/listners/tracker"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/googlepubsub"
)

func main() {
	logger.Print("Starting Tracker Click Server")

	//optional arguments
	cmd.Init()

	//start dependent dbs
	config.InitializeRedisSlave(2000)
	config.InitializeRedisTranxn(200)
	config.InitializeRedisMaster(200)
	config.InitialiseMasterES()

	server.Init(
		server.Name("go.micro.api.v1.click"),
	)

	server.Handle(
		server.NewHandler(
			new(listner.Click),
		),
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
