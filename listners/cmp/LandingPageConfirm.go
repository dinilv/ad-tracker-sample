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

type Lpc struct {
	LandingPageConfirmListener
}

type LandingPageConfirmListener interface {
	Track(context.Context, *api.Request, *api.Response) error
}

func (tracker *Lpc) Cmp(ctx context.Context, req *api.Request, rsp *api.Response) error {
	fmt.Println("Listener got the request for landing page confirm redirection:")

	startTime := time.Now()

	//rotation
	rotatedConfirm := false

	//cookies & request
	requestParams := make(map[string]string)
	requestParams[constants.API_TIME] = time.Now().UTC().String()

	//req & params
	helper.ParseGetRequest(requestParams, req)

	//escape path & host & headers
	requestParams[constants.HOST] = constants.CMP_HOST
	requestParams[constants.URL] = constants.AFF_LPC + requestParams[constants.URL]
	helper.ParseGetHeader(requestParams, req)
	requestParams[constants.CONFIRM_URL] = requestParams[constants.URL]
	requestParams[constants.URL] = ""

	//geo
	geo := new(model.GeoDetails)
	helper.GetGeo(req, geo)

	//mandatory checks
	offerID, offOK := requestParams[constants.OFFER_ID]
	affiliateID, affOK := requestParams[constants.AFF_ID]
	operator, operatorOK := requestParams[constants.OPERATOR]
	msisdn, msisdnOK := requestParams[constants.MSISDN]

	//missing affiliateID tag as rotated
	if affOK == false || len(affiliateID) == 0 || helper.ID_REGEX.MatchString(requestParams[constants.AFF_ID]) {
		rotatedConfirm = true
		requestParams[constants.RECV_OFFER_ID] = offerID
		requestParams[constants.RECV_AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "140"
		affiliateID = constants.TRACKER_MEDIA
	}
	if offOK == false || len(offerID) == 0 || helper.ID_REGEX.MatchString(requestParams[constants.OFFER_ID]) {
		rotatedConfirm = true
		//query for matching offer with geo
		rotationOfferID := helper.GetMOOfferByGeo(geo.CountryCode)
		requestParams[constants.OFFER_ID] = rotationOfferID
		requestParams[constants.RECV_OFFER_ID] = constants.TRACKER_OFFER
		requestParams[constants.RECV_AFF_ID] = requestParams[constants.AFF_ID]
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "141"
		offerID = rotationOfferID
		affiliateID = constants.TRACKER_MEDIA

	} else if operatorOK == false || len(operator) == 0 {
		rotatedConfirm = true
		//query for matching offer with geo
		rotationOfferID := helper.GetMOOfferByGeo(geo.CountryCode)
		requestParams[constants.OFFER_ID] = rotationOfferID
		requestParams[constants.RECV_OFFER_ID] = constants.TRACKER_OFFER
		requestParams[constants.RECV_AFF_ID] = requestParams[constants.AFF_ID]
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "148"
		offerID = rotationOfferID
		affiliateID = constants.TRACKER_MEDIA

	} else if msisdnOK == false || len(msisdn) == 0 || helper.ID_REGEX.MatchString(requestParams[constants.MSISDN]) {
		rotatedConfirm = true
		//query for matching offer with geo
		rotationOfferID := helper.GetMOOfferByGeo(geo.CountryCode)
		requestParams[constants.OFFER_ID] = rotationOfferID
		requestParams[constants.RECV_OFFER_ID] = constants.TRACKER_OFFER
		requestParams[constants.RECV_AFF_ID] = requestParams[constants.AFF_ID]
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "148"
		offerID = rotationOfferID
		affiliateID = constants.TRACKER_MEDIA

	} else if !dao.ValidateOperator(offerID, operator) {
		requestParams[constants.ACTIVITY] = "147"
		//rotate to another-MO-offer
	} else if !dao.ValidateSubscription(offerID, msisdn) {
		requestParams[constants.ACTIVITY] = "144"
		//rotate to another-MO-offer
	} else if !dao.ValidateMSISDN(offerID, msisdn) {
		requestParams[constants.ACTIVITY] = "145"
		//rotate to CPI offer
	}
	//pixel-ids detection
	if len(requestParams[constants.ADSAUCE_IMPRESSION_ID]) == 0 || len(requestParams[constants.ADSAUCE_BANNER_CLICK_ID]) == 0 ||
		len(requestParams[constants.ADSAUCE_LANDING_PAGE_VIEW_ID]) == 0 {
		rotatedConfirm = true
		requestParams[constants.ACTIVITY] = "146"
	}

	//if transaction-id already present in confirm
	receivedTransactionID, ok := requestParams[constants.TRANSACTION_ID]
	if ok {
		requestParams[constants.RECV_TRANXN_ID] = receivedTransactionID
	}
	if rotatedConfirm {
		//generate and use transactionID
		transactionID := constants.GenerateTransactionId(offerID, affiliateID)
		requestParams[constants.TRANSACTION_ID] = transactionID
	}

	template := dao.GetOperatorTemplateByOfferID(offerID, operator)
	//template check
	if len(template) == 0 {
		rotatedConfirm = true
		//query for matching offer with geo & offer
		rotationOfferID := helper.GetMOOfferFromStack(offerID, geo.CountryCode)
		requestParams[constants.OFFER_ID] = rotationOfferID
		requestParams[constants.RECV_OFFER_ID] = offerID
		requestParams[constants.RECV_AFF_ID] = requestParams[constants.AFF_ID]
		requestParams[constants.AFF_ID] = constants.TRACKER_MEDIA
		requestParams[constants.ACTIVITY] = "18"
		//generate new transactionID
		transactionID := constants.GenerateTransactionId(rotationOfferID, constants.TRACKER_MEDIA)
		requestParams[constants.TRANSACTION_ID] = transactionID
		template = dao.GetTemplateByOfferID(rotationOfferID)
	}

	redirect_url := helper.ReplaceURLParameters(requestParams, template)

	//process cookies and redirect to url
	if len(requestParams[constants.ADSAUCE_ID]) == 0 {
		userID := "uid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		lpcID := "lpcid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		values := []string{constants.ADSAUCE_ID + "=" + userID + ";expires=Fri, 31 Dec 9999 23:59:59 GMT;",
			constants.ADSAUCE_LANDING_PAGE_CONFIRM_ID + "=" + lpcID + ";max-age=1200;"}

		//save to message for further processing
		requestParams[constants.ADSAUCE_ID] = userID
		requestParams[constants.ADSAUCE_LANDING_PAGE_CONFIRM_ID] = lpcID

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

		lpcID := "lpcid_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		requestParams[constants.ADSAUCE_LANDING_PAGE_CONFIRM_ID] = lpcID
		values := []string{constants.ADSAUCE_LANDING_PAGE_CONFIRM_ID + "=" + lpcID + ";max-age=1200;"}

		header := map[string]*api.Pair{
			"Location": {
				Key:    "Location",
				Values: []string{redirect_url},
			}, "Set-Cookie": {
				Key:    "Set-Cookie",
				Values: values},
		}

		helper.Redirect(rsp, header)
	}
	requestParams[constants.CONFIRM_RED_URL] = redirect_url
	requestParams[constants.SESSION_IP] = geo.IP
	endTime := time.Now()
	requestParams[constants.TIME_TAKEN] = strconv.FormatFloat(endTime.Sub(startTime).Seconds(), 'f', 6, 64)
	fmt.Println("TIME Taken for complete activity:" + requestParams[constants.TIME_TAKEN])

	//put details to subscribe confirm for async processing
	if rotatedConfirm {
		fwdMap := map[string]string{}
		helper.DuplicateMap(fwdMap, requestParams)
		helper.Subscribe(constants.MORotatedTopic, fwdMap)
		requestParams[constants.ACTIVITY] = "61"
		helper.Subscribe(constants.LandingPageConfirmTopic, fwdMap)

	} else {
		requestParams[constants.ACTIVITY] = "61"
		fwdMap := map[string]string{}
		helper.DuplicateMap(fwdMap, requestParams)
		helper.Subscribe(constants.LandingPageConfirmTopic, fwdMap)
	}
	endTime = time.Now()
	fmt.Println("TIME Taken before API return:" + strconv.FormatFloat(endTime.Sub(startTime).Seconds(), 'f', 6, 64))
	return nil

}
