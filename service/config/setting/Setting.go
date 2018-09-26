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
	logger.Print("Starting Tracker Setting Microservice:")
	service := micro.NewService(
		micro.Name("go.micro.service.v1.tracker"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	//start redis
	config.InitializeRedisMaster(100)
	config.InitializeRedisTranxn(100)
	config.InitializeMongo()
	config.InitialiseMasterES()
	config.InitializeMongoSessionPool()

	service.Init()
	handler.RegisterTrackerHandler(service.Server(), new(handler.Tracker))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
