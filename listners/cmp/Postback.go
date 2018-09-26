package v1

import (
	"context"
	"fmt"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	helper "github.com/adcamie/adserver/helpers"
	logger "github.com/adcamie/adserver/logger"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
)

type Postback struct {
	PostbackListener
}

type PostbackListener interface {
	Tracker(context.Context, *api.Request, *api.Response) error
}

func (postback *Postback) Cmp(ctx context.Context, req *api.Request, rsp *api.Response) error {
	fmt.Print("Listener got the postback request:")

	requestParams := make(map[string]string)
	requestParams[constants.API_TIME] = time.Now().UTC().String()
	requestParams[constants.Processor] = constants.POSTBACK_API

	//process request headers & geo
	helper.ParseGetHeader(requestParams, req)

	//add method to logging
	requestParams[constants.Method] = req.Method
	requestParams[constants.URL] = requestParams[constants.URL] + constants.AFF_LSR

	//process all parameters
	helper.ParsePostbackReceived(requestParams, req)

	//check with whether transactionID exists, offert type with transactionID
	transactionID, transOK := requestParams[constants.TRANSACTION_ID]

	response := &model.PostbackRes{}

	if transOK == false || len(transactionID) == 0 {
		//postback is missing with details
		fmt.Println("TransactionID & OfferID are missing  in Parameters: Invalid postback.")
		//log to filtered postbacks
		requestParams[constants.ACTIVITY] = "35"
		helper.Subscribe(constants.FilteredTopic, requestParams)
		return errors.BadRequest("go.micro.api.v1.postback", "TransactionID is Not Present In Parameters.")
	} else if transOK == true && len(transactionID) > 0 {

		//check transactionID is already received or not for the day
		isConvertedTransaction := dao.ValidateConvertedTransactionID(transactionID)
		apiRecieved := dao.ValidateTransactionIDOnBackup(transactionID)
		if apiRecieved && !isConvertedTransaction && strings.Compare(requestParams[constants.Processor], constants.RETRY_JOB) != 0 {
			requestParams[constants.ACTIVITY] = "53"
			helper.Subscribe(constants.FilteredTopic, requestParams)
			return errors.BadRequest("go.micro.api.v1.tracker", "TransactionID already received. Try every transaction only once.")
		} else {
			dao.SavePostbackTransaction(transactionID)
		}

		//validate transactionID exists or not
		valid := dao.ValidateTransactionID(transactionID)
		if !valid {
			fmt.Println("TransactionID is invalid in postback.", transactionID)
			//log to filtered postbacks
			requestParams[constants.ACTIVITY] = "36"
			fwdMap := map[string]string{}
			helper.DuplicateMap(fwdMap, requestParams)
			helper.Subscribe(constants.FilteredTopic, fwdMap)

			if helper.TRANSACTION_REGEX.MatchString(transactionID) {
				//queue this for processing again after delay of 10 minutes
				requestParams[constants.RetryCount] = "1"
				requestParams[constants.ACTIVITY] = "51"
				helper.Subscribe(constants.DelayedPostbackTopic, requestParams)
				rsp.StatusCode = 200
				return nil
			} else {
				return errors.BadRequest("go.micro.api.v1.postback", "Transaction Id is Invalid.")
			}

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
			//check received_ip is whitelisted on advertiser or not
			valid, offerID := dao.ValidateAdveriserIPWithTransaction(transactionID, ip)
			if !valid {
				//log to filtered postbacks
				requestParams[constants.ACTIVITY] = "39"
				helper.Subscribe(constants.FilteredTopic, requestParams)
				return errors.BadRequest("go.micro.api.v1.postback", "Advertiser IP is not whitelisted.")

			}

			//check transaction id exists, checking for duplicate transaction or postevents
			isConverted, offerType, affiliateID := dao.ValidateTransactionIDForPostback(transactionID, offerID)
			//offer type without transactionID
			if strings.Compare(offerType, "8") == 0 {
				fmt.Println("Wrong Offer type received: Invalid postback.", offerType, requestParams[constants.OFFER_ID])
				//log to filtered postbacks
				requestParams[constants.ACTIVITY] = "43"
				helper.Subscribe(constants.FilteredTopic, requestParams)
				return errors.BadRequest("go.micro.api.v1.postback", "Wrong Offer type received.")

			}
			//offer type without postevents and already converted
			if isConverted && (strings.Compare(offerType, "3") == 0 || strings.Compare(offerType, "7") == 0) {
				fmt.Println("TransactionId Exists in converted postbacks and postevents not enabled")
				//log to filtered postbacks
				requestParams[constants.ACTIVITY] = "37"
				helper.Subscribe(constants.FilteredTopic, requestParams)
				return errors.BadRequest("go.micro.api.v1.postback", "Transaction Id is duplicate.")

			}

			setpostback := new(model.PostbackReq)
			setpostback.TransactionID = transactionID
			setpostback.AffiliateID = affiliateID
			setpostback.OfferID = offerID
			setpostback.OfferType = offerType
			setpostback.IsConverted = isConverted
			setpostback.GoalID = requestParams[constants.GOAL_ID]
			setpostback.ConversionIP = ip

			//for postback ping logs
			requestParams[constants.OFFER_ID] = offerID
			requestParams[constants.AFF_ID] = affiliateID

			request := client.NewJsonRequest("go.micro.service.v1.postback", "Postback.Setpostback", setpostback)
			if err := client.Call(ctx, request, response); err != nil {
				fmt.Print("Client Calling Error In Tracker Set Postback:", err, request, response)
				go logger.ErrorLogger(err.Error(), "PostBackAPI", "Calling Postback Service. Try 1")
				if err := client.Call(ctx, request, response); err != nil {
					fmt.Print("Client Calling Error In Tracker Set Postback:", err, request, response)
					go logger.ErrorLogger(err.Error(), "PostBackAPI", "Calling Postback Service. Try 2")
					if err := client.Call(ctx, request, response); err != nil {
						fmt.Print("Client Calling Error In Tracker Set Postback:", err, request, response)
						go logger.ErrorLogger(err.Error(), "PostBackAPI", "Calling Postback Service. Try 3")
						return err
					}
				}
			}

			fmt.Println("ConversionIP", setpostback.ConversionIP)
			fmt.Println("SessionIP", response.SessionIP)
		}
	}
	fmt.Println("Response Activity:", response.Activity)
	fmt.Println(response, "response", transactionID)

	helper.ParsePostbackResponse(requestParams, response)
	//return sucesss for postback always
	rsp.StatusCode = 200
	return nil
}
