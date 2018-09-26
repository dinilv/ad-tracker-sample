package v1

import (
	"time"

	"github.com/adcamie/adserver/common/v1"
	model "github.com/adcamie/adserver/db/model/v1"
)

func SubscribePostbackPing(msg map[string]string) *model.PostBackPingLog {

	postback := new(model.PostBackPingLog)
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()
	postback.Date = rounded
	postback.UTCDate = today
	processPostbackPingMessage(msg, postback)
	return postback
}

func processPostbackPingMessage(msg map[string]string, postback *model.PostBackPingLog) {
	for key, value := range msg {
		switch key {
		case common.TRANSACTION_ID:
			postback.TransactionID = value
		case common.AFF_ID:
			postback.AffiliateID = value
		case common.OFFER_ID:
			postback.OfferID = value
		case common.RESPONSE_CODE:
			postback.ResponseCode = value
		case common.RESPONSE_BODY:
			postback.Response = value
		case common.ERROR:
			postback.ErrorMessage = value
		case "redirect_url":
			postback.SentPostbackURL = value
		case common.TIME_TAKEN:
			postback.TimeTaken = value
		case common.API_TIME:
			postback.APITime = value
		}
	}
}
