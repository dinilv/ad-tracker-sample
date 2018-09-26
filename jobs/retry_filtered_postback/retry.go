package main

import (
	"log"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
	"github.com/dghubble/sling"
	"github.com/revel/cron"
	"gopkg.in/mgo.v2/bson"
)

var layout string

func main() {

	layout = "02/01/2006"

	RetryFailedTransactions()
	c := cron.New()
	c.AddFunc("0 40 * * * *", RetryFailedTransactions)
	c.Start()
	select {}
}

func RetryFailedTransactions() {

	retryMap := map[string]bool{}
	config.InitializeMongo()
	config.InitializeRedisTranxn(100)
	config.InitializeRedisBackup(100)

	//create request data of last 3  hour, consider cross over
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()
	hour := time.Now().UTC().Hour()
	hour = hour - 5
	startHour := hour
	startDate := rounded

	if hour < 0 {
		startDate = rounded.AddDate(0, 0, -(1))
		startHour = 23 + hour
	}

	dateFilter := map[string]interface{}{"$gte": startDate}
	hourFilter := map[string]interface{}{"$gte": startHour}

	filters := bson.M{"activity": 36, "transactionID": bson.M{"$regex": bson.RegEx{`gg15.*`, ""}}, "date": dateFilter, "hour": hourFilter}
	fields := bson.M{"receivedPostbackURL": 1, "transactionID": 1, "ip": 1}
	var i = 0
	for {
		i = i + 1
		log.Println("Start time:-", i, time.Now())
		results := dao.QueryAllLogsFromMongoWithOffsetForKeysWithoutSort(constants.MongoDB, "FilteredPostBackLog", 100, i, filters, fields)
		if len(results) == 0 {
			break
		}
		//retry with url if not a converted transaction
		for _, logging := range results {
			transactionID := logging["transactionID"].(string)
			if dao.ValidateTransactionID(transactionID) {
				if !dao.ValidateDelayedTransactionIDOnBackup(transactionID) {
					if !dao.ValidateConvertedTransactionID(transactionID) {
						if logging["receivedPostbackURL"] != nil {
							if len(logging["receivedPostbackURL"].(string)) != 0 {
								exist, ok := retryMap[transactionID]
								if !ok && !exist {
									sucess := bson.M{}
									urlToPing := logging["receivedPostbackURL"].(string) + "&ip=" + logging["ip"].(string) + "&processor=" + constants.RETRY_JOB
									sling.New().Get(urlToPing).ReceiveSuccess(sucess)
									retryMap[transactionID] = true
									dao.SaveRetryPostbackTransaction(transactionID)
								} else {
									log.Println("Already Retried")
								}
							}
						}
					}
				} else {
					log.Println("Converted Transaction")
				}
			} else {
				log.Println("Not a Valid transactionID")
				filtersForDelete := map[string]interface{}{"transactionID": transactionID}

				count := dao.QueryCountFromLogs(constants.MongoDB, constants.FailedTransactions, filtersForDelete)
				if count == 0 {
					dao.InsertToMongo(constants.MongoDB, constants.FailedTransactions, &model.Transactions{
						TransactionID: transactionID,
						UTCDate:       time.Now().UTC(),
						Date:          rounded})
				}
			}

		}
	}

}
