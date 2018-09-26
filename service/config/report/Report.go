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
	logger.Print("Starting Tracker Report Microservice:")
	service := micro.NewService(
		micro.Name("go.micro.service.v1.report"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	//initialise mongo
	config.InitializeMongo()
	config.InitializeMongoSessionPool()
	config.InitializeBigQuery()
	service.Init()
	handler.RegisterReportHandler(service.Server(), new(handler.Report))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
