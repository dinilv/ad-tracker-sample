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
	logger.Print("Starting Tracker Postback Server")
	cmd.Init()
	//start dbs
	config.InitializeRedisTranxn(1000)
	config.InitializeRedisBackup(1000)
	server.Init(
		server.Name("go.micro.api.v1.postback"),
	)
	server.Handle(
		server.NewHandler(
			new(listner.Postback),
		),
	)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
