package v1

import (
	"log"
	"strconv"

	"github.com/adcamie/adserver/common/v1"
	model "github.com/adcamie/adserver/db/model/v1"
	helpers "github.com/adcamie/adserver/helpers/v1"
	logger "github.com/adcamie/adserver/logger"
)

func SubscribeOperatorPageView(msg map[string]string) *model.TrackLog {

	track := new(model.TrackLog)
	//process message for common fields
	helpers.ProcessMessage(msg, track)

	//According to activity
	switch track.Activity {

	//request is received
	case 61:

		//Check offer_id,affiliate_id is there or not
		//filiterd if validation fails
		if len(track.OfferID) == 0 || track.OfferID == "" || len(track.AffiliateID) == 0 || track.AffiliateID == "" ||
			helpers.Regex.MatchString(track.OfferID) || helpers.Regex.MatchString(track.AffiliateID) {
			//not a valid impression
			track.Status = common.Filtered
			track.Comment = "offer_id or affiliate_id is missing or pattern for same is not proper."

		} else {
			track.Status = common.Unique
		}
	default:
		log.Print("Not Processed :")
		track.Status = common.Filtered
		track.Comment = "Not processed"
		go logger.ErrorLogger("On Not Processed", "LandingPageSubscriber:"+strconv.Itoa(track.Activity), "Switch case failed")

	}
	if recover() != nil {
		log.Println("Something Went Wrong :-()")
		go logger.ErrorLogger("On recover", "LandingPageSubscriber", "Switch case failed")
	}

	return track

}
