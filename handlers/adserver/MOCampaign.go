package v1

import (
	"log"
	"time"

	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

//Requests and response types for this handler
type AddMOCampaignReq struct {
	//Profiling details
	Name       string    `json:"name,omitempty"`
	CategoryId int       `json:"categoryId,omitempty"`
	BidType    int       `json:"bidType,omitempty"`
	Bid        float32   `json:"bid,omitempty"`
	Budget     float32   `json:"budget,omitempty"`
	StartDate  time.Time `json:"startDate,omitempty"`
	EndDate    time.Time `json:"endDate,omitempty"`
	Days       int       `json:"days,omitempty"`
	//targetting details
	NetworkTypes   []int `json:"networkTypes,omitempty"`
	CountryIds     []int `json:"countryIds,omitempty"`
	OsTypes        []int `json:"osTypes,omitempty"`
	InventoryTypes []int `json:"inventoryTypes,omitempty"`
	//Creatives
	CreativeIds []int `json:"creativeIds,omitempty"`
	Status      int
}

type AddMOCampaignRsp struct {
	Status string
	Id     string
}

type AddMOCampaignCreativeReq struct {
	CategoryId int
	Height     int    `json:"height,omitempty"`
	Width      int    `json:"width,omitempty"`
	CampaignID string `json:"campaignId,omitempty"`
	Source     string `json:"source,omitempty"`
	Status     int
	PixelUrl   string `json:"pixelUrl,omitempty"`
	Title      string `json:"title,omitempty"`
}

type AddMOCampaignCreativeRsp struct {
	Status string
	Id     string
}

type UpdateMOCampaignReq struct {
	Id     int32
	Status int32
	Name   string
}

type UpdateMOCampaignRsp struct {
	Status string
}

//Handlers and methods for API
type MOCampaignHandler interface {
	Add(context.Context, *AddMOCampaignReq, *AddMOCampaignRsp) error
	Update(context.Context, *UpdateMOCampaignReq, *UpdateMOCampaignRsp) error
	AddCreative(context.Context, *AddMOCampaignCreativeReq, *AddMOCampaignCreativeRsp) error
}

func RegisterMOCampaignHandler(s server.Server, hdlr MOCampaignHandler) {
	s.Handle(s.NewHandler(&MOCampaign{hdlr}))
}

type MOCampaign struct {
	MOCampaignHandler
}

func (mo *MOCampaign) Add(ctx context.Context, req *AddMOCampaignReq, rsp *AddMOCampaignRsp) error {
	log.Print("Received MOCampaign.Add request from registry")
	// Index a campaign creative
	put1, err := ESClient.Index().
		Index("MobileAds").
		Type("p1").
		BodyJson(req).
		Do(context.Background())
	if err != nil {
		// Handle error
		log.Print(err.Error())
	}

	//Create the meta data for ad serving & sabe in mocampaign ad index
	rsp.Status = "Success"
	rsp.Id = put1.Id
	return nil
}

func (mo *MOCampaign) Update(ctx context.Context, req *UpdateMOCampaignReq, rsp *UpdateMOCampaignRsp) error {
	log.Print("Received MOCampaign.Update request from registry")
	rsp.Status = "Success"
	return nil
}

func (mo *MOCampaign) AddCreative(ctx context.Context, req *AddMOCampaignCreativeReq, rsp *AddMOCampaignCreativeRsp) error {
	log.Print("Received MOCampaign.AddCreative request from registry")
	// Index a campaign creative
	put1, err := ESClient.Index().
		Index("mocampaign").
		Type("creative").
		BodyJson(req).
		Do(context.Background())
	if err != nil {
		// Handle error
		log.Print(err.Error())
	}
	log.Print(put1.Id)
	rsp.Status = "Success"
	rsp.Id = put1.Id
	return nil
}
