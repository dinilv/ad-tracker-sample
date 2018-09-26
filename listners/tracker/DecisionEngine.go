package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	model "github.com/adcamie/adserver/db/model"
	handler "github.com/adcamie/adserver/handlers/tracker"
	helper "github.com/adcamie/adserver/helpers"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
)

type Engine struct {
	DecisionEngineListener
}

type DecisionEngineListener interface {
	Events(context.Context, *api.Request, *api.Response) error
	List(context.Context, *api.Request, *api.Response) error
	Offerstack(context.Context, *api.Request, *api.Response) error
	Msisdn(context.Context, *api.Request, *api.Response) error
	Update(context.Context, *api.Request, *api.Response) error
	Uploadblocked(context.Context, *api.Request, *api.Response) error
	Uploadsubscriber(context.Context, *api.Request, *api.Response) error
}

func (engine *Engine) Events(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Listener got request to add events to offers:")
	//parse the json body to request
	log.Print("Req for decision is :", req)
	var events = new(model.AdcamieEvents)
	err := json.Unmarshal([]byte(req.Body), &events)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.decision.engine", "Error in parsing request body")
	}
	//add ip where events coming from
	ip := req.GetHeader()["X-Forwarded-For"].Values[0]
	events.IP = ip
	request := client.NewJsonRequest("go.micro.service.v1.engine", "Engine.Events", events)
	response := &handler.EngineRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In saving Events:", err)
		return err
	}
	rsp.StatusCode = 200
	return nil
}

func (engine *Engine) List(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to list offer events received:")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)
	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.decision.engine", "Token Authentication Failed Error:: Engine")
	}
	var getofferevent = new(handler.GetOfferEventListReq)
	getofferevent.Token = reqG["token"]
	offids, ok := reqG["offer_ids"]
	if ok {
		getofferevent.OfferIDs = strings.Split(offids, ",")
	}
	startdate, ok := reqG["start_date"]
	if ok {
		getofferevent.StartDate = startdate
	}
	enddate, ok := reqG["end_date"]
	if ok {
		getofferevent.EndDate = enddate
	}
	request := client.NewJsonRequest("go.micro.service.v1.engine", "Engine.Listevents", getofferevent)
	response := &handler.GetOfferEventListRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Event List :", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (engine *Engine) Msisdn(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to search MSISDN in stack:")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)
	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.decision.engine", "Token Authentication Failed Error:: Engine")
	}
	var reqM = new(handler.MsisdnReq)
	reqM.MSISDN = reqG["msisdn"]
	request := client.NewJsonRequest("go.micro.service.v1.engine", "Engine.MSISDN", reqM)
	response := &handler.MsisdnRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Msisdn validate :", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]interface{}{
		"data": response,
	})
	rsp.Body = string(b)
	return nil
}

func (engine *Engine) Offerstack(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to list offers in the rotation stack:")
	//parse get values
	reqG := make(map[string]string)
	helper.ParseGetRequestOnAPI(req, reqG)
	if validateAuthentication(reqG) {
		return errors.BadRequest("go.micro.service.v1.decision.engine", "Token Authentication Failed Error:: Engine")
	}
	var offerstackreq = new(handler.GetOfferStackReq)
	offerstackreq.Limit, _ = strconv.Atoi(reqG["limit"])
	offerstackreq.Page, _ = strconv.Atoi(reqG["page"])
	offerstackreq.Search, _ = reqG["search"]
	request := client.NewJsonRequest("go.micro.service.v1.engine", "Engine.Offerstack", offerstackreq)
	response := &handler.GetOfferStackRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Event Get Offer Stack:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]interface{}{
		"data":  response.Data,
		"count": response.Count,
	})
	rsp.Body = string(b)
	return nil
}

func (engine *Engine) Update(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got request to update  offerstack:")
	//parse the json body to request
	var updateOfferstack = new(handler.UpdateOfferStackReq)
	err := json.Unmarshal([]byte(req.Body), &updateOfferstack)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.decision.engine", "Error in parsing request body")
	}
	request := client.NewJsonRequest("go.micro.service.v1.engine", "Engine.Update", updateOfferstack)
	response := &handler.EngineRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In saving Events:", err)
		return err
	}
	rsp.StatusCode = 200
	return nil
}

func (engine *Engine) Uploadblocked(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Listener got the csv to save blocked msisdns for MO-Campaigns:")

	return nil
}

func (engine *Engine) Uploadsubscriber(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Listener got the request to list offers in the roataion stack:")
	return nil
}
