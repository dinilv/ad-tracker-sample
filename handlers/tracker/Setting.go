package v1

import (
	"context"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	helper "github.com/adcamie/adserver/helpers"
	"github.com/micro/go-micro/server"
	"github.com/olivere/elastic"
	"gopkg.in/mgo.v2/bson"
)

type SetOfferReq struct {
	OfferID          string    `json:"offerID,omitempty" bson:"offerID,omitempty"`
	OfferRefID       string    `json:"offerRefID,omitempty" bson:"offerRefID,omitempty"`
	OfferName        string    `json:"offerName,omitempty" bson:"offerName,omitempty"`
	Group            string    `json:"group,omitempty" bson:"group,omitempty"`
	Template         string    `json:"template,omitempty" bson:"template,omitempty"`
	Token            string    `json:"token,omitempty" bson:"token,omitempty"`
	Type             string    `json:"type,omitempty" bson:"type,omitempty"`
	WhiteListEnabled bool      `json:"whitelistEnabled,omitempty" bson:"whitelistEnabled,omitempty"`
	Advertiser       string    `json:"advertiser,omitempty" bson:"advertiser,omitempty"`
	CountryIDs       []string  `json:"countryIDs,omitempty" bson:"countryIDs,omitempty"`
	GoalID           string    `json:"goalID,omitempty" bson:"goalID,omitempty"`
	ServiceID        string    `json:"serviceID,omitempty" bson:"serviceID,omitempty"`
	CreatedAt        time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt        time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CreatedBy        int       `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	UpdatedBy        int       `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	IP               string    `json:"ip,omitempty" bson:"ip,omitempty"`
}
type SettingRes struct {
	Status string `json:"status,omitempty"`
	Id     string `json:"id,omitempty"`
}

type SetAffiliateReq struct {
	AffiliateID    string    `json:"affiliateID,omitempty" bson:"affiliateID,omitempty"`
	AffiliateName  string    `json:"affiliateName,omitempty" bson:"affiliateName,omitempty"`
	AffiliateRefID string    `json:"affiliateRefID,omitempty" bson:"affiliateRefID,omitempty"`
	Mqf            float64   `json:"mqf,omitempty" bson:"mqf,omitempty"`
	Token          string    `json:"token,omitempty"`
	MediaTemplate  string    `json:"mediaTemplate,omitempty" bson:"mediaTemplate,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CreatedBy      int       `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	UpdatedBy      int       `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	IP             string    `json:"ip,omitempty" bson:"ip,omitempty"`
}

type SetMqfReq struct {
	OfferID       string    `json:"offerID,omitempty" bson:"offerID,omitempty"`
	AffiliateID   string    `json:"affiliateID,omitempty" bson:"affiliateID,omitempty"`
	Mqf           float64   `json:"mqf,omitempty" bson:"mqf,omitempty"  `
	Token         string    `json:"token,omitempty" bson:"token,omitempty" `
	ResetCounter  bool      `json:"resetCounter,omitempty" bson:"resetCounter,omitempty" `
	MediaTemplate string    `json:"mediaTemplate,omitempty" bson:"mediaTemplate,omitempty"`
	Rotated       bool      `json:"rotated,omitempty" bson:"rotated,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	CreatedBy     int       `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	UpdatedBy     int       `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	IP            string    `json:"ip,omitempty" bson:"ip,omitempty"`
}

type GetOffersReq struct {
	Token string `json:"token,omitempty"`
}

type GetOffersRes struct {
	Data  []bson.M `json:"data,omitempty"`
	Count int      `json:"count,omitempty"`
}

type GetOfferReq struct {
	Token string `json:"token,omitempty"`
	Id    string `json:"id,omitempty"`
}

type DeleteOfferReq struct {
	Token      string `json:"token,omitempty"`
	OfferID    string `json:"offerID,omitempty" bson:"offerID,omitempty"`
	OfferRefID string `json:"offerRefID,omitempty" bson:"offerRefID,omitempty"`
}

type DeleteOfferRes struct {
	Status string `json:"status,omitempty"`
	Id     string `json:"id,omitempty"`
}

