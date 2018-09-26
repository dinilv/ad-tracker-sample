package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	constants "github.com/adcamie/adserver/common"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model"
	logger "github.com/adcamie/adserver/logger"
	"github.com/dghubble/sling"
	api "github.com/micro/micro/api/proto"
	"gopkg.in/mgo.v2/bson"
)

var topics []string
var host = ""

var TRANSACTION_REGEX = regexp.MustCompile("gg15\\d+")
var ID_REGEX = regexp.MustCompile("\\D+")

var brokers = map[string]*pubsub.Topic{}

func InitialiseTrackerBrokers() {

	//host
	host = "http://track.adcamie.com"

	//for backing up error pub/sub messages from google
	db.InitializeMongoBackup()

	topics = []string{constants.ImpressionTopic, constants.ClickTopic, constants.RotatedTopic,
		constants.PostbackTopic, constants.PostbackPingTopic, constants.DelayedPostbackTopic,
		constants.FilteredTopic}

	for _, topic := range topics {
		client, err := pubsub.NewClient(context.Background(), constants.ProjectName)
		if err != nil {
			logger.ErrorLogger(err.Error(), "GooglePubSub", "Topic: "+topic+" Initialization Error")
			panic(err.Error())
		}
		topicClient := client.Topic(topic)
		brokers[topic] = topicClient
	}
}

func InitialiseCMPBrokers() {

	//host
	host = "http://motrack.adcamie.com"

	//for backing up error pub/sub messages from google
	db.InitializeMongoBackup()

	topics = []string{
		constants.MOImpressionTopic, constants.BannerClickTopic, constants.RotatedTopic,
		constants.LandingPageViewTopic, constants.LandingPageConfirmTopic, constants.OperatorPageViewTopic,
		constants.ContentViewTopic, constants.PostbackTopic, constants.PostbackPingTopic,
		constants.DelayedPostbackTopic, constants.FilteredTopic}

	for _, topic := range topics {
		client, err := pubsub.NewClient(context.Background(), constants.ProjectName)
		if err != nil {
			logger.ErrorLogger(err.Error(), "GooglePubSub", "Topic: "+topic+" Initialization Error")
			panic(err.Error())
		}
		topicClient := client.Topic(topic)
		brokers[topic] = topicClient
	}
}

func ParseGetRequest(requestParams map[string]string, req *api.Request) {
	//for url recreation and convert url paramters to map
	var params string
	for key, get := range req.Get {
		if len(get.Values[0]) < 1000 && len(key) < 1000 {
			params = params + key + "=" + get.Values[0] + "&"
		}
		for _, val := range get.Values {
			if strings.Compare(key, "google_aid") == 0 {
				value := requestParams[key]
				key = "g_aid"
				requestParams[key] = value
			} else if len(val) < 1000 && len(key) < 1000 {
				requestParams[key] = val
			} else {
				fmt.Println("Length of URL key or value execeeded:", key, val)
			}
		}
	}

	//recreation of url, discard last character
	if len(params) > 0 {
		params = params[:len(params)-1]
	}
	requestParams[constants.URL] = requestParams[constants.URL] + params
	//cut down url to 1024 character limit
	if len(requestParams[constants.URL]) > 1000 {
		requestParams[constants.URL] = requestParams[constants.URL][0:1023]
	}
}

func ParsePostRequest(requestParams map[string]string, req *api.Request) {

	//for convert post url paramters to map
	for key, get := range req.GetPost() {
		for _, val := range get.Values {
			if strings.Compare(key, "google_aid") == 0 {
				value := requestParams[key]
				key = "g_aid"
				requestParams[key] = value
			} else if len(val) < 1000 && len(key) < 1000 {
				requestParams[key] = val
			} else {
				fmt.Println("Length of URL key & value execeeded:", key, val)
			}
		}
	}
}

