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
	logger.Print("Starting Tracker Logs Microservice:")
	service := micro.NewService(
		micro.Name("go.micro.service.v1.log"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	service.Init()
	handler.RegisterLogHandler(service.Server(), new(handler.Log))

	//start DBs
	config.InitializeMongo()
	config.InitializeMongoSessionPool()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
