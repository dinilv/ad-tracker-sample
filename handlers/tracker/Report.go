package v1

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	"github.com/micro/go-micro/server"
	"gopkg.in/mgo.v2/bson"
)

var layout, inputDateFormat string

func init() {
	layout = "2006-01-02"
	inputDateFormat = "02/01/2006"
}

type GetReportReq struct {
	Key              string   `json:"key,omitempty"`
	Token            string   `json:"token,omitempty"`
	SortField        string   `json:"sort,omitempty"`
	StartDate        string   `json:"start_date,omitempty"`
	EndDate          string   `json:"end_date,omitempty"`
	StartHour        string   `json:"start_hour,omitempty"`
	EndHour          string   `json:"end_hour,omitempty"`
	OfferIDs         []string `json:"offer_ids,omitempty"`
	RecvOfferIDs     []string `json:"recv_offer_ids,omitempty"`
	RecvAffiliateIDs []string `json:"recv_affiliate_ids,omitempty"`
	AffiliateIDs     []string `json:"affiliate_ids,omitempty"`
	Fields           []string `json:"fields,omitempty"`
	Type             string   `json:"type,omitempty"`
	FileName         string   `json:"filename,omitempty"`
	Page             int      `json:"page,omitempty"`
	Status           []string `json:"status,omitempty"`
	Transaction      string   `json:"transaction,omitempty"`
	MatchFilter      map[string]interface{}
	IDs              map[string]interface{}
	SumFields        map[string]interface{}
}

type GetReportRes struct {
	Report []model.APIReport
	Data   []bson.M
	Page   int
	Count  int
}

type ReportHandler interface {
	Basic(context.Context, *GetReportReq, *GetReportRes) error
	Rotation(context.Context, *GetReportReq, *GetReportRes) error
	Adcamie(context.Context, *GetReportReq, *GetReportRes) error
	AdcamieMetaReport(context.Context, *GetReportReq, *GetReportRes) error
	AdcamieMetaReportBigQuery(context.Context, *GetReportReq, *GetReportRes) error
	AdcamieOptimised(context.Context, *GetReportReq, *GetReportRes) error
	Dashboard(context.Context, *GetReportReq, *GetReportRes) error
}

type Report struct {
	ReportHandler
}

func RegisterReportHandler(s server.Server, hdlr ReportHandler) {
	log.Print("Getting Tracker Report Handler")
	s.Handle(s.NewHandler(&Report{hdlr}))
}

// For sorting basic reporting
type ByOfferId []model.APIReport

func (a ByOfferId) Len() int      { return len(a) }
func (a ByOfferId) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByOfferId) Less(i, j int) bool {
	if a[i].OfferID < a[j].OfferID {
		return true
	}
	if a[i].OfferID > a[j].OfferID {
		return false
	}
	return a[i].AffiliateID < a[j].AffiliateID
}

