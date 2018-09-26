package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	redis "github.com/go-redis/redis"
)

func main() {

	db.InitializeMongoSessionPool()

	//initialize redis
	RedisBackupClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "server@123",
		DB:       0,
		PoolSize: 100,
	})

	err := RedisBackupClient.Ping().Err()
	if err != nil {
		fmt.Println("Not able to connect to redis", err)
	}

	//for formatting date time
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()

	eachDayInMilliseconds := rounded.Unix()
	transactionIDTimeStamp := strconv.FormatInt(eachDayInMilliseconds, 10)[0:3]

	transactionPattern := "gg" + transactionIDTimeStamp + "*"

	//query with each pattern on redis live server
	results := RedisBackupClient.HScan(constants.Transactions, 1, transactionPattern, 1000)
	keys, cursor := results.Val()
	var transactions = []string{}
	for {
		//exit critieria
		log.Println("Length Of Keys", len(keys))
		if len(keys) == 0 || cursor == 0 {
			break
		}

		//get all redis keys and save value to redis backup
		for _, transactionID := range keys {
			log.Println("tranxn:-", transactionID)
			filters := map[string]interface{}{constants.TransactionID: transactionID}
			countTransactions := dao.GetCountFromMongoSession(constants.MongoDB, constants.PostBackLog, filters)
			if countTransactions == 0 {
				transactions = append(transactions, transactionID)
			}
		}
		//query next set of keys
		results = RedisBackupClient.HScan(constants.Transactions60, cursor, transactionPattern, 1000)
		keys, cursor = results.Val()
		log.Println("Next Iteration Start:", cursor)
	}

	log.Println("transactions", transactions)

	db.ShutdownMongoSessionPool()
}
