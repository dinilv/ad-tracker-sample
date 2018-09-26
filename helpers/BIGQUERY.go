package v1

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
	dao "github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	logger "github.com/adcamie/adserver/logger"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
)

func UploadToBigQueryTable(collection string) bool {
	log.Println(" Big Query Upload Start for collection :", collection)
	var filePath, backUpFilePath, tableName string
	switch collection {
	case "ClickLog":
		filePath = common.ClickLogPath
		backUpFilePath = common.ClickLogBackUpPath
		tableName = common.BigQueryClickLogTable
	case "PostBackLog":
		filePath = common.PostBackLogPath
		backUpFilePath = common.PostBackLogBackUPPath
		tableName = common.BigQueryPostbackLogTable
	case "PostEventLog":
		filePath = common.PostEventLogPath
		backUpFilePath = common.PostEventLogBackUpPath
		tableName = common.BigQueryPostEventLogTable
	}
	// create table name with month and year
	now := time.Now().UTC()
	monthID := common.MonthIds[now.Month().String()]
	year := now.Year()
	tableName = tableName + "_" + strconv.Itoa(year) + "_" + strconv.Itoa(monthID)

	input, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Println(err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while opening json file"+filePath)
		return false
	}
	//replacements
	input = bytes.Replace(input, []byte(`:{"$date"`), []byte(""), -1)
	input = bytes.Replace(input, []byte(`"},`), []byte(`",`), -1)

	if err = ioutil.WriteFile(common.IntermediateBufferPath+"buffer_3.json", input, 0666); err != nil {
		log.Println(err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while opening json file"+filePath)
		return false
	}

	f, err := os.Open(common.IntermediateBufferPath + "buffer_3.json")
	if err != nil {
		log.Println("Error while file opening", err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while opening json file"+filePath)
		return false
	}
	//To check whether the table already exist or not, if table not exist create the table
	if _, err := db.BigQueryClient.Dataset(common.BigQueryDataSet).Table(tableName).Metadata(context.Background()); err != nil {
		if errorDetail, ok := err.(*googleapi.Error); ok && errorDetail.Code == http.StatusNotFound {
			log.Println(errorDetail.Code)
			dao.CreateBigQueryTable(tableName, model.BigQueryTrackLogBackUp{})
		}
	}
	source := bigquery.NewReaderSource(f)
	source.FileConfig.SourceFormat = bigquery.JSON
	source.FileConfig.AutoDetect = false
	loader := db.BigQueryClient.Dataset(common.BigQueryDataSet).Table(tableName).LoaderFrom(source)
	loader.CreateDisposition = bigquery.CreateNever
	job, err := loader.Run(context.Background())
	if err != nil {
		//create table again with table ID and bigquery model
		//dao.CreateBigQueryTable(tableName, model.BigQueryTrackLogBackUp{})
		//job, err = loader.Run(context.Background())
		//still error log it
		//if err != nil {
		log.Println("Error while big query upload job running:", err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while big query upload job running")
		return false
		//}
	}
	status, err := job.Wait(context.Background())
	if err != nil {
		log.Println("Error while big query upload job waiting:", err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while big query upload job waiting")
		return false
	}
	if err := status.Err(); err != nil {
		log.Println("Error while big query upload job status check:", err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while big query upload job status check")
		return false
	}
	time.Sleep(2 * time.Minute)
	err = os.Rename(common.IntermediateBufferPath+"buffer_3.json", backUpFilePath+time.Now().UTC().String()+".json")
	if err != nil {
		log.Println(err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while moving buffer json to backup folder")
		return false
	}
	err = os.Remove(common.IntermediateBufferPath + "buffer_2.json")
	if err != nil {
		log.Println(err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while removing intermediate buffer_2")
		return false
	}
	err = os.Remove(common.IntermediateBufferPath + "buffer_1.json")
	if err != nil {
		log.Println(err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while removing intermediate buffer_1")
		return false
	}
	err = os.Remove(filePath)
	if err != nil {
		log.Println(err)
		go logger.ErrorLogger(err.Error(), "BigQueryUpload", "Error while removing backup json file")
		return false
	}

	f.Close()
	return true
}
