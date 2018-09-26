package v1

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	helper "github.com/adcamie/adserver/helpers"
	"github.com/micro/go-micro/server"
	"gopkg.in/mgo.v2/bson"
)

type LogHandler interface {
	All(context.Context, *GetReportReq, *GetReportRes) error
	Totalcount(context.Context, *GetReportReq, *GetReportRes) error
	CSV(context.Context, *GetReportReq, *GetReportRes) error
	MO(context.Context, *GetReportReq, *GetReportRes) error
}

type Log struct {
	LogHandler
}

func RegisterLogHandler(s server.Server, hdlr LogHandler) {
	log.Print("Getting Tracker Setting Log Handler")
	s.Handle(s.NewHandler(&Log{hdlr}))
}

func (rlog *Log) CSV(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {
	log.Println("Received CSV service request." + req.FileName)

	//Create filters with mandatory/optional parameters
	log.Print("Request on CSV writing log :", req)
	var fieldsFormatted = ""
	for _, field := range req.Fields {
		if strings.Compare("geo", field) == 0 {
			fieldsFormatted = fieldsFormatted + "geo.countryname,"
		} else if strings.Compare("clickGeo", field) == 0 {
			fieldsFormatted = fieldsFormatted + "clickGeo.countryname,"
		} else {
			fieldsFormatted = fieldsFormatted + field + ","
		}

	}
	//remove last comma character
	fieldsFormatted = fieldsFormatted[:len(fieldsFormatted)-1]
	queries := createCommandLineQuery(req)

	//jsonify queries
	json, _ := json.Marshal(queries)
	jsonQuery := string(json)
	//need to handle 90 days window
	var arguments = []string{"--port", "27017", "--host", "10.148.0.4", "-d", "Tracker", "-c", req.Type, "--query", jsonQuery, "-f", fieldsFormatted, "--type=csv", "--out=" + req.FileName}
	helper.RunMongoExport(arguments)

	rsp.Data = nil
	return nil
}

func (rlog *Log) All(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {
	log.Print(" Get  Log from  tracker:")
	limit := 50
	sort := "utcdate"
	offset := req.Page
	fields := helper.ConvertToBson(req.Fields...)
	reporttype := req.Type
	matchFilter := make(map[string]interface{})

	//offer ids and affiliate ids filter
	if len(req.AffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.AffiliateIDs}
		matchFilter["affiliateID"] = affFilter
	}
	if len(req.OfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.OfferIDs}
		matchFilter["offerID"] = offFilter
	}
	if len(req.RecvAffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.RecvAffiliateIDs}
		matchFilter["receivedAffiliateID"] = affFilter
	}
	if len(req.RecvOfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.RecvOfferIDs}
		matchFilter["receivedOfferID"] = offFilter
	}
	if len(req.Status) > 0 {
		statusFilter := map[string]interface{}{"$in": req.Status}
		matchFilter["status"] = statusFilter
	}
	if len(req.Transaction) > 0 {
		matchFilter["transactionID"] = req.Transaction
	}
	if len(req.Key) > 0 {
		matchFilter["selectedOfferForRotation"] = bson.M{"$regex": bson.RegEx{req.Key + `*`, ""}}
	}

	//date filter applicable tracker logs & other collections ignore the same
	if strings.Compare(req.Type, constants.FailedTransactions) != 0 && strings.Compare(req.Type, constants.ExhaustedOffer) != 0 &&
		strings.Compare(req.Type, constants.RotationGEOStack) != 0 && strings.Compare(req.Type, constants.ExhaustedOfferAffiliate) != 0 &&
		strings.Compare(req.Type, constants.RotationGroupStack) != 0 {
		//add date filters
		start, _ := time.Parse(inputDateFormat, req.StartDate)
		end, _ := time.Parse(inputDateFormat, req.EndDate)
		dateFilter := map[string]interface{}{"$gte": start, "$lte": end}
		matchFilter["date"] = dateFilter
		//add hour filter if required
		if strings.Compare(req.StartDate, req.EndDate) == 0 {
			hour1, _ := strconv.Atoi(req.StartHour)
			hour2, _ := strconv.Atoi(req.EndHour)
			hourFilter := map[string]interface{}{"$gte": hour1, "$lte": hour2}
			matchFilter["hour"] = hourFilter
		}

	}

	log.Println("Filter for querying", matchFilter, fields)
	data, count := dao.QueryAllLogsFromMongoWithOffsetSession(constants.MongoDB, reporttype, limit, offset, sort, matchFilter, fields)

	//real count if not ClickLog
	if strings.Compare(constants.ClickLog, req.Type) != 0 {
		count = dao.GetCountFromMongoSession(constants.MongoDB, reporttype, matchFilter)
	}
	//only if imp tracker logs to show
	if strings.Compare(req.Type, constants.FailedTransactions) != 0 && strings.Compare(req.Type, constants.ExhaustedOffer) != 0 &&
		strings.Compare(req.Type, constants.RotationGEOStack) != 0 && strings.Compare(req.Type, constants.ExhaustedOfferAffiliate) != 0 &&
		strings.Compare(req.Type, constants.RotationGroupStack) != 0 {

		affiliateMap, offerMap := createNameMap()
		for _, each := range data {
			if each["offerID"] != nil {
				each["offerName"] = offerMap[each["offerID"].(string)]
			}
			if each["affiliateID"] != nil {
				each["affiliateName"] = affiliateMap[each["affiliateID"].(string)]
			}
		}
	}
	rsp.Data = data
	rsp.Count = count

	return nil

}

