package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	logger "github.com/adcamie/adserver/logger"
	"github.com/revel/cron"
)

var reg *regexp.Regexp
var topics []string
var brokers = map[string]*pubsub.Topic{}

func initialize() {

	reg = regexp.MustCompile("gg15\\d+")

	topics = []string{constants.ImpressionTopic, constants.ClickTopic, constants.RotatedTopic, constants.PostbackTopic, constants.PostbackPingTopic, constants.DelayedPostbackTopic, constants.FilteredTopic}

	for _, topic := range topics {
		client, err := pubsub.NewClient(context.Background(), constants.ProjectName)
		if err != nil {
			log.Println("Error", err.Error())
			logger.ErrorLogger(err.Error(), "GooglePubSub", "Topic: "+topic+" Initialization Error")
			panic(err.Error())
		}
		topicClient := client.Topic(topic)
		brokers[topic] = topicClient
	}
}
func main() {

	initialize()
	FailedTransaction()
	c := cron.New()
	c.AddFunc("* 45 * * * *", FailedTransaction)
	c.Start()
	select {}
}

func FailedTransaction() {
	config.InitializeMongoBackup()
	config.InitializeMongoSessionPool()

	MoveErrorTransaction()
	MoveFailedTransaction()

	config.ShutdownMongoBackup()
	config.ShutdownMongoSessionPool()
}
func MoveErrorTransaction() {

	var filters = map[string]interface{}{}
	var transactionFilters = map[string]interface{}{}
	var i = 1

	//ping postback log for last hour
	results := dao.QueryErrorTransactionLogsFromMongoBackupWithOffset(100, i, filters)
	for {
		for _, msg := range results {
			//remove google_aid
			_, ok := msg["google_aid"]
			if ok {
				value := msg["google_aid"]
				delete(msg, "google_aid")
				msg["g_aid"] = value
			}
			//subscribe accordingly
			switch msg[constants.Activity] {

			case "2":
				msg["_id"] = ""
				msg["User-Agent"] = ""
				subscribe(constants.ClickTopic, msg)
			case "14", "15", "16", "17", "18", "19", "20", "21", "22", "23":
				msg["_id"] = ""
				msg["User-Agent"] = ""
				subscribe(constants.RotatedTopic, msg)
			default:
				log.Print("Not Processed :")
				logger.ErrorLogger("On Not Processed", "FailedTransactionJob:"+msg[constants.Activity], "Switch case failed")

			}
			transactionFilters[constants.TRANSACTION_ID] = msg[constants.TRANSACTION_ID]
			dao.DeleteFromMongoBackup(constants.MongoDB, constants.ErrorTransaction, transactionFilters)
		}
		i = i + 1
		results = dao.QueryErrorTransactionLogsFromMongoBackupWithOffset(100, i, filters)
		if len(results) == 0 {
			break
		}
	}

}

