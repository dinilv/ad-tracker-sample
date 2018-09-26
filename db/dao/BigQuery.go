package dao

import (
	"context"
	"fmt"
	"log"

	"time"

	"cloud.google.com/go/bigquery"
	"github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
	helpers "github.com/adcamie/adserver/db/helper"
	"github.com/adcamie/adserver/db/model"
	logger "github.com/adcamie/adserver/logger"
	"google.golang.org/api/iterator"
)

func CreateBigQueryTable(tableID string, model interface{}) bool {
	log.Println("Table creation started")

	ctx := context.Background()
	schema, err := bigquery.InferSchema(model)
	if err != nil {
		log.Println("Error while table creation :", err)
		logger.ErrorLogger(err.Error(), "BigQuery", "Create Schema for BigQueryTable")
		return false
	}
	table := db.BigQueryDataSet.Table(tableID)
	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		log.Println("Error while table creation :", err)
		logger.ErrorLogger(err.Error(), "BigQuery", "Creation for BigQueryTable")
		return false
	}
	log.Println("Table creation ended")
	return true
}

func InsertRecordsToBigQuery(tableID string, records []interface{}) bool {
	fmt.Println("Record insertion started", len(records))
	ctx := context.Background()
	u := db.BigQueryDataSet.Table(tableID).Uploader()
	if err := u.Put(ctx, records); err != nil {
		fmt.Println("Error while record insertion :", err)
		logger.ErrorLogger(err.Error(), "BigQuery", "Record Insertion Error")
		return false
	}
	fmt.Println("Record insertion ended")
	return true
}

func DataAggregationOnBigQuery(hour int, date string, collection string) bool {
	log.Println("Aggregation for table started for collection :", collection)
	ctx := context.Background()
	var query *bigquery.Query
	queryString := helpers.BigQueryBuilder(date, hour, collection)
	query = db.BigQueryClient.Query(queryString)
	query.UseStandardSQL = false
	job, err := query.Run(ctx)
	if err != nil {
		fmt.Println("Error while query run in job :", err)
		go logger.ErrorLogger(err.Error(), "BigQuery", "Query Run Error")
		return false
	}

	status, err := job.Wait(ctx)
	if err != nil {
		fmt.Println("Error while query status in job :", err)
		go logger.ErrorLogger(err.Error(), "BigQuery", "Query Status Error")
		return false
	}

	if err := status.Err(); err != nil {
		fmt.Println("Error while query status in job :", err)
		go logger.ErrorLogger(err.Error(), "BigQuery", "Query Status Error")
		return false
	}

	it, err := job.Read(ctx)
	var results []interface{}
	for {
		var row v1.APIMetaDataReport
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Error while record read :", err)
			go logger.ErrorLogger(err.Error(), "BigQuery", "Record Reading Error")
			return false
		}
		//meta data for deletion
		row.UTCDate = time.Now().UTC()
		//meta data for report creation
		row.Collection = collection
		results = append(results, row)
	}
	if len(results) > 0 {
		if collection == "ClickCookieIDLog" {
			DeleteCollection(common.MongoDB, common.AdcamieClickCookieData)
			InsertManyToMongo(common.MongoDB, common.AdcamieClickCookieData, results)
		} else if collection == "PostbackCookieIDLog" {
			DeleteCollection(common.MongoDB, common.AdcamiePostbackCookieData)
			InsertManyToMongo(common.MongoDB, common.AdcamiePostbackCookieData, results)
		} else {
			InsertManyToMongo(common.MongoDB, common.AdcamieMetaReport, results)
		}

	}
	fmt.Println("Aggregation for table ended")
	return true
}
