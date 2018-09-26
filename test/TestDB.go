package main

import (
	"context"
	"fmt"
	"log"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/olivere/elastic"
)

func main() {
	//config.InitializeMongoBackup()
	//config.InitializeMongo()
	//config.CreatePool()
	//config.InitialiseSlaveES()

	fmt.Printf("Creating Client for elastic search\n")
	ESMasterClient, err := elastic.NewSimpleClient(elastic.SetURL("http://35.185.178.246:9200"))
	if err != nil {
		log.Println("Not able to connect to elastic search", err)
	}
	wildcard := elastic.NewWildcardQuery("offerID", "1*")
	res, _ := ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(wildcard).From(0*10).Size(10).Sort(constants.UpdatedAt, true).Pretty(false).Do(context.Background())
	log.Println(res.TotalHits())
}
