package v1

import (
	"context"
	"log"
	"strconv"
	"time"

	handler "github.com/adcamie/adserver/handlers/v1/tracker"
	listners "github.com/adcamie/adserver/listners/v1/tracker"
	"github.com/micro/go-micro/client"
)

var layout string

func init() {
	layout = "02/01/2006"
}

func CreateAdcamieReportFromMongo() {

	log.Println("Start Time:-", time.Now())
	log.Println("Adcamie Report")

	//create request data of last 2 hour, consider cross over
	params := map[string]string{}
	start_date := time.Now().UTC()
	now := time.Now().UTC().Hour()
	now = now - 2
	if now < 0 {
		params["start_hour"] = strconv.Itoa(now + 23)
		t := start_date.AddDate(0, 0, -1)
		params["start_date"] = t.Format(layout)
		end_date := time.Now().UTC()
		params["end_date"] = end_date.Format(layout)
	} else {
		params["start_hour"] = strconv.Itoa(now)
		params["start_date"] = start_date.Format(layout)
		params["end_date"] = start_date.Format(layout)

	}
	log.Println("Parameters on call")
	log.Println(params)
	getreport := listners.ValidateForReport(params)
	request := client.NewJsonRequest("go.micro.service.v1.report", "Report.Adcamie", getreport)
	response := &handler.GetReportRes{}

	//init report creation
	client.Call(context.Background(), request, response)

}
