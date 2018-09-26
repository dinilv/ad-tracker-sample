package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
	logger "github.com/adcamie/adserver/logger"
	redis "github.com/go-redis/redis"
	"github.com/olivere/elastic"
	"github.com/revel/cron"
)

var layout = "02/01/2006"
var redisKeyLimit int64 = 50000
var redisLiveLimit = 10
var redisBackup60Limit = 60
var redisBackup90Limit = 90
var redisBackupDeleteLimit = 150

func main() {

	ExecuteRedisDeleteJob()
	c := cron.New()
	c.AddFunc("0 0 17 * * *", ExecuteRedisDeleteJob)
	c.Start()
	select {}
}

func ExecuteRedisDeleteJob() {
	MovingKeysFromLiveRedis()
	RemoveKeysFromLiveRedis()

}

func MovingKeysFromLiveRedis() {

	//initialize redis
	RedisTranxnClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "server@123",
		DB:       0,
		PoolSize: 1000,
	})

	db.InitialiseBackupES()

	//for formatting date time
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()

	//create 15-60 days before key pattern one by one
	for i := redisLiveLimit; i < (redisBackup60Limit + 1); i++ {
		eachDay := rounded.AddDate(0, 0, -(i))
		eachDayInMilliseconds := eachDay.Unix()
		transactionIDTimeStamp := strconv.FormatInt(eachDayInMilliseconds, 10)[0:5]
		transactionPattern := "gg" + transactionIDTimeStamp + "*"
		log.Println("Transaction Pattern on Live to redisBackup60Limit :", transactionPattern)
		//query with each pattern on redis live server
		results := RedisTranxnClient.HScan(constants.Transactions, 1, transactionPattern, redisKeyLimit)
		keys, cursor := results.Val()
		for {
			//exit critieria
			log.Println("Length Of Keys", len(keys))
			if len(keys) == 0 {
				break
			}
			//for saving and removing keys
			pipelineTranxn := RedisTranxnClient.Pipeline()
			var redisLogs []interface{}
			bulk := db.ESBackupClient.Bulk().Index(constants.Tracker).Type(constants.RedisKeysBackup)
			//get all redis keys and save value to redis backup
			for _, transactionID := range keys {

				pipelineTranxn.HDel(constants.Transactions, transactionID)

				transactionValue := RedisTranxnClient.HGet(constants.Transactions, transactionID).Val()
				transactionsSplitted := strings.Split(transactionValue, constants.Separator)
				redisLog := &model.RedisTransactionBackup{}
				redisLog.TransactionID = transactionID
				redisLog.UTCDate = time.Now().UTC()
				if len(transactionsSplitted) > 0 {
					redisLog.OfferID = transactionsSplitted[0]
				}
				if len(transactionsSplitted) > 1 {
					redisLog.AffiliateID = transactionsSplitted[1]
				}
				if len(transactionsSplitted) > 2 {
					redisLog.AffiliateSub = transactionsSplitted[2]
				}
				if len(transactionsSplitted) > 3 {
					redisLog.AffiliateSub2 = transactionsSplitted[3]
				}
				if len(transactionsSplitted) > 0 && len(transactionID) > 0 {
					redisLogs = append(redisLogs, redisLog)
					idES := redisLog.TransactionID[2:len(redisLog.TransactionID)]
					bulkRequest := elastic.NewBulkIndexRequest().Id(idES).Doc(redisLog)
					bulk.Add(bulkRequest)
				}
			}
			//execute data saving and deleting
			err := dao.BulkInsertionToESBackup(bulk)
			if err == nil {
				_, err := pipelineTranxn.Exec()
				if err != nil {
					fmt.Println("Redis error while deleting redis-keys :", err.Error())
					go logger.ErrorLogger(err.Error(), "RedisTranxn", "Deleting redis keys on jobs.(15-60 Days)")
					break
				}

			}
			if err != nil {
				break
			}
			//close the pipelines
			defer pipelineTranxn.Close()

			//query next set of keys
			results = RedisTranxnClient.HScan(constants.Transactions, cursor, transactionPattern, redisKeyLimit)
			keys, cursor = results.Val()
			log.Println("Next Iteration Start:", cursor)
			if cursor == 0 {
				log.Println("Cursor is zero.")
				break
			}
		}
	}

}

func MoveWithinMongoBackup() {

	//initialize redis
	db.InitializeMongoBackup()

	//for formatting date time
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()

	//create 60-90 days before key pattern one by one
	for i := redisBackup60Limit; i < (redisBackup90Limit + 1); i++ {
		eachDay := rounded.AddDate(0, 0, -(i))
		eachDayInMilliseconds := eachDay.Unix()
		transactionIDTimeStamp := strconv.FormatInt(eachDayInMilliseconds, 10)[0:5]
		transactionPattern := "gg" + transactionIDTimeStamp + "*"
		log.Println("Transaction Pattern on redisBackup60Limit to redisBackup90Limit :", transactionPattern)
		//query with each pattern on mongo backup server
		results := db.RedisTranxnClient.HScan(constants.Transactions, 1, transactionPattern, redisKeyLimit)
		keys, cursor := results.Val()
		for {
			//exit critieria
			log.Println("Length Of Keys", len(keys))
			if len(keys) == 0 || cursor == 0 {
				break
			}
			var redisLogs []interface{}

			//execute data saving and deleting
			err := dao.InsertManyToMongoBackup(constants.MongoDB, constants.RedisTransactionBackup90, redisLogs)
			if err != nil {
				break
			}

			//query next set of keys
			results = db.RedisTranxnClient.HScan(constants.Transactions, cursor, transactionPattern, redisKeyLimit)
			keys, cursor = results.Val()
			log.Println("Next Iteration Start:", cursor)
		}
	}

}

