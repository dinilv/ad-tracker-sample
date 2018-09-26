package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	helper "github.com/adcamie/adserver/helpers/v1"
)

var argumentsHelp = []string{"-d", "Tracker", "-c", "ClickLog", "--sort", "{utcdate:-1}", "-f", "utcdate,date,hour,offerID,affiliateID,transactionID,sessionIP,geo.country,status,clickGeo.country", "--type=csv", "--out=" + "t.csv"}

var mongoExportPath = "/Users/dinilv/Downloads/platform/mongodb-osx-x86_64-3.2.9/bin/mongoexport"

func main3() {

	helper.RunMongoExport(argumentsHelp)

	log.Println("Start writing to CSV:-", time.Now().UTC())

	cmd := exec.Command(mongoExportPath, argumentsHelp...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err.Error())
	}
	log.Println("Done writing to CSV:-", time.Now().UTC())
}
