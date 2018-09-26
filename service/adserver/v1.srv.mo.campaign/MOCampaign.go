package main

import (
	"log"
	"time"

	handler "github.com/adcamie/adserver/handlers/v1/adserver"

	"github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.v1.mocampaign"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	// optionally setup command line usage
	service.Init()

	//register handler for the service
	handler.RegisterMOCampaignHandler(service.Server(), new(handler.MOCampaign))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
