package v1

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	helper "github.com/adcamie/adserver/helpers"
	_ "github.com/micro/go-plugins/broker/googlepubsub"
	api "github.com/micro/micro/api/proto"
)

type Click struct {
	ClickListener
}

type ClickListener interface {
	Setclick(context.Context, *api.Request, *api.Response) error
}

func (tracker *Click) Tracker(ctx context.Context, req *api.Request, rsp *api.Response) error {

	fmt.Println("Listener got the request to click redirection:")
	startTime := time.Now()
	requestParams := make(map[string]string)
	requestParams[common.API_TIME] = time.Now().UTC().String()
	geo := new(model.GeoDetails)
	rotatedClick := false

	//prcocess request, parameters & geo
	helper.ParseGetRequest(requestParams, req)

	//add escape path
	requestParams[common.URL] = common.AFF_C + requestParams[common.URL]
	helper.ParseGetHeader(requestParams, req)

	//prcocess geo
	helper.GetGeo(req, geo)
	requestParams[common.CLICK_URL] = requestParams[common.URL]
	requestParams[common.URL] = ""

	//mandatory checks
	offerID, offOK := requestParams[common.OFFER_ID]
	affiliateID, affOK := requestParams[common.AFF_ID]

	//missing affiliateID tag as rotated
	if affOK == false || len(affiliateID) == 0 || helper.ID_REGEX.MatchString(requestParams[common.AFF_ID]) {
		rotatedClick = true
		requestParams[common.RECV_OFFER_ID] = offerID
		requestParams[common.RECV_AFF_ID] = common.TRACKER_MEDIA
		requestParams[common.AFF_ID] = common.TRACKER_MEDIA
		requestParams[common.ACTIVITY] = "15"
		affiliateID = common.TRACKER_MEDIA
	}
	if offOK == false || len(offerID) == 0 || helper.ID_REGEX.MatchString(requestParams[common.OFFER_ID]) {
		rotatedClick = true
		//query for matching offer with geo
		rotationOfferID := helper.GetOfferByGeo(geo.CountryCode)
		requestParams[common.OFFER_ID] = rotationOfferID
		requestParams[common.RECV_OFFER_ID] = common.TRACKER_OFFER
		requestParams[common.RECV_AFF_ID] = requestParams[common.AFF_ID]
		requestParams[common.AFF_ID] = common.TRACKER_MEDIA
		requestParams[common.ACTIVITY] = "16"
		offerID = rotationOfferID
		affiliateID = common.TRACKER_MEDIA

	} else if dao.ValidateExhaustedOfferHash(offerID) {
		rotatedClick = true
		//query for matching offer with geo
		rotationOfferID := helper.GetOfferFromStack(offerID, geo.CountryCode)
		requestParams[common.OFFER_ID] = rotationOfferID
		requestParams[common.RECV_OFFER_ID] = offerID
		requestParams[common.RECV_AFF_ID] = requestParams[common.AFF_ID]
		requestParams[common.AFF_ID] = common.TRACKER_MEDIA
		requestParams[common.ACTIVITY] = "22"
		offerID = rotationOfferID
		affiliateID = common.TRACKER_MEDIA
		fmt.Println("Rotated due to Exhausted Offer")
	} else if dao.ValidateExhaustedOfferAffiliateHash(offerID, affiliateID) {
		rotatedClick = true
		//query for matching offer with geo & offer
		rotationOfferID := helper.GetOfferFromStack(offerID, geo.CountryCode)
		requestParams[common.OFFER_ID] = rotationOfferID
		requestParams[common.RECV_OFFER_ID] = offerID
		requestParams[common.RECV_AFF_ID] = requestParams[common.AFF_ID]
		requestParams[common.AFF_ID] = common.TRACKER_MEDIA
		requestParams[common.ACTIVITY] = "23"
		offerID = rotationOfferID
		affiliateID = common.TRACKER_MEDIA
		fmt.Println("Rotated due to Exhausted Offer & Media")
	} else if !dao.ValidateOfferCountry(offerID, geo.CountryCode) {
		rotatedClick = true
		//query for matching offer with geo & offer
		rotationOfferID := helper.GetOfferFromStack(offerID, geo.CountryCode)
		requestParams[common.OFFER_ID] = rotationOfferID
		requestParams[common.RECV_OFFER_ID] = offerID
		requestParams[common.RECV_AFF_ID] = requestParams[common.AFF_ID]
		requestParams[common.AFF_ID] = common.TRACKER_MEDIA
		requestParams[common.ACTIVITY] = "14"
		offerID = rotationOfferID
		affiliateID = common.TRACKER_MEDIA
		fmt.Println("Rotated due to Wrong Geo")
	}

	//if transaction-id already present in click
	receivedTransactionID, ok := requestParams[common.TRANSACTION_ID]
	if ok {
		requestParams[common.RECV_TRANXN_ID] = receivedTransactionID
	}

	//generate and use transactionID
	transactionID := common.GenerateTransactionId(offerID, affiliateID)
	requestParams[common.TRANSACTION_ID] = transactionID

	template := dao.GetTemplateByOfferID(offerID)
	//template check
	if len(template) == 0 {
		rotatedClick = true
		//query for matching offer with geo & offer
		rotationOfferID := helper.GetOfferFromStack(offerID, geo.CountryCode)
		requestParams[common.OFFER_ID] = rotationOfferID
		requestParams[common.RECV_OFFER_ID] = offerID
		requestParams[common.RECV_AFF_ID] = requestParams[common.AFF_ID]
		requestParams[common.AFF_ID] = common.TRACKER_MEDIA
		requestParams[common.ACTIVITY] = "18"
		//generate new transactionID
		transactionID := common.GenerateTransactionId(rotationOfferID, common.TRACKER_MEDIA)
		requestParams[common.TRANSACTION_ID] = transactionID
		template = dao.GetTemplateByOfferID(rotationOfferID)
	}
	redirect_url := helper.ReplaceURLParameters(requestParams, template)

	//process cookies and redirect to url
	if len(requestParams[common.ADSAUCE_ID]) == 0 {
		userid := "uid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		values := []string{common.ADSAUCE_ID + "=" + userid + ";expires=Fri, 31 Dec 9999 23:59:59 GMT;"}

		//save to message for further processing
		requestParams[common.ADSAUCE_ID] = userid

		//redirecting after setting cookie
		header := map[string]*api.Pair{
			"Location": {
				Key:    "Location",
				Values: []string{redirect_url},
			}, "Set-Cookie": {
				Key:    "Set-Cookie",
				Values: values},
		}
		helper.Redirect(rsp, header)
	} else {
		header := map[string]*api.Pair{
			"Location": {
				Key:    "Location",
				Values: []string{redirect_url},
			},
		}

		helper.Redirect(rsp, header)
	}
	requestParams[common.CLICK_RED_URL] = redirect_url
	requestParams[common.SESSION_IP] = geo.IP
	endTime := time.Now()
	requestParams[common.TIME_TAKEN] = strconv.FormatFloat(endTime.Sub(startTime).Seconds(), 'f', 6, 64)
	fmt.Println("TIME Taken for complete activity:" + requestParams[common.TIME_TAKEN])

	//put details to subscribe click for async processing
	if rotatedClick {
		fwdMap := map[string]string{}
		helper.DuplicateMap(fwdMap, requestParams)
		helper.Subscribe(common.RotatedTopic, fwdMap)
	} else {
		requestParams[common.ACTIVITY] = "2"
		fwdMap := map[string]string{}
		helper.DuplicateMap(fwdMap, requestParams)
		helper.Subscribe(common.ClickTopic, fwdMap)
	}
	endTime = time.Now()
	fmt.Println("TIME Taken before API return:" + strconv.FormatFloat(endTime.Sub(startTime).Seconds(), 'f', 6, 64))
	return nil
}
