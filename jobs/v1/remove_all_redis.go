package v1

import (
	"log"

	constants "github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
)

import model "github.com/adcamie/adserver/db/model/v1"

func SaveAllTransactionKeys() {
	InitialzeToExternalMongoDB()
	db.InitializeRedisMaster(10)
	results := db.RedisMasterClient.HScan(constants.Transactions, 1, "gg*", 80000)
	keys, cursor := results.Val()
	for {
		log.Println("Length Of Keys", len(keys))
		if len(keys) == 0 || cursor == 0 {
			break
		}
		var objects []interface{}
		log.Println("Length Of Keys", len(keys))
		for key, value := range keys {
			trans := new(model.RedisTransaction)
			trans.Transaction = value
			trans.Page = key
			objects = append(objects, trans)
		}
		MultipleInsertToMongo("Tracker", "RedisKeys", objects)
		results = db.RedisMasterClient.HScan(constants.Transactions, cursor, "gg*", 80000)
		keys, cursor = results.Val()
		log.Println("Next Iteration Start:", cursor)
	}
}
