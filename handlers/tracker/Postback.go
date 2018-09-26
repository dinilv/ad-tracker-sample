package v1

import (
	"context"
	"fmt"
	"strings"

	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	helper "github.com/adcamie/adserver/helpers"
	"github.com/micro/go-micro/server"
)

type PostbackHandler interface {
	Setpostback(context.Context, *model.PostbackReq, **model.PostbackRes) error
}

type Postback struct {
	PostbackHandler
}

func RegisterTrackerPostbackHandler(s server.Server, hdlr PostbackHandler) {
	fmt.Print("Getting Postback Setting Handler")
	s.Handle(s.NewHandler(&Postback{hdlr}))
}

func (postback *Postback) Track(ctx context.Context, req *model.PostbackReq, rsp *model.PostbackRes) error {

	fmt.Println("Setting Postback in tracker service for transaction id: ", req.TransactionID)
	rsp.Activity = 50

	//get click log from mongo
	filters := make(map[string]interface{})
	filters[constants.TransactionID] = req.TransactionID
	results := dao.QueryLatestClickFromMongo(filters)
	if len(results) > 0 {
		rsp.ClickURL = results[0].ClickURL
		rsp.ReceivedOfferID = results[0].ReceivedOfferID
		rsp.ReceivedAffiliateID = results[0].ReceivedAffiliateID
		rsp.Status = results[0].Status
		rsp.SessionIP = results[0].SessionIP
	} else {
		//if primary mongo is down take from backup mongo
		results = dao.QueryLatestClickFromMongoBackup(filters)
		if len(results) > 0 {
			rsp.ClickURL = results[0].ClickURL
			rsp.ReceivedOfferID = results[0].ReceivedOfferID
			rsp.ReceivedAffiliateID = results[0].ReceivedAffiliateID
			rsp.Status = results[0].Status
			rsp.SessionIP = results[0].SessionIP
		} else {
			//no mongo log available then take from redis
			aff_sub, aff_sub2 := dao.GetTransactionWithAffiliate(req.TransactionID)

			//from redis backup too
			if len(aff_sub) == 0 {
				valid, redisTransaction := dao.SearchRedisKeysFromESBackup(req.TransactionID)
				if valid {
					rsp.ClickURL = "http://track.adcamie.com/aff_c?aff_sub=" + redisTransaction.AffiliateSub + "&aff_sub2=" + redisTransaction.AffiliateSub2
				}
			} else {
				rsp.ClickURL = "http://track.adcamie.com/aff_c?aff_sub=" + aff_sub + "&aff_sub2=" + aff_sub2
			}

		}
	}

	fmt.Println("ConversionIP", req.ConversionIP)
	fmt.Println("SessionIP", rsp.SessionIP)

	//Check for fraud detection
	if strings.Compare(req.OfferType, "9") != 0 && strings.Compare(req.ConversionIP, rsp.SessionIP) == 0 {
		fmt.Println("ips", req.ConversionIP, rsp.SessionIP)
		fmt.Println("ConversionIP", req.ConversionIP)
		fmt.Println("SessionIP", rsp.SessionIP)
		rsp.Activity = 0
		return nil
	}

	//according to offer type & isConverted & status(rotated) decide what to do
	if req.IsConverted && strings.Compare(rsp.Status, constants.Rotated) != 0 {

		fmt.Println("Converted and Not rotated")
		switch req.OfferType {
		case "1":
			//MO CPA campaign(possibility of duplicate transaction)
			helper.ValidatePostbackForMedia(req, rsp)
			//check sent or un-sent conversion
			rsp.Activity = 4
			if rsp.IsPingRequired {
				helper.CreatePostback(req, rsp)
				rsp.Activity = 3
			}

		case "2", "4":
			rsp.Activity = 6
			//CPI with postevents: sent conversion or not
			if dao.ValidateSentTransactionID(req.TransactionID) {
				helper.CreatePostback(req, rsp)
				rsp.Activity = 5
			}
		}

	} else if !req.IsConverted && strings.Compare(rsp.Status, constants.Rotated) != 0 {
		fmt.Println("Not Converted and Not rotated")
		switch req.OfferType {

		case "1", "2", "3", "7", "9":
			//CPI with conversion for MSISDN-MO & CPI
			helper.ValidatePostbackForMedia(req, rsp)
			//check sent or un-sent conversion
			rsp.Activity = 4
			if rsp.IsPingRequired {
				helper.CreatePostback(req, rsp)
				rsp.Activity = 3
			}

		case "4":
			//CPI with custom conversion point: decide conversion or event. Query offerID
			filters = make(map[string]interface{})
			filters[constants.OfferID] = req.OfferID
			offers := dao.GetOfferFromMongo(filters)
			if len(offers) > 0 {
				offer := offers[0]
				if strings.Compare(offer.GoalID, req.GoalID) == 0 {
					//conversion is valid
					helper.ValidatePostbackForMedia(req, rsp)
					//check sent or un-sent conversion
					rsp.Activity = 4
					if rsp.IsPingRequired {
						helper.CreatePostback(req, rsp)
						rsp.Activity = 3

					}
				} else {
					//log to Events
					rsp.Activity = 6

				}
			}

		}

	} else if !req.IsConverted && strings.Compare(rsp.Status, constants.Rotated) == 0 {

		fmt.Println("Not Converted and rotated")
		//not converted rotated offer postback
		switch req.OfferType {

		case "1", "2", "3", "7":
			//CPI with conversion for MO & CPI
			helper.ValidateRotatedPostbackForMedia(req, rsp)

			//check sent or un-sent conversion
			rsp.Activity = 8
			if rsp.IsPingRequired {
				req.OfferID = rsp.ReceivedOfferID
				req.AffiliateID = rsp.ReceivedAffiliateID
				helper.CreatePostback(req, rsp)
				rsp.Activity = 7

			}

		case "4":
			//CPI with custom conversion point: decide conversion or event. Query offerID
			filters = make(map[string]interface{})
			filters[constants.OfferID] = req.OfferID
			offers := dao.GetOfferFromMongo(filters)
			if len(offers) > 0 {
				offer := offers[0]
				if strings.Compare(offer.GoalID, req.GoalID) == 0 {
					//conversion is valid
					helper.ValidatePostbackForMedia(req, rsp)
					//check sent or un-sent conversion
					rsp.Activity = 8
					if rsp.IsPingRequired {
						req.OfferID = rsp.ReceivedOfferID
						req.AffiliateID = rsp.ReceivedAffiliateID
						helper.CreatePostback(req, rsp)
						rsp.Activity = 7
					}
				} else {
					//log to rotated events
					rsp.Activity = 10

				}
			}

		}

	} else if req.IsConverted && strings.Compare(rsp.Status, constants.Rotated) == 0 {
		fmt.Println("Converted and  rotated")
		switch req.OfferType {
		case "1":
			//MO CPA campaign(possibility of duplicate transaction)
			helper.ValidateRotatedPostbackForMedia(req, rsp)
			//check sent or un-sent conversion
			rsp.Activity = 8
			if rsp.IsPingRequired {
				helper.CreatePostback(req, rsp)
				rsp.Activity = 7
			}

		case "2", "4", "7":
			rsp.Activity = 10
			//CPI/MO-TransactionID with postevents: sent conversion or not
			if dao.ValidateSentTransactionID(req.TransactionID) {
				helper.CreatePostback(req, rsp)
				rsp.Activity = 9
			}
		}
	}

	return nil
}
