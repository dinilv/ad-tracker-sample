package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	handler "github.com/adcamie/adserver/handlers/tracker"
	helper "github.com/adcamie/adserver/helpers"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
)

type Tracker struct {
	TrackerListener
}

type TrackerListener interface {
	Setoffer(context.Context, *api.Request, *api.Response) error
	Setaffiliate(context.Context, *api.Request, *api.Response) error
	Setmqf(context.Context, *api.Request, *api.Response) error
	Getoffers(context.Context, *api.Request, *api.Response) error
	Getoffer(context.Context, *api.Request, *api.Response) error
	Deleteoffer(context.Context, *api.Request, *api.Response) error
	Getaffiliates(context.Context, *api.Request, *api.Response) error
	Getaffiliate(context.Context, *api.Request, *api.Response) error
	Deleteaffiliate(context.Context, *api.Request, *api.Response) error
	Getrotatedmqfs(context.Context, *api.Request, *api.Response) error
	Getmqfs(context.Context, *api.Request, *api.Response) error
	Getmqf(context.Context, *api.Request, *api.Response) error
	Deletemqf(context.Context, *api.Request, *api.Response) error
	Gettransactions(context.Context, *api.Request, *api.Response) error
	Dopostback(context.Context, *api.Request, *api.Response) error
	Offerdropdown(context.Context, *api.Request, *api.Response) error
	Affiliatedropdown(context.Context, *api.Request, *api.Response) error
	Addaffiliate(context.Context, *api.Request, *api.Response) error
	Getadvertisers(context.Context, *api.Request, *api.Response) error
	Getadvertiser(context.Context, *api.Request, *api.Response) error
	Deleteadvertiser(context.Context, *api.Request, *api.Response) error
}

func (tracker *Tracker) Setoffer(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Listener got the request to set offer:")
	log.Print("req is :", req)

	//parse offer parameters to req body
	var offer = new(handler.SetOfferReq)
	log.Print("Offer is :", offer)
	err := json.Unmarshal([]byte(req.Body), &offer)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Error in parsing request body")
	}

	if offer.OfferID == "" {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Offer Id Not Present In Parameters:: Offer Setting")
	}

	if offer.OfferName == "" {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Offer Name Not Present In Parameters:: Offer Setting")
	}

	if offer.Type == "" {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Offer Type Not Present In Parameters:: Offer Setting")
	}

	if offer.Template == "" {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Template Not Present In Parameters:: Offer Setting")
	} else {
		template := strings.Replace(offer.Template, "{", "", -1)
		template = strings.Replace(template, "}", "", -1)
		_, err := url.Parse(template)
		if err != nil {
			fmt.Println(err)
			return errors.BadRequest("go.micro.service.v1.trackersetting", "Template is not in proper URL format:: Offer Setting")
		}

	}

	if offer.Token != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Authentication Failed Error:: Offer Setting")
	}

	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Setoffer", offer)
	response := &handler.SettingRes{}
	client.Call(ctx, request, response)

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Result": response.Status, "Id": response.Id,
	})
	rsp.Body = string(b)

	return nil
}

