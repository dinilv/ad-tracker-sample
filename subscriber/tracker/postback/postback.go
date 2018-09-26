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
	db.InitializeRedisBackup(100)
	db.InitializeMongo()
	db.InitializeMongoBackup()
	db.InitializeRedisTranxn(100)
	db.InitializeMongoSessionPool()
	db.InitialiseBackupES()

	//initialize subscriber options
	router.Initialize(10, 5, common.PostbackSub, InsertBatch)

	//schedule for 5 secs loop
	go router.ProcessingBatch()

	//create a listner for shutdown
	apiRouter := router.InitRoutes("postback")
	server := &http.Server{
		Addr:    "0.0.0.0:" + os.Args[1],
		Handler: apiRouter,
	}
	log.Println("Listening...")
	err := server.ListenAndServe()
	go logger.ErrorLogger(err.Error(), "GooglePubSub", "Server Creation Failed:")
}

func InsertBatch(messages []*pubsub.Message) {

	var postbackLogs []interface{}
	var postEventLogs []interface{}
	var sentPostbackTranxn []string
	var unSentPostbackTranxn []string
	var transactionIDs = map[string]bool{}

	for _, msg := range messages {
		//add messageID as parameter in request map for checking duplicate in mongo
		msg.Attributes[common.MessageID] = msg.ID
		msg.Attributes[common.SubscriptionTime] = time.Now().UTC().String()
		log := processor.SubscribePostback(msg.Attributes)
		transactionID := log.TransactionID
		//create transaction keys and pass on to redis for batch
		switch log.Activity {
		//postbacks
		case 3, 4, 7, 8:
			//check for postback log for better be safe with duplicate messages
			filters := map[string]interface{}{common.MessageID: msg.ID}
			countMessages := dao.GetCountFromMongoSession(common.MongoDB, common.PostBackLog, filters)
			filters = map[string]interface{}{common.TransactionID: transactionID}
			countTransactions := dao.GetCountFromMongoSession(common.MongoDB, common.PostBackLog, filters)
			exist, ok := transactionIDs[transactionID]
			validBackup := dao.ValidateSubscriberTransactionIDOnBackup(transactionID)
			//handling with transactionID
			if countMessages == 0 && countTransactions == 0 && !ok && !exist && !validBackup {
				transactionIDs[transactionID] = true
				//append after all duplicate checks
				postbackLogs = append(postbackLogs, log)
				dao.SavePostbackTransactionOnSubscriber(transactionID)
				switch log.Status {
				case common.Sent, common.RotatedSent:
					//save tranxn and increment count
					sentPostbackTranxn = append(sentPostbackTranxn, log.TransactionID+common.ObjectSeparator+log.OfferID+common.Separator+log.AffiliateID)
				case common.UnSent, common.RotatedUnSent:
					//save tranxn
					unSentPostbackTranxn = append(unSentPostbackTranxn, log.TransactionID+common.ObjectSeparator+log.OfferID+common.Separator+log.AffiliateID)
				}
			}
			if countMessages > 0 {
				dao.InsertToMongo(common.MongoDB, common.DuplicateMesaagePostBackLog, log)
			}

			if countTransactions > 0 {
				dao.InsertToMongo(common.MongoDB, common.DuplicatePostBackLog, log)
			}
			if ok || exist {
				dao.InsertToMongo(common.MongoDB, common.DuplicatePostBackLogOnSubscriber, log)
			}
			if validBackup {
				dao.InsertToMongo(common.MongoDB, common.DuplicatePostBackLogOnRedis, log)
			}

			//postevents
		case 11, 12, 13:
			filters := map[string]interface{}{common.MessageID: msg.ID}
			countMessages := dao.GetCountFromMongoSession(common.MongoDB, common.PostBackLog, filters)
			//handling without transactonID
			if countMessages == 0 {
				//append after all duplicate checks
				postbackLogs = append(postbackLogs, log)
				switch log.Status {
				case common.Sent, common.RotatedSent:
					//save tranxn and increment count
					sentPostbackTranxn = append(sentPostbackTranxn, log.TransactionID+common.ObjectSeparator+log.OfferID+common.Separator+log.AffiliateID)
				case common.UnSent, common.RotatedUnSent:
					//save tranxn
					unSentPostbackTranxn = append(unSentPostbackTranxn, log.TransactionID+common.ObjectSeparator+log.OfferID+common.Separator+log.AffiliateID)
				}
			}

			if countMessages > 0 {
				dao.InsertToMongo(common.MongoDB, common.DuplicateMesaagePostBackLog, log)
			}

		//postevents
		case 5, 6, 9, 10:
			postEventLogs = append(postEventLogs, log)

		}
	}
	if len(postbackLogs) > 0 {
		dao.SaveSentPostBackBatch(sentPostbackTranxn)
		dao.SaveUnSentPostBackBatch(unSentPostbackTranxn)
		dao.InsertManyToMongo(common.MongoDB, common.PostBackLog, postbackLogs)
		dao.InsertManyToMongoBackup(common.MongoDB, common.PostBackLog, postbackLogs)
	}
	if len(postEventLogs) > 0 {
		dao.InsertManyToMongo(common.MongoDB, common.PostEventLog, postEventLogs)
		dao.InsertManyToMongoBackup(common.MongoDB, common.PostEventLog, postEventLogs)
	}
}
