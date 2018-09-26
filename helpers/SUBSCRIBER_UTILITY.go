package v1

import (
	"log"
	"strconv"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	"github.com/dghubble/sling"
)

func ConvertMinuteToReportMinute(minute int) int {

	if minute >= 0 && minute <= 14 {
		return 0
	} else if minute > 14 && minute <= 29 {
		return 15
	} else if minute > 29 && minute <= 44 {
		return 30
	} else {
		return 45
	}
}

func ProcessMessage(msg map[string]string, track *model.TrackLog) {
	activity, _ := strconv.Atoi(msg[constants.ACTIVITY])
	track.Activity = activity

	//handling ip on url as priority
	if len(msg[constants.IP]) > 0 {
		track.IP = msg[constants.IP]
	}

	log.Println("Async Handler For Tracker Received message:- ", activity)

	//process request headers and parameters
	for key, value := range msg {
		ProcessURlParameters(key, value, track)
	}
	//type
	if len(track.OfferID) > 0 {
		track.OfferType = dao.GetOfferTypeOnMaster(track.OfferID)
	}
	//Time stamp detailmsgs
	utc := time.Now().UTC()
	track.UTCDate = utc
	track.Date = time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	track.Year = utc.Year()
	track.Hour = utc.Hour()
	track.Minute = ConvertMinuteToReportMinute(utc.Minute())
	track.Month = utc.Month().String()
	_, track.Week = utc.ISOWeek()
	track.MonthID = constants.MonthIds[utc.Month().String()]

	//Location and geography details
	ip := track.IP
	geo := new(model.GeoDetails)
	if activity == 2 {
		sling.New().Get(GetFreeGeoIPClick(ip)).ReceiveSuccess(geo)
	} else {
		sling.New().Get(GetFreeGeoIP(ip)).ReceiveSuccess(geo)
	}
	track.Geo = *geo

}

