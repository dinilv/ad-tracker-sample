package main

import (
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	helper "github.com/adcamie/adserver/helpers/v1"
	"github.com/robfig/cron"
)

var layout = "02/01/2006"

func main() {

	CheckConversion()
	c := cron.New()
	c.AddFunc("0 55 * * * *", CheckConversion)
	c.Start()
	select {}
}

func CheckConversion() {

	config.InitializeMongoBackup()

	//for formatting date time
	var filters = map[string]interface{}{}
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()
	hour := time.Now().UTC().Hour() - 1

	dateFilter := map[string]interface{}{"$lte": rounded, "$gte": rounded}
	hourFilter := map[string]interface{}{"$lte": hour, "$gte": hour}

	filters["date"] = dateFilter
	filters["hour"] = hourFilter

	//ping postback log for last hour
	results := dao.QueryTrackerLogsFromMongoBackupWithOffset(constants.MongoDB, constants.PostBackLog, 1, 1, filters)
	//if not a single entry initiate alert
	if len(results) == 0 {
		//initiate alert
		sms := helper.NewSMS(4, constants.Camiens)
		sms.Send()
	}
}
