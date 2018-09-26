package dao

import (
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	mongo "github.com/adcamie/adserver/db/config"
	model "github.com/adcamie/adserver/db/model"
	logger "github.com/adcamie/adserver/logger"
	"gopkg.in/mgo.v2/bson"
)

func InsertToMongoSession(db string, col string, obj interface{}) {
	c := mongo.GetMongoSession().DB(db).C(col)
	err := c.Insert(obj)
	if err != nil {
		fmt.Print("Error while mongo insertion : ", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Insertion Error")
	}
}

func InsertManyToMongo(db string, col string, objects []interface{}) {
	fmt.Println("Inserting docs:", col, len(objects))
	c := mongo.GetMongoSession().DB(db).C(col)
	err := c.Insert(objects...)
	if err != nil {
		fmt.Print("Problem while mulitple insertion :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Multiple Insertion Error:"+"Length-"+strconv.Itoa(len(objects)))
	}

}

func QueryFromMongoSession(db string, col string, filters map[string]interface{}) []model.TrackerModel {
	var results []model.TrackerModel
	err := mongo.GetMongoSession().DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while mulitple query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query Error")
	}
	return results
}

func QueryTrackerLogFromMongoSession(db string, col string, filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := mongo.GetMongoSession().DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while tracker log query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query Tracker Log Error")
	}
	return results
}

func GetAllMessageFromMongoSession(coll string) []model.Message {
	var results []model.Message
	var filters map[string]interface{}
	err := mongo.GetMongoSession().DB(constants.TrackerTemp).C(coll).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while get all query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Get All Messages Error")
	}
	return results
}
func GetAllFromMongoSession(db string, col string, filters map[string]interface{}) ([]bson.M, int) {
	var results []bson.M
	count, _ := mongo.GetMongoSession().DB(db).C(col).Count()
	err := mongo.GetMongoSession().DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while get all query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Get All Error")
	}
	return results, count
}

func GetCountFromMongoSession(db string, col string, filters map[string]interface{}) int {
	count, err := mongo.GetMongoSession().DB(db).C(col).Find(filters).Count()
	if err != nil {
		fmt.Print("Problem while get count query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Get count Error")
	}
	return count
}

func GetAllSortedFromMongoSession(db string, col string, sortField string, filters map[string]interface{}) ([]bson.M, int) {
	var results []bson.M
	count, err := mongo.GetMongoSession().DB(db).C(col).Count()
	if err != nil {
		fmt.Print("Problem while get all sorted query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Get All Count Error")
	}
	err = mongo.GetMongoSession().DB(db).C(col).Find(filters).Sort("-" + sortField).All(&results)
	if err != nil {
		fmt.Print("Problem while get all sorted query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Get All Sorted Error")
	}
	return results, count
}

func DeleteRecordFromMongoSession(db string, col string, filters map[string]interface{}) {
	err := mongo.GetMongoSession().DB(db).C(col).Remove(filters)
	if err != nil {
		fmt.Print("Problem while delete record query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Delete Error")
	}
}

// Mongo function to get all mongo data
func QueryAllFromMongoSession(db string, col string, filters map[string]interface{}) []bson.M {
	var results []bson.M
	err := mongo.GetMongoSession().DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while getting all mongo data query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query All Error")
	}
	return results

}

func QueryAllOffersSession() []model.Offer {
	var results []model.Offer
	err := mongo.GetMongoSession().DB(constants.MongoDB).C(constants.Offer).Find(make(map[string]interface{})).All(&results)
	if err != nil {
		fmt.Print("Problem while get all offer:", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query All Offers Error")
	}
	return results

}

func QueryAllAffiliatesSession() []model.Affiliate {
	var results []model.Affiliate
	err := mongo.GetMongoSession().DB(constants.MongoDB).C(constants.Affiliate).Find(make(map[string]interface{})).All(&results)
	if err != nil {
		fmt.Print("Problem while get all affiliate :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query All Affiliates Error")
	}
	return results
}

//Mongo Function To get paginated data
func QueryAllLogsFromMongoWithOffsetSession(db string, col string, limit int, offset int, sort string, filters map[string]interface{}, fields bson.M) ([]bson.M, int) {
	var results []bson.M
	count, err := mongo.GetMongoSession().DB(db).C(col).Find(filters).Count()
	if err != nil {
		fmt.Print("Problem while get all log with offset :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query All Logs With Offset Count Error")
	}
	err = mongo.GetMongoSession().DB(db).C(col).Find(filters).Select(fields).Sort("-" + sort).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Problem while get all log with offset :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query All Logs With Offset Count Error")
	}
	return results, count
}

