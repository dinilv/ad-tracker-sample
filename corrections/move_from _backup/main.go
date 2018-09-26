package main

import (
	"log"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	config.InitializeMongo()
	config.InitializeMongoBackup()
	config.InitializeMongoSessionPool()

	//create request data of last 3  hour, consider cross over
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day()-25, 0, 0, 0, 0, today.Location()).UTC()
	startHour := 0
	endHour := 23

	dateFilter := map[string]interface{}{"$gte": rounded}
	hourFilter := map[string]interface{}{"$gte": startHour, "$lte": endHour}

	filters := bson.M{"date": dateFilter, "hour": hourFilter}
	var i = 0
	for {
		log.Println("Start time:-", i, time.Now())
		results := dao.QueryTrackerLogsFromMongoBackupWithOffset(constants.MongoDB, constants.PostBackLog, 100, i, filters)
		if len(results) == 0 {
			break
		}
		i = i + 1
		//retry with url if not a converted transaction
		for _, logging := range results {
			emptyfilters := map[string]interface{}{constants.TransactionID: logging.TransactionID}
			countTransactions := dao.GetCountFromMongoSession(constants.MongoDB, constants.PostBackLog, emptyfilters)
			if countTransactions == 0 {
				logging.UTCDate = today
				logging.Date = rounded
				logging.Hour = today.Hour()
				logging.Processor = "missing_backup_job"
				dao.InsertToMongo(constants.MongoDB, constants.PostBackLog, logging)
				log.Println(logging)

			}
		}
	}
	config.ShutdownMongo()
	config.ShutdownMongoBackup()
	config.ShutdownMongoSessionPool()
}
