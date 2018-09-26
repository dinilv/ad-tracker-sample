package v1

import (
	"log"
	"strconv"

	"github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
	helpers "github.com/adcamie/adserver/helpers/v1"
	logger "github.com/adcamie/adserver/logger"
)

func SubscribeClick(msg map[string]string) *model.TrackLog {

	track := new(model.TrackLog)
	track.OfferType = dao.GetOfferTypeOnTranxn(track.OfferID)

	//for common fields
	helpers.ProcessMessage(msg, track)
	//According to activity
	switch track.Activity {

	case 2:
		track.Status = common.Unique

	default:
		log.Print("Not Processed :")
		track.Status = common.Unique
		track.Comment = common.NOT_PROCESSED
		go logger.ErrorLogger("On Not Processed", "ClickSubscriber:"+strconv.Itoa(track.Activity), "Switch case failed")

	}
	if recover() != nil {
		log.Println("Something Went Wrong :-()")
		go logger.ErrorLogger("On recover", "ClickSubscriber", "On recover: Switch case failed")
	}
	return track
}
