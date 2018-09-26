package v1

import (
	"log"
	"time"

	"github.com/adcamie/adserver/db/dao"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongoSession *mgo.Session

func ConfigCollection() {

	//configuration collections
	var configCollections = []string{"Advertiser", "Affiliate", "Offer", "OfferAffiliateMQF"}

	for _, collection := range configCollections {

		//delete first then insert
		MultipleDeleteExternalMongo("Tracker", collection, make(map[string]interface{}))

		var i = 0
		for {
			i = i + 1
			//get mongo click 500 per items
			results := dao.QueryAllFromMongoWithOffset("Tracker", collection, 20000, i, make(map[string]interface{}))
			if len(results) == 0 {
				break
			}
			MultipleInsertToMongo("Tracker", collection, results)

		}
	}

}

func MoveCollection(filters map[string]interface{}) {

	var moveCollections = []string{"PostBackLog", "PostEventLog"}
	for _, collection := range moveCollections {

		var i = 0
		for {
			i = i + 1
			results := dao.QueryAllFromMongoWithOffset("Tracker", collection, 20000, i, filters)
			log.Println("Length:", len(results))
			if len(results) == 0 {
				break
			}
			MultipleInsertToMongo("Tracker", collection, results)
		}
		//delete the collection data by filter
		MultipleDeleteToMongo("Tracker", collection, filters)
	}

}

func DeletionCollectionThirtyDays(filters map[string]interface{}) {

	var deletionCollections = []string{"ClickLog"}
	for _, collection := range deletionCollections {
		MultipleDeleteToMongo("Tracker", collection, filters)
	}

}

//deletion record in one day
func DeletionCollectionImmediate(filters map[string]interface{}) {

	var deletionCollections = []string{"ImpressionLog", "RotatedClickLog", "FilteredImpressionLog", "FilteredClickLog", "FilteredPostBackLog"}
	for _, collection := range deletionCollections {
		MultipleDeleteToMongo("Tracker", collection, filters)
	}

}

func DeletionCollectionFifteenDays(filters map[string]interface{}) {

	var deletionCollections = []string{"AdcamieEvents", "AdcamieReport"}
	for _, collection := range deletionCollections {
		MultipleDeleteToMongo("Tracker", collection, filters)
	}

}

func InitialzeToExternalMongoDB() {
	log.Println("External Mongo Conection Intialised")
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"104.198.98.165:27017"},
		Timeout:  60 * time.Second,
		Username: "adcamie",
		Password: "server@123",
	}

	mongoSession, _ = mgo.DialWithInfo(mongoDBDialInfo)
}

func MultipleInsertToMongo(db string, col string, objects []interface{}) {

	log.Println("Inserting docs:", col, len(objects))
	c := mongoSession.DB(db).C(col)
	err := c.Insert(objects[0 : len(objects)-1]...)
	if err != nil {
		log.Print(err.Error())
	}

}

func InsertToExternalMongo(db string, col string, obj interface{}) {

	c := mongoSession.DB(db).C(col)
	err := c.Insert(obj)
	if err != nil {
		log.Print(err.Error())
	}

}

func MultipleDeleteToMongo(db string, col string, filters map[string]interface{}) {

	log.Println("Deleting Internal docs:", col)
	dao.DeleteRecordFromMongo(db, col, filters)
}

func MultipleDeleteExternalMongo(db string, col string, filters map[string]interface{}) {

	log.Println("Deleting External docs:", col)
	mongoSession.DB(db).C(col).RemoveAll(filters)
}

func QueryExternalMongoForKeys(db string, col string, limit int, offset int, sort string, filters map[string]interface{}, fields bson.M) []bson.M {
	var results []bson.M
	mongoSession.DB(db).C(col).Find(filters).Select(fields).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	return results
}

func QueryExternalMongoForTransactions(db string, col string, limit int, offset int, filters map[string]interface{}, fields bson.M) []bson.M {
	var results []bson.M
	mongoSession.DB(db).C(col).Find(filters).Select(fields).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	return results
}
