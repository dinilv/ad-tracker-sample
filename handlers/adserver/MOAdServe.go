package v1

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	db "github.com/adcamie/adserver/db/config"
	model "github.com/adcamie/adserver/db/model/v1"
	"github.com/micro/go-micro/server"
	"github.com/olivere/elastic"
	"golang.org/x/net/context"
)

var ESClient = db.ESMasterClient

//Requests and response types for this handler
type SearchAdserveReq struct {
	//targetting details
	NetworkType   int
	CountryId     int
	OsType        int
	InventoryType int
	CategoryId    int
	Height        int
	Width         int
	CookieId      string
}

type SearchAdserveRsp struct {
	Type       string            `json:"type,omitempty"`
	Status     int32             `json:"status,omitempty"`
	Data       SearchAdResultRow `json:"data,omitempty"`
	CreativeID int               `json:"creativeId,omitempty"`
	CampaignID int               `json:"campaignId,omitempty"`
	ID         string            `json:"id,omitempty"`
	SearchTime int64             `json:"searchTime,omitempty"`
}

/*func (results *SearchAdserveRsp) AddResult(row SearchAdResultRow) []SearchAdResultRow {
	results.Data = append(results.Data, row)
	return results.Data
}
*/

type SearchAdResultRow struct {
	Tid          string `json:"tid,omitempty"`
	CID          int    `json:"cid,omitempty"`
	AID          int    `json:"aid,omitempty"`
	CreativeLink string `json:"clink,omitempty"`
	ClickLink    string `json:"clklink,omitempty"`
	PixelLink    string `json:"pxlink,omitempty"`
}

//Handlers and methods for API
type AdServeHandler interface {
	Search(context.Context, *SearchAdserveReq, *SearchAdserveRsp) error
}

func RegisterAdServeHandlerHandler(s server.Server, hdlr AdServeHandler) {
	s.Handle(s.NewHandler(&AdServe{hdlr}))
}

type AdServe struct {
	AdServeHandler
}

func (ad *AdServe) Search(ctx context.Context, req *SearchAdserveReq, rsp *SearchAdserveRsp) error {
	log.Print("Received AdSearch request from registry")

	//Rotate Adserving according to  cookie id filtering, according to Size, bid, ecpm and ctr

	//elastic
	//cid := elastic.NewTermQuery("CategoryId", req.CategoryId)
	hgtr := elastic.NewRangeQuery("height").Gte(req.Height - 20).Lte(req.Height + 20).Boost(2.0)
	wdtr := elastic.NewRangeQuery("width").Gte(req.Width - 20).Lte(req.Width + 20).Boost(1.9)
	filter := elastic.NewTermQuery("cookieIds", req.CookieId)
	cookie := elastic.NewBoostingQuery().Positive(hgtr).Boost(1.0).Negative(filter).NegativeBoost(.5)
	log.Print(cookie.Source())
	results, err := ESClient.Search().Index().Index("mobileads").Query(wdtr).Query(cookie).Pretty(false).Size(1).Do(context.Background())

	if err != nil {
		// Handle error
		log.Print(err.Error())
	}

	fmt.Printf("Query took %d milliseconds\n", results.TotalHits())

	var ttyp model.MOAd

	for _, item := range results.Each(reflect.TypeOf(ttyp)) {
		log.Print(item)
		if t, ok := item.(model.MOAd); ok {
			//Create tracking unique for campaign & creative
			tid := "th_" + strconv.FormatInt(time.Now().UnixNano(), 10)
			each := SearchAdResultRow{}
			each.Tid = tid
			each.CID = t.CampaignID
			each.AID = 1
			each.CreativeLink = t.CreativeLink
			each.ClickLink = t.ClickURL
			each.PixelLink = t.PixelURL
			rsp.Data = each
			rsp.CampaignID = t.CampaignID
			rsp.CreativeID = t.CreativeID
			rsp.Type = results.Hits.Hits[0].Type
			rsp.ID = results.Hits.Hits[0].Id
			rsp.SearchTime = results.TookInMillis

		}
	}

	rsp.Status = 200

	return nil
}
