package main

import (
	"github.com/adcamie/adserver/jobs/v1"
	"github.com/jasonlvhit/gocron"
)

func main() {
	gocron.Every(1).Hour().Do(runJobs)
	<-gocron.Start()
}

func runJobs() {
	v1.CreateAdcamieReportFromMongo()
}
