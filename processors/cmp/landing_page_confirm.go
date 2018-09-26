package v1

import (
	"log"
	"os"
	"strconv"

	constants "github.com/adcamie/adserver/common/v1"
	model "github.com/adcamie/adserver/db/model/v1"
	helpers "github.com/adcamie/adserver/helpers/v1"
	logger "github.com/adcamie/adserver/logger"
)

func SubscribeLandingPageConfirm(msg map[string]string) *model.TrackLog {

	track := new(model.TrackLog)

	//process message for common fields
	helpers.ProcessMessage(msg, track)

	//host details
	host, _ := os.Hostname()
	track.Host = host

	//According to activity
	switch track.Activity {

	//request is received
	case 11:
		track.Status = constants.Unique
	case 128:
		track.Status = constants.Filtered
		track.Comment = constants.ErrorIDFormat
	case 129:
		track.Status = constants.Filtered
		track.Comment = constants.WRONG_GEO
	case 130:
		track.Status = constants.Filtered
		track.Comment = "Invalid carrier for offer"
	case 131:
		track.Status = constants.Filtered
		track.Comment = "Received Invalid MSISDN on MO-Offer ."
	case 132:
		track.Status = constants.Filtered
		track.Comment = "Offer already Subscribed by this MSISDN."
	case 133:
		track.Status = constants.Filtered
		track.Comment = "MSISDN received on LPC is in Black list."
	default:
		track.Status = constants.Filtered
		track.Comment = constants.NotProcessed
		go logger.ErrorLogger("On Not Processed", "LandingPageViewSubscriber:"+strconv.Itoa(track.Activity), "Switch case failed")

	}
	if recover() != nil {
		log.Println("Something Went Wrong :-()")
		go logger.ErrorLogger("On recover", "LandingPageViewSubscriber", "Switch case failed")
	}

	return track

}
