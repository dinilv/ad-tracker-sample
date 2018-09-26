package dao

import (
	"context"
	"fmt"
	"reflect"

	constants "github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
	model "github.com/adcamie/adserver/db/model"
	logger "github.com/adcamie/adserver/logger"
	"github.com/olivere/elastic"
)

func BulkInsertionToESBackup(bulk *elastic.BulkService) error {
	_, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Print("Search query ES Error :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticBackupMaster", "Bulk Insertion ES Error:"+err.Error())
	}
	return err
}

func SearchRedisKeysFromESBackup(transactionID string) (bool, model.RedisTransactionBackup) {

	idFilter := elastic.NewTermQuery(constants.TransactionID, transactionID)
	results, err := db.ESBackupClient.Search().Index().Index(constants.Tracker).Type(constants.RedisKeysBackup).Query(idFilter).Pretty(false).Do(context.Background())
	if err != nil {
		fmt.Print("Search query ES Error :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Search query ES Error:"+err.Error())
	}
	var ttyp model.RedisTransactionBackup
	if results != nil || results.TotalHits() > 0 {
		for _, item := range results.Each(reflect.TypeOf(ttyp)) {
			if t, ok := item.(model.RedisTransactionBackup); ok {
				fmt.Println("Got Offer:" + t.OfferID)
				return true, t
			}
		}
	}
	return false, ttyp
}

func DeleteFromESBackup(index string, dbtype string, query elastic.Query) {
	_, err := db.ESBackupClient.DeleteByQuery().Index(index).Type(dbtype).Query(query).Do(context.Background())
	if err != nil {
		fmt.Print("ElasticSearchMaster error on deletion :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Deletion Error:"+err.Error())
	}

}
