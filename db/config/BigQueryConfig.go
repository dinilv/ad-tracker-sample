package config

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
	"github.com/adcamie/adserver/common"
	logger "github.com/adcamie/adserver/logger"
)

var BigQueryClient *bigquery.Client
var BigQueryDataSet *bigquery.Dataset

func InitializeBigQuery() {
	log.Println("Intializing big-query-config")
	ctx := context.Background()
	projectID := common.ProjectName
	var err error
	BigQueryClient, err = bigquery.NewClient(ctx, projectID)
	if err != nil {
		fmt.Println("Could not create bigquery client:", err)
		go logger.ErrorLogger(err.Error(), "BigQuery", "Client Creation Error")
	} else {
		BigQueryDataSet = BigQueryClient.Dataset(common.BigQueryDataSet)
	}

}
