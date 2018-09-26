package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
)

var countryIDs []string
var countriesMap = make(map[string]string)

func GetOfferByGeo(countryID string) string {
	//check offer availble for that geo
	offerID := dao.GetRotationOfferByGEO(countryID)
	//take from all if not available
	if len(offerID) == 0 {
		offerID = dao.GetRotationOfferByGEO(common.ALL)
		//save offer to redis-master for future requests
		dao.SaveRotationOfferByGEO(countryID, offerID)
	}

	return offerID
}
func GetOfferFromStack(offerID string, countryID string) string {
	//check its readily available or not
	rotationOfferID := dao.GetRotationOfferByOfferID(offerID, countryID)
	group := ""
	//if not found search on ES
	if len(rotationOfferID) == 0 {
		group = dao.GetOfferGroup(offerID)
		rotationOfferID = dao.SearchOfferStack(offerID, group, countryID)
		//save offer to redis-master for future requests
		dao.SaveRotationOfferByGroup(offerID, countryID, rotationOfferID)
	}
	return rotationOfferID
}
func GetMOOfferFromStack(offerID string, countryID string) string {
	//check its readily available or not
	rotationOfferID := dao.GetRotationOfferByOfferID(offerID, countryID)
	group := ""
	//if not found search on ES
	if len(rotationOfferID) == 0 {
		group = dao.GetOfferGroup(offerID)
		rotationOfferID = dao.SearchMOOfferStack(offerID, group, countryID)
		//save offer to redis-master for future requests
		dao.SaveRotationOfferByGroup(offerID, countryID, rotationOfferID)
	}
	return rotationOfferID
}

func GetMOOfferByGeo(countryID string) string {
	//check offer availble for that geo
	offerID := dao.GetRotationMOOfferByGEO(countryID)
	//take from all if not available
	if len(offerID) == 0 {
		offerID = dao.SearchMOOfferStackWithGeo(common.ALL)
		//save offer to redis-master for future requests
		dao.SaveRotationOfferByGEO(countryID, offerID)
	}

	return offerID
}

func intializeDropDown() {

	resp, err := http.Get("http://country.io/names.json")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	jsonCountriesData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal([]byte(jsonCountriesData), &countriesMap)
	if err != nil {
		println(err)
		return
	}
	log.Println(countriesMap)
}
