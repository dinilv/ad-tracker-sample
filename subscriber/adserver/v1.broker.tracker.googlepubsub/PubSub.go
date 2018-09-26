package main

import (
	"log"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-plugins/broker/googlepubsub"
	"github.com/micro/go-micro/broker"
	subscriber "github.com/adcamie/adserver/subscribers/tracker"
	"github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
)

func main() {

	cmd.Init()
	config.InitializeRedis(6000)
	brokerOptions := googlepubsub.ProjectID("adcamie007")
	broker := googlepubsub.NewBroker(brokerOptions)
	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}
	log.Println("Broker is :",broker.String())
	log.Println("Broker is :",broker.Options())
	log.Println("Broker address :",broker.Address())

	//create 25 subscribers to impression Queue
	count := 1
	for i := 0; i < count; i++ {
		CreateSubscriber(common.ImpressionTopic, "test-queue",broker)

	}
	//create 50 subscribers to click Queue
	/*count = 50
	for i := 0; i < count; i++ {
		CreateSubscriber(common.ClickTopic, common.ClickQueue)
	}
	//create 10 subscribers to postback Queue
	count = 10
	for i := 0; i < count; i++ {
		CreateSubscriber(common.PostbackTopic, common.PostbackQueue)
	}
	//create 3 subscribers to filtered Queue
	count = 3
	for i := 0; i < count; i++ {
		CreateSubscriber(common.FilteredTopic, common.FilteredQueue)
	}
	//create 10 subscribers to rotated Queue
	count = 30
	for i := 0; i < count; i++ {
		CreateSubscriber(common.RotatedTopic, common.RotatedQueue)
	}*/
	select {}
}

func CreateSubscriber(topic string, queue string, pubsubBroker broker.Broker)  {

	result, err := pubsubBroker.Subscribe(topic, func(p broker.Publication) error {
		go subscriber.SubscribeImpression(p.Message().Header)
		return nil
	}, )
	if err != nil {
		log.Println(err, topic, queue)
	}
	log.Println(topic,queue,result)
}

