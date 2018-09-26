package main

import (
	"log"
	"strconv"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	"github.com/olivere/elastic"
	"github.com/revel/cron"
)

var FiveDaysDeletionCollections, FifteeenDaysDeletionCollections, NinetyDaysDeletionCollections, FiveDaysCreatedAtDeletionCollection, FiveDaysDateDeleetionCollection []string
var redisBackup90Limit = 90
var redisBackupDeleteLimit = 150

func init() {
	FiveDaysDeletionCollections = []string{"quota_usage", "GooglePubSub", "AdcamieClickCookieData", "AdcamiePostbackCookieData", "DuplicateMessagePostBackLog", "DuplicatePostBackLogOnSubscriber", "ImpressionLog", "RotatedClickLog", "FilteredImpressionLog", "FilteredClickLog", "FilteredPostBackLog"}
	FifteeenDaysDeletionCollections = []string{"ManualPostbackLog", "AdcamieReport", "FailedTransactions", "PostBackPingLog", "APIReport", "AdcamieMetaReport", "DuplicatePostBackLog", "DuplicateMessagePostBackLog", "DuplicatePostBackLogOnSubscriber"}
	NinetyDaysDeletionCollections = []string{"PostBackLog", "FailedTransactions", "PostEventLog"}
	FiveDaysCreatedAtDeletionCollection = []string{"AdcamieEvents"}
	FiveDaysDateDeleetionCollection = []string{"AdcamieMetaReport"}
}
func main() {
	DeleteKeysFromESBackup()
	c := cron.New()
	c.AddFunc("0 30 0 * * *", ExceuteDeletions)
	c.Start()
	select {}
}
func ExceuteDeletions() {
	//resolve start and end time slot
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()
	fifteenDaysBefore := rounded.AddDate(0, 0, -(15))
	fiveDaysBefore := rounded.AddDate(0, 0, -(5))
	thirtyDaysBefore := rounded.AddDate(0, 0, -(30))
	ninteyDaysBefore := rounded.AddDate(0, 0, -(90))

	//deletion form backup
	dateFilter := map[string]interface{}{"$lt": rounded}
	filters := make(map[string]interface{})
	filters["utcdate"] = dateFilter
	DeleteBackupCollections(filters)

	//five days collections
	dateFilter = map[string]interface{}{"$lt": fiveDaysBefore}
	filters["utcdate"] = dateFilter
	DeletionCollectionFiveDays(filters)

	//fifteen days collections
	dateFilter = map[string]interface{}{"$lt": fifteenDaysBefore}
	filters["utcdate"] = dateFilter
	DeletionCollectionFifteenDays(filters)

	//thirteen days collections
	dateFilter = map[string]interface{}{"$lt": thirtyDaysBefore}
	filters["utcdate"] = dateFilter
	DeletionCollectionThirtyDays(filters)

	//ninety days collections
	dateFilter = map[string]interface{}{"$lt": ninteyDaysBefore}
	filters["utcdate"] = dateFilter
	DeletionCollectionNinetyDays(filters)

	//deleting redis keys on ES backup
	DeleteKeysFromESBackup()
}

//delete backup collection for postbacks and clicks
func DeleteBackupCollections(dateFilter map[string]interface{}) {
	config.InitializeMongoBackup()
	dao.DeleteFromMongoBackup(constants.MongoDB, constants.ClickLog, dateFilter)
	dao.DeleteFromMongoBackup(constants.MongoDB, constants.PostBackLog, dateFilter)
	dao.DeleteFromMongoBackup(constants.MongoDB, constants.PostEventLog, dateFilter)
}

//delete records more than 5 days
func DeletionCollectionFiveDays(filters map[string]interface{}) {
	config.InitializeMongo()
	for _, collection := range FiveDaysDeletionCollections {
		dao.DeleteFromMongo(constants.MongoDB, collection, filters)
	}

}

//delete records more than 15 days
func DeletionCollectionFifteenDays(filters map[string]interface{}) {
	config.InitializeMongo()
	for _, collection := range FifteeenDaysDeletionCollections {
		dao.DeleteFromMongo(constants.MongoDB, collection, filters)
	}

}

//delete records more than 30 days
func DeletionCollectionThirtyDays(filters map[string]interface{}) {
	config.InitializeMongo()
	dao.DeleteFromMongo(constants.MongoDB, constants.ClickLog, filters)

}

//delete records more than 90 days
func DeletionCollectionNinetyDays(filters map[string]interface{}) {
	config.InitializeMongo()
	for _, collection := range NinetyDaysDeletionCollections {
		dao.DeleteFromMongo(constants.MongoDB, collection, filters)
	}

}

func DeleteKeysFromESBackup() {

	//initialize mongo
	config.InitialiseBackupES()

	//for formatting date time
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()

	//create more than 90 days before key pattern one by one
	for i := redisBackup90Limit; i < (redisBackupDeleteLimit + 1); i++ {

		eachDay := rounded.AddDate(0, 0, -(i))
		eachDayInMilliseconds := eachDay.Unix()
		transactionIDTimeStamp := strconv.FormatInt(eachDayInMilliseconds, 10)[0:5]
		transactionPattern := "gg" + transactionIDTimeStamp + "*"
		log.Println("Delete Query Pattern:", transactionPattern)
		regexQuery := elastic.NewRegexpQuery(constants.TransactionID, transactionPattern)
		dao.DeleteFromESBackup(constants.Tracker, constants.RedisKeysBackup, regexQuery)
		log.Println("Done Delete Query Pattern:", transactionPattern)
	}

}