func (tracker *Tracker) Setaffiliate(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to set affiliate:")

	//parse affiliate parameters to req body
	var affiliate = new(handler.SetAffiliateReq)
	err := json.Unmarshal([]byte(req.Body), &affiliate)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Error in parsing request body")
	}
	if affiliate.Token != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.tracker", "Token Authentication Failed Error : : Affiliate Setting")
	}
	if affiliate.AffiliateID == "" {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Affiliate Id Not Present In Parameters:: Affiliate Setting")
	}
	if affiliate.AffiliateName == "" {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Affiliate Name Not Present In Parameters:: Affiliate Setting")
	}
	if affiliate.Mqf == 0.0 {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "MQF  Required :: Affiliate Setting")
	}
	if affiliate.MediaTemplate != "" {
		//validate for proper url
		template := strings.Replace(affiliate.MediaTemplate, "{", "", -1)
		template = strings.Replace(template, "}", "", -1)
		_, err := url.Parse(template)
		if err != nil {
			fmt.Println(err)
			return errors.BadRequest("go.micro.service.v1.trackersetting", "Template is not in proper URL format:: Affiliate Setting")
		}

	}

	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Setaffiliate", affiliate)
	response := &handler.SettingRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Affiliate Setting:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Result": response.Status, "Id": response.Id,
	})
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Setmqf(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to set mqf offer affiliate")

	//parse set offer affiliate mqf parameters to req body
	var mqf = new(handler.SetMqfReq)
	err := json.Unmarshal([]byte(req.Body), &mqf)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Error in parsing request body")
	}
	if mqf.Token != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Authentication Failed Error:: Offer Affiliate MQF Setting")
	}
	if mqf.Mqf == 0.0 {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "MQF  Required : : Offer Affiliate MQF Setting")
	}

	if mqf.AffiliateID == "" {
		return errors.BadRequest("go.micro.service.v1.tracker", "Affiliate Id Required : : Offer Affiliate MQF Setting")
	}
	if mqf.OfferID == "" {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Offer Id  Required : : Offer Affiliate MQF Setting")
	}
	if mqf.MediaTemplate != "" {
		//validate for proper url
		template := strings.Replace(mqf.MediaTemplate, "{", "", -1)
		template = strings.Replace(template, "}", "", -1)
		_, err := url.Parse(template)
		if err != nil {
			fmt.Println(err)
			return errors.BadRequest("go.micro.service.v1.tracker", "Template is not in proper URL format:: Offer Affiliate MQF Setting")
		}

	}

	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Setmqf", mqf)
	response := &handler.SettingRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Affiliate MQF Setting:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Result": response.Status, "Id": response.Id,
	})
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Getoffers(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get offers :")
	token_, ok := req.Get["token"]
	if ok == false || strings.Compare(token_.Values[0], constants.PassToken) != 0 {
		log.Print("Token Not Present In Parameters:" + token_.Values[0])
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Not Present In Parameters:: Get offer")
	}
	var getoffers = new(handler.GetOffersReq)
	//getoffers.Page = pageno
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Getoffers", getoffers)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Offers:", err)
		return err
	}
	rsp.StatusCode = 200
	log.Print(response.Data)
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Getoffer(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get one offer :")
	token_, ok := req.Get["token"]
	if ok == false || strings.Compare(token_.Values[0], constants.PassToken) != 0 {
		log.Print("Token Not Present In Parameters:" + token_.Values[0])
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Not Present In Parameters:: Get one offer")
	}
	id, ok := req.Get["id"]
	if ok == false {
		log.Print("Id Not Present In Parameters:" + id.Values[0])
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Id Not Present In Parameters:: Get one offer")
	}
	var getoffer = new(handler.GetOfferReq)
	getoffer.Token = token_.Values[0]
	getoffer.Id = id.Values[0]
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Getoffer", getoffer)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Offer:", err)
		return err
	}
	rsp.StatusCode = 200
	log.Print(response.Data)
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Deleteoffer(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to delete one offer :")
	var deleteoffer = new(handler.DeleteOfferReq)
	log.Print("Offer is :", deleteoffer)
	err := json.Unmarshal([]byte(req.Body), &deleteoffer)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Error in parsing request body")
	}
	if deleteoffer.Token != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Authentication Failed Error:: Offer Deletion")
	}
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Deleteoffer", deleteoffer)
	response := &handler.DeleteOfferRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Offer:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Result": response.Status, "Id": response.Id,
	})
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Getaffiliates(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get affiliates :")
	token_, ok := req.Get["token"]
	if ok == false || strings.Compare(token_.Values[0], constants.PassToken) != 0 {
		log.Print("Token Not Present In Parameters:")
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Not Present In Parameters:: Get Affiliates")
	}
	var getoffers = new(handler.GetOffersReq)
	getoffers.Token = token_.Values[0]
	//getoffers.Page = pageno
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Getaffiliates", getoffers)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Affiliate:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Getaffiliate(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get affiliates :")
	token_, ok := req.Get["token"]
	if ok == false || strings.Compare(token_.Values[0], constants.PassToken) != 0 {
		log.Print("Token Not Present In Parameters:")
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Not Present In Parameters:: Get Affiliates")
	}
	id, ok := req.Get["id"]
	if ok == false {
		log.Print("Id Not Present In Parameters:" + id.Values[0])
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Id Not Present In Parameters:: Get one offer")
	}
	var getoffer = new(handler.GetOfferReq)
	getoffer.Token = token_.Values[0]
	getoffer.Id = id.Values[0]
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Getaffiliate", getoffer)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Affiliate:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Deleteaffiliate(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to delete one affiliate :")
	var deleteaffiliate = new(handler.DeleteAffiliateReq)
	log.Print("Affiliate is :", deleteaffiliate)
	err := json.Unmarshal([]byte(req.Body), &deleteaffiliate)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Error in parsing request body")
	}
	if deleteaffiliate.Token != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Authentication Failed Error:: Affiliate Deletion")
	}
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Deleteaffiliate", deleteaffiliate)
	response := &handler.DeleteAffiliateRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Offer:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Result": response.Status, "Id": response.Id,
	})
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Getmqfs(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get mqfs :")
	token_, ok := req.Get["token"]
	if ok == false || strings.Compare(token_.Values[0], constants.PassToken) != 0 {
		log.Print("Token Not Present In Parameters:")
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Not Present In Parameters:: Get Mqfs")
	}
	var getoffers = new(handler.GetOffersReq)
	getoffers.Token = token_.Values[0]
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Getmqfs", getoffers)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Mqfs:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Getrotatedmqfs(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get rotated mqfs :")
	token_, ok := req.Get["token"]
	if ok == false || strings.Compare(token_.Values[0], constants.PassToken) != 0 {
		log.Print("Token Not Present In Parameters:")
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Not Present In Parameters:: Get Mqfs")
	}
	var getoffers = new(handler.GetOffersReq)
	getoffers.Token = token_.Values[0]
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Getrotatedmqfs", getoffers)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Mqfs:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Getmqf(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get mqf :")
	token_, ok := req.Get["token"]
	if ok == false || strings.Compare(token_.Values[0], constants.PassToken) != 0 {
		log.Print("Token Not Present In Parameters:")
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Not Present In Parameters:: Get Mqf")
	}
	offerid, ok := req.Get["offer_id"]
	if ok == false {
		log.Print("Offer Id Not Present In Parameters:" + offerid.Values[0])
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Offer Id Not Present In Parameters:: Get mqf")
	}
	affilaiteid, ok := req.Get["affiliate_id"]
	if ok == false {
		log.Print("Affiliate Id Not Present In Parameters:" + affilaiteid.Values[0])
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Affiliate Id Not Present In Parameters:: Get mqf")
	}
	var getoffer = new(handler.GetMQFReq)
	getoffer.Token = token_.Values[0]
	getoffer.Offerid = offerid.Values[0]
	getoffer.Affiliateid = affilaiteid.Values[0]
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Getmqf", getoffer)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Mqfs:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Deletemqf(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to delete one mqf :")
	var deletemqf = new(handler.DeleteMQFReq)
	log.Print("MQF is :", deletemqf)
	err := json.Unmarshal([]byte(req.Body), &deletemqf)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Error in parsing request body")
	}
	if deletemqf.Token != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Authentication Failed Error:: MQF Deletion")
	}
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Deletemqf", deletemqf)
	response := &handler.DeleteMQFRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Offer:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Result": response.Status, "Id": response.Id,
	})
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Gettransactions(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get transaction :")
	//parse get values

	reqG := make(map[string]string)
	for _, get := range req.Get {
		for _, val := range get.Values {
			reqG[get.Key] = val
		}
	}
	if reqG["token"] != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Authentication Failed Error:: Transaction Id")
	}
	var gettransactions = new(handler.GetTransactionReq)
	gettransactions.Token = reqG["token"]
	affids, ok := reqG["affiliate_ids"]
	if ok {
		gettransactions.AffiliateIDs = strings.Split(affids, ",")
	}
	offids, ok := reqG["offer_ids"]
	if ok {
		gettransactions.OfferIDs = strings.Split(offids, ",")
	}
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Gettransactions", gettransactions)
	response := &handler.GetTransactionRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Transactions :", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Dopostback(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to do manual postback :")
	//parse postback parameters to req body
	var postback = new(handler.ManualPostbackReq)
	err := json.Unmarshal([]byte(req.Body), &postback)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Error in parsing request body:: Manual Postback")
	}

	if postback.TransactionID == "" {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "TransactionID Id Not Present In Parameters:: Manual Postback")
	}

	if postback.Token != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Authentication Failed Error:: Manual Postback")
	}

	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Dopostback", postback)
	response := &handler.ManualPostbackRes{}
	client.Call(ctx, request, response)
	if len(response.URL) > 0 {
		helper.Ping(response.URL, map[string]string{constants.TRANSACTION_ID: response.TransactionID})
		rsp.StatusCode = 200
		b, _ := json.Marshal(map[string]string{
			"url": response.URL,
		})
		rsp.Body = string(b)

	} else {
		rsp.StatusCode = 500
		b, _ := json.Marshal(map[string]string{
			"message": response.Status,
		})
		rsp.Body = string(b)

	}

	return nil
}