func QueryAllTrackerLogsFromMongoSession(db string, col string, sort string, filters map[string]interface{}, fields bson.M) []bson.M {
	var results []bson.M
	err := mongo.GetMongoSession().DB(db).C(col).Find(filters).Select(fields).Sort("-" + sort).All(&results)
	if err != nil {
		fmt.Print("Problem while get all log  :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query All Tracker Logs")
	}
	return results

}

func GetOfferFromMongoSession(filters map[string]interface{}) []model.Offer {
	var results []model.Offer
	log.Println("Filters for get ofer mongo", filters)
	err := mongo.GetMongoSession().DB(constants.MongoDB).C(constants.Offer).Find(filters).Limit(1).All(&results)
	if err != nil {
		fmt.Print("Problem while getting offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Get Offer Error")
	}
	return results
}

func GetAffiliateFromMongoSession(filters map[string]interface{}) (float64, string) {
	var results []model.Affiliate
	err := mongo.GetMongoSession().DB(constants.MongoDB).C(constants.Affiliate).Find(filters).Limit(1).All(&results)
	if err != nil {
		fmt.Print("Problem while get all affiliate  :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Get Affiliate Error")
	}
	if len(results) > 0 {
		return results[0].Mqf, results[0].MediaTemplate
	}
	return 0.7, ""
}

func QueryLatestClickFromMongoSession(filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := mongo.GetMongoSession().DB(constants.MongoDB).C(constants.ClickLog).Find(filters).Limit(1).Sort("-utcdate").All(&results)
	if err != nil {
		fmt.Print("Problem while query latest click  :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Get Latest Click")
	}
	return results
}

//Aggregation queries for reporting
func GetInterceptorReportForAPISession(db string, col string, matchQuery map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}) []bson.M {
	coll := mongo.GetMongoSession().DB(db).C(col)
	var results []bson.M
	idQuery := map[string]interface{}{"_id": idFields}
	for k, v := range sumFields {
		idQuery[k] = v
	}

	pipeline := []bson.M{
		{"$match": matchQuery},
		{"$group": idQuery},
	}
	err := coll.Pipe(pipeline).All(&results)
	if err != nil {
		fmt.Print("Problem while aggregation on report  :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Aggreagtion on Report Error")
	}

	return results
}

//for moving data
func QueryAllFromMongoSessionWithOffset(db string, col string, limit int, offset int, filters map[string]interface{}) []interface{} {
	var results []interface{}
	err := mongo.GetMongoSession().DB(db).C(col).Find(filters).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query All Offset Error")
	}
	return results
}

func QueryFailedTransactionLogsFromMongoSessionWithOffset(limit int, offset int, filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := mongo.GetMongoSession().DB(constants.MongoDB).C(constants.FailedTransactions).Find(filters).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Query Failed Transactions Log Offset Error")
	}
	return results
}
func QueryOldestFilteredLogFromMongoSession(filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := mongo.GetMongoSession().DB(constants.MongoDB).C(constants.FilteredPostBackLog).Find(filters).Limit(1).Sort("utcdate").All(&results)
	if err != nil {
		fmt.Print("Problem while query latest click  :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Get Oldest Postback Log")
	}
	return results
}

func QueryAdcamieReportFromMongoSession(db string, col string, filters map[string]interface{}) []model.APIReport {
	var results []model.APIReport
	err := mongo.GetMongoSession().DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Mongo query adcamie report error :", err.Error())
		go logger.ErrorLogger(err.Error(), "Mongo", "Query Adcamie Report Error")
	}
	return results

}
func GetInterceptorReportForAPIMongoSession(db string, col string, matchQuery map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}) []bson.M {
	coll := mongo.GetMongoSession().DB(db).C(col)
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

func QueryTrackerLogsFromMongoSessionWithOffset(db string, col string, limit int, offset int, filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := mongo.GetMongoSession().DB(db).C(col).Find(filters).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackup", "Query Track Log Offset Error")
	}
	return results
}

func DeleteCollection(db string, col string) {
	fmt.Println("Delete collection :", col)
	var selector interface{}
	result, err := mongo.GetMongoSession().DB(db).C(col).RemoveAll(selector)
	if err != nil {
		fmt.Print("Problem while collection deletion :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Collection Deletion Error")
	}
	fmt.Println("Result is :", result)

}

func DeleteFromMongoSession(db string, col string, filters map[string]interface{}) {
	_, err := mongo.GetMongoSession().DB(db).C(col).RemoveAll(filters)
	if err != nil {
		fmt.Print("Problem while deletion  :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoSessionPool", "Deletion Error")
	}
}
