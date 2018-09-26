package v1

import (
	"context"
	"fmt"
	"strconv"
	"time"

	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	helper "github.com/adcamie/adserver/helpers"
	api "github.com/micro/micro/api/proto"
)

type Click struct {
	BannerClickListener
}

type BannerClickListener interface {
	Track(context.Context, *api.Request, *api.Response) error
}

func (click *Click) Cmp(ctx context.Context, req *api.Request, rsp *api.Response) error {

	fmt.Println("Listener got the request for banner click redirection:")
	startTime := time.Now()

	//rotation
	rotatedClick := false

	//cookies & request
	requestParams := make(map[string]string)
	requestParams[constants.API_TIME] = time.Now().UTC().String()

	//req & params
	helper.ParseGetRequest(requestParams, req)

	//escape path & host & headers
	requestParams[constants.HOST] = constants.CMP_HOST
	requestParams[constants.URL] = constants.AFF_BC + requestParams[constants.URL]
	helper.ParseGetHeader(requestParams, req)
	requestParams[constants.CLICK_URL] = requestParams[constants.URL]
	requestParams[constants.URL] = ""

	//geo
	geo := new(model.GeoDetails)
	helper.GetGeo(req, geo)

	//mandatory checks
	offerID, offOK := requestParams[constants.OFFER_ID]
	affiliateID, affOK := requestParams[constants.AFF_ID]

	//missing affiliateID tag as rotated
	if affOK == false || len(affiliateID) == 0 || helper.ID_REGEX.MatchString(requestParams[constants.AFF_ID]) {
		rotatedClick = true
		requestParams[constants.RECV_OFFER_ID] = offerID
		requestParams[constants.RECV_AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "15"
		affiliateID = constants.TRACKER_MEDIA
	}
	if offOK == false || len(offerID) == 0 || helper.ID_REGEX.MatchString(requestParams[constants.OFFER_ID]) {
		rotatedClick = true
		//query for matching offer with geo
		rotationOfferID := helper.GetMOOfferByGeo(geo.CountryCode)
		requestParams[constants.OFFER_ID] = rotationOfferID
		requestParams[constants.RECV_OFFER_ID] = constants.TRACKER_OFFER
		requestParams[constants.RECV_AFF_ID] = requestParams[constants.AFF_ID]
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "16"
		offerID = rotationOfferID
		affiliateID = constants.TRACKER_MEDIA

	} else if dao.ValidateExhaustedOfferHash(offerID) {
		rotatedClick = true
		//query for matching offer with geo
		rotationOfferID := helper.GetMOOfferFromStack(offerID, geo.CountryCode)
		requestParams[constants.OFFER_ID] = rotationOfferID
		requestParams[constants.RECV_OFFER_ID] = offerID
		requestParams[constants.RECV_AFF_ID] = requestParams[constants.AFF_ID]
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "22"
		offerID = rotationOfferID
		affiliateID = constants.TRACKER_MEDIA
		fmt.Println("Rotated due to Exhausted Offer")
	} else if dao.ValidateExhaustedOfferAffiliateHash(offerID, affiliateID) {
		rotatedClick = true
		//query for matching offer with geo & offer
		rotationOfferID := helper.GetMOOfferFromStack(offerID, geo.CountryCode)
		requestParams[constants.OFFER_ID] = rotationOfferID
		requestParams[constants.RECV_OFFER_ID] = offerID
		requestParams[constants.RECV_AFF_ID] = requestParams[constants.AFF_ID]
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "23"
		offerID = rotationOfferID
		affiliateID = constants.TRACKER_MEDIA
		fmt.Println("Rotated due to Exhausted Offer & Media")
	} else if !dao.ValidateOfferCountry(offerID, geo.CountryCode) {
		rotatedClick = true
		//query for matching offer with geo & offer
		rotationOfferID := helper.GetMOOfferFromStack(offerID, geo.CountryCode)
		requestParams[constants.OFFER_ID] = rotationOfferID
		requestParams[constants.RECV_OFFER_ID] = offerID
		requestParams[constants.RECV_AFF_ID] = requestParams[constants.AFF_ID]
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "14"
		offerID = rotationOfferID
		affiliateID = constants.TRACKER_MEDIA
		fmt.Println("Rotated due to Wrong Geo")
	} else if !dao.ValidateBlackListIP(geo.IP, offerID) {
	}
	//impression-id detection
	if len(requestParams[constants.ADSAUCE_IMPRESSION_ID]) == 0 {
		rotatedClick = true
		requestParams[constants.ACTIVITY] = "124"
	}

	//if transaction-id already present in click
	_, ok := requestParams[constants.TRANSACTION_ID]
	if !ok {
		//generate and use transactionID
		requestParams[constants.TRANSACTION_ID] = constants.GenerateTransactionId(offerID, affiliateID)
	}

	template := dao.GetTemplateByOfferID(offerID)
	//template check
	if len(template) == 0 {
		rotatedClick = true
		//query for matching offer with geo & offer
		rotationOfferID := helper.GetOfferFromStack(offerID, geo.CountryCode)
		requestParams[constants.OFFER_ID] = rotationOfferID
		requestParams[constants.RECV_OFFER_ID] = offerID
		requestParams[constants.RECV_AFF_ID] = requestParams[constants.AFF_ID]
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "18"
		template = dao.GetTemplateByOfferID(rotationOfferID)
	}
	redirectURL := helper.ReplaceURLParameters(requestParams, template)

	//process cookies and redirect to url
	if len(requestParams[constants.ADSAUCE_ID]) == 0 {
		userID := "uid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		bannerClickID := "bcid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		values := []string{constants.ADSAUCE_ID + "=" + userID + ";expires=Fri, 31 Dec 9999 23:59:59 GMT;",
			constants.ADSAUCE_BANNER_CLICK_ID + "=" + bannerClickID + ";max-age=1200;"}

		//save to message for further processing
		requestParams[constants.ADSAUCE_ID] = userID
		requestParams[constants.ADSAUCE_BANNER_CLICK_ID] = bannerClickID

		//redirecting after setting cookie
		header := map[string]*api.Pair{
			"Location": {
				Key:    "Location",
				Values: []string{redirectURL},
			}, "Set-Cookie": {
				Key:    "Set-Cookie",
				Values: values},
		}
		helper.Redirect(rsp, header)
	} else {

		bannerClickID := "bcid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		requestParams[constants.ADSAUCE_BANNER_CLICK_ID] = bannerClickID
		values := []string{constants.ADSAUCE_BANNER_CLICK_ID + "=" + bannerClickID + ";max-age=1200;"}

		header := map[string]*api.Pair{
			"Location": {
				Key:    "Location",
				Values: []string{redirectURL},
			}, "Set-Cookie": {
				Key:    "Set-Cookie",
				Values: values},
		}

		helper.Redirect(rsp, header)
	}
	requestParams[constants.CLICK_RED_URL] = redirectURL
	requestParams[constants.SESSION_IP] = geo.IP
	endTime := time.Now()
	requestParams[constants.TIME_TAKEN] = strconv.FormatFloat(endTime.Sub(startTime).Seconds(), 'f', 6, 64)
	fmt.Println("TIME Taken for complete activity:" + requestParams[constants.TIME_TAKEN])

	//put details to subscribe click for async processing
	if rotatedClick {
		fwdMap := map[string]string{}
		helper.DuplicateMap(fwdMap, requestParams)
		helper.Subscribe(constants.MORotatedTopic, fwdMap)
	} else {
		requestParams[constants.ACTIVITY] = "2"
		fwdMap := map[string]string{}
		helper.DuplicateMap(fwdMap, requestParams)
		helper.Subscribe(constants.BannerClickTopic, fwdMap)
	}
	endTime = time.Now()
	fmt.Println("TIME Taken before API return:" + strconv.FormatFloat(endTime.Sub(startTime).Seconds(), 'f', 6, 64))
	return nil
}