func ParseGetHeader(requestParams map[string]string, req *api.Request) {
	for _, header := range req.GetHeader() {
		for _, val := range header.Values {

			switch header.Key {

			case "Cookie":
				cookies := strings.Split(val, ";")
				for _, cookie := range cookies {
					cookieStrings := strings.Split(cookie, "=")
					if len(cookieStrings) > 1 {
						if len(cookieStrings[0]) < 1000 && len(cookieStrings[1]) < 1000 {
							requestParams[strings.Trim(cookieStrings[0], " ")] = cookieStrings[1]
						}
					}
				}
			case "Host", "User-Agent", "Content-Type", "method", "body", "X-Forwarded-For", "X-Requested-With", "Referer":
				//checking character limitation
				if len(val) < 1000 && len(header.Key) < 1000 {
					requestParams[header.Key] = val
				} else {
					fmt.Println("Length of header key execeeded:", val)
				}
			}
		}
	}

	//create host url
	if len(requestParams["Host"]) == 0 {
		full_url := host + requestParams[constants.URL]
		requestParams[constants.URL] = full_url
	} else {
		full_url := "http://" + requestParams["Host"] + requestParams[constants.URL]
		requestParams[constants.URL] = full_url
	}
	//full url cut down if length is more than 1024
	if len(requestParams[constants.URL]) > 1000 {
		requestParams[constants.URL] = requestParams[constants.URL][0:1020]
	}
}

