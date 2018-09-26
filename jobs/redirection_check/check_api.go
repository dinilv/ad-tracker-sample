package main

import (
	"log"
	"net/http"

	constants "github.com/adcamie/adserver/common/v1"
	helper "github.com/adcamie/adserver/helpers/v1"
	logger "github.com/adcamie/adserver/logger"
	"github.com/robfig/cron"
)

var TestClickURL, TestPostbackURL string

func init() {
	TestClickURL = "http://track.adcamie.com/aff_c?offer_id=100&aff_id=100"
	TestPostbackURL = "http://track.adcamie.com/aff_lsr?offer_id=100&aff_id=100"
}

func main() {

	c := cron.New()
	c.AddFunc("0 5 * * * *", CheckClickRedirection)
	c.AddFunc("0 10 * * * *", CheckPostback)
	c.Start()
	select {}
}

func CheckClickRedirection() {

	//Check for Click
	response, err := http.Get(TestClickURL)
	if err != nil {
		//log error
		go logger.ErrorLogger(err.Error(), "CheckAPIJob", "Pinging Click URL failed")
		//send sms
		sms := helper.NewSMS(2, constants.Camiens)
		sms.Send()
	} else {
		defer response.Body.Close()
		log.Println(response.StatusCode, "Status")
		//check response code
		if response.StatusCode != 200 {
			//not success initiate alert
			sms := helper.NewSMS(2, constants.Camiens)
			sms.Send()
		}
	}

}

func CheckPostback() {

	//Check for Postback
	response, err := http.Get(TestPostbackURL)
	if err != nil {
		//log error
		go logger.ErrorLogger(err.Error(), "CheckAPIJob", "Pinging PostBack URL failed")
		//send sms
		sms := helper.NewSMS(3, constants.Camiens)
		sms.Send()
	} else {
		defer response.Body.Close()
		//check response code
		if response.StatusCode != 200 {
			//not success initiate alert
			sms := helper.NewSMS(3, constants.Camiens)
			sms.Send()
		}
	}
}