type DeleteAffiliateReq struct {
	Token          string `json:"token,omitempty"`
	AffiliateID    string `json:"affiliateID,omitempty" bson:"affiliateID,omitempty"`
	AffiliateRefID string `json:"affiliateRefID,omitempty" bson:"affiliateRefID,omitempty"`
}

type DeleteAffiliateRes struct {
	Status string `json:"status,omitempty"`
	Id     string `json:"id,omitempty"`
}

type GetMQFReq struct {
	Token       string `json:"token,omitempty"`
	Offerid     string `json:"offerID,omitempty"`
	Affiliateid string `json:"affiliateID,omitempty"`
}

type DeleteMQFReq struct {
	Token       string `json:"token,omitempty"`
	Offerid     string `json:"offerID,omitempty"`
	Affiliateid string `json:"affiliateID,omitempty"`
}

type DeleteMQFRes struct {
	Status string `json:"status,omitempty"`
	Id     string `json:"id,omitempty"`
}

type GetOfferRes struct {
	Data []bson.M `json:"data,omitempty"`
}

type GetTransactionReq struct {
	Token        string   `json:"token,omitempty"`
	OfferIDs     []string `json:"offer_ids,omitempty"`
	AffiliateIDs []string `json:"affiliate_ids,omitempty"`
}

type GetTransactionRes struct {
	Data []bson.M `json:"data"`
}

