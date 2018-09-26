package v1

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	handler "github.com/adcamie/adserver/handlers/tracker"
	helper "github.com/adcamie/adserver/helpers"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
)

var layout, inputDateFormat string

func init() {
	layout = "02/01/2006"
	inputDateFormat = "02/01/2006"
}

type Report struct {
	ReportListener
}

type ReportListener interface {
	Basic(context.Context, *api.Request, *api.Response) error
	Rotation(context.Context, *api.Request, *api.Response) error
	Adcamie(context.Context, *api.Request, *api.Response) error
	Optimised(context.Context, *api.Request, *api.Response) error
	Dashboard(context.Context, *api.Request, *api.Response) error
}

func (report *Report) Basic(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get report :")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)
	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.tracker.report", "Token Authentication Failed Error:: Reports")
	}
	var getreport = ValidateForReport(reqG)

	request := client.NewJsonRequest("go.micro.service.v1.report", "Report.Basic", getreport)
	response := &handler.GetReportRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Basic Report :", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]interface{}{
		"data": response.Report,
	})
	rsp.Body = string(b)
	return nil
}

func (report *Report) Adcamie(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get adcamie report :")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)
	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.tracker.report", "Token Authentication Failed Error:: Reports")
	}
	var getreport = ValidateForReport(reqG)

	request := client.NewJsonRequest("go.micro.service.v1.report", "Report.Adcamie", getreport)
	response := &handler.GetReportRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Adcamie Report :", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]interface{}{
		"data": response.Report,
	})
	rsp.Body = string(b)
	return nil
}

func (report *Report) Optimised(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Listener got the request to get adcamie optimised report :")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)
	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.tracker.report", "Token Authentication Failed Error:: Reports")
	}
	var getreport = ValidateForReport(reqG)

	request := client.NewJsonRequest("go.micro.service.v1.report", "Report.AdcamieOptimised", getreport)
	response := &handler.GetReportRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Adcamie Report :", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]interface{}{
		"data": response.Report,
	})
	rsp.Body = string(b)
	return nil
}

func (report *Report) Rotation(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get rotation report :")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)
	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.report", "Token Authentication Failed Error:: Reports")
	}

	var getreport = ValidateForReport(reqG)

	request := client.NewJsonRequest("go.micro.service.v1.report", "Report.Rotation", getreport)
	response := &handler.GetReportRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Rotaion Get Report:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]interface{}{
		"data": response.Report,
	})
	rsp.Body = string(b)
	return nil
}

func validateAuthentication(reqG map[string]string) bool {

	//validations
	if strings.Compare(reqG[constants.Token], constants.PassToken) == 0 {
		return false
	}
	return true
}

func ValidateForReport(reqG map[string]string) *handler.GetReportReq {

	getreport := new(handler.GetReportReq)

	start, ok := reqG["start_date"]
	if ok {
		getreport.StartDate = start
	} else {
		t := time.Now().UTC()
		getreport.StartDate = t.Format(layout)
	}
	end, ok := reqG["end_date"]
	if ok {
		getreport.EndDate = end
	} else {
		t := time.Now().UTC()
		getreport.EndDate = t.Format(layout)
	}

	if strings.Compare(getreport.Type, constants.PostBackPingLog) == 0 {
		end, _ := time.Parse(inputDateFormat, getreport.EndDate)
		rounded := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location()).UTC()
		dayBefore := rounded.AddDate(0, 0, 1)
		getreport.EndDate = dayBefore.Format(layout)
	}

	startH, ok := reqG["start_hour"]
	if ok {
		getreport.StartHour = startH
	} else {
		getreport.StartHour = "0"
	}
	endH, ok := reqG["end_hour"]
	if ok {
		getreport.EndHour = endH
	} else {
		getreport.EndHour = "23"
	}
	affids, ok := reqG["affiliate_ids"]
	if ok {
		affids_splitted := strings.Split(affids, ",")
		if len(affids_splitted) > 0 && len(affids_splitted[0]) > 0 {
			getreport.AffiliateIDs = affids_splitted
		}
	}
	offids, ok := reqG["offer_ids"]
	if ok {
		offids_splitted := strings.Split(offids, ",")
		if len(offids_splitted) > 0 && len(offids_splitted[0]) > 0 {
			getreport.OfferIDs = offids_splitted
		}
	}
	status, ok := reqG["status"]
	if ok {
		status_splitted := strings.Split(status, ",")
		if len(status_splitted) > 0 && len(status_splitted[0]) > 0 {
			getreport.Status = status_splitted
		}
	}
	recv_offids, ok := reqG["recv_offer_ids"]
	if ok {
		recv_offids_splitted := strings.Split(recv_offids, ",")
		if len(recv_offids_splitted) > 0 && len(recv_offids_splitted[0]) > 0 {
			getreport.RecvOfferIDs = recv_offids_splitted
		}
	}
	recv_affids, ok := reqG["recv_affiliate_ids"]
	if ok {
		recv_affids_splitted := strings.Split(recv_affids, ",")
		if len(recv_affids_splitted) > 0 && len(recv_affids_splitted[0]) > 0 {
			getreport.RecvAffiliateIDs = recv_affids_splitted
		}
	}
	fields, ok := reqG["fields"]
	if ok {
		fields_splitted := strings.Split(fields, ",")
		if len(fields_splitted) > 0 && len(fields_splitted[0]) > 0 {
			getreport.Fields = fields_splitted
		}

	}
	reporttype, ok := reqG["type"]
	if ok {
		getreport.Type = reporttype
	}

	transaction, ok := reqG["transaction"]
	if ok {
		if len(transaction) > 0 {
			getreport.Transaction = transaction
		}
	}

	offset, ok := reqG["page"]
	if ok {
		pageno, _ := strconv.Atoi(offset)
		getreport.Page = pageno
	} else {
		getreport.Page = 1
	}

	key, ok := reqG["key"]
	if ok {
		if len(key) > 0 {
			getreport.Key = key
		}
	}
	return getreport
}
