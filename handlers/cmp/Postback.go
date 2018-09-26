package v1

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/dao"
	helper "github.com/adcamie/adserver/helpers/v1"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/googlepubsub"
)

type PostbackReq struct {
	IsConverted   bool
	AffiliateID   string
	OfferID       string
	OfferType     string
	TransactionId string
	ConversionIP  string
	SessionIP     string
	GoalID        string
	ClickURL      string
}

type PostbackRes struct {
	Status              string
	ReceivedOfferID     string
	ReceivedAffiliateID string
	Url                 string
	ClickURL            string
	Activity            int
	IsPingRequired      bool
	SessionIP           string
}

type PostbackHandler interface {
	Setpostback(context.Context, *PostbackReq, *PostbackRes) error
}

type Postback struct {
	PostbackHandler
}

func RegisterTrackerPostbackHandler(s server.Server, hdlr PostbackHandler) {
	fmt.Print("Getting Postback Setting Handler")
	s.Handle(s.NewHandler(&Postback{hdlr}))
}

func (postback *Postback) Setpostback(ctx context.Context, req *PostbackReq, rsp *PostbackRes) error {

	fmt.Println("Setting Postback in tracker service for transaction id: ", req.TransactionId)
	rsp.Activity = 50

	//get click log from mongo
	filters := make(map[string]interface{})
	filters[constants.TransactionID] = req.TransactionId
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
			aff_sub, aff_sub2 := dao.GetTransactionWithAffiliate(req.TransactionId)

			//from redis backup too
			if len(aff_sub) == 0 {
				valid, redisTransaction := dao.SearchRedisKeysFromESBackup(req.TransactionId)
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
			validatePostbackForMedia(req, rsp)
			//check sent or un-sent conversion
			rsp.Activity = 4
			if rsp.IsPingRequired {
				createPostback(req, rsp)
				rsp.Activity = 3
			}

		case "2", "4":
			rsp.Activity = 6
			//CPI with postevents: sent conversion or not
			if dao.ValidateSentTransactionID(req.TransactionId) {
				createPostback(req, rsp)
				rsp.Activity = 5
			}
		}

	} else if !req.IsConverted && strings.Compare(rsp.Status, constants.Rotated) != 0 {
		fmt.Println("Not Converted and Not rotated")
		switch req.OfferType {

		case "1", "2", "3", "7", "9":
			//CPI with conversion for MSISDN-MO & CPI
			validatePostbackForMedia(req, rsp)
			//check sent or un-sent conversion
			rsp.Activity = 4
			if rsp.IsPingRequired {
				createPostback(req, rsp)
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
					validatePostbackForMedia(req, rsp)
					//check sent or un-sent conversion
					rsp.Activity = 4
					if rsp.IsPingRequired {
						createPostback(req, rsp)
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
			validateRotatedPostbackForMedia(req, rsp)

			//check sent or un-sent conversion
			rsp.Activity = 8
			if rsp.IsPingRequired {
				req.OfferID = rsp.ReceivedOfferID
				req.AffiliateID = rsp.ReceivedAffiliateID
				createPostback(req, rsp)
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
					validatePostbackForMedia(req, rsp)
					//check sent or un-sent conversion
					rsp.Activity = 8
					if rsp.IsPingRequired {
						req.OfferID = rsp.ReceivedOfferID
						req.AffiliateID = rsp.ReceivedAffiliateID
						createPostback(req, rsp)
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
			validateRotatedPostbackForMedia(req, rsp)
			//check sent or un-sent conversion
			rsp.Activity = 8
			if rsp.IsPingRequired {
				createPostback(req, rsp)
				rsp.Activity = 7
			}

		case "2", "4", "7":
			rsp.Activity = 10
			//CPI/MO-TransactionID with postevents: sent conversion or not
			if dao.ValidateSentTransactionID(req.TransactionId) {
				createPostback(req, rsp)
				rsp.Activity = 9
			}
		}
	}

	return nil
}

func validatePostbackForMedia(req *PostbackReq, rsp *PostbackRes) {

	//Check for offer affiliate mqf
	mqf := dao.GetMQFByOfferAndAffiliate(req.OfferID, req.AffiliateID)
	if len(mqf) == 0 {
		//check for global affiliate mqf
		mqf = dao.GetMQFByAffiliate(req.AffiliateID)
		if len(mqf) == 0 {
			mqf = "0.7"
		}
	}
	mqfFloat, _ := strconv.ParseFloat(mqf, 64)
	total_conversion_count, sent_conversion_count := dao.GetConversionData(req.OfferID, req.AffiliateID)
	total_conversion_float, _ := strconv.ParseFloat(total_conversion_count, 64)
	sent_conversion_float, _ := strconv.ParseFloat(sent_conversion_count, 64)

	if total_conversion_float == 0.0 {
		fmt.Println("Total Conversion is Zero")
		rsp.IsPingRequired = true
	} else if (sent_conversion_float / total_conversion_float) <= mqfFloat {
		fmt.Println("sent.total", sent_conversion_float, total_conversion_float)
		rsp.IsPingRequired = true
	} else {
		fmt.Println("Conversion MQF is less")
		rsp.IsPingRequired = false
	}

}

func validateRotatedPostbackForMedia(req *PostbackReq, rsp *PostbackRes) {

	mqf := 0.3
	total_conversion_count, sent_conversion_count := dao.GetRotatedConversionData(req.OfferID)
	total_conversion_float, _ := strconv.ParseFloat(total_conversion_count, 64)
	sent_conversion_float, _ := strconv.ParseFloat(sent_conversion_count, 64)

	if (sent_conversion_float / total_conversion_float) <= mqf {
		fmt.Println("sent.total", sent_conversion_float, total_conversion_float)
		rsp.IsPingRequired = true
	} else if total_conversion_float == 0.0 {
		fmt.Println("Total Conversion is Zero")
		rsp.IsPingRequired = true
	} else {
		fmt.Println("Conversion MQF is less")
		rsp.IsPingRequired = false
	}

}

func createPostback(req *PostbackReq, rsp *PostbackRes) {

	postbackURLTemplate := dao.GetOfferAffiliatePostbackTemplate(req.OfferID, req.AffiliateID)
	//Check offer affiliate specific template is available or not
	if len(postbackURLTemplate) == 0 || postbackURLTemplate == "" {
		postbackURLTemplate = dao.GetTemplateByAffiliateID(req.AffiliateID)
		//Check affiliate specific template is available or not
		if len(postbackURLTemplate) != 0 || postbackURLTemplate != "" {
			postbackURLTemplate = helper.ReplaceTemplateParameters(rsp.ClickURL, postbackURLTemplate, req.TransactionId)
		}
	} else {
		postbackURLTemplate = helper.ReplaceTemplateParameters(rsp.ClickURL, postbackURLTemplate, req.TransactionId)
	}

	rsp.Url = postbackURLTemplate

}
