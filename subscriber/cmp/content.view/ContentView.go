package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
	db.InitializeRedisTranxn(10)
	db.InitializeBigQuery()
	db.InitializeMongoSessionPool()

	//initialize subscriber options
	router.Initialize(200, 5, common.ContentViewSub, InsertBatch)

	//schedule for 5 secs loop
	go router.ProcessingBatch()
	//create a listner for shutdown
	apiRouter := router.InitRoutes("content-view")
	server := &http.Server{
		Addr:    "0.0.0.0:" + os.Args[1],
		Handler: apiRouter,
	}
	log.Println("Listening...")
	err := server.ListenAndServe()
	go logger.ErrorLogger(err.Error(), "GooglePubSub", "Server Creation Failed:")
}

func InsertBatch(messages []*pubsub.Message) {

	log.Println("Inside batch processing")
	var trackerLog []interface{}
	var filteredTrackerLog []interface{}

	fmt.Println("Size of message received :", len(messages))
	for _, msg := range messages {
		//add message parameters in request map
		msg.Attributes[common.MessageID] = msg.ID
		msg.Attributes[common.SubscriptionTime] = time.Now().UTC().String()

		//send to processor for processing
		logs := processor.SubscribeContentView(msg.Attributes)
		if strings.Compare(logs.Status, common.Filtered) == 0 {
			filteredTrackerLog = append(filteredTrackerLog, logs)
		} else {
			trackerLog = append(trackerLog, logs)
		}

	}
	//insert to filtered collection
	if len(filteredTrackerLog) > 0 {
		dao.InsertManyToMongo(common.MongoDB, common.FilteredContentViewLog, filteredTrackerLog)
	}
	//insert to content view log
	if len(trackerLog) > 0 {
		dao.InsertManyToMongo(common.MongoDB, common.ContentViewLog, trackerLog)
	}
}