func ProcessURlParameters(key string, value string, track *model.TrackLog) {

	switch key {

	//related to request header
	case "User-Agent":
		track.UserAgent = value
	case "Content-Type":
		track.ContentType = value
	case "method":
		track.Method = value
	case "body":
		track.RequestBody = value
	case "X-Forwarded-For":
		//split for load-balancer adding ip
		ipsplitted := strings.Split(value, ",")
		ip := ipsplitted[0]
		if len(track.IP) > 0 {
			track.SessionIP = track.IP
			track.ConversionIP = track.IP
		} else {
			track.IP = ip
			track.SessionIP = ip
			track.ConversionIP = ip
		}
	case "X-Requested-With":
		track.RequestedPackage = value
	case "Referer":
		track.Referer = value
		//custom variables passed
	case "url":
		track.ReceivedPostbackURL = value
	case "redirect_url":
		track.SentPostbackURL = value
	case "click_url":
		track.ClickURL = value
	case "click_red_url":
		track.ClickRedirectURL = value
	//performance monitoring
	case "time_taken":
		track.ResponseTime, _ = strconv.ParseFloat(value, 64)
	case constants.SubscriberName:
		track.SubscriberName = value
	case constants.Processor:
		track.Processor = value
	case constants.API_TIME:
		track.APITime = value
	case constants.BROKER_TIME:
		track.BrokerTime = value
	case constants.GO_ROUTINE_TIME:
		track.GoRoutineTime = value
	case constants.SubscriptionTime:
		track.SubscriptionTime = value
	case constants.RetryCount:
		track.RetryCount = value
	case constants.ErrorMessage:
		track.ErrorMessage = value
	case constants.FilteredCount:
		track.FilteredCount = value

	//related to user
	case constants.ADSAUCE_ID:
		track.CookieID = value
	case "click_date":
		timeClick, _ := time.Parse(time.RFC1123Z, value)
		track.ClickDate = timeClick

	//mandatory parameters
	case constants.OFFER_ID:
		track.OfferID = value
	case constants.AFF_ID:
		track.AffiliateID = value
	case constants.TRANSACTION_ID:
		track.TransactionID = value
	case constants.MessageID:
		track.MessageID = value

	//special varibales
	case constants.RECV_OFFER_ID:
		track.ReceivedOfferID = value
	case constants.RECV_AFF_ID:
		track.ReceivedAffiliateID = value
	case constants.RECV_TRANXN_ID:
		track.ReceivedTransactionID = value
	case constants.ROTATED_CLICK:
		track.Status = constants.Rotated
	case constants.OFFER_TYPE:
		track.OfferType = value

	//variables
	case "seq":
		track.Sequence, _ = strconv.Atoi(value)
	case "msisdn":
		track.MSISDN = value
	case "operator":
		track.Operator = value
	case "service_code":
		track.ServiceCode = value
	case "service_id":
		track.ServiceID = value
	case "cp_code":
		track.CPCode = value

	//HO Optional parameters
	case "advertiser_id":
		track.AdvertiserID = value
	case "advertiser_ref":
		track.AdvertiserRefID = value
	case "adv_sub":
		track.AdvertiserSub = value
	case "adv_sub1":
		track.AdvertiserSub1 = value
	case "adv_sub2":
		track.AdvertiserSub2 = value
	case "adv_sub3":
		track.AdvertiserSub3 = value
	case "adv_sub4":
		track.AdvertiserSub4 = value
	case "adv_sub5":
		track.AdvertiserSub5 = value
	case "aff_sub":
		track.AffiliateSub = value
	case "aff_sub1":
		track.AffiliateSub1 = value
	case "aff_sub2":
		track.AffiliateSub2 = value
	case "aff_sub3":
		track.AffiliateSub3 = value
	case "aff_sub4":
		track.AffiliateSub4 = value
	case "aff_sub5":
		track.AffiliateSub5 = value
	case "aff_sub6":
		track.AffiliateSub6 = value
	case "aff_sub7":
		track.AffiliateSub7 = value
	case "aff_sub8":
		track.AffiliateSub8 = value
	case "aff_sub9":
		track.AffiliateSub9 = value
	case "aff_sub10":
		track.AffiliateSub10 = value
	case "affiliate_name":
		track.AffiliateName = value
	case "affiliate_ref":
		track.AffiliateRef = value
	case "currency":
		track.Currency = value
	case "date":
		track.ConvertedDate = value
	case "datetime":
		track.ConvertedDateTime = value
	case "file_name":
		track.CreativeFile = value
	case "goal_id":
		track.GoalID = value
	case "ip":
		track.IP = value
		track.SessionIP = value
		track.ConversionIP = value
	case "payout":
		track.AffiliatePayout = value
	case "revenue":
		track.AffiliateRevenue = value
	case "goal_ref":
		track.GoalRef = value
	case "offer_file_id":
		track.OfferFileID = value
	case "offer_name":
		track.OfferName = value
	case "offer_url_id":
		track.OfferURLID = value
	case "ran":
		track.Ran = value
	case "sale_amount":
		track.SaleAmount = value
	case "session_ip":
		track.SessionIP = value
	case "source":
		track.Source = value
	case "time":
		track.ConvertedTime = value
	case "g_aid":
		track.GoogleAID = value
	case "android_id":
		track.AndroidID = value
	case "android_id_md5":
		track.AndroidIDMD5 = value
	case "android_id_sha1":
		track.AndroidIDSHA1 = value
	case "device_brand":
		track.DeviceBrand = value
	case "device_id":
		track.DeviceID = value
	case "device_id_md5":
		track.DeviceIDMD5 = value
	case "device_id_sha1":
		track.DeviceIDSHA1 = value
	case "device_model":
		track.DeviceModel = value
	case "device_os":
		track.DeviceOS = value
	case "device_os_version":
		track.DeviceOSVersion = value
	case "ios_ifa":
		track.IOSIfa = value
	case "ios_ifa_md5":
		track.IOSIfaMD5 = value
	case "ios_ifa_sha1":
		track.IOSIfaSHA1 = value
	case "ios_ifv":
		track.IOSIfv = value
	case "mac_address":
		track.MacAddress = value
	case "mac_address_md5":
		track.MacAddressMD5 = value
	case "mac_address_sha1":
		track.MacAddressSHA1 = value
	case "windows_aid":
		track.WindowsAID = value
	case "windows_aid_md5":
		track.WindowsAIDMD5 = value
	case "windows_aid_sha1":
		track.WindowsSHA1 = value
	case "odin":
		track.ODIN = value
	case "open_udid":
		track.OpenUDID = value
	case "unid":
		track.UNID = value
	case "user_id":
		track.AppUserID = value

	//PostEventLog
	case "event_status":
		track.EventStatus = value
	case "amount":
		amount, _ := strconv.ParseFloat(value, 32)
		track.Amount = amount
	case "conversion_unique_id":
		track.ConversionUniqueID = value

	}

}

