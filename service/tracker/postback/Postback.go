package main

import (
	"log"
	"time"

	logger "github.com/Sirupsen/logrus"
	"github.com/adcamie/adserver/db/config"
	handler "github.com/adcamie/adserver/handlers/v1/tracker"
	"github.com/micro/go-micro"
)

func main() {
	logger.Print("Starting Tracker Postback Microservice:")
	service := micro.NewService(
		micro.Name("go.micro.service.v1.postback"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	//start DBs
	config.InitializeRedisTranxn(100)
	config.InitializeMongo()
	config.InitializeMongoBackup()
	config.InitialiseBackupES()

	service.Init()
	handler.RegisterTrackerPostbackHandler(service.Server(), new(handler.Postback))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
