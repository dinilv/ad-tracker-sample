package v1

import (
	"fmt"
	"log"

	"github.com/olivere/elastic"
)

func ElasticSearchHealthCheck() bool {

	fmt.Printf("Creating Client for test elastic search\n")
	_, err := elastic.NewSimpleClient(elastic.SetURL("http://10.148.0.3:9200"))
	if err != nil {
		log.Println("Not able to connect to elastic search", err)
		return true
	}
	return false
}
