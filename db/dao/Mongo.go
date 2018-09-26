package dao

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/config"
	model "github.com/adcamie/adserver/db/model"
	logger "github.com/adcamie/adserver/logger"
	"gopkg.in/mgo.v2/bson"
)

func InsertToMongo(db string, col string, obj interface{}) {

	c := config.MongoSession.DB(db).C(col)
	err := c.Insert(obj)
	if err != nil {
		fmt.Print("Mongo Insertion error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Insertion Error")
	}
}

func UpdateToMongo(db string, col string, upsert bson.M, filters bson.M) {
	c := config.MongoSession.DB(db).C(col)
	err := c.Update(filters, upsert)
	if err != nil {
		fmt.Print("Mongo Update error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Update Error")
	}
}

func QueryFromMongo(db string, col string, filters map[string]interface{}) []model.TrackerModel {
	var results []model.TrackerModel
	err := config.MongoSession.DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Mongo Query error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query From Mongo Error")
	}
	return results
}

func QueryTrackerLogFromMongo(db string, col string, filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := config.MongoSession.DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Mongo query tracker log error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query Tracker Log Error")
	}
	return results
}

func GetAllFromMongo(db string, col string, filters map[string]interface{}) ([]bson.M, int) {
	var results []bson.M
	count, err := config.MongoSession.DB(db).C(col).Count()
	if err != nil {
		fmt.Print("Mongo get all error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Count Error")
	}
	err = config.MongoSession.DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Mongo get all error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Get All Mongo Error")
	}
	return results, count
}
func GetAllSortedFromMongo(db string, col string, sortField string, filters map[string]interface{}) ([]bson.M, int) {
	var results []bson.M
	count, err := config.MongoSession.DB(db).C(col).Count()
	if err != nil {
		fmt.Print("Mongo get all sorted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Count Error")
	}
	err = config.MongoSession.DB(db).C(col).Find(filters).Sort("-" + sortField).All(&results)
	if err != nil {
		fmt.Print("Mongo get all sorted error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Get All Mongo Sorted Error")
	}
	return results, count
}

func DeleteRecordFromMongo(db string, col string, filters map[string]interface{}) {
	_, err := config.MongoSession.DB(db).C(col).RemoveAll(filters)
	if err != nil {
		fmt.Print("Mongo delete record error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Remove Error")
	}
}

// Mongo function to get all mongo data
func QueryAllFromMongo(db string, col string, filters map[string]interface{}) []bson.M {
	var results []bson.M
	err := config.MongoSession.DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query All Error")
	}
	return results
}

func QueryAdcamieReportFromMongo(db string, col string, filters map[string]interface{}) []model.APIReport {
	var results []model.APIReport
	err := config.MongoSession.DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Mongo query adcamie report error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query Adcamie Report Error")
	}
	return results

}

func QueryAllFromMongoWithOffset(db string, col string, limit int, offset int, filters map[string]interface{}) []interface{} {
	var results []interface{}
	err := config.MongoSession.DB(db).C(col).Find(filters).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query All Offset Error")
	}
	return results
}

func QueryAllOffers() []model.Offer {
	var results []model.Offer
	err := config.MongoSession.DB(constants.MongoDB).C(constants.Offer).Find(make(map[string]interface{})).All(&results)
	if err != nil {
		fmt.Print("Mongo query all offer error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query All Offers Error")
	}
	return results
}

func QueryMongoForKeys(db string, col string, limit int, offset int, filters map[string]interface{}, fields bson.M) []bson.M {
	var results []bson.M
	err := config.MongoSession.DB(db).C(col).Find(filters).Select(fields).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query for keys error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query Mongo For Keys Error")
	}
	return results
}

func QueryAllAffiliates() []model.Affiliate {
	var results []model.Affiliate
	err := config.MongoSession.DB(constants.MongoDB).C(constants.Affiliate).Find(make(map[string]interface{})).All(&results)
	if err != nil {
		fmt.Print("Mongo query all affiliate error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query All Affiliates Error")
	}
	return results

}

