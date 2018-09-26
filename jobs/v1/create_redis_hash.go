package v1

import (
	"log"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	"gopkg.in/mgo.v2/bson"
)

func CreateRedisHash() {

	log.Println("Start Time:-", time.Now())
	log.Println("Click Log")

	//recent 30 days log
	now := time.Now()
	then := now.Add(-3 * 24 * 60 * time.Minute)
	start := map[string]interface{}{"$gte": then}
	//hour := map[string]interface{}{"$gte": 14}
	var filters = make(map[string]interface{})
	filters["utcdate"] = start
	//filters["hour"] = hour
	log.Println("Filters", filters)

	var fields = bson.M{"transactionID": 1, "offerID": 1, "affiliateID": 1}

	log.Println("Postback Log")
	i := 0
	for {
		pipeline := db.RedisMasterClient.Pipeline()
		i = i + 1
		results := dao.QueryMongoForKeys("Tracker", "PostBackLog", 40000, i, filters, fields)
		log.Println("Length:", len(results))
		if len(results) == 0 {
			break
		}
		for _, postbackLog := range results {
			//handling nil
			if postbackLog[constants.TransactionID] != nil && postbackLog[constants.Status] != nil {
				//extract transactionID
				transactionID := postbackLog[constants.TransactionID].(string)
				status := postbackLog[constants.Status].(string)
				//extract sent or unsent
				if strings.Compare(constants.Sent, status) == 0 || strings.Compare(constants.RotatedSent, status) == 0 {
					pipeline.HSet(constants.SentTransactionIDHash, transactionID, "1")
				}
				//add to converted transactionIDHash
				pipeline.HSet(constants.ConvertedTransactionIDHash, transactionID, "1")
			}
		}
		_, err := pipeline.Exec()
		defer pipeline.Close()
		log.Println("error", err)
	}

	log.Println("End Time:-", time.Now())

}

func CreateRedisHashNew() {

	log.Println("Start Time:-", time.Now())

	log.Println("For Redis Click Log")

	//recent 30 days log
	now := time.Now()
	then := now.Add(-.3 * 24 * 60 * time.Minute)
	start := map[string]interface{}{"$gte": then}
	var filters = make(map[string]interface{})
	filters["utcdate"] = start

	log.Println("Filters", filters)

	//redis keys to delete
	var fields = bson.M{"transactionID": 1, "offerID": 1, "affiliateID": 1}

	//take click log
	var i = 0
	for {
		pipeline := db.RedisMasterClient.Pipeline()
		i = i + 1
		results := dao.QueryMongoForKeys("Tracker", "ClickLog", 40000, i, filters, fields)
		log.Println("Length:", len(results))
		if len(results) == 0 {
			break
		}
		for _, clickLog := range results {
			//handling nil
			if clickLog[constants.TransactionID] != nil && clickLog[constants.OfferID] != nil && clickLog[constants.AffiliateID] != nil {
				//extract transactionID & add to transactionIDHash
				transactionID := clickLog[constants.TransactionID].(string)
				offerID := clickLog[constants.OfferID].(string)
				affiliateID := clickLog[constants.AffiliateID].(string)
				pipeline.HSet(constants.Transactions, transactionID, offerID+constants.Separator+affiliateID)
			}
		}
		_, err := pipeline.Exec()
		defer pipeline.Close()
		log.Println("error", err)
		log.Println("End Time:-", time.Now())
	}
}
