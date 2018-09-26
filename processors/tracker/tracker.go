package v1

import (
	"log"
	"strconv"
	"strings"

	"github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
	helpers "github.com/adcamie/adserver/helpers/v1"
)

func SubscribeTracking(msg map[string]string) error {

	track := new(model.TrackLog)
	activity, _ := strconv.Atoi(msg[common.ACTIVITY])
	track.Activity = activity
	log.Println("Async Handler For Tracker Received message:- ", activity)

	//According to activity
	switch activity {

	//request is received
	case 1:

		//Check offer_id,affiliate_id is there or not
		//filiterd if validation fails
		if len(track.OfferID) == 0 || track.OfferID == "" || len(track.AffiliateID) == 0 || track.AffiliateID == "" ||
			helpers.Regex.MatchString(track.OfferID) || helpers.Regex.MatchString(track.AffiliateID) {
			//not a valid impression
			track.Status = common.Filtered
			track.Comment = "offer_id or affiliate_id is missing or pattern for same is not proper."
			dao.InsertToMongoSession(common.MongoDB, common.FilteredImpressionLog, &track)

		} else {
			track.Status = common.Unique
			dao.InsertToMongoSession(common.MongoDB, common.ImpressionLog, &track)
			var trackerLog []interface{}
			trackerLog = append(trackerLog, track)

		}

	case 2:

		if strings.Compare(track.OfferType, "8") != 0 {
			dao.SaveClick(track.TransactionID, track.OfferID, track.AffiliateID, track.AffiliateSub, track.AffiliateSub2)
			//check the user_id, offer_id, affiliate_id  exists
			if strings.Compare(track.Status, common.Rotated) == 0 {
				track.Status = common.Rotated
				track.Comment = "Rotated click due to wrong geo or offer rotation."
			} else {
				track.Status = common.Unique
			}
			//save to db
			dao.InsertToMongoSession(common.MongoDB, common.ClickLog, &track)
			var trackerLog []interface{}
			trackerLog = append(trackerLog, track)
		} else {
			log.Println("Log not saved due to silent campaign type")
		}

	//on valid conversions
	case 3:
		log.Print("On receiving sent conversion: ")
		track.Status = common.Sent
		track.Comment = "Conversion happened on  this transaction Id"
		helpers.CopyTransaction(track.TransactionID, track)
		dao.SaveSentPostBacks(track.TransactionID, track.OfferID, track.AffiliateID)
		dao.InsertToMongoSession(common.MongoDB, common.PostBackLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)
	case 4:
		log.Print("On receiving un-sent conversion: ")
		track.Status = common.UnSent
		track.Comment = "Conversion happened on this transaction Id, but not forwarded to media due to MQF"
		helpers.CopyTransaction(track.TransactionID, track)
		dao.SaveUnSentPostBacks(track.TransactionID, track.OfferID, track.AffiliateID)
		dao.InsertToMongoSession(common.MongoDB, common.PostBackLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)
	case 5:
		log.Print("On receiving post events: ")
		track.Status = common.Sent
		helpers.CopyTransaction(track.TransactionID, track)
		track.Comment = "Post Events forwarded to media"
		dao.InsertToMongoSession(common.MongoDB, common.PostEventLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)
	case 6:
		log.Print("On receiving post events un sent : ")
		helpers.CopyTransaction(track.TransactionID, track)
		track.Status = common.UnSent
		track.Comment = "Post Events Not forwarded to media"
		dao.InsertToMongoSession(common.MongoDB, common.PostEventLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)
	case 7:
		log.Print("On receiving rotated offer sent conversion: ")
		track.Status = common.RotatedSent
		track.Comment = "Conversion forwarded to media"
		helpers.CopyTransaction(track.TransactionID, track)
		dao.SaveSentPostBacks(track.TransactionID, track.OfferID, common.TRACKER_MEDIA)
		dao.InsertToMongoSession(common.MongoDB, common.PostBackLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)

	case 8:
		log.Print("On receiving rotated offer un-sent conversion: ")
		track.Status = common.RotatedUnSent
		track.Comment = "Conversion forwarded to media"
		helpers.CopyTransaction(track.TransactionID, track)
		dao.SaveUnSentPostBacks(track.TransactionID, track.OfferID, common.TRACKER_MEDIA)
		dao.InsertToMongoSession(common.MongoDB, common.PostBackLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)

	case 9:
		log.Print("On receiving rotated post events un sent : ")
		track.Status = common.RotatedSent
		track.Comment = "Event forwarded to media"
		helpers.CopyTransaction(track.TransactionID, track)
		dao.InsertToMongoSession(common.MongoDB, common.PostEventLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)

	case 10:
		log.Print("On receiving rotated post events un sent : ")
		track.Status = common.RotatedUnSent
		track.Comment = "Event forwarded to media"
		helpers.CopyTransaction(track.TransactionID, track)
		dao.InsertToMongoSession(common.MongoDB, common.PostEventLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)

	case 11:
		log.Print("On receiving sent postback without transaction ID : ")
		track.Status = common.Sent
		track.Comment = "Conversion happened on this offer and media without transactionID."
		dao.InsertToMongoSession(common.MongoDB, common.PostBackLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)

	case 12:
		log.Print("On unsent postback without transaction ID : ")
		track.Status = common.UnSent
		track.Comment = "Conversion happened on this offer and media, but not forwarded due to MQF."
		dao.InsertToMongoSession(common.MongoDB, common.PostBackLog, &track)
		var trackerLog []interface{}
		trackerLog = append(trackerLog, track)

	//rotated Click Events
	case 14:
		log.Print("On receiving MO-Click, invalid MSISDN as transactionID:")
		track.Status = common.Rotated
		track.Comment = "Received Invalid MSISDN on MO-Offer Click."
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)
	case 15:
		log.Print("On receiving MO-Click, already Subscribed Offer:")
		track.Status = common.Rotated
		track.Comment = "Offer already Subscribed by this MSISDN."
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)
	case 16:
		log.Print("On receiving Click of Exhausted Offer:")
		track.Status = common.Rotated
		track.Comment = "Rotation is enabled due to: manual intervention/daily budget/total budget/daily cap/campaign paused."
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)

	case 17:
		log.Print("On receiving Click without OfferID:")
		track.Status = common.Rotated
		track.Comment = "offer_id is Missing or offer_id is not in proper format On Click."
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)

	case 18:
		log.Print("On receiving Wrong Geo Click:")
		track.Status = common.Rotated
		track.Comment = "Wrong Geo Click"
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)

	//filtered click events
	case 19:
		log.Print("On receiving Duplicate transaction ID on Click:")
		track.Status = common.Rotated
		track.Comment = "Duplicate Transaction ID"
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)
	case 21:
		log.Print("On receiving Filtered Events:")
		track.Status = common.Rotated
		track.Comment = "aff_id is missing in URL or aff_id is not proper"
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)
	case 22:
		log.Print("On receiving Filtered Events:")
		track.Status = common.Rotated
		track.Comment = "Template is missing for the offer,possibly offer is not set in tracker."
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)
	case 23:
		log.Print("On receiving Blocked MSISDN Click:")
		track.Status = common.Rotated
		track.Comment = "MSISDN received on click is in blocked list."
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)

	case 24:
		log.Print("On receiving Template not found:")
		track.Status = common.Rotated
		track.Comment = "No templates found for offer."
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)

	case 25:
		log.Print("On receiving Exhausted Offer & Media :")
		track.Status = common.Rotated
		track.Comment = "Rotation Enbaled on Offer & Media."
		dao.InsertToMongoSession(common.MongoDB, common.RotatedClickLog, &track)

	//filtered postbacks
	case 35:
		log.Print("On receiving Filtered Events:")
		track.Status = common.Filtered
		track.Comment = "TransactionID and OfferID are missing, possibly invalid postback to tracker."
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)

	case 36:
		log.Print("On receiving Filtered Events:")
		track.Status = common.Filtered
		track.Comment = "Invalid Transaction Id received in postback."
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)

	case 37:
		log.Print("On receiving Filtered Events:")
		track.Status = common.Filtered
		track.Comment = "Postback already received for this transaction id"
		helpers.CopyTransaction(track.TransactionID, track)
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)

	case 38:
		log.Print("On receiving Fraud Events:")
		track.Status = common.Fraud
		track.Comment = "Session IP is same as Conversion IP."
		helpers.CopyTransaction(track.TransactionID, track)
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)
	case 39:
		log.Print("On receiving Wrong Advertiser IP:")
		track.Status = common.Fraud
		track.Comment = "Received Postback from Non-Advertiser IP for Offer."
		helpers.CopyTransaction(track.TransactionID, track)
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)

	case 40:
		log.Print("On receiving Filtered Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Offer Id is missing, possibly invalid offer media postback(without transactionID)  to tracker."
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)

	case 41:
		log.Print("On receiving Filtered Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Affiliate Id is missing, possibly invalid offer media postback(without transactionID)  to tracker."
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)

	case 42:
		log.Print("On receiving Filtered Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Template is missing, possibly offer-media & media template is not set in tracker."
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)

	case 43:
		log.Print("On receiving Filtered Postbacks:")
		track.Status = common.Filtered
		track.Comment = "Wrong Offer type for offer media postback is received."
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)
	case 50:
		log.Print("On no actions in Postbacks:")
		track.Status = common.Filtered
		track.Comment = "No actions taken on postback. Possible Wrong Offer Type"
		dao.InsertToMongoSession(common.MongoDB, common.FilteredPostBackLog, &track)
	}
	if recover() != nil {
		log.Print("Something Went Wrong :-()")
	}

	return nil
}
