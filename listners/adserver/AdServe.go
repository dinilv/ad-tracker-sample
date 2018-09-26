package v1

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"strconv"

	handler "github.com/adcamie/adserver/handlers/v1/adserver"
	message "github.com/adcamie/adserver/messages/proto/v1"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	api "github.com/micro/micro/api/proto"

	"golang.org/x/net/context"
)

type Adserve struct {
	AdserveListner
}

type AdserveListner interface {
	Search(context.Context, *api.Request, *api.Response) error
	Impression(context.Context, *api.Request, *api.Response) error
	Click(context.Context, *api.Request, *api.Response) error
}

func RegisterAdserveListner(s server.Server, hdlr AdserveListner) {
	s.Handle(s.NewHandler(&Adserve{hdlr}))
}

func (ad *Adserve) Search(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Received Adserve Search Req + Cookie-id:-", ctx)

	//Handling user cookies to find the user is repeated or not and set cookie for future uses
	cookieid := ""
	//Save the request headers/parameters,response and cookie to subscriber for tracking
	reqG := make(map[string]string)
	for _, get := range req.Get {
		for _, val := range get.Values {
			reqG[get.Key] = val
		}
	}
	for _, header := range req.GetHeader() {
		for _, val := range header.Values {
			if strings.Compare("Cookie", header.Key) == 0 {
				cookies := strings.Split(val, "=")
				cookieid = cookies[1]
				log.Print("Received Cookie-id:-" + cookieid)
			}
			reqG[header.Key] = val
		}
	}

	//Logging request first
	requestId := strconv.FormatInt(time.Now().UnixNano(), 10)
	reqG["RequestID"] = requestId
	pub := client.NewPublication("go.micro.sub.moadserve", &message.Message{Reqh: reqG, TrackingActivity: 1})
	client.Publish(ctx, pub)

	//Check for available ads search
	sr := new(handler.SearchAdserveReq)
	sr.CategoryId, _ = strconv.Atoi(req.Get["cid"].Values[0])
	sr.Height, _ = strconv.Atoi(req.Get["hgt"].Values[0])
	sr.Width, _ = strconv.Atoi(req.Get["wdt"].Values[0])
	sr.CookieId = cookieid
	request := client.NewJsonRequest("go.micro.service.v1.adserve", "AdServe.Search", sr)

	response := &handler.SearchAdserveRsp{}

	if err := client.Call(ctx, request, response); err != nil {
		return err
	}

	rsp.StatusCode = response.Status
	b, _ := json.Marshal(response.Data)
	rsp.Body = string(b)

	//set cookie if its empty
	if strings.Compare(cookieid, "") == 0 {
		log.Print("cookieid" + cookieid)
		cookieid = "th_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		values := []string{"adsauce_id=" + cookieid + "; expires=Fri, 31 Dec 9999 23:59:59 GMT;"}
		rsp.Header = map[string]*api.Pair{"Set-Cookie": {Key: "Set-Cookie", Values: values}}
	}

	reqG = make(map[string]string)
	reqG["id"] = response.ID
	reqG["type"] = response.Type
	reqG["tid"] = response.Data.Tid
	reqG["RequestID"] = requestId
	reqG["CookieID"] = cookieid
	reqG["CampaignID"] = strconv.Itoa(response.CampaignID)
	reqG["CreativeID"] = strconv.Itoa(response.CreativeID)
	reqG["searchTime"] = strconv.FormatInt(response.SearchTime, 10)
	pub = client.NewPublication("go.micro.sub.moadserve", &message.Message{Reqh: reqG, TrackingActivity: 2})
	client.Publish(ctx, pub)

	return nil
}

func (ad *Adserve) Impression(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received Ad Impression Track Req")

	//Save the request response and cookie to subscriber for tracking
	reqG := make(map[string]string)
	for _, get := range req.Get {
		for _, val := range get.Values {
			reqG[get.Key] = val
		}
	}
	for _, header := range req.GetHeader() {
		for _, val := range header.Values {
			reqG[header.Key] = val
		}
	}
	pub := client.NewPublication("go.micro.sub.moadserve", &message.Message{Reqh: reqG, TrackingActivity: 3})
	client.Publish(ctx, pub)

	rsp.StatusCode = 200

	return nil
}

func (ad *Adserve) Click(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received Click Track Req")

	//Save the request response and cookie to subscriber for tracking
	reqG := make(map[string]string)
	for _, get := range req.Get {
		for _, val := range get.Values {
			reqG[get.Key] = val
		}
	}
	for _, header := range req.GetHeader() {
		for _, val := range header.Values {
			reqG[header.Key] = val
		}
	}
	pub := client.NewPublication("go.micro.sub.moadserve", &message.Message{Reqh: reqG, TrackingActivity: 4})
	client.Publish(ctx, pub)

	rsp.StatusCode = 200

	return nil
}

func NewAuthHandler() server.HandlerWrapper {

	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			log.Print("Inside Auth Handler", req.Request())
			return h(ctx, req, rsp)
		}
	}
}