func (tracker *Tracker) Addadvertiser(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to add advertiser")
	//parse data to  req body
	var advertiser = new(handler.SetAdvertiserReq)
	err := json.Unmarshal([]byte(req.Body), &advertiser)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.tracker", "Error in parsing request body while advertiser creation")
	}
	if advertiser.Token != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.tracker", "Token Authentication Failed Error:: Advertiser Creation")
	}
	if advertiser.AdvertiserName == "" {
		return errors.BadRequest("go.micro.service.v1.tracker", "Advertiser Name Required : : Advertiser Creation")
	}

	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Addadvertiser", advertiser)
	response := &handler.SettingRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Advertiser Creation:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Result": response.Status, "Id": response.Id,
	})
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Getadvertisers(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get advertisers :")
	token_, ok := req.Get["token"]
	if ok == false || strings.Compare(token_.Values[0], constants.PassToken) != 0 {
		log.Print("Token Not Present In Parameters:" + token_.Values[0])
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Not Present In Parameters:: Get advertisers")
	}
	var getadvertisers = new(handler.GetOffersReq)
	//getoffers.Page = pageno
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Getadvertisers", getadvertisers)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Advertisers:", err)
		return err
	}
	rsp.StatusCode = 200
	log.Print(response.Data)
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Getadvertiser(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to get one of the advertiser :")
	token_, ok := req.Get["token"]
	if ok == false || strings.Compare(token_.Values[0], constants.PassToken) != 0 {
		log.Print("Token Not Present In Parameters:" + token_.Values[0])
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Not Present In Parameters:: Get one advertiser")
	}
	name, ok := req.Get["id"]
	if ok == false {
		log.Print("Name Not Present In Parameters:" + name.Values[0])
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Name Not Present In Parameters:: Get one advertiser")
	}
	var getadvertiser = new(handler.GetOfferReq)
	getadvertiser.Token = token_.Values[0]
	getadvertiser.Id = name.Values[0]
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Getadvertiser", getadvertiser)
	response := &handler.GetOffersRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Advertiser:", err)
		return err
	}
	rsp.StatusCode = 200
	log.Print(response.Data)
	b, _ := json.Marshal(response)
	rsp.Body = string(b)
	return nil
}

func (tracker *Tracker) Deleteadvertiser(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Listener got the request to delete one advertiser :")
	var deleteadvertiser = new(handler.DeleteAdvertiserReq)
	log.Print("Advertiser is :", deleteadvertiser)
	err := json.Unmarshal([]byte(req.Body), &deleteadvertiser)
	if err != nil {
		fmt.Println("Whoops:", err)
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Error in parsing request body")
	}
	if deleteadvertiser.Token != constants.PassToken {
		return errors.BadRequest("go.micro.service.v1.trackersetting", "Token Authentication Failed Error:: Advertiser Deletion")
	}
	request := client.NewJsonRequest("go.micro.service.v1.tracker", "Tracker.Deleteadvertiser", deleteadvertiser)
	response := &handler.DeleteAdvertiserRes{}
	if err := client.Call(ctx, request, response); err != nil {
		log.Print("Client Calling Error In Tracker Get Offer:", err)
		return err
	}
	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Result": response.Status, "Id": response.Id,
	})
	rsp.Body = string(b)
	return nil
}
