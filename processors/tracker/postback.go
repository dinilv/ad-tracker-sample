package v1

import (
	"log"
	"strconv"

	"github.com/adcamie/adserver/common/v1"
	model "github.com/adcamie/adserver/db/model/v1"
	helpers "github.com/adcamie/adserver/helpers/v1"
	logger "github.com/adcamie/adserver/logger"
)

func SubscribePostback(msg map[string]string) *model.TrackLog {

	track := new(model.TrackLog)

	//process message for common fields
	helpers.ProcessMessage(msg, track)

	//handle mo postback
	if len(track.CPCode) > 0 {
		track.ServiceID = track.ServiceCode
		track.ServiceCode = track.CPCode + track.ServiceCode
	}

	switch track.Activity {
	case 3:
		log.Print("On receiving sent conversion: ")
		track.Status = common.Sent
		track.Comment = "Conversion happened on  this transaction Id"
		helpers.CopyTransaction(track.TransactionID, track)
	case 4:
		log.Print("On receiving un-sent conversion: ")
		track.Status = common.UnSent
		track.Comment = "Conversion happened on this transaction Id, but not forwarded to media due to MQF"
		helpers.CopyTransaction(track.TransactionID, track)

	case 5:
		log.Print("On receiving post events: ")
		track.Status = common.Sent
		helpers.CopyTransaction(track.TransactionID, track)
		track.Comment = "Post Events forwarded to media"

	case 6:
		log.Print("On receiving post events un sent : ")
		helpers.CopyTransaction(track.TransactionID, track)
		track.Status = common.UnSent
		track.Comment = "Post Events Not forwarded to media"

	case 7:
		log.Print("On receiving rotated offer sent conversion: ")
		track.Status = common.RotatedSent
		track.Comment = "Conversion forwarded to media"
		helpers.CopyTransaction(track.TransactionID, track)
	case 8:
		log.Print("On receiving rotated offer un-sent conversion: ")
		track.Status = common.RotatedUnSent
		track.Comment = "Conversion forwarded to media"
		helpers.CopyTransaction(track.TransactionID, track)

	case 9:
		log.Print("On receiving rotated post events sent : ")
		track.Status = common.RotatedSent
		track.Comment = "Event forwarded to media"
		helpers.CopyTransaction(track.TransactionID, track)

	case 10:
		log.Print("On receiving rotated post events un sent : ")
		track.Status = common.RotatedUnSent
		track.Comment = "Event forwarded to media"
		helpers.CopyTransaction(track.TransactionID, track)

	case 11:
		log.Print("On receiving sent postback without transaction ID : ")
		track.Status = common.Sent
		track.Comment = "Conversion happened on this offer and media without transactionID."

	case 12:
		log.Print("On unsent postback without transaction ID : ")
		track.Status = common.UnSent
		track.Comment = "Conversion happened on this offer and media, but not forwarded due to MQF."
	case 13:
		log.Print("On unsent postback without transaction ID & without media template: ")
		track.Status = common.UnSent
		track.Comment = "Conversion happened on this offer and media, but not forwarded beacuse no media template found."
	default:
		log.Print("Not Processed :")
		track.Status = common.UnSent
		track.Comment = "Not processed"
		go logger.ErrorLogger("On Not Processed", "PostbackSubscriber:"+strconv.Itoa(track.Activity), "Switch case failed")

	}
	if recover() != nil {
		log.Println("Something Went Wrong :-()")
		go logger.ErrorLogger("On recover", "PostbackSubscriber", "Switch case failed")
	}

	return track
}
