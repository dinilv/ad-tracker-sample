package v1

import (
	"context"
	"reflect"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	"github.com/micro/go-micro/server"
	"github.com/olivere/elastic"
	"gopkg.in/mgo.v2/bson"
)

type GetOfferStackReq struct {
	Token  string `json:"token,omitempty"`
	Search string `json:"search,omitempty"`
	Page   int    `json:"from,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}
type UpdateOfferStackReq struct {
	Token       string `json:"token,omitempty"`
	OfferID     string `json:"offerID"`
	Status      string `json:"status"`
	AvoidStatus string `json:"avoidStatus"`
}
type MsisdnReq struct {
	MSISDN string `json:"msisdn,omitempty"`
}
type MsisdnRes struct {
	MSISDN   string         `json:"msisdn,omitempty" bson:"msisdn,omitempty"`
	Operator string         `json:"operator,omitempty" bson:"operator,omitempty"`
	Blocked  string         `json:"blocked,omitempty" bson:"blocked,omitempty"`
	Offers   []OfferDetails `json:"offers,omitempty" bson:"offerIDs,omitempty"`
}

type OfferDetails struct {
	OfferID   string `json:"offerID"`
	OfferName string `json:"offerName"`
}

type GetOfferStackRes struct {
	Data  []model.OfferStack `json:"data,omitempty"`
	Count int64              `json:"count,omitempty"`
}

type GetOfferEventListReq struct {
	Token     string   `json:"token,omitempty"`
	OfferIDs  []string `json:"offer_ids,omitempty"`
	StartDate string   `json:"start_date,omitempty"`
	EndDate   string   `json:"end_date,omitempty"`
}

type GetOfferEventListRes struct {
	Data  []bson.M `json:"data,omitempty"`
	Count int      `json:"count,omitempty"`
}

type EngineRes struct {
	Status string `json:"status,omitempty"`
	Id     string `json:"id,omitempty"`
}

type EngineHandler interface {
	Events(context.Context, *model.AdcamieEvents, *EngineRes) error
	Listevents(context.Context, *GetOfferEventListReq, *GetOfferEventListRes) error
	Offerstack(context.Context, *GetOfferStackReq, *GetOfferStackRes) error
	MSISDN(context.Context, *MsisdnReq, *MsisdnRes) error
	Update(context.Context, *UpdateOfferStackReq, *EngineRes) error
}

type Engine struct {
	EngineHandler
}

func RegisterEngineHandler(s server.Server, hdlr EngineHandler) {
	log.Print("Getting Decision Engine Setting Handler")
	s.Handle(s.NewHandler(&Engine{hdlr}))
}

func (engine *Engine) Events(ctx context.Context, req *model.AdcamieEvents, rsp *EngineRes) error {
	log.Print("Add events to tracker:")

	bulk := db.ESMasterClient.Bulk().Index(constants.OfferStack).Type(constants.Offers)

	for _, event := range req.Events {
		//save to offer events
		event.IP = req.IP
		event.CreatedAt = time.Now().UTC()
		event.UTCDate = time.Now().UTC()
		dao.InsertToMongoSession(constants.MongoDB, constants.AdcamieEvents, event)

		//save/update to offer stack with script
		var scripts []string
		params := make(map[string]interface{})

		//updated at
		scripts = append(scripts, "ctx._source.updatedAt=params.updatedAt")
		params[constants.UpdatedAt] = time.Now().UTC()

		if len(event.Status) > 0 {
			scripts = append(scripts, "ctx._source.status=params.status")
			params[constants.Status] = event.Status
			//update to redis db for exhausted offer or not
			if strings.Compare(event.Status, constants.NOT_ROTATING) == 0 {
				//remove offer_id from exhausted offer
				dao.RemoveOfferInExhausted(event.OfferID)
				dao.RemoveOfferInExhaustedHash(event.OfferID)
				dao.DeleteFromMongoSession(constants.MongoDB, constants.ExhaustedOffer, map[string]interface{}{"offerID": event.OfferID})
			} else {
				//add to exhausted offer_id and rotation starts
				dao.SaveExhaustedOffer(event.OfferID)
				dao.SaveExhaustedOfferMap(event.OfferID)
				dao.InsertToMongoSession(constants.MongoDB, constants.ExhaustedOffer, &model.RotationStack{
					OfferID:   event.OfferID,
					AddedDate: time.Now().UTC(),
					Event:     event.EventName,
				})
			}
		}
		if len(event.AvoidStatus) > 0 {
			scripts = append(scripts, "ctx._source.avoidStatus=params.avoidStatus")
			params[constants.AvoidStatus] = event.AvoidStatus
		}
		if len(event.OfferType) > 0 {
			scripts = append(scripts, "ctx._source.offerType=params.offerType")
			params[constants.OfferType] = event.OfferType
		}
		if len(event.CampaignID) > 0 {
			scripts = append(scripts, "ctx._source.campaignID=params.campaignID")
			params[constants.CampaignID] = event.CampaignID
		}
		if len(event.EventID) > 0 {
			scripts = append(scripts, "ctx._source.lastEventID=params.eventID")
			params[constants.EventID] = event.EventID
		}
		if len(event.EventName) > 0 {
			scripts = append(scripts, "ctx._source.lastEvent=params.eventName")
			params[constants.EventName] = event.EventName
		}
		if len(event.Comment) > 0 {
			scripts = append(scripts, "ctx._source.comment=params.comment")
			params[constants.Comment] = event.Comment
		}
		if len(event.Group) > 0 {
			scripts = append(scripts, "ctx._source.group=params.group")
			params[constants.Group] = event.Group
		}
		if len(event.CountryIDs) > 0 {
			scripts = append(scripts, "ctx._source.countryIDs=params.countryIDs")
			params[constants.CountryIDs] = event.CountryIDs
		}
		if event.Bid != 0 {
			scripts = append(scripts, "ctx._source.bid=params.bid")
			params[constants.Bid] = event.Bid
		}
		if event.ECPC != 0 {
			scripts = append(scripts, "ctx._source.ecpc=params.ecpc")
			params[constants.ECPC] = event.ECPC
		}
		if event.Clicks != 0 {
			scripts = append(scripts, "ctx._source.clicks=params.clicks")
			params[constants.Clicks] = event.Clicks
		}
		if event.Conversions != 0 {
			scripts = append(scripts, "ctx._source.conversions=params.conversions")
			params[constants.Conversions] = event.Conversions
		}
		//create upsert obj
		obj := &model.OfferStack{
			OfferID:     event.OfferID,
			OfferType:   event.OfferType,
			LastEventID: event.EventID,
			LastEvent:   event.EventName,
			Bid:         event.Bid,
			Status:      event.Status,
			AvoidStatus: event.AvoidStatus,
			Group:       event.Group,
			ECPC:        event.ECPC,
			Clicks:      event.Clicks,
			Conversions: event.Conversions,
			CountryIDs:  event.CountryIDs,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		script := strings.Join(scripts, ";")
		escript := elastic.NewScript(script).Params(params).Lang("painless")
		bulk.Add(elastic.NewBulkUpdateRequest().Id(event.OfferID).Script(escript).ScriptedUpsert(false).Upsert(obj))
	}

	res, err := bulk.Do(context.Background())
	if err != nil {
		log.Println(err, res, bulk.NumberOfActions())
	}
	rsp.Status = "Success"
	return nil
}

func (engine *Engine) Offerstack(ctx context.Context, req *GetOfferStackReq, rsp *GetOfferStackRes) error {
	log.Print(" Get Offers Stack from  tracker:", req.Limit, req.Page)
	var res *elastic.SearchResult
	//check search keyword is empty or not
	if len(req.Search) > 0 {
		wildcard := elastic.NewWildcardQuery("offerID", req.Search+"*")
		res = dao.SearchFromES(constants.OfferStack, constants.Offers, req.Page-1, req.Limit, wildcard)
	} else {
		query := elastic.NewMatchAllQuery()
		res = dao.SearchFromES(constants.OfferStack, constants.Offers, req.Page-1, req.Limit, query)
	}

	var ofstack []model.OfferStack
	var ttyp model.OfferStack
	if res != nil || res.TotalHits() > 0 {
		for _, item := range res.Each(reflect.TypeOf(ttyp)) {
			log.Print(item)
			if t, ok := item.(model.OfferStack); ok {
				ofstack = append(ofstack, t)
			}
		}
	}
	rsp.Data = ofstack
	rsp.Count = res.Hits.TotalHits
	return nil
}

func (engine *Engine) Listevents(ctx context.Context, req *GetOfferEventListReq, rsp *GetOfferEventListRes) error {
	log.Print(" Get Event list for offers from  tracker:")
	matchFilter := make(map[string]interface{})
	if len(req.OfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.OfferIDs}
		matchFilter["offerID"] = offFilter
	}
	start, _ := time.Parse(inputDateFormat, req.StartDate)
	end, _ := time.Parse(inputDateFormat, req.EndDate)
	start = start.UTC()
	end = end.UTC()
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	dateFilter := map[string]interface{}{"$gte": start, "$lte": end}
	matchFilter["createdAt"] = dateFilter
	data, count := dao.GetAllSortedFromMongoSession(constants.MongoDB, constants.AdcamieEvents, constants.CreatedAt, matchFilter)
	rsp.Data = data
	rsp.Count = count
	return nil
}

func (engine *Engine) MSISDN(ctx context.Context, req *MsisdnReq, rsp *MsisdnRes) error {
	log.Print(" Get MSISDN Details from tracker:")
	return nil
}

func (engine *Engine) Update(ctx context.Context, req *UpdateOfferStackReq, rsp *EngineRes) error {
	log.Print("Update Offer Stack Details:", req)

	//save to offer events
	event := model.AdcamieEvent{}
	if len(req.Status) > 0 {
		event.OfferID = req.OfferID
		event.EventID = "20"
		event.EventName = "Manual Intervention For Rotation Status"
		event.Comment = "Manual Intervention from Tracker UI, status-" + req.Status
		event.Status = req.Status
		event.CreatedAt = time.Now().UTC()

	}
	if len(req.AvoidStatus) > 0 {
		event.OfferID = req.OfferID
		event.EventID = "20"
		event.EventName = "Manual Intervention For Avoid Rotation Status"
		event.Comment = "Manual Intervention from Tracker UI, status-" + req.AvoidStatus
		event.AvoidStatus = req.AvoidStatus
		event.CreatedAt = time.Now().UTC()

	}
	dao.InsertToMongoSession(constants.MongoDB, constants.AdcamieEvents, event)

	//update to offer stack with script
	script := elastic.NewScript("ctx._source.status=params.status;ctx._source.updatedAt=params.updatedAt;ctx._source.lastEventID=params.eventID;ctx._source.lastEventName=params.eventName").Params(map[string]interface{}{"status": req.Status, "updatedAt": time.Now().UTC(), "eventID": "20", "eventName": "Manaul Intervention"}).Lang("painless")
	if len(req.AvoidStatus) > 0 {
		script = elastic.NewScript("ctx._source.avoidStatus=params.avoidStatus;ctx._source.updatedAt=params.updatedAt;ctx._source.lastEventID=params.eventID;ctx._source.lastEventName=params.eventName").Params(map[string]interface{}{"avoidStatus": req.AvoidStatus, "updatedAt": time.Now().UTC(), "eventID": "21", "eventName": "Manaul Intervention for avoid status"}).Lang("painless")
	}
	dao.UpdateToESByScript(constants.OfferStack, constants.Offers, req.OfferID, script)
	//save to redis server
	if len(req.Status) > 0 {
		//update to redis db for exhausted offer or not
		if strings.Compare(event.Status, constants.NOT_ROTATING) == 0 {
			//remove offer_id from exhausted offer
			dao.RemoveOfferInExhausted(event.OfferID)
			dao.RemoveOfferInExhaustedHash(event.OfferID)
			dao.DeleteFromMongoSession(constants.MongoDB, constants.ExhaustedOffer, map[string]interface{}{"offerID": event.OfferID})
		} else {
			//add to exhausted offer_id and rotation starts
			dao.SaveExhaustedOffer(event.OfferID)
			dao.SaveExhaustedOfferMap(event.OfferID)
			dao.InsertToMongoSession(constants.MongoDB, constants.ExhaustedOffer, &model.RotationStack{
				OfferID:   event.OfferID,
				AddedDate: time.Now().UTC(),
				Event:     event.EventName,
			})
		}
	}
	rsp.Status = "Success"
	return nil
}
