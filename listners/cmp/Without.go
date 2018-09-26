package v1

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	helper "github.com/adcamie/adserver/helpers"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
)

func (postback *Postback) Without(ctx context.Context, req *api.Request, rsp *api.Response) error {

	fmt.Print("Listener got the postback request:")

	requestParams := make(map[string]string)
	requestParams[constants.API_TIME] = time.Now().UTC().String()
	requestParams[constants.Processor] = constants.POSTBACK_API

	//process request headers & geo
	helper.ParseGetHeader(requestParams, req)

	//add method to logging
	requestParams[constants.Method] = req.Method
	requestParams[constants.URL] = requestParams[constants.URL] + constants.AFF_L

	//process received postback
	helper.ParsePostbackReceived(requestParams, req)

	//check with whether offerID exists, offert type without transactionID
	offerID, offOK := requestParams[constants.OFFER_ID]

	//get offer type
	offerType := dao.GetOfferTypeOnMaster(requestParams[constants.OFFER_ID])

	response := &model.PostbackRes{}
	//check for offer media postback
	if offOK == true && len(offerID) > 0 {

		//Check Affiliate ID is received or not
		affiliateID, affOK := requestParams[constants.AFF_ID]
		if affOK == false || len(affiliateID) == 0 {
			fmt.Print("Affiliate ID is missing  in Parameters: Invalid postback.")
			//log to filtered postbacks
			requestParams[constants.ACTIVITY] = "41"
			helper.Subscribe(constants.FilteredTopic, requestParams)
			return errors.BadRequest("go.micro.api.v1.postback", "Affiliate ID is Not Present In Parameters.")
		}

		//validate offer type
		if strings.Compare(offerType, "8") != 0 || strings.Compare(offerType, "10") != 0 {
			fmt.Print("Wrong Offer type received: Invalid postback.", offerType, requestParams[constants.OFFER_ID])
			//log to filtered postbacks
			requestParams[constants.ACTIVITY] = "43"
			helper.Subscribe(constants.FilteredTopic, requestParams)
			return errors.BadRequest("go.micro.api.v1.postback", "Wrong Offer type received.")

		} else {

			//check ip is coming as parameter in url
			ip := "0.0.0.0"
			if len(requestParams[constants.IP]) > 0 {
				ip = requestParams[constants.IP]
			} else {
				//split for load-balancer adding ip
				ips := strings.Split(requestParams["X-Forwarded-For"], ",")
				ip = ips[0]
			}
			//check received_ip is whitelisted on advertiser or not & split ip since lb adds its own
			valid := dao.ValidateAdveriserWithOfferID(requestParams[constants.OFFER_ID], ip)
			if !valid {
				//log to filtered postbacks
				requestParams[constants.ACTIVITY] = "39"
				helper.Subscribe(constants.FilteredTopic, requestParams)
				return errors.BadRequest("go.micro.api.v1.postback", "Advertiser IP is not whitelisted.")

			}

			setPostbackReq := new(model.PostbackReq)
			setPostbackReq.AffiliateID = requestParams[constants.AFF_ID]
			setPostbackReq.OfferID = requestParams[constants.OFFER_ID]
			setPostbackReq.ClickURL = requestParams[constants.URL]
			setPostbackRsp := new(model.PostbackRes)
			//validate media MQF
			helper.ValidatePostbackForMedia(setPostbackReq, setPostbackRsp)
			//check sent or un-sent conversion
			setPostbackRsp.Activity = 12
			if setPostbackRsp.IsPingRequired {
				setPostbackRsp.ClickURL = setPostbackReq.ClickURL
				helper.CreatePostback(setPostbackReq, setPostbackRsp)
				if len(setPostbackRsp.URL) == 0 {
					setPostbackRsp.Activity = 42
				} else {
					setPostbackRsp.Activity = 11
				}
				log.Println("Printing Offer Media at the last:", setPostbackRsp.Activity)
			}

		}
	}

	helper.ParsePostbackResponse(requestParams, response)
	//return sucesss for postback always
	rsp.StatusCode = 200
	return nil
}
