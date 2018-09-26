package dao

import (
	"context"
	"fmt"
	"log"
	"reflect"

	constants "github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
	model "github.com/adcamie/adserver/db/model"
	logger "github.com/adcamie/adserver/logger"
	"github.com/olivere/elastic"
)

func SaveToES(index string, dbtype string, obj interface{}) {

	put1, err := db.ESMasterClient.Index().Index(index).Type(dbtype).BodyJson(obj).Pretty(false).Do(context.Background())
	if err != nil {
		fmt.Print("ElasticSearchMaster error on saving :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Saving Error")
	}
	fmt.Println("Id on Saving:-", put1.Id)
}

func UpdateToESByQuery(index string, dbtype string, query elastic.Query, script *elastic.Script) {

	update, err := db.ESMasterClient.UpdateByQuery(index).Type(dbtype).Query(query).Script(script).Do(context.Background())
	if err != nil {
		fmt.Print("ElasticSearchMaster error on update query :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Update Query Error")
	}
	fmt.Println("Id on Update:-", update)
}
func UpdateToESByScript(index string, dbtype string, id string, script *elastic.Script) {

	update, err := db.ESMasterClient.Update().Index(index).Type(dbtype).Id(id).Script(script).Do(context.Background())
	if err != nil {
		fmt.Print("ElasticSearchMaster error on update script :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Update Script Error")
	}
	fmt.Println("Id on Update:-", update)
}

func UpsertToES(index string, dbtype string, id string, script *elastic.Script, obj interface{}) {
	_, err := db.ESMasterClient.Update().Index(index).Type(dbtype).Id(id).Script(script).ScriptedUpsert(false).Upsert(obj).Do(context.Background())
	if err != nil {
		fmt.Print("ElasticSearchMaster error on upserting :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Upsert ES Error")
	}
}

func SearchFromES(index string, dbtype string, page int, limit int, query elastic.Query) *elastic.SearchResult {
	res, err := db.ESMasterClient.Search().Index().Index(index).Type(dbtype).Query(query).From((page)*limit).Size(limit).Sort(constants.UpdatedAt, true).Pretty(false).Do(context.Background())
	if err != nil {
		fmt.Print("Search query ES Error :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Search query ES Error")
	}
	return res
}

func DeleteFromES(index string, dbtype string, query elastic.Query) {
	_, err := db.ESMasterClient.DeleteByQuery().Index(index).Type(dbtype).Query(query).Do(context.Background())
	if err != nil {
		fmt.Print("ElasticSearchMaster error on deletion :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Deletion Error")
	}

}

func BulkUpdateToES(index string, dbtype string, requests *elastic.BulkUpdateRequest) {

	update, err := db.ESMasterClient.Bulk().Index(index).Type(dbtype).Add(requests).Do(context.Background())
	if err != nil {
		fmt.Print("ElasticSearchMaster error on bulk update :", err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Bulk Update Error")
	}
	fmt.Println("On bulk update:-", update)
}

func SearchOfferStackWithGeo(countryID string) string {

	if len(countryID) == 0 {
		countryID = "ALL"
	}
	//test group Filter
	TESTGroupFilter := elastic.NewTermQuery("group", "TEST").Boost(2.0)
	testGroupFilter := elastic.NewTermQuery("group", "test").Boost(2.0)
	TestGroupFilter := elastic.NewTermQuery("group", "Test").Boost(2.0)
	//status
	statusFilter := elastic.NewTermQuery("status", constants.NOT_ROTATING).Boost(5.0)
	//only non-mo offer to be rotated to
	offerTypeFilter := elastic.NewTermsQuery("offerType", "2", "3", "4", "7").Boost(2.0)
	//matching geo
	countryFilter := elastic.NewTermQuery("countryIDs", countryID).Boost(5.0)
	//all geo
	allFilter := elastic.NewTermQuery("countryIDs", "ALL").Boost(4.5)
	//high eCPC
	ecpcFilter := elastic.NewRangeQuery("eCPC").Gte(.01).Boost(2.0)
	//avoid status
	avoidFilter := elastic.NewTermQuery("avoidStatus", constants.TRUE).Boost(5.0)
	//minimum clicks to consider
	clicksFilter := elastic.NewRangeQuery("clicks").Gte(5000).Boost(2.0)

	queryBool := elastic.NewBoolQuery().Must(statusFilter, countryFilter, offerTypeFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
	results, err := db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
	if err != nil {
		fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  GEO Error. . Step 1")
		return "105"
	}
	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, allFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  GEO Error. . Step 2")
			return "105"
		}
	}
	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  GEO Error. . Step 3")
			return "105"
		}
	}
	var ttyp model.OfferStack
	if results.TotalHits() > 0 {
		fmt.Println("On Parsing total hits")
		for _, item := range results.Each(reflect.TypeOf(ttyp)) {
			log.Print(item)
			if t, ok := item.(model.OfferStack); ok {
				fmt.Println("Selected Offer:", t.OfferID)
				return t.OfferID
			}
		}
	}
	return "105"
}

