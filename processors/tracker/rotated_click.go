package v1

import (
	"log"
	"strconv"

	"github.com/adcamie/adserver/common/v1"
	model "github.com/adcamie/adserver/db/model/v1"
	helpers "github.com/adcamie/adserver/helpers/v1"
	logger "github.com/adcamie/adserver/logger"
)

func SubscribeRotatedClick(msg map[string]string) *model.TrackLog {

	track := new(model.TrackLog)
	//for common fields
	helpers.ProcessMessage(msg, track)
	track.Status = common.Rotated

	//According to activity
	switch track.Activity {

	//wrong geo
	case 14:
		log.Println("On receiving Wrong Geo Click:")
		track.Comment = "Wrong Geo Click"

	//mandatory parameters
	case 15:
		log.Print("On receiving Click without proper affiliate ID:")
		track.Comment = "aff_id is Missing or aff_id is not in proper format On Click."
	case 16:
		log.Println("On receiving Click without OfferID:")
		track.Comment = "offer_id is Missing or offer_id is not in proper format On Click."
	case 17:
		log.Println("On receiving MO-Click without MSISDN:")
		track.Comment = "MSISDN is Missing or MSISDN is not in proper format On Click."
	case 18:
		log.Println("Template missing")
		track.Comment = "Template is missing."

	//mo-click validations
	case 19:
		log.Println("On receiving MO-Click, invalid MSISDN as transactionID:")
		track.Comment = "Received Invalid MSISDN on MO-Offer Click."
	case 20:
		log.Println("On receiving MO-Click, already Subscribed Offer:")
		track.Comment = "Offer already Subscribed by this MSISDN."
	case 21:
		log.Print("On receiving Blocked MSISDN Click:")
		track.Comment = "MSISDN received on click is in blocked list."
	case 124:
		log.Print("On receiving click without impression-id")
		track.Comment = common.ERROR_IMPRESSION_COOKIE_ID

	//exhausted
	case 22:
		log.Println("On receiving Click of Exhausted Offer:")
		track.Comment = "Rotated Offer due to: manual intervention/daily budget/total budget/daily cap/campaign paused."
	case 23:
		log.Println("On receiving Exhausted Offer & Media :")
		track.Comment = "Rotated Offer & Media due to: manual intervention/daily budget/total budget/daily cap/media paused."

	default:
		log.Println("Not Processed :")
		track.Comment = "Not processed"
		go logger.ErrorLogger("On Not Processed", "RotatedClickSubscriber:-"+strconv.Itoa(track.Activity), "Switch case failed")

	}
	if recover() != nil {
		log.Println("Something Went Wrong :-()")
		go logger.ErrorLogger("On recover", "RotatedClickSubscriber", "On recover. Switch case failed")
	}

	return track
}
