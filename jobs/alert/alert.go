package main

import (
	constants "github.com/adcamie/adserver/common/v1"
	helper "github.com/adcamie/adserver/helpers/v1"
	"github.com/adcamie/adserver/jobs/v1"
	"github.com/jasonlvhit/gocron"
)

func main() {
	runAlertJobs()
	gocron.Every(5).Minutes().Do(runAlertJobs)
	<-gocron.Start()
	select {}
}
func runAlertJobs() {

	//check Mongo master
	alert := v1.MongoMasterHealthCheck()
	if alert {
		//send sms
		sms := helper.NewSMS(1, constants.Camiens)
		sms.Send()
	}

	//check mongo backup
	alert = v1.MongoBackupHealthCheck()
	if alert {
		//send sms
		sms := helper.NewSMS(6, constants.Camiens)
		sms.Send()
	}

	//check redis master
	alert = v1.RedisMasterHealthCheck()
	if alert {
		//send sms
		sms := helper.NewSMS(0, constants.Camiens)
		sms.Send()
	}
	//check redis tranxn
	alert = v1.RedisTranxnHealthCheck()
	if alert {
		//send sms
		sms := helper.NewSMS(8, constants.Camiens)
		sms.Send()
	}

	//check redis backup
	alert = v1.RedisBackupHealthCheck()
	if alert {
		//send sms
		sms := helper.NewSMS(7, constants.Camiens)
		sms.Send()
	}

	//check es backup
	alert = v1.ElasticSearchHealthCheck()
	if alert {
		//send sms
		sms := helper.NewSMS(9, constants.Camiens)
		sms.Send()
	}

}
