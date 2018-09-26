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

func SubscribeLandingPageView(msg map[string]string) *model.TrackLog {

	track := new(model.TrackLog)

	//process message for common fields
	helpers.ProcessMessage(msg, track)

	//host details
	host, _ := os.Hostname()
	track.Host = host

	//According to activity
	switch track.Activity {

	//request is received
	case 101:
		track.Status = constants.Unique
	case 125:
		track.Status = constants.Filtered
		track.Comment = constants.ErrorIDFormat
	case 126:
		track.Status = constants.Filtered
		track.Comment = constants.WRONG_GEO
	case 127:
		track.Status = constants.Filtered
		track.Comment = constants.WRONG_GEO
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