func SearchOfferStack(offerID string, group string, countryID string) string {

	if len(countryID) == 0 {
		countryID = "ALL"
	}
	//offer type filter
	offerTypeFilter := elastic.NewTermsQuery("offerType", "2", "3", "4", "7").Boost(5.0)
	//matching geo
	countryFilter := elastic.NewTermQuery("countryIDs", countryID).Boost(5.0)
	//all geo
	allFilter := elastic.NewTermQuery("countryIDs", "ALL").Boost(4.5)
	//offerID Filter
	offerIDFilter := elastic.NewTermQuery("offerID", offerID).Boost(2.0)
	//test group Filter
	testGroupFilter := elastic.NewTermQuery("group", "test").Boost(2.0)
	TestGroupFilter := elastic.NewTermQuery("group", "Test").Boost(2.0)
	TESTGroupFilter := elastic.NewTermQuery("group", "TEST").Boost(2.0)
	//status
	statusFilter := elastic.NewTermQuery("status", constants.NOT_ROTATING).Boost(5.0)
	//same group
	groupFilter := elastic.NewTermQuery("group", group).Boost(5.0)
	//high eCPC
	ecpcFilter := elastic.NewRangeQuery("eCPC").Gte(.01).Boost(2.0)
	//minimum clicks to consider
	clicksFilter := elastic.NewRangeQuery("clicks").Gte(5000).Boost(2.0)
	//avoid status
	avoidFilter := elastic.NewTermQuery("avoidStatus", constants.TRUE).Boost(5.0)

	queryBool := elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, groupFilter, countryFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
	results, err := db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
	if err != nil {
		fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search Error. Step 0")
		return "105"
	}
	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, groupFilter, countryFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search Error. Step 1")
			return "105"
		}
	}

	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, groupFilter, allFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search Error. . Step 2")
			return "105"
		}
	}
	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, countryFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  Error. . Step 3")
			return "105"
		}
	}

	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, allFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  Error. . Step 4")
			return "105"
		}
	}
	var ttyp model.OfferStack
	if results != nil || results.TotalHits() > 0 {
		for _, item := range results.Each(reflect.TypeOf(ttyp)) {
			if t, ok := item.(model.OfferStack); ok {
				fmt.Println("Selected Offer:" + t.OfferID)
				return t.OfferID
			}
		}
	}
	return "105"
}