//Mongo Function To get paginated data
func QueryAllLogsFromMongoWithOffset(db string, col string, limit int, offset int, sort string, filters map[string]interface{}, fields bson.M) ([]bson.M, int) {
	var results []bson.M
	err := config.MongoSession.DB(db).C(col).Find(filters).Select(fields).Sort("-" + sort).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query All Logs from mongo with Offset Error")
	}
	return results, 50
}

func QueryCountFromLogs(db string, col string, filters map[string]interface{}) int {
	count, err := config.MongoSession.DB(db).C(col).Find(filters).Count()
	if err != nil {
		fmt.Print("Mongo query count logs error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query Count Logs Error")
	}
	return count
}

func QueryAllLogsFromMongoWithOffsetForKeys(db string, col string, limit int, offset int, sort string, filters map[string]interface{}, fields bson.M) []bson.M {
	var results []bson.M
	err := config.MongoSession.DB(db).C(col).Find(filters).Select(fields).Sort("-" + sort).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all logs with offset for keys error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query All Logs Offset Keys Error")
	}
	return results
}

func QueryAllLogsFromMongoWithOffsetForKeysWithoutSort(db string, col string, limit int, offset int, filters bson.M, fields bson.M) []bson.M {
	var results []bson.M
	err := config.MongoSession.DB(db).C(col).Find(filters).Select(fields).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all logs with offset for keys error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query All Logs Offset Keys Error")
	}
	return results
}
func QueryAllTrackerLogsFromMongo(db string, col string, sort string, filters map[string]interface{}, fields bson.M) []bson.M {
	var results []bson.M
	err := config.MongoSession.DB(db).C(col).Find(filters).Select(fields).Sort("-" + sort).All(&results)
	if err != nil {
		fmt.Print("Mongo query all tracker logs error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query All Tracker Logs Error")
	}
	return results

}

func GetOfferFromMongo(filters map[string]interface{}) []model.Offer {
	var results []model.Offer
	log.Println("Filters for get ofer mongo", filters)
	err := config.MongoSession.DB(constants.MongoDB).C(constants.Offer).Find(filters).Limit(1).All(&results)
	if err != nil {
		fmt.Print("Mongo query all offer error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Get Offer Error")
	}
	return results
}

func GetAffiliateFromMongo(filters map[string]interface{}) (float64, string) {
	var results []model.Affiliate
	err := config.MongoSession.DB(constants.MongoDB).C(constants.Affiliate).Find(filters).Limit(1).All(&results)
	if err != nil {
		fmt.Print("Mongo query all affilaite error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Get Affiliate Error")
	}
	if len(results) > 0 {
		return results[0].Mqf, results[0].MediaTemplate
	}
	return 0.7, ""
}

func QueryLatestClickFromMongo(filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := config.MongoSession.DB(constants.MongoDB).C(constants.ClickLog).Find(filters).Limit(1).All(&results)
	if err != nil {
		fmt.Print("Mongo query latest click error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Latest Click Query Error")
	}
	return results
}

func DeleteFromMongo(db string, col string, filters map[string]interface{}) {
	_, err := config.MongoSession.DB(db).C(col).RemoveAll(filters)
	if err != nil {
		fmt.Print("Mongo query deletion error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Deletion Error")
	}
}

//Aggregation queries for reporting
func GetInterceptorReportForAPI(db string, col string, matchQuery map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}) []bson.M {
	coll := config.MongoSession.DB(db).C(col)
	var results []bson.M
	idQuery := map[string]interface{}{"_id": idFields}
	sortQuery := map[string]interface{}{"hour": -1, "offerID": -1, "affiliateID": -1}
	for k, v := range sumFields {
		idQuery[k] = v
	}
	pipeline := []bson.M{
		{"$match": matchQuery},
		{"$group": idQuery},
		{"$sort": sortQuery},
	}
	err := coll.Pipe(pipeline).All(&results)
	if err != nil {
		fmt.Print("Mongo query interceptor report error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Aggreagtion on Report Error")
	}
	return results
}
