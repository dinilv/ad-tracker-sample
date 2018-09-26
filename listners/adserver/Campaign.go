package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	db "github.com/adcamie/adserver/db/config"
	model "github.com/adcamie/adserver/db/model/v1"
	handler "github.com/adcamie/adserver/handlers/v1/adserver"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
	api "github.com/micro/micro/api/proto"
	"github.com/olivere/elastic"
)

var ESClient = db.ESMasterClient

type Mocampaign struct {
	MocampaignListner
}

type MocampaignListner interface {
	//Campaign creation part on our Ad-server-UI
	Add(context.Context, *api.Request, *api.Response) error
	Update(context.Context, *api.Request, *api.Response) error
	Addcreative(context.Context, *api.Request, *api.Response) error
	// Ad-serving and reporting
	Moad(context.Context, *api.Request, *api.Response) error
	List(context.Context, *api.Request, *api.Response) error
	Analytics(context.Context, *api.Request, *api.Response) error
	Log(context.Context, *api.Request, *api.Response) error
}

func RegisterMOCampaignListner(s server.Server, hdlr MocampaignListner) {
	s.Handle(s.NewHandler(&Mocampaign{hdlr}))
}

func (mo *Mocampaign) Add(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received MO Campaign Add Req")
	name, ok := req.Get["name"]
	if !ok || len(name.Values) == 0 {
		return errors.BadRequest("go.micro.service.v1.mocampaign", "Name cannot be blank")
	}

	request := client.NewJsonRequest("go.micro.service.v1.mocampaign", "MOCampaign.Add", &handler.AddMOCampaignReq{
		Name: strings.Join(name.Values, " "),
	})

	response := &handler.AddMOCampaignRsp{}

	if err := client.Call(ctx, request, response); err != nil {
		return err
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Add MO Campaign API Listnener": response.Status,
	})
	rsp.Body = string(b)

	return nil
}

func (mo *Mocampaign) Update(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received MO Campaign Update Req")

	name, ok := req.Get["id"]
	if !ok || len(name.Values) == 0 {
		return errors.BadRequest("go.micro.service.v1.mocampaign", "ID cannot be blank")
	}

	request := client.NewJsonRequest("go.micro.service.v1.mocampaign", "MOCampaign.Update", &handler.UpdateMOCampaignReq{
		Name: strings.Join(name.Values, " "),
	})
	response := &handler.UpdateMOCampaignRsp{}

	if err := client.Call(ctx, request, response); err != nil {
		return err
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Update MO API Listnener": response.Status,
	})
	rsp.Body = string(b)

	return nil
}

func (mo *Mocampaign) Addcreative(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received MO Campaign Creative Add Req")

	var cr = new(handler.AddMOCampaignCreativeReq)
	err := json.Unmarshal([]byte(req.Body), &cr)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	request := client.NewJsonRequest("go.micro.service.v1.mocampaign", "MOCampaign.AddCreative", cr)
	response := &handler.UpdateMOCampaignRsp{}

	if err := client.Call(ctx, request, response); err != nil {
		return err
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"Update MO API Listnener": response.Status,
	})
	rsp.Body = string(b)

	return nil
}

func (mo *Mocampaign) Moad(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received MO Campaign Add MOAD Req")

	var moad = new(model.MOAd)
	err := json.Unmarshal([]byte(req.Body), &moad)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	cookieIds := []string{"0"}
	moad.CookieIds = cookieIds
	moad.Impressions = 0
	moad.Clicks = 0
	put1, err := ESClient.Index().
		Index("mobileads").
		Type("p" + strconv.Itoa(moad.Priority)).BodyJson(moad).
		Do(context.Background())
	if err != nil {
		// Handle error
		log.Print(err.Error())
	}
	fmt.Println(put1.Id + "Id")
	_, err = ESClient.Update().Index("mobileads").Type("p" + strconv.Itoa(moad.Priority)).Id(put1.Id).Doc(map[string]interface{}{"impressions": 0, "clicks": 0}).DocAsUpsert(true).Do(context.Background())
	rsp.StatusCode = 200
	rsp.Body = "{}"

	return nil
}

func (mo *Mocampaign) Analytics(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received MO Campaign Analytics Req")
	rsp.StatusCode = 200

	return nil
}

func (mo *Mocampaign) List(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received MO Campaign List Req")

	all := elastic.NewMatchAllQuery()
	res, err := ESClient.Search().Index().Index("mobileads").Type("p" + req.Get["priority"].Values[0]).Query(all).Pretty(false).Do(context.Background())
	if err != nil {
		// Handle error
		log.Print(err)
	}

	b, _ := json.Marshal(res.Hits.Hits)
	rsp.Body = string(b)
	rsp.StatusCode = 200

	return nil
}

func (mo *Mocampaign) Log(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received MO Campaign Logs Req")

	all := elastic.NewMatchAllQuery()
	res, err := ESClient.Search().Index().Index("tracking").Type(req.Get["type"].Values[0]).Query(all).Pretty(false).Do(context.Background())
	if err != nil {
		// Handle error
		log.Print(err)
	}

	b, _ := json.Marshal(res.Hits.Hits)
	rsp.Body = string(b)
	rsp.StatusCode = 200

	return nil
}
