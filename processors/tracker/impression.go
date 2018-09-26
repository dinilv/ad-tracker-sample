package v1

import (
	"os"
	"strconv"

	constants "github.com/adcamie/adserver/common/v1"
	model "github.com/adcamie/adserver/db/model/v1"
	helpers "github.com/adcamie/adserver/helpers/v1"
	logger "github.com/adcamie/adserver/logger"
)

func SubscribeImpression(msg map[string]string) *model.TrackLog {

	track := new(model.TrackLog)
	//process message for constants fields
	helpers.ProcessMessage(msg, track)

	//host details
	host, _ := os.Hostname()
	track.Host = host

	//According to activity
	switch track.Activity {

	//request is received
	case 1:
		//Check offer_id,affiliate_id is there or not
		if len(track.OfferID) == 0 || track.OfferID == "" || len(track.AffiliateID) == 0 || track.AffiliateID == "" ||
			helpers.Regex.MatchString(track.OfferID) || helpers.Regex.MatchString(track.AffiliateID) {
			track.Activity = 120
			track.Status = constants.Filtered
			track.Comment = constants.ErrorIDFormat
		} else {
			track.Status = constants.Unique
		}
	case 120:
		track.Status = constants.Filtered
		track.Comment = constants.ErrorIDFormat
	case 121:
		track.Status = constants.Filtered
		track.Comment = constants.EXHAUSETD_OFFER
	case 122:
		track.Status = constants.Filtered
		track.Comment = constants.WRONG_GEO
	default:
		track.Status = constants.Filtered
		track.Comment = constants.NOT_PROCESSED
		go logger.ErrorLogger("On Not Processed", "ImpressionSubscriber:"+strconv.Itoa(track.Activity), "Switch case failed")

	}
	if recover() != nil {
		go logger.ErrorLogger("On recover", "ImpressionSubscriber", "Switch case failed")
	}

	return track

}
