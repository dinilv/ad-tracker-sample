package main

import (
	"context"
	"log"
	"time"

	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	handler "github.com/adcamie/adserver/handlers/v1/tracker"
	helper "github.com/adcamie/adserver/helpers/v1"
	"github.com/micro/go-micro/client"
	"github.com/robfig/cron"
)

var MoveCollections = []string{"PostBackLog", "PostEventLog", "ClickLog"}

//var MoveCollections = []string{"PostBackLog"}
var AggregateCollections = []string{"PostBackLogWithoutTransaction", "PostBackLog", "PostEventLog", "ClickLog"}

//var CookieCollections = []string{"ClickCookieIDLog", "PostbackCookieIDLog"}
var OperationCount = 0
var layoutBigQuery string
var layoutReport string

func init() {
	layoutBigQuery = "2006-01-02"
	layoutReport = "02/01/2006"
}

type Quota struct {
	SlotType  string    `bson:"slotType,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty"`
	Usage     int       `bson:"usage,omitempty"`
}

func main() {
	MoveCollection()
	c := cron.New()
	c.AddFunc("0 15 * * * *", MoveCollection)
	c.Start()
	select {}
}

//big query has limit of 1k per day per table
func MoveCollection() {
	log.Println("Starting moving collection")
	config.InitializeMongo()
	config.InitializeBigQuery()
	config.InitializeMongoSessionPool()
	//moving one by one collections
	for _, collection := range MoveCollections {
		log.Println("Starting moving collection :", collection)
		helper.UploadToBigQueryTable(collection)
	}
	//delay for bigquery  report upload
	/*	time.Sleep(5 * time.Minute)
		today := time.Now().UTC()
		rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()
		hour := time.Now().UTC().Hour() - 1
		if hour < 0 {
			rounded = rounded.AddDate(0, 0, -(1))
			hour = 23
		}
		hourString := strconv.Itoa(hour)

		for _, collection := range AggregateCollections {
			CreateBigQueryReport(hour, rounded.Format(layoutBigQuery), collection)
		}
		//start creating report from bigquery meta data collection
		CreateAdcamieReportFromBigQuery(hourString, rounded.Format(layoutReport))

		//time.Sleep(2 * time.Minute)

		/*for _, collection := range CookieCollections {
			CreateBigQueryReport(hour, rounded.Format(layoutBigQuery), collection)
		}
	*/

	config.ShutdownMongo()
	config.ShutdownMongoSessionPool()
}

func CreateBigQueryReport(hour int, date string, collection string) {
	config.InitializeBigQuery()
	config.InitializeMongoSessionPool()
	log.Println("Aggregation for big query table started for collection :", collection)
	dao.DataAggregationOnBigQuery(hour, date, collection)
	config.ShutdownMongoSessionPool()
}

func CreateAdcamieReportFromBigQuery(start_hour string, start_date string) {
	log.Println("Start Time:-", time.Now())
	log.Println("Adcamie BigQuery Report")

	getreport := new(handler.GetReportReq)
	getreport.StartHour = start_hour
	getreport.StartDate = start_date
	request := client.NewJsonRequest("go.micro.service.v1.report", "Report.AdcamieMetaReportBigQuery", getreport)
	response := &handler.GetReportRes{}

	//init report creation
	client.Call(context.Background(), request, response)
}