func SearchMOOfferStackWithGeo(countryID string) string {

	if len(countryID) == 0 {
		countryID = "ALL"
	}
	//test group Filter
	TESTGroupFilter := elastic.NewTermQuery("group", "TEST").Boost(2.0)
	testGroupFilter := elastic.NewTermQuery("group", "test").Boost(2.0)
	TestGroupFilter := elastic.NewTermQuery("group", "Test").Boost(2.0)
	//status
	statusFilter := elastic.NewTermQuery("status", constants.NOT_ROTATING).Boost(5.0)
	//only non-mo offer to be rotated to
	offerTypeFilter := elastic.NewTermsQuery("offerType", "1", "7").Boost(2.0)
	//matching geo
	countryFilter := elastic.NewTermQuery("countryIDs", countryID).Boost(5.0)
	//all geo
	allFilter := elastic.NewTermQuery("countryIDs", "ALL").Boost(4.5)
	//high eCPC
	ecpcFilter := elastic.NewRangeQuery("eCPC").Gte(.01).Boost(2.0)
	//avoid status
	avoidFilter := elastic.NewTermQuery("avoidStatus", constants.TRUE).Boost(5.0)
	//minimum clicks to consider
	clicksFilter := elastic.NewRangeQuery("clicks").Gte(5000).Boost(2.0)

	queryBool := elastic.NewBoolQuery().Must(statusFilter, countryFilter, offerTypeFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
	results, err := db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
	if err != nil {
		fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  GEO Error. . Step 1")
		return "105"
	}
	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, allFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  GEO Error. . Step 2")
			return "105"
		}
	}
	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  GEO Error. . Step 3")
			return "105"
		}
	}
	var ttyp model.OfferStack
	if results.TotalHits() > 0 {
		fmt.Println("On Parsing total hits")
		for _, item := range results.Each(reflect.TypeOf(ttyp)) {
			log.Print(item)
			if t, ok := item.(model.OfferStack); ok {
				fmt.Println("Selected Offer:", t.OfferID)
				return t.OfferID
			}
		}
	}
	return "105"
}

func SearchMOOfferStack(offerID string, group string, countryID string) string {

	if len(countryID) == 0 {
		countryID = "ALL"
	}
	//offer type filter
	offerTypeFilter := elastic.NewTermsQuery("offerType", "1", "7").Boost(5.0)
	//matching geo
	countryFilter := elastic.NewTermQuery("countryIDs", countryID).Boost(5.0)
	//all geo
	allFilter := elastic.NewTermQuery("countryIDs", "ALL").Boost(4.5)
	//offerID Filter
	offerIDFilter := elastic.NewTermQuery("offerID", offerID).Boost(2.0)
	//test group Filter
	testGroupFilter := elastic.NewTermQuery("group", "test").Boost(2.0)
	TestGroupFilter := elastic.NewTermQuery("group", "Test").Boost(2.0)
	TESTGroupFilter := elastic.NewTermQuery("group", "TEST").Boost(2.0)
	//status
	statusFilter := elastic.NewTermQuery("status", constants.NOT_ROTATING).Boost(5.0)
	//same group
	groupFilter := elastic.NewTermQuery("group", group).Boost(5.0)
	//high eCPC
	ecpcFilter := elastic.NewRangeQuery("eCPC").Gte(.01).Boost(2.0)
	//minimum clicks to consider
	clicksFilter := elastic.NewRangeQuery("clicks").Gte(5000).Boost(2.0)
	//avoid status
	avoidFilter := elastic.NewTermQuery("avoidStatus", constants.TRUE).Boost(5.0)

	queryBool := elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, groupFilter, countryFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
	results, err := db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
	if err != nil {
		fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
		go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search Error. Step 0")
		return "105"
	}
	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, groupFilter, countryFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search Error. Step 1")
			return "105"
		}
	}

	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, groupFilter, allFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search Error. . Step 2")
			return "105"
		}
	}
	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, countryFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  Error. . Step 3")
			return "105"
		}
	}

	if results.TotalHits() == 0 {
		queryBool = elastic.NewBoolQuery().Must(statusFilter, offerTypeFilter, allFilter).Should(ecpcFilter, clicksFilter).MustNot(avoidFilter, offerIDFilter, testGroupFilter, TestGroupFilter, TESTGroupFilter)
		results, err = db.ESMasterClient.Search().Index().Index(constants.OfferStack).Type(constants.Offers).Query(queryBool).Pretty(false).Size(1).Do(context.Background())
		if err != nil {
			fmt.Println("ElasticSearchMaster offer search stack error :" + err.Error())
			go logger.ErrorLogger(err.Error(), "ElasticSearchMaster", "Offer Stack Search  Error. . Step 4")
			return "105"
		}
	}
	var ttyp model.OfferStack
	if results != nil || results.TotalHits() > 0 {
		for _, item := range results.Each(reflect.TypeOf(ttyp)) {
			if t, ok := item.(model.OfferStack); ok {
				fmt.Println("Selected Offer:" + t.OfferID)
				return t.OfferID
			}
		}
	}
	return "105"
}