func ParseMultiPartReq(contentType string, requestParams map[string]string, req *api.Request) {

	var splitContentType = strings.Split(contentType, "boundary=")

	var boundary = splitContentType[len(splitContentType)-1]
	mr := multipart.NewReader(strings.NewReader(req.Body), boundary)
	for {
		p, err := mr.NextPart()
		if err != nil {
			fmt.Println(err)
			break
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(p)

		//get form value
		var formName = p.FormName()
		var formValue = buf.String()
		//assign to req map
		requestParams[formName] = formValue

	}

}

func ParseJSONReq(requestParams map[string]string, req *api.Request) {

	//for supporting JSON post request
	var postReq = new(model.PostBackReq)
	json.Unmarshal([]byte(req.Body), &postReq)

	//pass parameters to request map
	if len(postReq.AffiliateID) > 0 {
		requestParams[constants.AFF_ID] = postReq.AffiliateID
	}
	if len(postReq.OfferID) > 0 {
		requestParams[constants.OFFER_ID] = postReq.OfferID
	}
	if len(postReq.TransactionID) > 0 {
		requestParams[constants.TRANSACTION_ID] = postReq.TransactionID
	}

}

func GetGeo(req *api.Request, geo *model.GeoDetails) {
	//get geo details for filtering and rotation
	ip := req.GetHeader()["X-Forwarded-For"].Values[0]
	if len(ip) == 0 {
		ip = req.GetHeader()["X-Real-Ip"].Values[0]
		if len(ip) == 0 {
			log.Println("NO IP headers worked, you idiot. :-(")
		}
	}
	//split ip since lb adds its own
	ips := strings.Split(ip, ",")
	sling.New().Get(GetFreeGeoIP(ips[0])).ReceiveSuccess(geo)

}

func Redirect(rsp *api.Response, header map[string]*api.Pair) {
	log.Println("Redirection started")
	startTime := time.Now()
	rsp.StatusCode = int32(303)
	rsp.Header = header
	endTime := time.Now()
	log.Println("Time Taken For redirectin:-", endTime.Sub(startTime).Seconds())

}

func Ping(url string, msg map[string]string) {

	log.Println("Pinging started")
	if len(url) > 0 {
		startTime := time.Now()
		sucess := bson.M{}
		rsp, _ := sling.New().Get(url).ReceiveSuccess(sucess)
		endTime := time.Now()
		timeTaken := endTime.Sub(startTime).Seconds()
		log.Println("Time Taken For Postback Ping:-", timeTaken)
		//add for Logging
		var rspBytes []byte
		rsp.Body.Read(rspBytes)
		msg[constants.RESPONSE_BODY] = string(rspBytes)
		msg[constants.RESPONSE_CODE] = strconv.Itoa(rsp.StatusCode)
		timeTakenString := strconv.FormatFloat(timeTaken, 'E', -1, 64)
		msg[constants.TIME_TAKEN] = timeTakenString
		rsp.Body.Close()
	} else {
		log.Println("Received no length URL for Ping")
		msg[constants.RESPONSE_BODY] = "0"
		msg[constants.ERROR] = "0"
		msg[constants.RESPONSE_CODE] = "0"
	}
	//publish to subscriber
	Subscribe(constants.PostbackPingTopic, msg)
}

func Subscribe(topic string, msg map[string]string) {
	fmt.Println("Transaction Received on Subscriber:", topic, "-:-", msg[constants.ACTIVITY], msg[constants.TRANSACTION_ID], msg[constants.OFFER_ID], msg[constants.AFF_ID])
	result := brokers[topic].Publish(context.Background(), &pubsub.Message{Attributes: msg})
	serverID, err := result.Get(context.Background())
	if err != nil {
		fmt.Println("Error in Publishing to Google Pub/Sub", serverID, err, err.Error(), msg)
		msg[constants.ErrorMessage] = err.Error()
		//remove google_aid
		_, ok := msg["google_aid"]
		if ok {
			value := msg["google_aid"]
			delete(msg, "google_aid")
			msg["g_aid"] = value
		}
		//remove user-agent
		msg["User-Agent"] = ""
		SubscribeError(topic, msg)
	}
}

func SubscribeError(topic string, msg map[string]string) {
	fmt.Println("Transaction Received on Subscriber:-", msg[constants.TRANSACTION_ID], msg[constants.OFFER_ID], msg[constants.AFF_ID])
	result := brokers[topic].Publish(context.Background(), &pubsub.Message{Attributes: msg})
	serverID, err := result.Get(context.Background())
	if err != nil {
		fmt.Println("Error in Publishing to Google Pub/Sub", serverID, err, err.Error(), msg)
		msg[constants.ErrorMessage] = err.Error()
		dao.InsertToMongoBackup(constants.MongoDB, constants.ErrorTransaction, msg)
	}
}
func DuplicateMap(fwdMap map[string]string, request map[string]string) {
	for k := range request {
		fwdMap[k] = request[k]
	}
}

func ValidatePostbackForMedia(req *model.PostbackReq, rsp *model.PostbackRes) {

	//Check for offer affiliate mqf
	mqf := dao.GetMQFByOfferAndAffiliate(req.OfferID, req.AffiliateID)
	if len(mqf) == 0 {
		//check for global affiliate mqf
		mqf = dao.GetMQFByAffiliate(req.AffiliateID)
		if len(mqf) == 0 {
			mqf = "0.7"
		}
	}
	mqfFloat, _ := strconv.ParseFloat(mqf, 64)
	totalConversionCount, sentConversionCount := dao.GetConversionData(req.OfferID, req.AffiliateID)
	totalConversionFloat, _ := strconv.ParseFloat(totalConversionCount, 64)
	sentConversionFloat, _ := strconv.ParseFloat(sentConversionCount, 64)

	if totalConversionFloat == 0.0 {
		fmt.Println("Total Conversion is Zero")
		rsp.IsPingRequired = true
	} else if (sentConversionFloat / totalConversionFloat) <= mqfFloat {
		fmt.Println("sent.total", sentConversionFloat, totalConversionFloat)
		rsp.IsPingRequired = true
	} else {
		fmt.Println("Conversion MQF is less")
		rsp.IsPingRequired = false
	}

}

func ValidateRotatedPostbackForMedia(req *model.PostbackReq, rsp *model.PostbackRes) {

	mqf := 0.3
	totalConversionCount, sentConversionCount := dao.GetRotatedConversionData(req.OfferID)
	totalConversionFloat, _ := strconv.ParseFloat(totalConversionCount, 64)
	sentConversionFloat, _ := strconv.ParseFloat(sentConversionCount, 64)

	if (sentConversionFloat / totalConversionFloat) <= mqf {
		fmt.Println("sent.total", sentConversionFloat, totalConversionFloat)
		rsp.IsPingRequired = true
	} else if totalConversionFloat == 0.0 {
		fmt.Println("Total Conversion is Zero")
		rsp.IsPingRequired = true
	} else {
		fmt.Println("Conversion MQF is less")
		rsp.IsPingRequired = false
	}

}

func CreatePostback(req *model.PostbackReq, rsp *model.PostbackRes) {

	postbackURLTemplate := dao.GetOfferAffiliatePostbackTemplate(req.OfferID, req.AffiliateID)
	//Check offer affiliate specific template is available or not
	if len(postbackURLTemplate) == 0 || postbackURLTemplate == "" {
		postbackURLTemplate = dao.GetTemplateByAffiliateID(req.AffiliateID)
		//Check affiliate specific template is available or not
		if len(postbackURLTemplate) != 0 || postbackURLTemplate != "" {
			postbackURLTemplate = ReplaceTemplateParameters(rsp.ClickURL, postbackURLTemplate, req.TransactionID)
		}
	} else {
		postbackURLTemplate = ReplaceTemplateParameters(rsp.ClickURL, postbackURLTemplate, req.TransactionID)
	}

	rsp.URL = postbackURLTemplate

}

func ParsePostbackResponse(requestParams map[string]string, response *model.PostbackRes) {
	//check redirection is needed or not
	switch response.Activity {

	case 0:
		//process to fraud postback
		requestParams[constants.ACTIVITY] = "38"
		Subscribe(constants.FilteredTopic, requestParams)

	case 3:
		//process to sent conversions
		requestParams[constants.REDIRECT_URL] = response.URL
		requestParams[constants.ACTIVITY] = "3"
		Subscribe(constants.PostbackTopic, requestParams)
		Ping(response.URL, requestParams)

	case 4:
		//process this to unsent conversions
		requestParams[constants.ACTIVITY] = "4"
		Subscribe(constants.PostbackTopic, requestParams)

	case 5:
		//process to sent post Events
		requestParams[constants.REDIRECT_URL] = response.URL
		requestParams[constants.ACTIVITY] = "5"
		Subscribe(constants.PostbackTopic, requestParams)
		Ping(response.URL, requestParams)

	case 6:
		//process this to unsent post events
		requestParams[constants.ACTIVITY] = "6"
		Subscribe(constants.PostbackTopic, requestParams)

	case 7:
		//process to rotated sent conversion
		requestParams[constants.REDIRECT_URL] = response.URL
		requestParams[constants.ACTIVITY] = "7"
		Subscribe(constants.PostbackTopic, requestParams)
		Ping(response.URL, requestParams)

	case 8:
		//process to rotated un-sent conversion
		requestParams[constants.ACTIVITY] = "8"
		Subscribe(constants.PostbackTopic, requestParams)

	case 9:
		//process to rotated sent postevents
		requestParams[constants.REDIRECT_URL] = response.URL
		requestParams[constants.ACTIVITY] = "9"
		Subscribe(constants.PostbackTopic, requestParams)
		Ping(response.URL, requestParams)

	case 10:
		//process to rotated un-sent postevents
		requestParams[constants.ACTIVITY] = "10"
		Subscribe(constants.PostbackTopic, requestParams)

	case 11:
		//process to sent conversions without transactionID
		requestParams[constants.REDIRECT_URL] = response.URL
		requestParams[constants.ACTIVITY] = "11"
		Subscribe(constants.PostbackTopic, requestParams)
		Ping(response.URL, requestParams)

	case 12:
		//process this to unsent conversions without transactionID
		requestParams[constants.ACTIVITY] = "12"
		Subscribe(constants.PostbackTopic, requestParams)

	case 42:
		//process this to unsent conversions without transactionID & without media postback template
		requestParams[constants.ACTIVITY] = "13"
		Subscribe(constants.PostbackTopic, requestParams)

	default:
		//no actions taken on postback possible wrong offer types
		requestParams[constants.ACTIVITY] = "50"
		Subscribe(constants.FilteredTopic, requestParams)

	}

	if len(response.URL) == 0 && response.Activity != 0 && response.Activity != 12 && response.Activity != 10 &&
		response.Activity != 8 && response.Activity != 6 && response.Activity != 4 {
		//process to filter postback log for template error
		fwdMap := make(map[string]string)
		DuplicateMap(fwdMap, requestParams)
		fwdMap[constants.ACTIVITY] = "42"
		Subscribe(constants.FilteredTopic, fwdMap)
	}

}

func ParsePostbackReceived(requestParams map[string]string, req *api.Request) {

	//by default GET method
	if len(req.Method) == 0 {
		ParseGetRequest(requestParams, req)
	}
	//based on method extract parameters
	if strings.Compare(req.Method, "GET") == 0 {
		ParseGetRequest(requestParams, req)
	} else {

		//add body to logging
		requestParams[constants.Body] = req.Body

		//content header
		var contentType = requestParams["Content-Type"]

		//for POST switch with content-type
		if strings.HasPrefix(contentType, "multipart/form-data") {
			ParseMultiPartReq(contentType, requestParams, req)
		} else if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
			//x-wwww-urlencoded
			ParsePostRequest(requestParams, req)
		} else if strings.HasPrefix(contentType, "application/json") {
			//json body
			ParseJSONReq(requestParams, req)
		}

	}
}