type ManualPostbackReq struct {
	OfferID       string    `json:"offerID,omitempty" bson:"offerID,omitempty"`
	AffiliateID   string    `json:"affiliateID,omitempty" bson:"affiliateID,omitempty"`
	Token         string    `json:"token,omitempty" bson:"token,omitempty" `
	TransactionID string    `json:"transactionID,omitempty" bson:"transactionID,omitempty"`
	CreatedAt     time.Time `json:"utcdate,omitempty" bson:"utcdate,omitempty"`
	Hour          int       `json:"hour,omitempty" bson:"hour,omitempty"`
	CreatedBy     int       `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	IP            string    `json:"ip,omitempty" bson:"ip,omitempty"`
	URL           string    `json:"url,omitempty" bson:"url,omitempty"`
	Date          time.Time `json:"date,omitempty" bson:"date,omitempty"`
}
type ManualPostbackRes struct {
	TransactionID string
	URL           string `json:"url,omitempty" bson:"url,omitempty"`
	Status        string `json:"status,omitempty" bson:"status,omitempty"`
}

type SetAdvertiserReq struct {
	AdvertiserName     string    `json:"advertiserName,omitempty" bson:"advertiserName,omitempty"`
	AdvertiserIP       []string  `json:"advertiserIPs,omitempty" bson:"advertiserIPs,omitempty"`
	AdvertiserTemplate string    `json:"advertiserTemplate,omitempty" bson:"advertiserTemplate,omitempty"`
	Token              string    `json:"token,omitempty" bson:"token,omitempty" `
	CreatedAt          time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt          time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type DeleteAdvertiserReq struct {
	Token          string `json:"token,omitempty"`
	AdvertiserName string `json:"advertiserName,omitempty" bson:"advertiserName,omitempty"`
}

type DeleteAdvertiserRes struct {
	Status string `json:"status,omitempty"`
	Id     string `json:"id,omitempty"`
}

type TrackerHandler interface {
	Setoffer(context.Context, *SetOfferReq, *SettingRes) error
	Setaffiliate(context.Context, *SetAffiliateReq, *SettingRes) error
	Setmqf(context.Context, *SetMqfReq, *SettingRes) error
	Getoffers(context.Context, *GetOffersReq, *GetOffersRes) error
	Getoffer(context.Context, *GetOfferReq, *GetOfferRes) error
	Deleteoffer(context.Context, *DeleteOfferReq, *DeleteOfferRes) error
	Getaffiliates(context.Context, *GetOffersReq, *GetOffersRes) error
	Getaffiliate(context.Context, *GetOfferReq, *GetOfferRes) error
	Deleteaffiliate(context.Context, *DeleteAffiliateReq, *DeleteAffiliateRes) error
	Getmqfs(context.Context, *GetOffersReq, *GetOffersRes) error
	Getmqf(context.Context, *GetMQFReq, *GetOfferRes) error
	Deletemqf(context.Context, *DeleteMQFReq, *DeleteMQFRes) error
	Gettransactions(context.Context, *GetTransactionReq, *GetTransactionRes) error
	Dopostback(context.Context, *ManualPostbackReq, *ManualPostbackRes) error
	Getmanuallog(context.Context, *GetOfferReq, *GetOfferRes) error
	Addadvertiser(context.Context, *SetAdvertiserReq, *SettingRes) error
	Getadvertisers(context.Context, *GetOffersReq, *GetOffersRes) error
	Getadvertiser(context.Context, *GetOfferReq, *GetOfferRes) error
	Deleteadvertiser(context.Context, *DeleteAdvertiserReq, *DeleteAdvertiserRes) error
}

type Tracker struct {
	TrackerHandler
}

func RegisterTrackerHandler(s server.Server, hdlr TrackerHandler) {
	log.Print("Getting Tracker Setting Handler")
	s.Handle(s.NewHandler(&Tracker{hdlr}))
}

func (tracker *Tracker) Setoffer(ctx context.Context, req *SetOfferReq, rsp *SettingRes) error {
	log.Print("Setting Offer in tracker setting for offer id: ", req.OfferID)
	if req.WhiteListEnabled {
		dao.SaveOffer(req.OfferID, req.Type, req.Template, req.Group, req.Advertiser, req.CountryIDs)
	} else {
		dao.SaveOffer(req.OfferID, req.Type, req.Template, req.Group, "", req.CountryIDs)
	}

	//delete if exists
	filters := make(map[string]interface{})
	filters[constants.OfferID] = req.OfferID
	dao.DeleteFromMongoSession(constants.MongoDB, constants.Offer, filters)
	req.CreatedAt = time.Now().UTC()
	dao.InsertToMongoSession(constants.MongoDB, constants.Offer, req)
	//insert or update to offer stack
	var scripts []string
	params := make(map[string]interface{})
	scripts = append(scripts, "ctx._source.updatedAt=params.updatedAt")
	params[constants.UpdatedAt] = time.Now().UTC()
	scripts = append(scripts, "ctx._source.status=params.status")
	params[constants.Status] = constants.NOT_ROTATING
	scripts = append(scripts, "ctx._source.avoidStatus=params.avoidStatus")
	params[constants.AvoidStatus] = constants.FALSE
	scripts = append(scripts, "ctx._source.offerType=params.offerType")
	params[constants.OfferType] = req.Type
	scripts = append(scripts, "ctx._source.lastEventID=params.eventID")
	params[constants.EventID] = "0"
	scripts = append(scripts, "ctx._source.lastEvent=params.eventName")
	params[constants.EventName] = "Offer added/edited in Tracker UI/NIS UI."
	scripts = append(scripts, "ctx._source.group=params.group")
	params[constants.Group] = req.Group
	scripts = append(scripts, "ctx._source.countryIDs=params.countryIDs")
	params[constants.CountryIDs] = req.CountryIDs
	script := strings.Join(scripts, ";")
	escript := elastic.NewScript(script).Params(params).Lang("painless")
	//create upsert obj
	obj := &model.OfferStack{
		OfferID:     req.OfferID,
		OfferType:   req.Type,
		Status:      constants.NOT_ROTATING,
		AvoidStatus: constants.FALSE,
		LastEventID: "0",
		LastEvent:   "Offer added/edited in Tracker UI/NIS UI.",
		Group:       req.Group,
		CountryIDs:  req.CountryIDs,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	dao.UpsertToES(constants.OfferStack, constants.Offers, req.OfferID, escript, obj)
	rsp.Status = "Success"
	rsp.Id = req.OfferID
	return nil
}

func (tracker *Tracker) Setaffiliate(ctx context.Context, req *SetAffiliateReq, rsp *SettingRes) error {
	log.Print("Setting Affiliate in tracker setting for affiliate id: ", req.AffiliateID)

	dao.SaveAffiliate(req.AffiliateID, strconv.FormatFloat(req.Mqf, 'f', 2, 64), req.MediaTemplate)

	//delete if exists
	filters := make(map[string]interface{})
	filters[constants.AffiliateID] = req.AffiliateID
	dao.DeleteFromMongoSession(constants.MongoDB, constants.Affiliate, filters)
	req.CreatedAt = time.Now().UTC()
	dao.InsertToMongoSession(constants.MongoDB, constants.Affiliate, req)
	rsp.Status = "Success"
	rsp.Id = req.AffiliateID
	return nil
}

func (tracker *Tracker) Setmqf(ctx context.Context, req *SetMqfReq, rsp *SettingRes) error {
	log.Print("Setting Mqf Offer Affiliate in tracker setting for affiliate id: ", req)
	dao.SaveOfferAffiliateMQF(req.OfferID, req.AffiliateID, strconv.FormatFloat(req.Mqf, 'f', 2, 64), req.MediaTemplate, req.Rotated, req.ResetCounter)
	//delete if exists
	filters := make(map[string]interface{})
	filters[constants.OfferID] = req.OfferID
	filters[constants.AffiliateID] = req.AffiliateID
	dao.DeleteFromMongoSession(constants.MongoDB, constants.OfferAffiliateMQF, filters)
	req.CreatedAt = time.Now().UTC()
	dao.InsertToMongoSession(constants.MongoDB, constants.OfferAffiliateMQF, req)
	//handle rotation stack in mongo
	if !req.Rotated {
		dao.DeleteFromMongoSession(constants.MongoDB, constants.ExhaustedOffer, map[string]interface{}{"offerID": req.OfferID, "affiliateID": req.AffiliateID})
	} else {
		dao.InsertToMongoSession(constants.MongoDB, constants.ExhaustedOfferAffiliate, &model.RotationStack{
			OfferID:     req.OfferID,
			AffiliateID: req.AffiliateID,
			AddedDate:   time.Now().UTC(),
			Event:       "Saving Offer Affiliate Index",
		})
	}

	rsp.Status = "Success"
	rsp.Id = req.OfferID
	return nil
}

func (tracker *Tracker) Getoffers(ctx context.Context, req *GetOffersReq, rsp *GetOffersRes) error {
	log.Print(" Get Offers from  tracker:")
	matchFilter := make(map[string]interface{})
	data, count := dao.GetAllSortedFromMongo(constants.MongoDB, constants.Offer, constants.CreatedAt, matchFilter)
	rsp.Data = data
	rsp.Count = count
	return nil
}

func (tracker *Tracker) Getoffer(ctx context.Context, req *GetOfferReq, rsp *GetOfferRes) error {
	log.Print(" Get Offers from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter[constants.OfferID] = req.Id
	data, _ := dao.GetAllFromMongoSession(constants.MongoDB, constants.Offer, matchFilter)
	rsp.Data = data
	return nil
}

func (tracker *Tracker) Deleteoffer(ctx context.Context, req *DeleteOfferReq, rsp *DeleteOfferRes) error {
	log.Print(" Delete Offer from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter[constants.OfferID] = req.OfferID
	matchFilter[constants.OfferRefID] = req.OfferRefID
	dao.DeleteRecordFromMongoSession(constants.MongoDB, constants.Offer, matchFilter)
	rsp.Id = req.OfferID
	rsp.Status = "Success"
	return nil
}

func (tracker *Tracker) Getaffiliates(ctx context.Context, req *GetOffersReq, rsp *GetOffersRes) error {
	log.Print(" Get Affiliates from  tracker:")
	matchFilter := make(map[string]interface{})
	data, count := dao.GetAllSortedFromMongoSession(constants.MongoDB, constants.Affiliate, constants.CreatedAt, matchFilter)
	rsp.Data = data
	rsp.Count = count
	return nil

}

func (tracker *Tracker) Getaffiliate(ctx context.Context, req *GetOfferReq, rsp *GetOfferRes) error {
	log.Print(" Get Affiliate from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter[constants.AffiliateID] = req.Id
	data, _ := dao.GetAllFromMongoSession(constants.MongoDB, constants.Affiliate, matchFilter)
	rsp.Data = data
	return nil

}

func (tracker *Tracker) Deleteaffiliate(ctx context.Context, req *DeleteAffiliateReq, rsp *DeleteAffiliateRes) error {
	log.Print(" Delete Affiliate from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter[constants.AffiliateID] = req.AffiliateID
	matchFilter[constants.AffiliateRefID] = req.AffiliateRefID
	dao.DeleteRecordFromMongoSession(constants.MongoDB, constants.Affiliate, matchFilter)
	rsp.Id = req.AffiliateID
	rsp.Status = "Success"
	return nil
}

func (tracker *Tracker) Getmqfs(ctx context.Context, req *GetOffersReq, rsp *GetOffersRes) error {
	log.Print(" Get mqfs from  tracker:")
	matchFilter := make(map[string]interface{})
	data, count := dao.GetAllFromMongoSession(constants.MongoDB, constants.OfferAffiliateMQF, matchFilter)
	rsp.Data = data
	rsp.Count = count
	return nil
}

func (tracker *Tracker) Getrotatedmqfs(ctx context.Context, req *GetOffersReq, rsp *GetOffersRes) error {
	log.Print(" Get mqfs from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter["rotated"] = true
	data, count := dao.GetAllFromMongoSession(constants.MongoDB, constants.OfferAffiliateMQF, matchFilter)
	rsp.Data = data
	rsp.Count = count
	return nil
}

func (tracker *Tracker) Getmqf(ctx context.Context, req *GetMQFReq, rsp *GetOfferRes) error {
	log.Print(" Get mqf from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter[constants.OfferID] = req.Offerid
	matchFilter[constants.AffiliateID] = req.Affiliateid
	data, _ := dao.GetAllFromMongoSession(constants.MongoDB, constants.OfferAffiliateMQF, matchFilter)
	rsp.Data = data
	return nil
}

func (tracker *Tracker) Deletemqf(ctx context.Context, req *DeleteMQFReq, rsp *DeleteMQFRes) error {
	log.Print(" Delete MQF from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter[constants.OfferID] = req.Offerid
	matchFilter[constants.AffiliateID] = req.Affiliateid
	dao.DeleteRecordFromMongoSession(constants.MongoDB, constants.OfferAffiliateMQF, matchFilter)
	rsp.Id = req.Offerid
	rsp.Status = "Success"
	return nil
}

func (tracker *Tracker) Gettransactions(ctx context.Context, req *GetTransactionReq, rsp *GetTransactionRes) error {
	log.Print(" Get Transaction Log from  tracker:")
	sort := "utcdate"
	field := []string{"transactionID"}
	fields := helper.ConvertToBson(field...)
	limit := 10
	offset := 1
	log.Print(fields)
	matchFilter := make(map[string]interface{})
	//offer ids and affiliate ids filter
	if len(req.AffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.AffiliateIDs}
		matchFilter[constants.AffiliateID] = affFilter
	}
	if len(req.OfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.OfferIDs}
		matchFilter[constants.OfferID] = offFilter
	}
	data, _ := dao.QueryAllLogsFromMongoWithOffsetSession(constants.MongoDB, constants.ClickLog, limit, offset, sort, matchFilter, fields)
	rsp.Data = data
	return nil
}

func (tracker *Tracker) Dopostback(ctx context.Context, req *ManualPostbackReq, rsp *ManualPostbackRes) error {
	log.Print("Do manual postback from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter[constants.TransactionID] = req.TransactionID
	clickLog := dao.QueryLatestClickFromMongoSession(matchFilter)

	if len(clickLog) > 0 {
		url := createManualPostback(clickLog[0].OfferID, clickLog[0].AffiliateID, req.TransactionID, clickLog[0].ClickURL)
		rsp.URL = url
		rsp.TransactionID = req.TransactionID
		req.URL = url
		req.OfferID = clickLog[0].OfferID
		req.AffiliateID = clickLog[0].AffiliateID
		req.CreatedAt = time.Now().UTC()
		req.Hour = time.Now().UTC().Hour()
		req.Date = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		dao.InsertToMongo(constants.MongoDB, constants.ManualPostbackLog, req)
	} else {
		rsp.URL = ""
		rsp.Status = "Invalid transactionID"
	}
	return nil
}

func (tracker *Tracker) Getmanuallog(ctx context.Context, req *GetOfferReq, rsp *GetOfferRes) error {
	log.Print(" Get manual postback logs from  tracker:")
	matchFilter := make(map[string]interface{})
	data, _ := dao.GetAllFromMongoSession(constants.MongoDB, constants.ManualPostbackLog, matchFilter)
	affiliateMap, offerMap := createNameMap()
	for _, each := range data {
		each["offerName"] = offerMap[each["offerID"].(string)]
		each["affiliateName"] = affiliateMap[each["affiliateID"].(string)]
	}
	rsp.Data = data
	return nil
}

func createManualPostback(offerID string, affiliateID string, transactionID string, clickURL string) string {

	postbackURLTemplate := dao.GetOfferAffiliatePostbackTemplate(offerID, affiliateID)
	//Check offer affiliate specific template is available or not
	if len(postbackURLTemplate) == 0 || postbackURLTemplate == "" {
		postbackURLTemplate = dao.GetTemplateByAffiliateID(affiliateID)
		//Check affiliate specific template is available or not
		if len(postbackURLTemplate) == 0 || postbackURLTemplate == "" {
			postbackURLTemplate = "http://tk.adcamie.com/aff_lsr?transaction_id=" + transactionID
		} else {
			postbackURLTemplate = helper.ReplaceTemplateParameters(clickURL, postbackURLTemplate, transactionID)
		}
	} else {
		postbackURLTemplate = helper.ReplaceTemplateParameters(clickURL, postbackURLTemplate, transactionID)
	}

	return postbackURLTemplate

}

func (tracker *Tracker) Addadvertiser(ctx context.Context, req *SetAdvertiserReq, rsp *SettingRes) error {
	log.Print("Advertiser adding in tracker: ", req.AdvertiserName)
	//delete if exists
	filters := make(map[string]interface{})
	filters[constants.AdvertiserName] = req.AdvertiserName
	dao.DeleteFromMongoSession(constants.MongoDB, constants.Advertiser, filters)
	req.CreatedAt = time.Now().UTC()
	req.UpdatedAt = time.Now().UTC()
	dao.InsertToMongo(constants.MongoDB, constants.Advertiser, req)
	rsp.Status = "Success"
	rsp.Id = req.AdvertiserName
	return nil
}

func (tracker *Tracker) Getadvertisers(ctx context.Context, req *GetOffersReq, rsp *GetOffersRes) error {
	log.Print(" Get Advertisers from  tracker:")
	matchFilter := make(map[string]interface{})
	data, count := dao.GetAllFromMongoSession(constants.MongoDB, constants.Advertiser, matchFilter)
	rsp.Data = data
	rsp.Count = count
	return nil
}

func (tracker *Tracker) Getadvertiser(ctx context.Context, req *GetOfferReq, rsp *GetOfferRes) error {
	log.Print(" Get Advertiser from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter[constants.AdvertiserName] = req.Id
	data, _ := dao.GetAllFromMongoSession(constants.MongoDB, constants.Advertiser, matchFilter)
	rsp.Data = data
	return nil
}

func (tracker *Tracker) Deleteadvertiser(ctx context.Context, req *DeleteAdvertiserReq, rsp *DeleteAdvertiserRes) error {
	log.Print(" Delete Advertiser from  tracker:")
	matchFilter := make(map[string]interface{})
	matchFilter[constants.AdvertiserName] = req.AdvertiserName
	dao.DeleteRecordFromMongoSession(constants.MongoDB, constants.Advertiser, matchFilter)
	rsp.Id = req.AdvertiserName
	rsp.Status = "Success"
	return nil
}
