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

type Opv struct {
	OperatorPageViewListener
}

type OperatorPageViewListener interface {
	Track(context.Context, *api.Request, *api.Response) error
}

func (tracker *Opv) Cmp(ctx context.Context, req *api.Request, rsp *api.Response) error {

	log.Print("Listener got the request to track impression:")

	//cookies & request
	requestParams := make(map[string]string)
	requestParams[constants.API_TIME] = time.Now().UTC().String()

	//req & params
	helper.ParseGetRequest(requestParams, req)

	//escape path & host & headers
	requestParams[constants.HOST] = constants.CMP_HOST
	requestParams[constants.URL] = constants.AFF_OPV + requestParams[constants.URL]
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
		requestParams[constants.ACTIVITY] = "151"
	} else if affOK == false || len(affiliateID) == 0 || helper.ID_REGEX.MatchString(requestParams[constants.AFF_ID]) {
		requestParams[constants.ACTIVITY] = "151"
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
	} else if dao.ValidateExhaustedOfferHash(offerID) {
		requestParams[constants.ACTIVITY] = "152"
	} else {
		//valid
		if len(requestParams[constants.ADSAUCE_ID]) == 0 {
			userID := "uid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
			opvID := "opvid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
			values := []string{constants.ADSAUCE_ID + "=" + userID + ";expires=Fri, 31 Dec 9999 23:59:59 GMT;",
				constants.ADSAUCE_OPERATOR_PAGE_VIEW_ID + "=" + opvID + ";max-age=1200;"}
			//add to response
			rsp.Header = map[string]*api.Pair{"Set-Cookie": {Key: "Set-Cookie", Values: values}}
			//add to message
			requestParams[constants.ADSAUCE_ID] = userID
			requestParams[constants.ADSAUCE_OPERATOR_PAGE_VIEW_ID] = opvID
		} else {
			//set impression cookie for fraud detection
			opvID := "opvid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
			values := []string{constants.ADSAUCE_OPERATOR_PAGE_VIEW_ID + "=" + opvID + ";max-age=1200;"}
			//add to response
			rsp.Header = map[string]*api.Pair{"Set-Cookie": {Key: "Set-Cookie", Values: values}}
			//add to message
			requestParams[constants.ADSAUCE_OPERATOR_PAGE_VIEW_ID] = opvID
		}
		requestParams[constants.ACTIVITY] = "62"
	}

	//pixel-ids detection
	if len(requestParams[constants.ADSAUCE_IMPRESSION_ID]) == 0 && len(requestParams[constants.ADSAUCE_BANNER_CLICK_ID]) == 0 {
		requestParams[constants.ACTIVITY] = "153"
	}

	//async processing
	helper.Subscribe(constants.OperatorPageViewTopic, requestParams)

	//sucesss
	rsp.StatusCode = 200
	return nil
}