func CopyTransaction(transactionID string, track *model.TrackLog) {

	//find from click log
	filters := make(map[string]interface{})
	filters[constants.TransactionID] = transactionID
	clickLog := dao.QueryLatestClickFromMongoSession(filters)
	//click log is only valid for 30 days
	if len(clickLog) > 0 {
		//copy all details from click to conversion
		track.ReceivedOfferID = clickLog[0].ReceivedOfferID
		track.ReceivedAffiliateID = clickLog[0].ReceivedAffiliateID
		track.OfferID = clickLog[0].OfferID
		track.AffiliateID = clickLog[0].AffiliateID
		track.OfferType = clickLog[0].OfferType
		track.ClickURL = clickLog[0].ClickURL
		track.ClickDate = clickLog[0].UTCDate
		track.ClickRedirectURL = clickLog[0].ClickRedirectURL
		track.SessionIP = clickLog[0].SessionIP
		track.ClickGeo = clickLog[0].Geo
		track.CookieID = clickLog[0].CookieID
		track.UserAgent = clickLog[0].UserAgent
		track.Referer = clickLog[0].Referer
		track.RequestedPackage = clickLog[0].RequestedPackage
		track.AdvertiserID = clickLog[0].AdvertiserID
		track.AdvertiserRefID = clickLog[0].AdvertiserRefID
		track.AdvertiserSub = clickLog[0].AdvertiserSub
		track.AffiliateSub = clickLog[0].AffiliateSub
		track.AffiliateSub2 = clickLog[0].AffiliateSub2
		track.AffiliateSub3 = clickLog[0].AffiliateSub3
		track.AffiliateSub4 = clickLog[0].AffiliateSub4
		track.AffiliateSub5 = clickLog[0].AffiliateSub5
		track.AffiliateName = clickLog[0].AffiliateName
		track.AffiliateRef = clickLog[0].AffiliateRef
	} else {
		//check on backup Mongo
		clickLog := dao.QueryLatestClickFromMongoBackup(filters)
		if len(clickLog) > 0 {
			//copy all details from click to conversion
			track.ReceivedOfferID = clickLog[0].ReceivedOfferID
			track.ReceivedAffiliateID = clickLog[0].ReceivedAffiliateID
			track.OfferID = clickLog[0].OfferID
			track.AffiliateID = clickLog[0].AffiliateID
			track.OfferType = clickLog[0].OfferType
			track.ClickURL = clickLog[0].ClickURL
			track.ClickDate = clickLog[0].UTCDate
			track.ClickRedirectURL = clickLog[0].ClickRedirectURL
			track.SessionIP = clickLog[0].SessionIP
			track.ClickGeo = clickLog[0].Geo
			track.CookieID = clickLog[0].CookieID
			track.UserAgent = clickLog[0].UserAgent
			track.Referer = clickLog[0].Referer
			track.RequestedPackage = clickLog[0].RequestedPackage
			track.AdvertiserID = clickLog[0].AdvertiserID
			track.AdvertiserRefID = clickLog[0].AdvertiserRefID
			track.AdvertiserSub = clickLog[0].AdvertiserSub
			track.AffiliateSub = clickLog[0].AffiliateSub
			track.AffiliateSub2 = clickLog[0].AffiliateSub2
			track.AffiliateSub3 = clickLog[0].AffiliateSub3
			track.AffiliateSub4 = clickLog[0].AffiliateSub4
			track.AffiliateSub5 = clickLog[0].AffiliateSub5
			track.AffiliateName = clickLog[0].AffiliateName
			track.AffiliateRef = clickLog[0].AffiliateRef

		} else {
			//if click not found take from transaction add offerID & AffiliateID
			track.OfferID, track.AffiliateID = dao.GetTransaction(transactionID)
			if strings.Compare(track.OfferID, constants.OfferDefault) == 0 && strings.Compare(track.AffiliateID, constants.AffiliateDefault) == 0 {
				//take from mongo backup transactions
				valid, redisBackupModel := dao.SearchRedisKeysFromESBackup(transactionID)
				if valid {
					track.OfferID = redisBackupModel.OfferID
					track.AffiliateID = redisBackupModel.AffiliateID
				}
			}
		}
	}

}
