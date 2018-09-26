package dao

import (
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/config"
	model "github.com/adcamie/adserver/db/model"
	logger "github.com/adcamie/adserver/logger"
	"gopkg.in/mgo.v2/bson"
)

func InsertToMongoBackup(db string, col string, obj interface{}) {
	c := config.GetMongoBackupSession().DB(db).C(col)
	err := c.Insert(obj)
	if err != nil {
		fmt.Print("Error while mongo insertion : ", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Insertion Error")
	}
}

func InsertManyToMongoBackup(db string, col string, objects []interface{}) error {
	fmt.Println("Inserting docs:", col, len(objects))
	c := config.GetMongoBackupSession().DB(db).C(col)
	err := c.Insert(objects...)
	if err != nil {
		fmt.Print("Problem while mulitple insertion :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Multiple Insertion Error:"+"Length-"+strconv.Itoa(len(objects)))
	}
	return err
}

func QueryFromMongoBackup(db string, col string, filters map[string]interface{}) []model.TrackerModel {
	var results []model.TrackerModel
	err := config.GetMongoBackupSession().DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while mulitple query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Query Error")
	}
	return results
}

func QueryTrackerLogFromMongoBackup(db string, col string, filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := config.GetMongoBackupSession().DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while tracker log query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Query Tracker Log Error")
	}
	return results
}

func GetAllMessageFromMongoBackup(coll string) []model.Message {
	var results []model.Message
	var filters map[string]interface{}
	err := config.GetMongoBackupSession().DB(constants.TrackerTemp).C(coll).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while get all query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Get All Messages Error")
	}
	return results
}
func GetAllFromMongoBackup(db string, col string, filters map[string]interface{}) ([]bson.M, int) {
	var results []bson.M
	count, _ := config.GetMongoBackupSession().DB(db).C(col).Count()
	err := config.GetMongoBackupSession().DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while get all query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Get All Error")
	}
	return results, count
}

func GetAllSortedFromMongoBackup(db string, col string, sortField string, filters map[string]interface{}) ([]bson.M, int) {
	var results []bson.M
	count, err := config.GetMongoBackupSession().DB(db).C(col).Count()
	if err != nil {
		fmt.Print("Problem while get all sorted query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Get All Count Error")
	}
	err = config.GetMongoBackupSession().DB(db).C(col).Find(filters).Sort("-" + sortField).All(&results)
	if err != nil {
		fmt.Print("Problem while get all sorted query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Get All Sorted Error")
	}
	return results, count
}

func DeleteRecordFromMongoBackup(db string, col string, filters map[string]interface{}) {
	err := config.GetMongoBackupSession().DB(db).C(col).Remove(filters)
	if err != nil {
		fmt.Print("Problem while delete record query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Delete Error")
	}
}

// Mongo function to get all mongo data
func QueryAllFromMongoBackup(db string, col string, filters map[string]interface{}) []bson.M {
	var results []bson.M
	err := config.GetMongoBackupSession().DB(db).C(col).Find(filters).All(&results)
	if err != nil {
		fmt.Print("Problem while getting all mongo data query :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Query All Error")
	}
	return results

}

func QueryAllOffersMongoBackup() []model.Offer {
	var results []model.Offer
	err := config.GetMongoBackupSession().DB(constants.MongoDB).C(constants.Offer).Find(make(map[string]interface{})).All(&results)
	if err != nil {
		fmt.Print("Problem while get all offer:", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Query All Offers Error")
	}
	return results

}

func QueryAllAffiliatesMongoBackup() []model.Affiliate {
	var results []model.Affiliate
	err := config.GetMongoBackupSession().DB(constants.MongoDB).C(constants.Affiliate).Find(make(map[string]interface{})).All(&results)
	if err != nil {
		fmt.Print("Problem while get all affiliate :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Query All Affiliates Error")
	}
	return results
}

//Mongo Function To get paginated data
func QueryAllLogsFromMongoWithOffsetMongoBackup(db string, col string, limit int, offset int, sort string, filters map[string]interface{}, fields bson.M) ([]bson.M, int) {
	var results []bson.M
	count, err := config.GetMongoBackupSession().DB(db).C(col).Find(filters).Count()
	if err != nil {
		fmt.Print("Problem while get all log with offset :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Query All Logs With Offset Count Error")
	}
	err = config.GetMongoBackupSession().DB(db).C(col).Find(filters).Select(fields).Sort("-" + sort).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Problem while get all log with offset :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Query All Logs With Offset Count Error")
	}
	return results, count
}

