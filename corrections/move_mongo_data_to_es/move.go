package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
	model "github.com/adcamie/adserver/db/model/v1"
	logger "github.com/adcamie/adserver/logger"
	"github.com/olivere/elastic"
	mgo "gopkg.in/mgo.v2"
)

var offerIDs = []string{}
var completedOffers = []string{}
var ESBackupClient *elastic.Client
var MongoClient *mgo.Session

func init() {

	ESBackupClient, _ = elastic.NewSimpleClient(elastic.SetURL("http://10.148.0.5:9200"))

	MongoClient, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Username: "adcamie",
		Password: "gs#adcamie2017@nov",
		Database: "admin",
	})
	MongoClient.SetSocketTimeout(50 * time.Minute)
	if err != nil {
		log.Println("Not able to connect to mongo", err)
		logger.ErrorLogger(err.Error(), "MongoBackup", "MongoBackup Creation failed")
		go logger.ErrorLogger(err.Error(), "MongoBackup", "Client Creation Error")

	}
}
func main() {
	intializeOfferIDs()
	db.InitializeMongoBackup()
	//loop over offerIDs
	for _, offerID := range offerIDs {
		//for offset incremental
		i := 1
		results := QueryRedisTransactionsFromMongoBackupWithOffset(constants.MongoDB, constants.RedisTransactionBackup, 100000, i, map[string]interface{}{"offerID": offerID})
		bulk := ESBackupClient.Bulk().Index("tracker").Type("redis_keys_backup")

		for {

			log.Println("offerID", offerID)
			log.Println("Length of results:", len(results))
			//check results is size zero
			if len(results) == 0 {
				completedOffers = append(completedOffers, offerID)
				break
			}
			//loop through results
			for _, log := range results {
				if len(log.TransactionID) > 0 {
					bulkRequest := elastic.NewBulkIndexRequest().Id(log.TransactionID).Doc(log)
					bulk.Add(bulkRequest)
				}
			}
			//save bulkRequest
			err := commitToES(bulk)
			if err != nil {
				log.Println(completedOffers)
				log.Println("Breaking as error occureed at:", offerID)
			}
			//increase Offset
			i++
			//query Logs
			results = QueryRedisTransactionsFromMongoBackupWithOffset(constants.MongoDB, constants.RedisTransactionBackup, 100000, i, map[string]interface{}{"offerID": offerID})
		}
	}

}

func commitToES(bulk *elastic.BulkService) error {
	log.Println("Before bulk")
	_, err := bulk.Do(context.Background())
	if err != nil {
		log.Println(err, bulk.NumberOfActions())
		return err
	}

	log.Println("After  Bulk")
	bulk = ESBackupClient.Bulk().Index("tracker").Type("redis_keys_backup")
	return nil
}

func intializeOfferIDs() {

	db.InitializeMongoBackup()
	//query distinct countryIDS from Mongo
	var results = []string{}
	err := db.GetBackupSession().DB(constants.MongoDB).C(constants.RedisTransactionBackup).Find(nil).Distinct(constants.OfferID, &results)
	log.Println("Results", results, err)
	//find results for each country and save to Redis-Master
	for i := 1700; i < 2010; i++ {
		offerID := strconv.Itoa(i)
		if len(offerID) == 4 && strings.Compare("1756", offerID) == 0 && strings.Compare("1925", offerID) == 0 {
			offerIDs = append(offerIDs, offerID)
		}
	}
}

func QueryRedisTransactionsFromMongoBackupWithOffset(db string, col string, limit int, offset int, filters map[string]interface{}) []model.RedisTransactionBackup {
	var results []model.RedisTransactionBackup
	err := MongoClient.DB(db).C(col).Find(filters).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackup", "Query Redis Trasnactions Log Offset Error")
	}
	return results
}
