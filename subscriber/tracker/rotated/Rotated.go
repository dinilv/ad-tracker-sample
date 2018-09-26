package main

import (
	"fmt"
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
	router.Initialize(100, 10, common.RotatedSub, InsertBatch)

	//schedule for 5 secs loop
	go router.ProcessingBatch()

	//create a listner for shutdown
	apiRouter := router.InitRoutes("rotated")
	server := &http.Server{
		Addr:    "0.0.0.0:" + os.Args[1],
		Handler: apiRouter,
	}
	log.Println("Listening...")
	err := server.ListenAndServe()
	go logger.ErrorLogger(err.Error(), "GooglePubSub", "Server Creation Failed:")
}

func InsertBatch(messages []*pubsub.Message) {
	var trackerLog []interface{}
	var trackerLogBatch []string
	for _, msg := range messages {

		fmt.Println("TransactionID:-", msg.Attributes[common.TRANSACTION_ID])
		//add messageID as parameter in request map for checking duplicate in mongo
		msg.Attributes[common.MessageID] = msg.ID
		msg.Attributes[common.SubscriptionTime] = time.Now().UTC().String()

		log := processor.SubscribeRotatedClick(msg.Attributes)
		//create transaction keys and pass on to redis for batch
		trackerLogBatch = append(trackerLogBatch, log.TransactionID+common.ObjectSeparator+log.OfferID+common.Separator+log.AffiliateID+common.Separator+log.AffiliateSub+common.Separator+log.AffiliateSub2)
		trackerLog = append(trackerLog, log)

	}
	if len(trackerLog) > 0 {
		dao.SaveClickBatch(trackerLogBatch)
		dao.InsertManyToMongo(common.MongoDB, common.ClickLog, trackerLog)
		dao.InsertManyToMongoBackup(common.MongoDB, common.ClickLog, trackerLog)
	}
}
