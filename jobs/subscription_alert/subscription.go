package main

import (
	"log"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	helper "github.com/adcamie/adserver/helpers/v1"
	"github.com/robfig/cron"
)

//format      2017-12-15 11:55:02.874446009 +0000 UTC"
var layout = "2006-01-02 15:04:05.000000000 +0000 UTC"

func main() {

	c := cron.New()
	c.AddFunc("0 50 * * * *", CheckSubscription)
	c.Start()
	select {}

}

func CheckSubscription() {
	config.InitializeMongoBackup()

	//for formatting date time
	var filters = map[string]interface{}{}
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()
	hour := time.Now().UTC().Hour()

	dateFilter := map[string]interface{}{"$lte": rounded, "$gte": rounded}
	hourFilter := map[string]interface{}{"$lte": hour, "$gte": hour}

	filters["date"] = dateFilter
	filters["hour"] = hourFilter

	//take one click log for last hour
	results := dao.QueryTrackerLogsFromMongoBackupWithOffset(constants.MongoDB, constants.ClickLog, 1, 1, filters)
	//take the lateset click log
	if len(results) > 0 {
		//convert API time and Subscription TIme
		converterAPITime, _ := time.Parse(layout, results[0].APITime)
		convertedSubscriptionTime, _ := time.Parse(layout, results[0].SubscriptionTime)
		log.Println("converterted API time", converterAPITime)
		log.Println("converterted Subscription time", convertedSubscriptionTime)
		duration := convertedSubscriptionTime.Sub(converterAPITime)
		log.Println("duration", duration.Seconds())
		//Check it more than 6 minutes gap
		if duration.Seconds() > 360 {
			//send sms
			sms := helper.NewSMS(5, constants.Camiens)
			sms.Send()
		}
	}

	//take one postback log for last hour
	results = dao.QueryTrackerLogsFromMongoBackupWithOffset(constants.MongoDB, constants.PostBackLog, 1, 1, filters)
	//take the lateset click log
	if len(results) > 0 {
		//convert API time and Subscription TIme
		converterAPITime, _ := time.Parse(layout, results[0].APITime)
		convertedSubscriptionTime, _ := time.Parse(layout, results[0].SubscriptionTime)
		log.Println("converterted API time", converterAPITime)
		log.Println("converterted Subscription time", convertedSubscriptionTime)
		duration := convertedSubscriptionTime.Sub(converterAPITime)
		log.Println("duration", duration.Seconds())
		//Check it more than 4 minutes gap
		if duration.Seconds() > 360 {
			//send sms
			sms := helper.NewSMS(5, constants.Camiens)
			sms.Send()
		}
	}

}
