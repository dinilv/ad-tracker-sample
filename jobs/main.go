package main

import (
	"log"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"

	"github.com/adcamie/adserver/jobs/v1"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	v1.SaveAllTransactionKeys()
	v1.InitialzeToExternalMongoDB()

}

func runJobs() {

	//intialize DBs
	config.InitializeRedisMaster(1000)
	v1.InitialzeToExternalMongoDB()
	config.InitializeMongo()

	//30 days before log
	now := time.Now()
	then := now.Add(-30 * 24 * 60 * time.Minute)
	start := map[string]interface{}{"$lte": then}
	var filters = make(map[string]interface{})
	filters["utcdate"] = start
	log.Println("Filters", filters)

	//clear cookies keys in redis
	v1.RemoveCookieRedisKeys()

	//clear transactions till latest 30
	v1.RemoveClickRedisKeys(filters)

	//move collections to be backed up as whole and part
	v1.ConfigCollection()

	//move only postbacks log
	v1.MoveCollection(filters)

	//delete click logs
	v1.DeletionCollectionThirtyDays(filters)

	//15 days before log
	now = time.Now()
	then = now.Add(-14 * 24 * 60 * time.Minute)
	start = map[string]interface{}{"$lte": then}
	filters = make(map[string]interface{})
	filters["utcdate"] = start

	//delete events
	v1.DeletionCollectionFifteenDays(filters)

	//1 day before log
	now = time.Now()
	then = now.Add(-1 * 24 * 60 * time.Minute)
	start = map[string]interface{}{"$lte": then}
	filters = make(map[string]interface{})
	filters["utcdate"] = start

	//delete filtered/rotation events
	v1.DeletionCollectionImmediate(filters)

}

func CookieJobToDeleteOldKeys() {

	v1.InitialzeToExternalMongoDB()
	config.InitializeRedisMaster(1000)
	//create filter
	patternArray := []string{"/^gg1507.*/", "/^gg1508.*/", "/^gg1509.*/"}
	notIn := map[string]interface{}{"$nin": patternArray}
	var filters = make(map[string]interface{})
	filters["transaction"] = notIn
	var fields = bson.M{"transaction": 1}
	i := 1
	// query from mongo
	results := v1.QueryExternalMongoForTransactions("Tracker", "RedisKeys", 40000, i, filters, fields)
	log.Println("Result Length:-", len(results))

	for {
		if len(results) == 0 {
			break
		}
		pipeline := config.RedisTranxnClient.Pipeline()
		//extract results
		for _, key := range results {
			transactionID := key["transaction"].(string)
			pipeline.HDel(constants.Transactions, transactionID)
			pipeline.HDel(constants.ConvertedTransactionIDHash, transactionID)
			pipeline.HDel(constants.SentTransactionIDHash, transactionID)
		}
		_, err := pipeline.Exec()
		defer pipeline.Close()
		log.Println("Error", err)
		i = i + 1
		results = v1.QueryExternalMongoForTransactions("Tracker", "RedisKeys", 40000, i, filters, fields)
		log.Println("Result Length:-", len(results))
	}
}