func (rlog *Log) Totalcount(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {
	log.Print(" Get  count from  tracker:")

	reporttype := req.Type
	matchFilter := make(map[string]interface{})

	//offer ids and affiliate ids filter
	if len(req.AffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.AffiliateIDs}
		matchFilter["affiliateID"] = affFilter
	}
	if len(req.OfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.OfferIDs}
		matchFilter["offerID"] = offFilter
	}
	if len(req.RecvAffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.RecvAffiliateIDs}
		matchFilter["receivedAffiliateID"] = affFilter
	}
	if len(req.RecvOfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.RecvOfferIDs}
		matchFilter["receivedOfferID"] = offFilter
	}
	if len(req.Status) > 0 {
		statusFilter := map[string]interface{}{"$in": req.Status}
		matchFilter["status"] = statusFilter
	}
	if len(req.Transaction) > 0 {
		matchFilter["transactionID"] = req.Transaction
	}
	if len(req.Key) > 0 {
		matchFilter["key"] = map[string]interface{}{"$regex": "*" + req.Key + "*"}
	}
	//date filter applicable tracker logs & other collections ignore the same
	if strings.Compare(req.Type, constants.FailedTransactions) != 0 &&
		strings.Compare(req.Type, constants.RotationGEOStack) != 0 &&
		strings.Compare(req.Type, constants.RotationGroupStack) != 0 {
		//add date filters
		start, _ := time.Parse(inputDateFormat, req.StartDate)
		end, _ := time.Parse(inputDateFormat, req.EndDate)
		dateFilter := map[string]interface{}{"$gte": start, "$lte": end}
		matchFilter["date"] = dateFilter
		//add hour filter if required
		if strings.Compare(req.StartDate, req.EndDate) == 0 {
			hour1, _ := strconv.Atoi(req.StartHour)
			hour2, _ := strconv.Atoi(req.EndHour)
			hourFilter := map[string]interface{}{"$gte": hour1, "$lte": hour2}
			matchFilter["hour"] = hourFilter
		}

	}

	count := 5000
	//real count if not ClickLog
	if strings.Compare(constants.ClickLog, req.Type) != 0 {
		count = dao.GetCountFromMongoSession(constants.MongoDB, reporttype, matchFilter)
	}
	rsp.Count = count
	return nil

}

func (rlog *Log) MO(ctx context.Context, req *GetReportReq, rsp *GetReportRes) error {
	log.Print(" Get MO Log from  tracker:")
	limit := 50
	sort := "utcdate"
	offset := req.Page
	fields := helper.ConvertToBson(req.Fields...)
	reporttype := req.Type
	matchFilter := make(map[string]interface{})
	start, _ := time.Parse(inputDateFormat, req.StartDate)
	end, _ := time.Parse(inputDateFormat, req.EndDate)
	dateFilter := map[string]interface{}{"$gte": start, "$lte": end}
	matchFilter["date"] = dateFilter
	//offer ids and affiliate ids filter
	if len(req.AffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.AffiliateIDs}
		matchFilter["affiliateID"] = affFilter
	}
	if len(req.OfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.OfferIDs}
		matchFilter["offerID"] = offFilter
	}
	if len(req.Transaction) > 0 {
		matchFilter["transactionID"] = req.Transaction
	}
	matchFilter["offerType"] = "1"
	if strings.Compare(req.StartDate, req.EndDate) == 0 {
		hour1, _ := strconv.Atoi(req.StartHour)
		hour2, _ := strconv.Atoi(req.EndHour)
		hourFilter := map[string]interface{}{"$gte": hour1, "$lte": hour2}
		matchFilter["hour"] = hourFilter
	}
	data, count := dao.QueryAllLogsFromMongoWithOffsetSession(constants.MongoDB, reporttype, limit, offset, sort, matchFilter, fields)
	affiliateMap, offerMap := createNameMap()
	for _, each := range data {
		each["offerName"] = offerMap[each["offerID"].(string)]
		each["affiliateName"] = affiliateMap[each["affiliateID"].(string)]
	}
	rsp.Data = data
	rsp.Count = count

	return nil

}

func createCommandLineQuery(req *GetReportReq) map[string]interface{} {

	queries := make(map[string]interface{})
	start, _ := time.Parse(inputDateFormat, req.StartDate)
	end, _ := time.Parse(inputDateFormat, req.EndDate)

	dateFilter := map[string]interface{}{"$gte": map[string]interface{}{"$date": start}, "$lte": map[string]interface{}{"$date": end}}
	queries["date"] = dateFilter
	//offer ids and affiliate ids filter
	if len(req.AffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.AffiliateIDs}
		queries["affiliateID"] = affFilter
	}
	if len(req.OfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.OfferIDs}
		queries["offerID"] = offFilter
	}
	if len(req.RecvAffiliateIDs) > 0 {
		affFilter := map[string]interface{}{"$in": req.RecvAffiliateIDs}
		queries["receivedAffiliateID"] = affFilter
	}
	if len(req.RecvOfferIDs) > 0 {
		offFilter := map[string]interface{}{"$in": req.RecvOfferIDs}
		queries["receivedOfferID"] = offFilter
	}
	if len(req.Status) > 0 {
		statusFilter := map[string]interface{}{"$in": req.Status}
		queries["status"] = statusFilter
	}
	if len(req.Transaction) > 0 {
		queries["transactionID"] = req.Transaction
	}

	return queries

}
