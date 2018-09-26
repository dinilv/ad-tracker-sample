package v1

import (
	"log"
	"strconv"

	"github.com/adcamie/adserver/common/v1"
	model "github.com/adcamie/adserver/db/model/v1"
	helpers "github.com/adcamie/adserver/helpers/v1"
	logger "github.com/adcamie/adserver/logger"
)

func SubscribeFiltered(msg map[string]string) *model.TrackLog {

	track := new(model.TrackLog)
	//for common fields
	helpers.ProcessMessage(msg, track)

	//According to activity
	switch track.Activity {

	//filtered postbacks
	case 35:
		log.Print("On receiving Filtered Events:")
		track.Status = common.Filtered
		track.Comment = "TransactionID and OfferID are missing, possibly invalid postback to tracker."

	case 36:
		log.Print("On receiving Filtered Events:")
		track.Status = common.Filtered
		track.Comment = "Invalid Transaction Id received in postback."

	case 37:
		log.Print("On receiving Filtered Events:")
		track.Status = common.Filtered
		track.Comment = "Postback already received for this transaction id"
		helpers.CopyTransaction(track.TransactionID, track)

	case 38:
		log.Print("On receiving Fraud Events:")
		track.Status = common.Fraud
		track.Comment = "Session IP is same as Conversion IP."
		helpers.CopyTransaction(track.TransactionID, track)
	case 39:
		log.Print("On receiving Wrong Advertiser IP:")
		track.Status = common.Fraud
		track.Comment = "Received Postback from Non-Advertiser IP for Offer."
		helpers.CopyTransaction(track.TransactionID, track)

	case 40:
		log.Print("On receiving Filtered Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Offer Id is missing, possibly invalid offer media postback(without transactionID)  to tracker."

	case 41:
		log.Print("On receiving Filtered Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Affiliate Id is missing, possibly invalid offer media postback(without transactionID)  to tracker."

	case 42:
		log.Print("On receiving Filtered Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Template is missing, possibly offer-media & media template is not set in tracker."

	case 43:
		log.Print("On receiving Filtered Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Wrong Offer type for offer media postback is received."

	case 50:
		log.Print("On no actions in Postbacks:")
		track.Status = common.Filtered
		track.Comment = "No actions taken on postback. Possible Wrong Offer Type"
	case 51:
		log.Print("On retry of transactionID in Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Retrying filtered transactionIDs"
	case 52:
		log.Print("On retry of transactionID in Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Maximum Retry limit is over."
	case 53:
		log.Print("On Duplicate  transactionID in Postback API:")
		track.Status = common.Filtered
		track.Comment = "Duplicate  transactionID in Postback API"
	case 54:
		log.Print("On Duplicate  transactionID on Delayed Postback:")
		track.Status = common.Filtered
		track.Comment = " Duplicate transactionID on Delayed Postback Subscriber"
	case 55:
		log.Print("On Duplicate  transactionID on RetryPostback:")
		track.Status = common.Filtered
		track.Comment = "Duplicate transactionID: Retry job has already converted."
	case 56:
		log.Print("On Duplicate  transactionID on DelayPostback:")
		track.Status = common.Filtered
		track.Comment = "Duplicate transactionID: Delay job has already converted."
	default:
		log.Print("Not Processed :")
		track.Status = common.Filtered
		track.Comment = "Not processed"
		go logger.ErrorLogger("On Not Processed", "FilteredSubscriber:"+strconv.Itoa(track.Activity), "Switch case failed")

	}
	if recover() != nil {
		log.Println("Something Went Wrong :-()")
		go logger.ErrorLogger("On recover", "FilteredSubscriber", "Switch case failed")
	}

	return track
}