func RemoveKeysFromLiveRedis() {

	//initialize redis
	RedisTranxnClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "server@123",
		DB:       0,
		PoolSize: 1000,
	})

	//for formatting date time
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()

	//create 90 days before key pattern one by one
	for i := redisBackup90Limit; i < (redisBackupDeleteLimit + 1); i++ {
		eachDay := rounded.AddDate(0, 0, -(i))
		eachDayInMilliseconds := eachDay.Unix()
		transactionIDTimeStamp := strconv.FormatInt(eachDayInMilliseconds, 10)[0:5]
		transactionPattern := "gg" + transactionIDTimeStamp + "*"
		log.Println("Transaction Pattern on Live to redisBackup60Limit :", transactionPattern)
		//query with each pattern on redis live server
		results := RedisTranxnClient.HScan(constants.Transactions, 1, transactionPattern, redisKeyLimit)
		keys, cursor := results.Val()
		for {
			//exit critieria
			log.Println("Length Of Keys", len(keys))
			if len(keys) == 0 || cursor == 0 {
				break
			}
			//for saving and removing keys
			pipelineTranxn := RedisTranxnClient.Pipeline()

			//get all redis keys and save value to redis backup
			for _, transactionID := range keys {

				pipelineTranxn.HDel(constants.ConvertedTransactionIDHash, transactionID)
				pipelineTranxn.HDel(constants.SentTransactionIDHash, transactionID)
			}

			//execute deleting
			_, err := pipelineTranxn.Exec()
			if err != nil {
				fmt.Println("Redis error while deleting redis-keys(Converted & Sent) :", err.Error())
				logger.ErrorLogger(err.Error(), "RedisTranxn", "Deleting redis-keys(Converted & Sent) :")
				break

			}
			if err != nil {
				break
			}
			//close the pipelines
			defer pipelineTranxn.Close()

			//query next set of keys
			results = RedisTranxnClient.HScan(constants.Transactions, cursor, transactionPattern, redisKeyLimit)
			keys, cursor = results.Val()
			log.Println("Next Iteration Start:", cursor)
		}
	}
}

func SaveToMongo() {

	//initialize redis
	RedisTranxnClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "server@123",
		DB:       0,
		PoolSize: 1000,
	})

	db.InitializeMongoBackup()
	transactionPattern := "*"

	//query with each pattern on redis live server
	results := RedisTranxnClient.HScan(constants.ConvertedTransactionIDHash, 1, transactionPattern, redisKeyLimit)
	keys, cursor := results.Val()
	for {
		//exit critieria
		log.Println("Length Of Keys", len(keys))
		if len(keys) == 0 || cursor == 0 {
			break
		}

		var redisLogs []interface{}

		//get all redis keys and save value to redis backup
		for _, transactionID := range keys {
			redisLog := &model.RedisTransactionBackup{}
			redisLog.TransactionID = transactionID
			redisLogs = append(redisLogs, redisLog)

		}
		//execute data saving and deleting
		err := dao.InsertManyToMongoBackup(constants.MongoDB, "ConvertedTransactions", redisLogs)
		if err != nil {
			fmt.Println("Redis error while deleting redis-keys :", err.Error())
			logger.ErrorLogger(err.Error(), "RedisTranxn", "Saving Converted redis keys on jobs.")
		}
		//query next set of keys
		results = RedisTranxnClient.HScan(constants.ConvertedTransactionIDHash, cursor, transactionPattern, redisKeyLimit)
		keys, cursor = results.Val()
		log.Println("Next Iteration Start:", cursor)

	}

	//query with each pattern on redis live server
	results = RedisTranxnClient.HScan(constants.SentTransactionIDHash, 1, transactionPattern, redisKeyLimit)
	keys, cursor = results.Val()
	for {
		//exit critieria
		log.Println("Length Of Keys", len(keys))
		if len(keys) == 0 || cursor == 0 {
			break
		}

		var redisLogs []interface{}

		//get all redis keys and save value to redis backup
		for _, transactionID := range keys {
			redisLog := &model.RedisTransactionBackup{}
			redisLog.TransactionID = transactionID
			redisLogs = append(redisLogs, redisLog)

		}
		//execute data saving and deleting
		err := dao.InsertManyToMongoBackup(constants.MongoDB, "SentTransactions", redisLogs)
		if err != nil {
			fmt.Println("Redis error while deleting redis-keys :", err.Error())
			logger.ErrorLogger(err.Error(), "RedisTranxn", "Saving sent transactions redis keys on jobs.")
		}
		//query next set of keys
		results = RedisTranxnClient.HScan(constants.SentTransactionIDHash, cursor, transactionPattern, redisKeyLimit)
		keys, cursor = results.Val()
		log.Println("Next Iteration Start:", cursor)

	}

}
