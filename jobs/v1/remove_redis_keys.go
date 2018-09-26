package v1

import (
	"log"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	"github.com/olivere/elastic"
	"gopkg.in/mgo.v2/bson"
)

func RemoveClickRedisKeys(filters map[string]interface{}) {

	//redis keys to delete
	var fields = bson.M{"transactionID": 1, "cookieID": 1, "offerID": 1, "affiliateID": 1}
	var i = 0
	for {
		log.Println("Click Start Time:-", time.Now())
		i = i + 1
		//get mongo click 10000 per items
		results := dao.QueryMongoForKeys("Tracker", "ClickLog", 40000, i, filters, fields)
		log.Println("Click Results:-", len(results))
		if len(results) == 0 {
			break
		}
		pipeline := config.RedisMasterClient.Pipeline()
		for _, clickLog := range results {
			if clickLog[constants.TransactionID] != nil {
				transactionID := clickLog[constants.TransactionID].(string)
				//sent transaction ids
				pipeline.HDel(constants.SentTransactionIDHash, transactionID)
				//converted transaction ids
				pipeline.HDel(constants.ConvertedTransactionIDHash, transactionID)
				//transaction ids
				pipeline.HDel(constants.Transactions, transactionID)
			}
		}

		//remove 40000 keys
		_, err := pipeline.Exec()
		defer pipeline.Close()
		log.Println("Error", err)

		log.Println("Click End Time:-", time.Now())
	}

}

//delete cookies
func RemoveCookieRedisKeys() {
	pipeline := config.RedisMasterClient.Pipeline()
	pipeline.Del(constants.ClickCookieHash)
	pipeline.Del(constants.ImpressionCookieHash)
	_, err := pipeline.Exec()
	defer pipeline.Close()
	log.Println("Error", err)

	log.Println("Deleted Cookies End Time:-", time.Now())
}

func DeleteESKeys() {
	offerFilter := elastic.NewTermQuery("offerID", "129999")
	dao.DeleteFromES(constants.OfferStack, constants.Offers, offerFilter)
}