func (report *Report) AdcamieMetaReport(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {
	return nil
}

func (report *Report) Basic(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {

	//get all offers and affiliates
	affiliateMap, offerMap := createNameMap()
	matchFilter, idFields := parseReportRequest(req)

	//for reports to collect all, imps clicks & conversions
	reportMap := make(map[string]model.APIReport)

	//sum for impressions
	sumFields := make(map[string]interface{})
	sumFields["impressions"] = map[string]interface{}{"$sum": 1}
	getImpressionReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	sumFields = make(map[string]interface{})
	sumFields["clicks"] = map[string]interface{}{"$sum": 1}
	getClickReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	sumFields = make(map[string]interface{})
	sumFields["conversions"] = map[string]interface{}{"$sum": 1}
	getConversionReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	sumFields = make(map[string]interface{})
	sumFields["events"] = map[string]interface{}{"$sum": 1}
	getPostEventReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	var reports []model.APIReport
	for _, v := range reportMap {
		reports = append(reports, v)
	}
	sort.Sort(ByOfferId(reports))
	rsp.Report = reports
	return nil
}

func (report *Report) Rotation(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {

	//get all offers and affiliates
	affiliateNameMap, offerNameMap := createNameMap()
	matchFilter, idFields := parseReportRequest(req)

	recvOfferID := map[string]interface{}{constants.ReceivedOfferID: "$" + constants.ReceivedOfferID}
	idFields[constants.ReceivedOfferID] = recvOfferID

	receivedAffiliateId := map[string]interface{}{constants.ReceivedAffiliateID: "$" + constants.ReceivedAffiliateID}
	idFields[constants.ReceivedAffiliateID] = receivedAffiliateId

	matchFilter["status"] = constants.Rotated

	//for reports to collect rotated clicks & conversions
	rotationReportMap := make(map[string]model.APIReport)

	sumFields := make(map[string]interface{})
	sumFields["clicks"] = map[string]interface{}{"$sum": 1}
	getRotatedClickReport(matchFilter, idFields, sumFields, rotationReportMap, affiliateNameMap, offerNameMap)

	sumFields = make(map[string]interface{})
	sumFields["conversions"] = map[string]interface{}{"$sum": 1}
	statusFilter := map[string]interface{}{"$in": []string{constants.RotatedSent, constants.RotatedUnSent}}
	matchFilter["status"] = statusFilter
	getRotatedConversionReport(matchFilter, idFields, sumFields, rotationReportMap, affiliateNameMap, offerNameMap)

	var reports []model.APIReport
	for _, v := range rotationReportMap {
		reports = append(reports, v)
	}

	rsp.Report = reports

	return nil
}

func (report *Report) Adcamie(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {

	//get all offers and affiliates
	affiliateMap, offerMap := createRefMap()
	matchFilter, idFields := parseReportRequest(req)

	date := map[string]interface{}{"date": "$date"}
	idFields["date"] = date

	hour := map[string]interface{}{"hour": "$hour"}
	idFields["hour"] = hour

	//for reports to collect all, imps clicks & conversions
	reportMap := make(map[string]model.APIReport)

	sumFields := make(map[string]interface{})
	sumFields["impressions"] = map[string]interface{}{"$sum": 1}
	getImpressionReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	sumFields = make(map[string]interface{})
	sumFields["clicks"] = map[string]interface{}{"$sum": 1}
	getClickReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	sumFields = make(map[string]interface{})
	sumFields["conversions"] = map[string]interface{}{"$sum": 1}
	getConversionReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	sumFields = make(map[string]interface{})
	sumFields["events"] = map[string]interface{}{"$sum": 1}
	getPostEventReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	//add rotated clicks
	receivedOfferID := map[string]interface{}{constants.OfferID: "$" + constants.ReceivedOfferID}
	idFields[constants.OfferID] = receivedOfferID

	receivedAffiliateID := map[string]interface{}{constants.AffiliateID: "$" + constants.ReceivedAffiliateID}
	idFields[constants.AffiliateID] = receivedAffiliateID

	matchFilter["status"] = constants.Rotated
	sumFields = make(map[string]interface{})
	sumFields["clicks"] = map[string]interface{}{"$sum": 1}
	getRotatedClickReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	//for rotated conversions
	sumFields = make(map[string]interface{})
	sumFields["conversions"] = map[string]interface{}{"$sum": 1}

	statusFilter := map[string]interface{}{"$in": []string{constants.RotatedSent, constants.RotatedUnSent}}
	matchFilter["status"] = statusFilter
	getRotatedConversionReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	//rotated conversions forwarded
	fwdOfferID := map[string]interface{}{constants.FwdOfferID: "$" + constants.OfferID}
	idFields[constants.FwdOfferID] = fwdOfferID
	getRotatedConversionForwardReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	//rotated clicks forwarded
	matchFilter["status"] = constants.Rotated
	sumFields = make(map[string]interface{})
	sumFields["clicks"] = map[string]interface{}{"$sum": 1}
	getRotatedClicksForwardReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	//rotations conversions received
	matchFilter, idFields = parseReportRequest(req)
	idFields["date"] = date
	idFields["hour"] = hour

	matchFilter["status"] = statusFilter

	sumFields = make(map[string]interface{})
	sumFields["conversions"] = map[string]interface{}{"$sum": 1}

	receivedOfferID = map[string]interface{}{constants.ReceivedOfferID: "$" + constants.ReceivedOfferID}
	idFields[constants.ReceivedOfferID] = receivedOfferID

	receivedAffiliateID = map[string]interface{}{constants.ReceivedAffiliateID: "$" + constants.ReceivedAffiliateID}
	idFields[constants.ReceivedAffiliateID] = receivedAffiliateID

	getRotatedConversionReceivedReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	//rotated clicks received
	matchFilter["status"] = constants.Rotated
	sumFields = make(map[string]interface{})
	sumFields["clicks"] = map[string]interface{}{"$sum": 1}

	getRotatedClicksReceivedReport(matchFilter, idFields, sumFields, reportMap, affiliateMap, offerMap)

	var reports []model.APIReport
	for _, v := range reportMap {
		reports = append(reports, v)
		//save to adcamie report datarow
		v.UtcDate, _ = time.Parse(layout, v.Date)
		deleteFilter := make(map[string]interface{})
		deleteFilter["offer_id"] = v.OfferID
		deleteFilter["affiliate_id"] = v.AffiliateID
		deleteFilter["date"] = v.Date
		deleteFilter["hour"] = v.Hour
		dao.DeleteFromMongoSession(constants.MongoDB, constants.AdcamieReport, deleteFilter)
		dao.InsertToMongoSession(constants.MongoDB, constants.AdcamieReport, v)
	}
	rsp.Report = reports

	return nil
}

func (report *Report) AdcamieOptimised(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {

	log.Println("Inside optimised Adcamie report")
	matchFilter := parseReportRequestForAdcamie(req)
	results := dao.QueryAdcamieReportFromMongoSession(constants.MongoDB, constants.AdcamieReport, matchFilter)
	rsp.Report = results
	return nil
}

func (report *Report) AdcamieMetaReportBigQuery(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {

	log.Println("Inside Bigquery Adcamie report")
	//get all offers and affiliates
	affiliateMap, offerMap := createRefMap()
	startDate, _ := time.Parse(inputDateFormat, req.StartDate)
	matchFilter := make(map[string]interface{})
	matchFilter["hour"] = req.StartHour
	matchFilter["date"] = startDate
	results := dao.QueryAdcamieReportFromMongoSession(constants.MongoDB, constants.AdcamieMetaReport, matchFilter)
	//create report map
	reportMap := make(map[string]model.APIReport)
	for _, apiReport := range results {
		key := apiReport.OfferID + "_" + apiReport.AffiliateID
		data, ok := reportMap[key]
		if ok {

		} else {
			data.AffiliateRefID = affiliateMap[apiReport.AffiliateID]
			data.OfferRefID = offerMap[apiReport.OfferID]
		}
	}
	for _, v := range reportMap {
		deleteFilter := make(map[string]interface{})
		deleteFilter["offer_id"] = v.OfferID
		deleteFilter["affiliate_id"] = v.AffiliateID
		deleteFilter["date"] = v.Date
		deleteFilter["hour"] = v.Hour
		dao.DeleteFromMongoSession(constants.MongoDB, constants.AdcamieBigQueryReport, deleteFilter)
		dao.InsertToMongoSession(constants.MongoDB, constants.AdcamieBigQueryReport, v)
	}

	return nil
}

func processReportAgg(v bson.M, affiliateMap map[string]string, offerMap map[string]string) (string, string, model.APIReport) {

	log.Println("Report Agg:-", v)

	idsMap := v["_id"].(bson.M)

	hour := 0
	if idsMap[constants.Hour] != nil {
		hourIDBson := idsMap[constants.Hour].(bson.M)
		hour = hourIDBson[constants.Hour].(int)
	}
	formattedDate := "01/01/2017"
	if idsMap[constants.Date] != nil {
		dateIDBson := idsMap[constants.Date].(bson.M)
		date := dateIDBson[constants.Date].(time.Time)
		formattedDate = date.Format(layout)
	}
	country := ""
	if idsMap["geo"] != nil {
		countryIDBson := idsMap["geo"].(bson.M)
		if countryIDBson["geo"] != nil {
			country = countryIDBson["geo"].(string)
		}

	}
	affiliateSub3 := ""
	if idsMap[constants.AffiliateSub3] != nil {
		affIDBson := idsMap[constants.AffiliateSub3].(bson.M)
		if affIDBson[constants.AffiliateSub3] != nil {
			affiliateSub3 = affIDBson[constants.AffiliateSub3].(string)
		}
	}
	status := ""
	if idsMap[constants.Status] != nil {
		statusIDBson := idsMap[constants.Status].(bson.M)
		status = statusIDBson[constants.Status].(string)
	}
	// Ids
	offerID := ""
	if idsMap[constants.OfferID] != nil {
		offerIDBson := idsMap[constants.OfferID].(bson.M)
		if offerIDBson[constants.OfferID] != nil {
			offerID = offerIDBson[constants.OfferID].(string)
		}
	}
	affiliateID := ""
	if idsMap[constants.AffiliateID] != nil {
		affiliateIDBson := idsMap[constants.AffiliateID].(bson.M)
		if affiliateIDBson[constants.AffiliateID] != nil {
			affiliateID = affiliateIDBson[constants.AffiliateID].(string)
		}
	}
	recvOfferID := ""
	if idsMap[constants.ReceivedOfferID] != nil {
		recvOfferIDBson := idsMap[constants.ReceivedOfferID].(bson.M)
		if recvOfferIDBson[constants.ReceivedOfferID] != nil {
			recvOfferID = recvOfferIDBson[constants.ReceivedOfferID].(string)
		}
	}
	recvAffiliateID := ""
	if idsMap[constants.ReceivedAffiliateID] != nil {
		recvAffiliateIDBson := idsMap[constants.ReceivedAffiliateID].(bson.M)
		if recvAffiliateIDBson[constants.ReceivedAffiliateID] != nil {
			recvAffiliateID = recvAffiliateIDBson[constants.ReceivedAffiliateID].(string)
		}
	}
	fwdOfferID := ""
	if idsMap[constants.FwdOfferID] != nil {
		fwdOfferIDBson := idsMap[constants.FwdOfferID].(bson.M)
		if fwdOfferIDBson[constants.FwdOfferID] != nil {
			fwdOfferID = fwdOfferIDBson[constants.FwdOfferID].(string)
		}
	}

	fwdAffiliateID := ""
	if idsMap[constants.FwdAffiliateID] != nil {
		fwdAffiliateIDBson := idsMap[constants.FwdAffiliateID].(bson.M)
		if fwdAffiliateIDBson[constants.FwdAffiliateID] != nil {
			fwdAffiliateID = fwdAffiliateIDBson[constants.FwdAffiliateID].(string)
		}
	}

	//create key for map
	keyToReportMap := strconv.Itoa(hour) + "_" + formattedDate + "_" + affiliateID + "_" + offerID + "_" +
		recvAffiliateID + "_" + recvOfferID + "_" + country + "_" + affiliateSub3

	//create report object
	report := model.APIReport{
		Hour:              hour,
		Date:              formattedDate,
		Country:           country,
		AffiliateSub3:     affiliateSub3,
		AffiliateID:       affiliateID,
		AffiliateName:     affiliateMap[affiliateID],
		AffiliateRefID:    affiliateMap[affiliateID],
		OfferID:           offerID,
		FwdOfferID:        fwdOfferID,
		FwdAffiliateID:    fwdAffiliateID,
		OfferName:         offerMap[offerID],
		OfferRefID:        offerMap[offerID],
		RecvOfferID:       recvOfferID,
		RecvOfferName:     offerMap[recvOfferID],
		RecvAffiliateID:   recvAffiliateID,
		RecvAffiliateName: affiliateMap[recvAffiliateID],
	}

	return keyToReportMap, status, report
}

func processRotatedReportAgg(v bson.M, affiliateMap map[string]string, offerMap map[string]string) (string, string, model.APIReport) {

	log.Println("Results:-", v)
	idsMap := v["_id"].(bson.M)

	hour := 0
	if idsMap[constants.Hour] != nil {
		hourIDBson := idsMap[constants.Hour].(bson.M)
		hour = hourIDBson[constants.Hour].(int)
	}
	formattedDate := "01/01/2017"
	if idsMap[constants.Date] != nil {
		dateIDBson := idsMap[constants.Date].(bson.M)
		date := dateIDBson[constants.Date].(time.Time)
		formattedDate = date.Format(layout)
	}
	status := ""
	if idsMap[constants.Status] != nil {
		statusIDBson := idsMap[constants.Status].(bson.M)
		status = statusIDBson[constants.Status].(string)
	}
	offerID := ""
	if idsMap[constants.OfferID] != nil {
		offerIDBson := idsMap[constants.OfferID].(bson.M)
		if offerIDBson[constants.OfferID] != nil {
			offerID = offerIDBson[constants.OfferID].(string)
		}
	}
	affiliateID := ""
	if idsMap[constants.AffiliateID] != nil {
		affiliateIDBson := idsMap[constants.AffiliateID].(bson.M)
		if affiliateIDBson[constants.AffiliateID] != nil {
			affiliateID = affiliateIDBson[constants.AffiliateID].(string)
		}
	}
	recvOfferID := ""
	if idsMap[constants.ReceivedOfferID] != nil {
		recvOfferIDBson := idsMap[constants.ReceivedOfferID].(bson.M)
		if recvOfferIDBson[constants.ReceivedOfferID] != nil {
			recvOfferID = recvOfferIDBson[constants.ReceivedOfferID].(string)
		}
	}
	recvAffiliateID := ""
	if idsMap[constants.ReceivedAffiliateID] != nil {
		recvAffiliateIDBson := idsMap[constants.ReceivedAffiliateID].(bson.M)
		if recvAffiliateIDBson[constants.ReceivedAffiliateID] != nil {
			recvAffiliateID = recvAffiliateIDBson[constants.ReceivedAffiliateID].(string)
		}
	}

	//create key for map
	keyToReportMap := strconv.Itoa(hour) + "_" + formattedDate + "_" + affiliateID + "_" + offerID + "_" +
		"" + "_" + "" + "_" + "" + "_" + ""

	//create report object
	report := model.APIReport{
		Hour:              hour,
		Date:              formattedDate,
		AffiliateID:       affiliateID,
		AffiliateName:     affiliateMap[affiliateID],
		AffiliateRefID:    affiliateMap[affiliateID],
		OfferID:           offerID,
		OfferName:         offerMap[offerID],
		OfferRefID:        offerMap[offerID],
		RecvOfferID:       recvOfferID,
		RecvOfferName:     offerMap[recvOfferID],
		RecvAffiliateID:   recvAffiliateID,
		RecvAffiliateName: affiliateMap[recvAffiliateID],
	}

	return keyToReportMap, status, report
}

func createAllMap() (map[string]string, map[string]string, map[string]string, map[string]string) {

	//get all offers and affiliates
	offerResults := dao.QueryAllOffers()
	affiliateResults := dao.QueryAllAffiliates()

	offerNameMap := make(map[string]string)
	offerRefMap := make(map[string]string)
	for _, v := range offerResults {
		offerNameMap[v.OfferID] = strings.Replace(v.OfferName, ",", " ", -1)
		offerRefMap[v.OfferID] = v.OfferRefID
	}
	affiliateNameMap := make(map[string]string)
	affiliateRefMap := make(map[string]string)
	for _, v := range affiliateResults {
		affiliateNameMap[v.AffiliateID] = strings.Replace(v.AffiliateName, ",", " ", -1)
		affiliateRefMap[v.AffiliateID] = v.AffiliateRefID
	}
	return affiliateNameMap, offerNameMap, affiliateRefMap, offerRefMap

}

func createNameMap() (map[string]string, map[string]string) {

	//get all offers and affiliates
	offerResults := dao.QueryAllOffers()
	affiliateResults := dao.QueryAllAffiliates()

	offerRefMap := make(map[string]string)
	for _, v := range offerResults {
		offerRefMap[v.OfferID] = strings.Replace(v.OfferName, ",", " ", -1)
	}

	affiliateRefMap := make(map[string]string)
	for _, v := range affiliateResults {
		affiliateRefMap[v.AffiliateID] = strings.Replace(v.AffiliateName, ",", " ", -1)
	}
	return affiliateRefMap, offerRefMap

}

func createRefMap() (map[string]string, map[string]string) {

	//get all offers and affiliates
	offerResults := dao.QueryAllOffers()
	affiliateResults := dao.QueryAllAffiliates()

	offerRefMap := make(map[string]string)
	for _, v := range offerResults {
		offerRefMap[v.OfferID] = v.OfferRefID
	}

	affiliateRefMap := make(map[string]string)
	for _, v := range affiliateResults {
		affiliateRefMap[v.AffiliateID] = v.AffiliateRefID
	}
	return affiliateRefMap, offerRefMap

}

func createOfferMap() map[string]string {

	//get all offers and affiliates
	offerResults := dao.QueryAllOffers()
	offerRefMap := make(map[string]string)
	for _, v := range offerResults {
		offerRefMap[v.OfferID] = strings.Replace(v.OfferName, ",", " ", -1)
	}

	return offerRefMap

}

func parseReportRequest(req *GetReportReq) (map[string]interface{}, map[string]interface{}) {

	//Create filters with mandatory/optional parameters
	matchFilter := make(map[string]interface{})
	if strings.Compare(req.StartDate, req.EndDate) == 0 {
		hour1, _ := strconv.Atoi(req.StartHour)
		hour2, _ := strconv.Atoi(req.EndHour)
		hourFilter := map[string]interface{}{"$gte": hour1, "$lte": hour2}
		matchFilter["hour"] = hourFilter
	}
	start, _ := time.Parse(inputDateFormat, req.StartDate)
	end, _ := time.Parse(inputDateFormat, req.EndDate)
	dateFilter := map[string]interface{}{"$gte": start, "$lte": end}
	matchFilter["date"] = dateFilter

	//offer ids and affiliate ids filter
	if len(req.AffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.AffiliateIDs}
		matchFilter[constants.AffiliateID] = affFilter
	}
	if len(req.OfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.OfferIDs}
		matchFilter[constants.OfferID] = offFilter
	}
	if len(req.RecvOfferIDs) > 0 {
		recvOffFilter := map[string]interface{}{"$in": req.RecvOfferIDs}
		matchFilter[constants.ReceivedOfferID] = recvOffFilter
	}
	if len(req.RecvAffiliateIDs) > 0 {
		recvAffFilter := map[string]interface{}{"$in": req.RecvAffiliateIDs}
		matchFilter[constants.ReceivedAffiliateID] = recvAffFilter
	}

	idFields := make(map[string]interface{})

	offerId := map[string]interface{}{"offerID": "$offerID"}
	idFields["offerID"] = offerId

	affiliateId := map[string]interface{}{"affiliateID": "$affiliateID"}
	idFields["affiliateID"] = affiliateId

	for _, field := range req.Fields {

		if strings.Compare(field, "geo") == 0 {
			fieldId := map[string]interface{}{"geo": "$geo.countryname"}
			idFields[field] = fieldId
		} else if strings.Compare(field, "clickGeo") == 0 {
			fieldId := map[string]interface{}{"geo": "$clickGeo.countryname"}
			idFields[field] = fieldId
		} else {
			fieldId := map[string]interface{}{field: "$" + field}
			idFields[field] = fieldId

		}
	}

	status := map[string]interface{}{"status": "$status"}
	idFields["status"] = status

	return matchFilter, idFields

}

func parseReportRequestForAdcamie(req *GetReportReq) map[string]interface{} {

	//Create filters with mandatory/optional parameters
	matchFilter := make(map[string]interface{})
	if strings.Compare(req.StartDate, req.EndDate) == 0 {
		hour1, _ := strconv.Atoi(req.StartHour)
		hour2, _ := strconv.Atoi(req.EndHour)
		hourFilter := map[string]interface{}{"$gte": hour1, "$lte": hour2}
		matchFilter["hour"] = hourFilter
	}
	start, _ := time.Parse(inputDateFormat, req.StartDate)
	end, _ := time.Parse(inputDateFormat, req.EndDate)
	dateFilter := map[string]interface{}{"$gte": start, "$lte": end}
	matchFilter["utcdate"] = dateFilter

	//offer ids and affiliate ids filter
	if len(req.AffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.AffiliateIDs}
		matchFilter[constants.AFFILIATE_ID] = affFilter
	}
	if len(req.OfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.OfferIDs}
		matchFilter[constants.OFFER_ID] = offFilter
	}

	return matchFilter

}

func getImpressionReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {

	impAggs := dao.GetInterceptorReportForAPIMongoSession(constants.MongoDB, constants.ImpressionLog, matchFilter, idFields, sumFields)

	//process impression report
	for _, v := range impAggs {

		impressions := v["impressions"].(int)
		keyToReportMap, status, reportData := processReportAgg(v, affiliateMap, offerMap)

		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]
		if ok {
			//total of unique and duplicate
			data.Impressions = impressions + data.Impressions
			if strings.Compare(status, constants.Unique) == 0 {
				data.UniqueImpressions = impressions
			} else {
				data.DuplicateImpressions = impressions
			}
			reportMap[keyToReportMap] = data

		} else {
			reportData.Impressions = impressions
			if strings.Compare(status, constants.Unique) == 0 {
				reportData.UniqueImpressions = impressions
			} else {
				reportData.DuplicateImpressions = impressions
			}
			reportMap[keyToReportMap] = reportData

		}

	}

}

func getRotatedClickReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {

	log.Println("On Rotations Clicks")

	clickAggs := dao.GetInterceptorReportForAPIMongoSession(constants.MongoDB, constants.ClickLog, matchFilter, idFields, sumFields)
	//parse aggregation click results back to struct and add to map
	for _, v := range clickAggs {

		clicks := v["clicks"].(int)
		keyToReportMap, _, reportData := processReportAgg(v, affiliateMap, offerMap)
		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]
		if ok {
			//total of unique,duplicate & rotated
			data.RotatedClicksFwd = clicks
			reportMap[keyToReportMap] = data
		} else {
			reportData.RotatedClicksFwd = clicks
			//check offerID is available or Not
			if len(reportData.OfferID) > 0 {
				reportMap[keyToReportMap] = reportData
			}

		}
	}

}

func getClickReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {

	clickAggs := dao.GetInterceptorReportForAPISession(constants.MongoDB, constants.ClickLog, matchFilter, idFields, sumFields)

	//parse aggregation click results back to struct and add to map
	for _, v := range clickAggs {

		clicks := v["clicks"].(int)
		keyToReportMap, status, reportData := processReportAgg(v, affiliateMap, offerMap)
		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]
		if ok {
			//total of unique,duplicate & rotated
			data.Clicks = clicks + data.Clicks
			if strings.Compare(status, constants.Unique) == 0 {
				data.UniqueClicks = clicks
			} else if strings.Compare(status, constants.Duplicate) == 0 {
				data.DuplicateClicks = clicks
			} else if strings.Compare(status, constants.Rotated) == 0 {
				data.UniqueClicks = clicks
			}
			reportMap[keyToReportMap] = data
		} else {
			reportData.Clicks = clicks
			if strings.Compare(status, constants.Unique) == 0 {
				//add entry to map
				reportData.UniqueClicks = clicks
			} else if strings.Compare(status, constants.Duplicate) == 0 {
				//add entry to map
				reportData.DuplicateClicks = clicks
			} else {
				//add entry to map
				reportData.UniqueClicks = clicks
			}
			reportMap[keyToReportMap] = reportData
		}
	}

}

func getRotatedClicksReceivedReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {

	log.Println("On Rotated Received Clicks")

	clickAggs := dao.GetInterceptorReportForAPISession(constants.MongoDB, constants.ClickLog, matchFilter, idFields, sumFields)

	//process rotated clicks forwarded
	for _, v := range clickAggs {

		clicks := v["clicks"].(int)
		keyToReportMap, _, reportData := processRotatedReportAgg(v, affiliateMap, offerMap)

		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]

		if ok {
			var rotatedData model.RotationDetails
			var rotatedDataAvailable bool
			var rotations []model.RotationDetails
			for _, datarotateed := range data.ReceivedRotations {
				if (strings.Compare(datarotateed.OfferID, reportData.RecvOfferID) == 0) && (strings.Compare(datarotateed.AffiliateID, reportData.RecvAffiliateID) == 0) {
					rotatedData = datarotateed
					rotatedDataAvailable = true
				} else {
					rotations = append(rotations, datarotateed)
				}
			}
			if !rotatedDataAvailable {
				var fwdRotations = model.RotationDetails{
					OfferID:           reportData.RecvOfferID,
					OfferRefID:        offerMap[reportData.RecvOfferID],
					AffiliateID:       reportData.RecvAffiliateID,
					AffiliateRefID:    affiliateMap[reportData.RecvAffiliateID],
					Clicks:            clicks,
					Conversions:       0,
					SentConversions:   0,
					UnSentConversions: 0,
				}
				rotations = append(rotations, fwdRotations)
			} else {
				//add to specific rotated data
				rotatedData.Clicks = rotatedData.Clicks + clicks
				rotations = append(rotations, rotatedData)
			}

			data.ReceivedRotations = rotations
			reportMap[keyToReportMap] = data

		} else {
			log.Println("Entry Doesnt Exist")
		}

	}
}

func getRotatedClicksForwardReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {

	log.Println("On Rotatated Forward Clicks")

	clickAggs := dao.GetInterceptorReportForAPISession(constants.MongoDB, constants.ClickLog, matchFilter, idFields, sumFields)

	//process rotated conversions forwarded
	for _, v := range clickAggs {

		clicks := v["clicks"].(int)
		keyToReportMap, _, reportData := processReportAgg(v, affiliateMap, offerMap)

		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]

		if ok {
			var rotatedData model.RotationDetails
			var rotatedDataAvailable bool
			var rotations []model.RotationDetails
			for _, datarotateed := range data.Rotations {
				if strings.Compare(datarotateed.OfferID, reportData.FwdOfferID) == 0 {
					rotatedData = datarotateed
					rotatedDataAvailable = true
				} else {
					rotations = append(rotations, datarotateed)
				}
			}
			if !rotatedDataAvailable {
				var fwdRotations = model.RotationDetails{
					OfferID:           reportData.FwdOfferID,
					OfferRefID:        offerMap[reportData.FwdOfferID],
					AffiliateID:       constants.TRACKER_MEDIA,
					AffiliateRefID:    constants.TRACKER_MEDIA,
					Clicks:            clicks,
					Conversions:       0,
					SentConversions:   0,
					UnSentConversions: 0,
				}
				rotations = append(rotations, fwdRotations)
			} else {
				//add to specific rotated data
				rotatedData.Clicks = rotatedData.Clicks + clicks
				rotations = append(rotations, rotatedData)
			}
			data.Rotations = rotations
			reportMap[keyToReportMap] = data

		} else {
			log.Println("Entry doesnt exist")

		}

	}

}

func getConversionReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {

	convAggs := dao.GetInterceptorReportForAPISession(constants.MongoDB, constants.PostBackLog, matchFilter, idFields, sumFields)

	//process conversion report
	for _, v := range convAggs {

		conversions := v["conversions"].(int)
		keyToReportMap, status, reportData := processReportAgg(v, affiliateMap, offerMap)

		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]
		if ok {
			data.Conversions = data.Conversions + conversions
			if strings.Compare(constants.Sent, status) == 0 {
				data.SentConversions = conversions
			} else if strings.Compare(constants.UnSent, status) == 0 {
				data.UnSentConversions = conversions
			} else if strings.Compare(constants.RotatedSent, status) == 0 {
				data.SentConversions = conversions
			} else if strings.Compare(constants.RotatedUnSent, status) == 0 {
				data.UnSentConversions = conversions
			}
			reportMap[keyToReportMap] = data

		} else {
			reportData.Conversions = conversions
			if strings.Compare(status, constants.Sent) == 0 {
				//add entry to map
				reportData.SentConversions = conversions
			} else if strings.Compare(status, constants.UnSent) == 0 {
				//add entry to map
				reportData.UnSentConversions = conversions
			} else if strings.Compare(status, constants.RotatedSent) == 0 {
				//add entry to map
				reportData.SentConversions = conversions
			} else if strings.Compare(status, constants.RotatedUnSent) == 0 {
				//add entry to map
				reportData.UnSentConversions = conversions
			}
			reportMap[keyToReportMap] = reportData

		}

	}
}

func getRotatedConversionReceivedReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {

	log.Println("On Rotated Received Conversions")

	convAggs := dao.GetInterceptorReportForAPISession(constants.MongoDB, constants.PostBackLog, matchFilter, idFields, sumFields)

	//process rotated conversions forwarded
	for _, v := range convAggs {

		conversions := v["conversions"].(int)
		keyToReportMap, status, reportData := processRotatedReportAgg(v, affiliateMap, offerMap)

		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]

		if ok {
			var rotatedData model.RotationDetails
			var rotatedDataAvailable bool
			var rotations []model.RotationDetails
			for _, datarotateed := range data.ReceivedRotations {
				if (strings.Compare(datarotateed.OfferID, reportData.RecvOfferID) == 0) && (strings.Compare(datarotateed.AffiliateID, reportData.RecvAffiliateID) == 0) {
					rotatedData = datarotateed
					rotatedDataAvailable = true
				} else {
					rotations = append(rotations, datarotateed)
				}
			}
			if strings.Compare(constants.RotatedSent, status) == 0 {
				if !rotatedDataAvailable {
					var fwdRotations = model.RotationDetails{
						OfferID:           reportData.RecvOfferID,
						OfferRefID:        offerMap[reportData.RecvOfferID],
						AffiliateID:       reportData.RecvAffiliateID,
						AffiliateRefID:    affiliateMap[reportData.RecvAffiliateID],
						Conversions:       conversions,
						SentConversions:   conversions,
						UnSentConversions: 0,
					}
					rotations = append(rotations, fwdRotations)
				} else {
					//add to specific rotated data
					rotatedData.Conversions = rotatedData.Conversions + conversions
					rotatedData.SentConversions = rotatedData.SentConversions + conversions
					rotations = append(rotations, rotatedData)
				}

			} else if strings.Compare(constants.RotatedUnSent, status) == 0 {
				if !rotatedDataAvailable {
					var fwdRotations = model.RotationDetails{
						OfferID:           reportData.RecvOfferID,
						OfferRefID:        offerMap[reportData.RecvOfferID],
						AffiliateID:       reportData.RecvAffiliateID,
						AffiliateRefID:    affiliateMap[reportData.RecvAffiliateID],
						Conversions:       conversions,
						SentConversions:   0,
						UnSentConversions: conversions,
					}
					rotations = append(rotations, fwdRotations)
				} else {
					//add to specific rotated data
					rotatedData.Conversions = rotatedData.Conversions + conversions
					rotatedData.UnSentConversions = rotatedData.UnSentConversions + conversions
					rotations = append(rotations, rotatedData)
				}

			}
			data.ReceivedRotations = rotations
			reportMap[keyToReportMap] = data

		} else {
			log.Println("Entry Doesnt Exist")
		}

	}

}

func getRotatedConversionForwardReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {

	log.Println("On Rotatated Forward Conversions")

	convAggs := dao.GetInterceptorReportForAPISession(constants.MongoDB, constants.PostBackLog, matchFilter, idFields, sumFields)

	//process rotated conversions forwarded
	for _, v := range convAggs {

		conversions := v["conversions"].(int)
		keyToReportMap, status, reportData := processReportAgg(v, affiliateMap, offerMap)

		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]

		if ok {
			var rotatedData model.RotationDetails
			var rotatedDataAvailable bool
			var rotations []model.RotationDetails
			for _, datarotateed := range data.Rotations {
				if strings.Compare(datarotateed.OfferID, reportData.FwdOfferID) == 0 {
					rotatedData = datarotateed
					rotatedDataAvailable = true
				} else {
					rotations = append(rotations, datarotateed)
				}
			}
			if strings.Compare(constants.RotatedSent, status) == 0 {
				if !rotatedDataAvailable {
					var fwdRotations = model.RotationDetails{
						OfferID:           reportData.FwdOfferID,
						OfferRefID:        offerMap[reportData.FwdOfferID],
						AffiliateID:       constants.TRACKER_MEDIA,
						AffiliateRefID:    constants.TRACKER_MEDIA,
						Conversions:       conversions,
						SentConversions:   conversions,
						UnSentConversions: 0,
					}
					rotations = append(rotations, fwdRotations)
				} else {
					//add to specific rotated data
					rotatedData.Conversions = rotatedData.Conversions + conversions
					rotatedData.SentConversions = rotatedData.SentConversions + conversions
					rotations = append(rotations, rotatedData)
				}

			} else if strings.Compare(constants.RotatedUnSent, status) == 0 {
				if !rotatedDataAvailable {
					var fwdRotations = model.RotationDetails{
						OfferID:           reportData.FwdOfferID,
						OfferRefID:        offerMap[reportData.FwdOfferID],
						AffiliateID:       constants.TRACKER_MEDIA,
						AffiliateRefID:    constants.TRACKER_MEDIA,
						Conversions:       conversions,
						SentConversions:   0,
						UnSentConversions: conversions,
					}
					rotations = append(rotations, fwdRotations)
				} else {
					//add to specific rotated data
					rotatedData.Conversions = rotatedData.Conversions + conversions
					rotatedData.UnSentConversions = rotatedData.UnSentConversions + conversions
					rotations = append(rotations, rotatedData)
				}

			}
			data.Rotations = rotations
			reportMap[keyToReportMap] = data

		} else {
			log.Println("Entry doesnt exist")

		}

	}

}

func getRotatedConversionReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {

	log.Println("On Rotated Conversions")

	convAggs := dao.GetInterceptorReportForAPISession(constants.MongoDB, constants.PostBackLog, matchFilter, idFields, sumFields)
	//process rotated conversions forwarded report
	for _, v := range convAggs {

		conversions := v["conversions"].(int)
		keyToReportMap, status, reportData := processReportAgg(v, affiliateMap, offerMap)

		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]
		if ok {
			data.RotatedConversionsFwd = data.RotatedConversionsFwd + conversions
			if strings.Compare(constants.RotatedSent, status) == 0 {
				data.SentRotatedConversionsFwd = conversions
			} else if strings.Compare(constants.RotatedUnSent, status) == 0 {
				data.UnSentRotatedConversionsFwd = conversions
			}
			reportMap[keyToReportMap] = data

		} else {
			reportData.RotatedConversionsFwd = conversions
			if strings.Compare(status, constants.RotatedSent) == 0 {
				reportData.SentRotatedConversionsFwd = conversions
			} else if strings.Compare(status, constants.RotatedUnSent) == 0 {
				reportData.UnSentRotatedConversionsFwd = conversions
			}
			reportMap[keyToReportMap] = reportData

		}

	}

}

func getPostEventReport(matchFilter map[string]interface{}, idFields map[string]interface{}, sumFields map[string]interface{}, reportMap map[string]model.APIReport, affiliateMap map[string]string, offerMap map[string]string) {
	eventAggs := dao.GetInterceptorReportForAPISession(constants.MongoDB, constants.PostEventLog, matchFilter, idFields, sumFields)
	//process conversion report
	for _, v := range eventAggs {

		events := v["events"].(int)
		keyToReportMap, _, reportData := processReportAgg(v, affiliateMap, offerMap)

		//check map entry exists or Not
		data, ok := reportMap[keyToReportMap]
		if ok {
			data.Events = events
			reportMap[keyToReportMap] = data

		} else {
			//add entry to map
			reportData.Events = events
			reportMap[keyToReportMap] = reportData

		}

	}
}
