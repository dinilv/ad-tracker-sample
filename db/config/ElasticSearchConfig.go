package config

import (
	"fmt"
	"log"

	logger "github.com/adcamie/adserver/logger"
	"github.com/olivere/elastic"
)

var ESMasterClient, ESBackupClient *elastic.Client
var esURL = "http://0.0.0.0:9200"

func InitialiseMasterES() {
	fmt.Printf("Creating Client for elastic search\n")
	var err error
	ESMasterClient, err = elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		log.Println("Not able to connect to elastic search", err)
		go logger.ErrorLogger(err.Error(), "Elastic Search", "Client Creation Error")
	}
}

func InitialiseBackupES() {
	fmt.Printf("Creating Client for test elastic search\n")
	var err error
	ESBackupClient, err = elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		log.Println("Not able to connect to elastic search", err)
		go logger.ErrorLogger(err.Error(), "ElasticSearch", "Client Creation Error")
	}
}
