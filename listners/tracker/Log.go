package v1

import (
	"context"
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	handler "github.com/adcamie/adserver/handlers/tracker"
	helper "github.com/adcamie/adserver/helpers"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
)

type Log struct {
	LogListener
}

type LogListener interface {
	All(context.Context, *api.Request, *api.Response) error
	Totalcount(context.Context, *api.Request, *api.Response) error
	Csv(context.Context, *api.Request, *api.Response) error
	Mo(context.Context, *api.Request, *api.Response) error
}

//To get log reports
func (rlog *Log) All(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Println("Listener got the request to get all logs :")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)

	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.tracker.log", "Token Authentication Failed Error:: Reports")
	}
	var getreport = ValidateForReport(reqG)
	callOpts := client.WithRequestTimeout(10 * time.Minute)
	request := client.NewJsonRequest("go.micro.service.v1.log", "Log.All", getreport)
	response := &handler.GetReportRes{}
	if err := client.Call(ctx, request, response, callOpts); err != nil {
		log.Print("Client Calling Error In Report Log Setting:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

//To get mongodb count log reports
func (rlog *Log) Totalcount(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Println("Listener got the request to get mongo count :")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)
	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.tracker.log", "Token Authentication Failed Error:: Log")
	}
	var getreport = ValidateForReport(reqG)
	callOpts := client.WithRequestTimeout(10 * time.Minute)
	request := client.NewJsonRequest("go.micro.service.v1.log", "Log.Totalcount", getreport)
	response := &handler.GetReportRes{}
	if err := client.Call(ctx, request, response, callOpts); err != nil {
		log.Print("Client Calling Error In Mongo Counting :", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (rlog *Log) Csv(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Println("Listener got the request to download report :")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)
	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.tracker.log", "Token Authentication Failed Error:: Reports")
	}
	var getreport = ValidateForReport(reqG)
	getreport.FileName = constants.Listener_Path + reqG["filename"]
	getreport.SortField = "utcdate"
	request := client.NewJsonRequest("go.micro.service.v1.log", "Log.CSV", getreport)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Report Setting:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string][]bson.M{
		"data": response.Data,
	})
	rsp.Body = string(b)
	return nil
}

//To get mo-log
func (rlog *Log) Mo(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Println("Listener got the request to get all mo logs :")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)

	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.tracker.log", "Token Authentication Failed Error:: Reports")
	}
	var getreport = ValidateForReport(reqG)
	request := client.NewJsonRequest("go.micro.service.v1.log", "Log.MO", getreport)
	response := &handler.GetReportRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Report Log Setting:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}
