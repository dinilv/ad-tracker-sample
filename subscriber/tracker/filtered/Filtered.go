package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	logger "github.com/adcamie/adserver/logger"
	processor "github.com/adcamie/adserver/processors/tracker"
	router "github.com/adcamie/adserver/subscriber"
)

func main() {
	db.InitializeRedisTranxn(100)
	db.InitializeMongoBackup()
	db.InitializeMongoSessionPool()

	//initialize subscriber options
	router.Initialize(20, 10, common.FilteredSub, InsertBatch)

	//schedule for 5 secs loop
	go router.ProcessingBatch()

	//create a listner for shutdown
	apiRouter := router.InitRoutes("filtered")
	server := &http.Server{
		Addr:    "0.0.0.0:" + os.Args[1],
		Handler: apiRouter,
	}
	log.Println("Listening...")
	err := server.ListenAndServe()
	go logger.ErrorLogger(err.Error(), "GooglePubSub", "Server Creation Failed:")
}

func InsertBatch(messages []*pubsub.Message) {
	log.Println("On Insertion of batch", len(messages))
	var filteredPostbackLog []interface{}
	for _, msg := range messages {
		//add messageID as parameter in request map
		msg.Attributes[common.MessageID] = msg.ID
		msg.Attributes[common.SubscriptionTime] = time.Now().UTC().String()
		//send to processor for processing
		log := processor.SubscribeFiltered(msg.Attributes)
		filteredPostbackLog = append(filteredPostbackLog, log)
	}

	//insert to filtered collection
	if len(filteredPostbackLog) > 0 {
		dao.InsertManyToMongo(common.MongoDB, common.FilteredPostBackLog, filteredPostbackLog)
	}
}