func MoveFailedTransaction() {

	//create request data of last 3  hour, consider cross over
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()
	hour := time.Now().UTC().Hour()
	hour = hour - 3
	startHour := hour
	startDate := rounded

	if hour < 0 {
		startDate = rounded.AddDate(0, 0, -(1))
		startHour = 23 + hour
	}
	retryMap := map[string]bool{}
	filterDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startHour, 0, 0, 0, today.Location()).UTC()
	dateFilter := map[string]interface{}{"$gte": filterDate}
	log.Println(dateFilter)
	var filters = map[string]interface{}{}
	var i = 1

	//ping postback log for last hour
	dbResults := dao.QueryFailedTransactionLogsFromMongoSessionWithOffset(100, i, filters)
	for {
		for _, tranxn := range dbResults {
			_, ok := retryMap[tranxn.TransactionID]
			//check transactionID pattern
			if reg.MatchString(tranxn.TransactionID) && !ok {
				//add for retry ValidateTransactionID
				retryMap[tranxn.TransactionID] = true
				//query oldest transaction from filtered logs
				filters = map[string]interface{}{}
				filters[constants.TransactionID] = tranxn.TransactionID
				results := dao.QueryOldestFilteredLogFromMongoSession(filters)

				if len(results) > 0 {

					//split transactionID, add affiliateID and offeriD
					offerID := tranxn.TransactionID[12:16]
					affiliateID := tranxn.TransactionID[16:22]

					//check atleast one conversion happened for this offer,media and ip
					transactionFilter := map[string]interface{}{constants.OfferID: offerID, constants.AffiliateID: affiliateID, constants.IP: results[0].IP}
					conversionCount := dao.GetCountFromMongoSession(constants.MongoDB, constants.PostBackLog, transactionFilter)
					if conversionCount > 0 {
						transactionIDFilter := map[string]interface{}{constants.TransactionID: tranxn.TransactionID}
						postbackCount := dao.GetCountFromMongoSession(constants.MongoDB, constants.PostBackLog, transactionIDFilter)
						if postbackCount == 0 {
							//add date fields, offer details and comment
							now := time.Now().UTC()
							results[0].UTCDate = now
							results[0].Date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UTC()
							results[0].Hour = now.Hour()
							results[0].Activity = 4
							results[0].Status = constants.UnSent
							results[0].OfferID = offerID
							results[0].AffiliateID = affiliateID
							results[0].Processor = constants.FAILED_TRANSACTION_JOB
							results[0].Comment = "Processed by failed transaction log"
							dao.InsertToMongoSession(constants.MongoDB, constants.PostBackLog, results[0])
							dao.InsertToMongoBackup(constants.MongoDB, constants.PostBackLog, results[0])
							dao.DeleteFromMongoSession(constants.MongoDB, constants.FailedTransactions, transactionIDFilter)
						}
					} else {
						//check its rotated media or not
						rotatedAffiliateID := tranxn.TransactionID[16:18]
						if strings.Compare(rotatedAffiliateID, constants.TRACKER_MEDIA) == 0 {
							transactionIDFilter := map[string]interface{}{constants.TransactionID: tranxn.TransactionID}
							postbackCount := dao.GetCountFromMongoSession(constants.MongoDB, constants.PostBackLog, transactionIDFilter)
							if postbackCount == 0 {
								//add date fields, offer details and comment
								now := time.Now().UTC()
								results[0].UTCDate = now
								results[0].Date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UTC()
								results[0].Hour = now.Hour()
								results[0].Activity = 8
								results[0].Status = constants.RotatedUnSent
								results[0].OfferID = offerID
								results[0].AffiliateID = rotatedAffiliateID
								results[0].Processor = constants.FAILED_TRANSACTION_JOB
								results[0].Comment = "Processed by failed transaction log"
								dao.InsertToMongoSession(constants.MongoDB, constants.PostBackLog, results[0])
								dao.InsertToMongoBackup(constants.MongoDB, constants.PostBackLog, results[0])
								dao.DeleteFromMongoSession(constants.MongoDB, constants.FailedTransactions, transactionIDFilter)
							}
						}
					}
				}
				i = i + 1
				dbResults = dao.QueryFailedTransactionLogsFromMongoSessionWithOffset(100, i, filters)
				if len(results) == 0 {
					break
				}
			}
		}
	}
}

func subscribe(topic string, msg map[string]string) {
	fmt.Println("Transaction Received on Subscriber:-", msg[constants.TRANSACTION_ID], msg[constants.OFFER_ID], msg[constants.AFF_ID])
	retryMsg := map[string]string{}
	result := brokers[topic].Publish(context.Background(), &pubsub.Message{Attributes: msg})
	serverID, err := result.Get(context.Background())
	if err != nil {
		fmt.Println("Error in Publishing to Google Pub/Sub", serverID, err, err.Error(), msg)
		//remove no length keys
		for key, value := range msg {
			if len(key) != 0 && len(value) != 0 {
				retryMsg[key] = value
			}
		}
		subscribeRetry(topic, retryMsg)
	}
}

func subscribeRetry(topic string, msg map[string]string) {
	fmt.Println("Transaction Received on Subscriber:-", msg[constants.TRANSACTION_ID], msg[constants.OFFER_ID], msg[constants.AFF_ID])
	result := brokers[topic].Publish(context.Background(), &pubsub.Message{Attributes: msg})
	serverID, err := result.Get(context.Background())
	if err != nil {
		fmt.Println("Error in Publishing to Google Pub/Sub", serverID, err, err.Error(), msg)
		msg[constants.ErrorMessage] = err.Error()
		dao.InsertToMongoBackup(constants.MongoDB, constants.ErrorTransactionOnRetry, msg)
	}
}