func QueryAllTrackerLogsFromMongoBackup(db string, col string, sort string, filters map[string]interface{}, fields bson.M) []bson.M {
	var results []bson.M
	err := config.GetMongoBackupSession().DB(db).C(col).Find(filters).Select(fields).Sort("-" + sort).All(&results)
	if err != nil {
		fmt.Print("Problem while get all log  :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Query All Tracker Logs")
	}
	return results

}

func GetOfferFromMongoBackup(filters map[string]interface{}) []model.Offer {
	var results []model.Offer
	log.Println("Filters for get ofer mongo", filters)
	err := config.GetMongoBackupSession().DB(constants.MongoDB).C(constants.Offer).Find(filters).Limit(1).All(&results)
	if err != nil {
		fmt.Print("Problem while getting offer :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Get Offer Error")
	}
	return results
}

func GetAffiliateFromMongoBackup(filters map[string]interface{}) (float64, string) {
	var results []model.Affiliate
	err := config.GetMongoBackupSession().DB(constants.MongoDB).C(constants.Affiliate).Find(filters).Limit(1).All(&results)
	if err != nil {
		fmt.Print("Problem while get all affiliate  :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Get Affiliate Error")
	}
	if len(results) > 0 {
		return results[0].Mqf, results[0].MediaTemplate
	}
	return 0.7, ""
}

func QueryLatestClickFromMongoBackup(filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := config.GetMongoBackupSession().DB(constants.MongoDB).C(constants.ClickLog).Find(filters).Limit(1).Sort("-utcdate").All(&results)
	if err != nil {
		fmt.Print("Problem while query latest click  :", err.Error())
		logger.ErrorLogger(err.Error(), "MongoBackupPool", "Get Latest Click")
	}
	return results
}

func DeleteFromMongoBackup(db string, col string, filters map[string]interface{}) {
	count, err := config.GetMongoBackupSession().DB(db).C(col).RemoveAll(filters)
	fmt.Println("Count of removed:", count)
	if err != nil {
		fmt.Print("Problem while deletion  :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Deletion Error")
	}

}

//Aggregation queries for reporting
func GetInterceptorReportForAPIMongoBackup(db string, col string, matchQuery map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}) []bson.M {
	coll := config.GetMongoBackupSession().DB(db).C(col)
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
		go logger.ErrorLogger(err.Error(), "MongoBackupPool", "Aggreagtion on Report Error")
	}

	return results
}

//for moving data
func QueryAllFromMongoBackupWithOffset(db string, col string, limit int, offset int, filters map[string]interface{}) []interface{} {
	var results []interface{}
	err := config.GetMongoBackupSession().DB(db).C(col).Find(filters).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackup", "Query All Offset Error")
	}
	return results
}

func QueryTrackerLogsFromMongoBackupWithOffset(db string, col string, limit int, offset int, filters map[string]interface{}) []model.TrackLog {
	var results []model.TrackLog
	err := config.GetMongoBackupSession().DB(db).C(col).Find(filters).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackup", "Query Track Log Offset Error")
	}
	return results
}
func QueryRedisTransactionsFromMongoBackupWithOffset(db string, col string, limit int, offset int, filters map[string]interface{}) []model.RedisTransactionBackup {
	var results []model.RedisTransactionBackup
	err := config.GetMongoBackupSession().DB(db).C(col).Find(filters).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackup", "Query Redis Trasnactions Log Offset Error")
	}
	return results
}

func QueryErrorTransactionLogsFromMongoBackupWithOffset(limit int, offset int, filters map[string]interface{}) []map[string]string {
	var results []map[string]string
	err := config.GetMongoBackupSession().DB(constants.MongoDB).C(constants.ErrorTransaction).Find(filters).Limit(limit).Skip(limit * (offset - 1)).All(&results)
	if err != nil {
		fmt.Print("Mongo query all mongo with offset error :", err.Error())
		go logger.ErrorLogger(err.Error(), "MongoBackup", "Query Error Transactions Log Offset Error")
	}
	return results
}
