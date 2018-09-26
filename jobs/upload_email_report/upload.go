package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	helper "github.com/adcamie/adserver/helpers/v1"
	"github.com/adcamie/gcp-scripts/gcp"
	"github.com/robfig/cron"
)

var months = map[int]string{1: "january", 2: "february", 3: "march", 4: "april", 5: "may", 6: "june", 7: "july", 8: "august", 9: "septemebr", 10: "october", 11: "november", 12: "december"}
var jsonQuery string
var fieldsFormatted string

func init() {
	fieldsFormatted = "utcdate,date,hour,offerID,affiliateID,transactionID,status,receivedOfferID,receivedAffiliateID,clickDate,sessionIP,clickGeo.countryname,conversionIP,geo.countryname,advertiserID,goalID,affiliateSub,affiliateSub2,affiliateSub3,affiliateSub4,affiliatePayout"
}

func main() {
	//schedule to run at the start of the day
	c := cron.New()
	c.AddFunc("0 20 0 * * *", SendPostbackReport)
	c.Start()
	select {}
}

func SendPostbackReport() {

	//resolve start and end time slot
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()
	dayBefore := rounded.AddDate(0, 0, -(1))
	//create query
	dateFilter := map[string]interface{}{"$gte": map[string]interface{}{"$date": dayBefore}, "$lte": map[string]interface{}{"$date": dayBefore}}
	queries := make(map[string]interface{})
	queries["date"] = dateFilter
	//jsonify
	json, _ := json.Marshal(queries)
	jsonQuery = string(json)
	//create report folder by year and month
	year := dayBefore.Year()
	yearString := strconv.Itoa(year)
	month := dayBefore.Month().String()
	day := strconv.Itoa(dayBefore.Day())
	//create report name by day
	reportName := "daily_report_" + day + "_" + strings.ToLower(month) + "_" + yearString + ".csv"
	reportPath := "/tmp/" + reportName
	//access mongo export on backup db
	var arguments = []string{"--port", "27017", "--host", "10.148.0.4", "-d", "Tracker", "-c", "PostBackLog", "-q", jsonQuery, "-f", fieldsFormatted, "--type=csv", "--out=" + reportPath}
	helper.RunMongoExport(arguments)
	file, _ := os.Open(reportPath)
	//upload to cloud storage
	gcp.UploadToCloudStorage(file, constants.GCPReportFolder, yearString+"/"+month, reportName)
	//create email links
	emailLink := "https://storage.googleapis.com/adcamie-daily-reports/" + yearString + "/" + month + "/" + reportName
	//email links to jen,sharon,hannah,eko,janie,may,jeffy & dinil
	email := "Hi Adcamiens, \n\nReport Link: " + emailLink + "\n\n NB: This link is available for life time.\n\nWish you a Happy Office Day"
	subject := "Automated PostBack Report For " + yearString + "-" + month + "-" + day + " is ready."
	senders := []string{"dinil@adcamie.com", "jeffy@adcamie.com",
		"dilip@adcamie.com", "gabriel@adcamie.com", "jennifer@adcamie.com", "hannah@adcamie.com",
		"sharon@adcamie.com", "janie@adcamie.com", "may@adcamie.com", "eko@adcamie.com"}
	//senders = []string{"dinil@adcamie.com"}
	helper.SendEmail(subject, email, senders)

}
