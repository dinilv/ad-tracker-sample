package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
	message "github.com/adcamie/adserver/messages/proto/v1"
	"github.com/dghubble/sling"
	"github.com/olivere/elastic"
	"golang.org/x/net/context"
)

var ESClient *elastic.Client

func Initialise() {
	ESClient = config.ESMasterClient
}

type Moadserve struct{}

func (mo *Moadserve) Handle(ctx context.Context, msg *message.Message) error {

	log.Print("Async Handler Received message: ")

	ip := ""
	moadTrack := new(model.MOAdTracking)

	//Time stamp details
	utc := time.Now().UTC()
	moadTrack.UTCDate = utc
	moadTrack.Date = time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	moadTrack.Year = utc.Year()
	moadTrack.Hour = utc.Hour()
	moadTrack.Month = utc.Month().String()
	_, moadTrack.Week = utc.ISOWeek()
	moadTrack.MonthID = constants.MonthIds[utc.Month().String()]
	moadTrack.ResponseStatus = constants.AdWorkFlowIds[msg.TrackingActivity]

	for key, value := range msg.Reqh {

		switch key {

		//Related to creative & campaign
		case "wdt":
			moadTrack.Width, _ = strconv.Atoi(value)
		case "hgt":
			moadTrack.Height, _ = strconv.Atoi(value)
		case "cid":
			moadTrack.CategoryType, _ = strconv.Atoi(value)
		case "id":
			moadTrack.ID = value
		case "type":
			moadTrack.PriorityType = value
		case "CampaignID":
			moadTrack.CampaignID, _ = strconv.Atoi(value)
		case "CreativeID":
			moadTrack.CreativeID, _ = strconv.Atoi(value)
		//related to request header
		case "User-Agent":
			moadTrack.UserAgent = value
		case "X-Forwarded-For":
			//split for lb adding ip
			ipsplitted := strings.Split(value, ",")
			ip = ipsplitted[0]
		//related to user
		case "CookieID":
			moadTrack.CookieId = value
		case "RequestID":
			moadTrack.RequestID = value
		case "tid":
			moadTrack.TrackingID = value
		case "searchTime":
			moadTrack.ResponseTime, _ = strconv.Atoi(value)
		//pubisher
		case "appid":
			moadTrack.AppID = value
			details := strings.Split(value, "_")
			moadTrack.Version = details[0]
			moadTrack.Publisher = details[4]
			moadTrack.ContentSubCategory, _ = strconv.Atoi(details[6])
			moadTrack.ContentPageID, _ = strconv.Atoi(details[6])
		//Operator
		case "mid":
			moadTrack.MSISDN = value
		case "op":
			moadTrack.Operator = value
		}

	}

	//According to activity
	switch msg.TrackingActivity {

	//request is received
	case 1:
		log.Print("On receiving request:-" + moadTrack.RequestID)

		//Location and geography details
		geo := new(model.GeoDetails)
		sling.New().Get(constants.URL + ip).ReceiveSuccess(geo)
		moadTrack.Geo = *geo

		//save camapign details & creative if needed
		dao.SaveToES(constants.TrackingIndex, constants.RequestLogType, moadTrack)

	//Ad is served
	case 2:
		log.Print("On server response. Request_Id:-"+moadTrack.RequestID+", Tracking_id:-"+moadTrack.TrackingID, ", Response Time:-", moadTrack.ResponseTime)

		//Check tid is generated or not
		if moadTrack.TrackingID != "" {

			//if generated update the details to request_log
			query := elastic.NewTermQuery(constants.RequestID, moadTrack.RequestID)
			//search for records with same request id get the ID and update the same
			results, err := ESClient.Search().Index().Index(constants.TrackingIndex).Type(constants.RequestLogType).Query(query).Do(context.Background())
			log.Print("Initial Checking for Request Log:-", results.TotalHits())

			//Check whether the entry is updated or not
			if results.TotalHits() == 0 {
				time.Sleep(5 * time.Second)
			}
			//Check after 5 sec
			results, err = ESClient.Search().Index().Index(constants.TrackingIndex).Type(constants.RequestLogType).Query(query).Do(context.Background())
			log.Print("Checking for Request Log After 5 secs:-", results.TotalHits())
			if err != nil {
				// Handle error
				log.Print(err.Error())
			}
			//Updates & save to ES
			script := elastic.NewScript("ctx._source.responseStatus = status;ctx._source.trackingID = tid;ctx._source.campaignID = campaignID;ctx._source.creativeID = creativeId;ctx._source.cookieID = cookieId;ctx._source.rowID = rowID;ctx._source.pType = ptype;ctx._source.responseTime=responseTime").
				Params(map[string]interface{}{"status": moadTrack.ResponseStatus, "tid": moadTrack.TrackingID, "campaignID": moadTrack.CampaignID, "creativeId": moadTrack.CreativeID, "cookieId": moadTrack.CookieId, "rowID": moadTrack.ID, "ptype": moadTrack.PriorityType, "responseTime": moadTrack.ResponseTime})
			dao.UpdateToESByQuery(constants.TrackingIndex, constants.RequestLogType, query, script)

			//update CookieId
			script = elastic.NewScript("ctx._source.cookieIds += cookieId").Params(map[string]interface{}{"cookieId": moadTrack.CookieId})
			update, err := ESClient.Update().Index("mobileads").Type(moadTrack.PriorityType).Id(moadTrack.ID).Script(script).ScriptedUpsert(true).Do(context.Background())
			log.Print(update, "Udpate Cookie Id: Check nil or Not", moadTrack.ID)
			log.Print(update, err)

		}

	//Impression is received
	case 3:
		log.Print("On Impression saving:-" + moadTrack.TrackingID)
		//Location and geography details
		geo := new(model.GeoDetails)
		sling.New().Get(constants.URL + ip).ReceiveSuccess(geo)
		moadTrack.Geo = *geo

		//upate the impression count
		query := elastic.NewTermQuery("trackingID", moadTrack.TrackingID)
		results, _ := ESClient.Search().Index().Index(constants.TrackingIndex).Type(constants.RequestLogType).Query(query).Do(context.Background())
		id := ""
		ptype := "p1"
		var moadDBTrack model.MOAdTracking
		for _, hit := range results.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &moadDBTrack)
			if err != nil {
				// Deserialization failed
				log.Print(err.Error())
			}
			id = moadDBTrack.ID
			ptype = moadDBTrack.PriorityType
			moadTrack.CampaignID = moadDBTrack.CampaignID
			moadTrack.CreativeID = moadDBTrack.CreativeID
			moadTrack.ID = moadDBTrack.ID

		}

		log.Print("Id:-"+id, ptype, moadTrack.CampaignID, moadTrack.CreativeID, moadTrack.ID)
		//log.Print(id, ptype)
		responseScript := elastic.NewScript("ctx._source.impressions+=1").Params(map[string]interface{}{})
		updates := ESClient.Update().Index(constants.MobileAdsIndex).Type(ptype).Id(id).Script(responseScript).ScriptedUpsert(true)
		log.Print(updates, ptype)
		update, err := updates.Do(context.Background())

		//idQuery := elastic.NewIdsQuery(ptype).Ids(id).QueryName("idQuery")

		//ESClient.UpdateByQuery().Index(constants.MobileAdsIndex).Query(idQuery).Script(responseScript).Do(context.Background())
		if err != nil {
			// Handle error
			log.Print(err.Error())
		}
		fmt.Println("Id on Update:-", update)

		//log.Print("update", update)

		//if err != nil {
		// Deserialization failed
		//log.Print(err.Error())
		//}

		//Save to ES
		dao.SaveToES(constants.TrackingIndex, constants.ImpressionLogType, moadTrack)
		//Save to mongo
		db.MongoSession.DB(constants.TrackingIndex).C(constants.ImpressionLogType).Insert(moadTrack)

	//Click is received
	case 4:
		log.Print("On Click saving:-" + moadTrack.TrackingID)

		//Location and geography details
		geo := new(model.GeoDetails)
		sling.New().Get(constants.URL + ip).ReceiveSuccess(geo)
		moadTrack.Geo = *geo

		//upate the click count
		query := elastic.NewTermQuery("trackingID", moadTrack.TrackingID)
		results, _ := ESClient.Search().Index().Index("tracking").Type("request_log").Query(query).Do(context.Background())
		id := ""
		ptype := "p1"
		var moadDBTrack model.MOAdTracking
		for _, hit := range results.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &moadDBTrack)
			if err != nil {
				// Deserialization failed
				log.Print(err.Error())
			}
			id = moadDBTrack.ID
			ptype = moadDBTrack.PriorityType
			moadTrack.CampaignID = moadDBTrack.CampaignID
			moadTrack.CreativeID = moadDBTrack.CreativeID
			moadTrack.ID = moadDBTrack.ID

		}
		responseScript := elastic.NewScript("ctx._source.clicks+=1")
		updates := ESClient.Update().Index("mobileads").Type(ptype).Id(id).Script(responseScript).ScriptedUpsert(true)
		log.Print(updates, ptype)
		update, _ := updates.Do(context.Background())
		log.Print("update", update)

		dao.SaveToES(constants.TrackingIndex, constants.ClickLogType, moadTrack)

	}

	if recover() != nil {
		log.Print("Something Went Wrong :-()")
	}

	return nil
}
