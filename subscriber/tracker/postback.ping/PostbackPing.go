package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	logger "github.com/adcamie/adserver/logger"
	processor "github.com/adcamie/adserver/processors/tracker"
	router "github.com/adcamie/adserver/subscriber"
)

func main() {

	config.InitializeMongoSessionPool()

	//initialize subscriber options
	router.Initialize(10, 5, common.PostbackPingSub, InsertBatch)

	//schedule for 5 secs loop
	go router.ProcessingBatch()

	//create a listner for shutdown
	apiRouter := router.InitRoutes("postback-ping")
	server := &http.Server{
		Addr:    "0.0.0.0:" + os.Args[1],
		Handler: apiRouter,
	}
	log.Println("Listening...")
	err := server.ListenAndServe()
	go logger.ErrorLogger(err.Error(), "GooglePubSub", "Server Creation Failed:")
}

func InsertBatch(messages []*pubsub.Message) {

	var postbackPingLogs []interface{}
	for _, msg := range messages {
		msg.Attributes[common.MessageID] = msg.ID
		msg.Attributes[common.SubscriptionTime] = time.Now().UTC().String()
		log := processor.SubscribePostbackPing(msg.Attributes)
		postbackPingLogs = append(postbackPingLogs, log)
	}
	dao.InsertManyToMongo(common.MongoDB, common.PostBackPingLog, postbackPingLogs)

}
