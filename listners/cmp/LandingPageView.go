package v1

import (
	"context"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	helper "github.com/adcamie/adserver/helpers"
	api "github.com/micro/micro/api/proto"
)

type Lpv struct {
	LandingPageViewListener
}

type LandingPageViewListener interface {
	Track(context.Context, *api.Request, *api.Response) error
}

func (tracker *Lpv) Cmp(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Listener got the request to CMP Landing Page View:")

	//cookies & request
	requestParams := make(map[string]string)
	requestParams[constants.API_TIME] = time.Now().UTC().String()

	//req & params
	helper.ParseGetRequest(requestParams, req)

	//escape path & host & headers
	requestParams[constants.HOST] = constants.CMP_HOST
	requestParams[constants.URL] = constants.AFF_LPV + requestParams[constants.URL]
	helper.ParseGetHeader(requestParams, req)
	requestParams[constants.PING_URL] = requestParams[constants.URL]
	requestParams[constants.URL] = ""

	//geo
	geo := new(model.GeoDetails)
	helper.GetGeo(req, geo)

	//mandatory checks
	offerID, offOK := requestParams[constants.OFFER_ID]
	affiliateID, affOK := requestParams[constants.AFF_ID]

	if offOK == false || len(offerID) == 0 || helper.ID_REGEX.MatchString(requestParams[constants.OFFER_ID]) {
		requestParams[constants.ACTIVITY] = "130"
	} else if affOK == false || len(affiliateID) == 0 || helper.ID_REGEX.MatchString(requestParams[constants.AFF_ID]) {
		requestParams[constants.ACTIVITY] = "130"
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
	} else if dao.ValidateExhaustedOfferHash(offerID) {
		requestParams[constants.ACTIVITY] = "131"
	} else if !dao.ValidateOfferCountry(offerID, geo.CountryCode) {
		requestParams[constants.ACTIVITY] = "132"
	} else {
		//valid
		if len(requestParams[constants.ADSAUCE_ID]) == 0 {
			userID := "uid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
			landingPageViewID := "lpvid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
			values := []string{constants.ADSAUCE_ID + "=" + userID + ";expires=Fri, 31 Dec 9999 23:59:59 GMT;",
				constants.ADSAUCE_LANDING_PAGE_VIEW_ID + "=" + landingPageViewID + ";max-age=1200;"}
			//add to response
			rsp.Header = map[string]*api.Pair{"Set-Cookie": {Key: "Set-Cookie", Values: values}}
			//add to message
			requestParams[constants.ADSAUCE_ID] = userID
			requestParams[constants.ADSAUCE_LANDING_PAGE_VIEW_ID] = landingPageViewID
		} else {
			//set lpv cookie for fraud detection
			landingPageViewID := "lpvid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
			values := []string{constants.ADSAUCE_LANDING_PAGE_VIEW_ID + "=" + landingPageViewID + ";max-age=900;"}
			//add to response
			rsp.Header = map[string]*api.Pair{"Set-Cookie": {Key: "Set-Cookie", Values: values}}
			//add to message
			requestParams[constants.ADSAUCE_LANDING_PAGE_VIEW_ID] = landingPageViewID
		}
		requestParams[constants.ACTIVITY] = "60"
	}

	//pixel-ids detection
	if len(requestParams[constants.ADSAUCE_IMPRESSION_ID]) == 0 && len(requestParams[constants.ADSAUCE_BANNER_CLICK_ID]) == 0 {
		requestParams[constants.ACTIVITY] = "133"
	}

	//async processing
	helper.Subscribe(constants.LandingPageViewTopic, requestParams)

	//sucesss
	rsp.StatusCode = 200
	return nil
}
