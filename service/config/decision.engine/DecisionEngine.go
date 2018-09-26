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
	logger.Print("Starting Tracker Decision Engine Microservice:")
	service := micro.NewService(
		micro.Name("go.micro.service.v1.engine"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	//start dependent dbs
	config.InitializeRedisMaster(100)
	config.InitializeRedisTranxn(100)
	config.InitialiseMasterES()
	config.InitializeMongo()
	config.InitializeMongoSessionPool()
	service.Init()
	handler.RegisterEngineHandler(service.Server(), new(handler.Engine))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
