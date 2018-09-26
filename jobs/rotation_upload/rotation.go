package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	dao "github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
	"github.com/jasonlvhit/gocron"
)

var countryIDs []string
var countriesMap = make(map[string]string)

func main() {
	intializeDropDown()
	runRotationJobs()
	gocron.Every(30).Minutes().Do(runRotationJobs)
	<-gocron.Start()
	select {}

}

func runRotationJobs() {
	config.InitialiseMasterES()
	config.InitializeRedisMaster(100)
	config.InitializeRedisTranxn(100)
	intializeCountryIDs()
	config.InitializeMongo()
	activeDB := config.InitializeMongoSessionPool()
	if activeDB {
		//clear existing data
		config.RedisMasterClient.Del(constants.RotationOfferGroupStack, constants.RotationOfferGeoStack)
		dao.DeleteRecordFromMongo(constants.MongoDB, constants.RotationGEOStack, map[string]interface{}{})
		dao.DeleteRecordFromMongo(constants.MongoDB, constants.RotationGroupStack, map[string]interface{}{})
		AddCountryRotationForCPI()
		AddGroupRotationForCPI()

	}

	config.ShutdownMongo()
	config.ShutdownMongoSessionPool()
}
func AddCountryRotationForCPI() {
	//saving to Mongo
	config.InitializeMongo()
	//search for each country
	for _, countryID := range countryIDs {
		//search on ES
		offerID := dao.SearchOfferStackWithGeo(countryID)
		dao.SaveRotationOfferByGEO(countryID, offerID)
		if strings.Compare("ALL", countryID) != 0 {
			dao.InsertToMongo(constants.MongoDB, constants.RotationGEOStack, map[string]interface{}{"countryID": countryID, "country": countriesMap[countryID], "selectedOfferForRotation": offerID, "addedDate": time.Now().UTC()})
		} else {
			dao.InsertToMongo(constants.MongoDB, constants.RotationGEOStack, map[string]interface{}{"countryID": "WORLWIDE", "country": "WORLWIDE", "selectedOfferForRotation": offerID, "addedDate": time.Now().UTC()})
		}

	}

}
func AddGroupRotationForCPI() {
	//saving to Mongo
	config.InitializeMongo()

	//kepping already selected offer
	//repetitionMap := map[string]bool{}
	allPattern := "*"
	var redisKeyLimit int64 = 1000

	//retrieve all redis-keys from offer-geo stack in redis
	results := config.RedisTranxnClient.HScan(constants.RotationOfferGroupStack, 1, allPattern, redisKeyLimit)
	keys, cursor := results.Val()
	pipelineTranxn := config.RedisTranxnClient.Pipeline()
	for {
		//exit critieria
		log.Println("Length Of Keys", len(keys))
		if len(keys) == 0 {
			log.Println("Before break")
			break
		}
		log.Println("After break")

		pipeline := config.RedisMasterClient.Pipeline()

		var records []interface{}

		//find th length of the keys, create batch of 200
		for _, key := range keys {

			//split key to find offerID and geo
			log.Println("Key received:", key)
			splittedKeys := strings.Split(key, constants.Separator)
			if len(splittedKeys) > 1 {
				//get offerID
				offerID := splittedKeys[0]
				//get geo
				countryID := splittedKeys[1]
				//get group
				group := dao.GetOfferGroupOnTranxn(offerID)
				//run elastic query
				rotationOfferID := dao.SearchOfferStack(offerID, group, countryID)
				//save to redis
				pipeline.HSet(constants.RotationOfferGroupStack, offerID+constants.Separator+countryID, rotationOfferID)
				pipelineTranxn.HSet(constants.RotationOfferGroupStack, offerID+constants.Separator+countryID, rotationOfferID)
				if len(group) == 0 {
					group = "NO GROUP OFFER"
				}
				country := countriesMap[countryID]
				if len(countryID) == 0 {
					countryID = "NO COUNTRY FOUND"
					country = "NO COUNTRY FOUND"
				}
				//batch
				records = append(records, model.RotationGroupStack{RotatedFromOffer: offerID, Country: country, Group: group, SelectedOfferForRotation: rotationOfferID, AddedDate: time.Now().UTC()})
			}

		}

		_, err := pipeline.Exec()
		if err != nil {
			fmt.Println("Redis error while saving rotation group stack batch :", err.Error())
		}
		pipeline.Close()

		dao.InsertManyToMongo(constants.MongoDB, constants.RotationGroupStack, records)

		if cursor == 0 {
			log.Println("Cursor is zero.")
			break
		}

		//query next set of keys
		results = config.RedisTranxnClient.HScan(constants.RotationOfferGroupStack, cursor, allPattern, redisKeyLimit)
	}

	pipelineTranxn.Exec()

}
func AddRotationForMOCPA() {

	//retrieve CPA offers with rotation enabled status from MongoDB
	//get their group(bucket)
	//run elastic query and save to CPA-rotation bucket

}

func intializeCountryIDs() {

	config.InitializeMongoBackup()
	//query distinct countryIDS from Mongo
	var results = []string{}
	config.GetMongoBackupSession().DB(constants.MongoDB).C(constants.ClickLog).Find(nil).Distinct("geo.countrycode", &results)
	log.Println(results)
	//find results for each country and save to Redis-Master
	for _, countryID := range results {
		countryIDs = append(countryIDs, countryID)
	}
}

func intializeDropDown() {

	resp, err := http.Get("http://country.io/names.json")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	jsonCountriesData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(jsonCountriesData), &countriesMap)

	if err != nil {
		println(err)
		return
	}
	log.Println(countriesMap)
}
