package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	common "github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
	handler "github.com/adcamie/adserver/handlers/v1/tracker"
	logger "github.com/adcamie/adserver/logger"
	router "github.com/adcamie/adserver/subscriber"
	"github.com/dghubble/sling"
	"github.com/micro/go-micro/client"
	"gopkg.in/mgo.v2/bson"
)

var reg *regexp.Regexp
var topics []string
var brokers = map[string]*pubsub.Topic{}

func init() {

	reg = regexp.MustCompile("gg15\\d+")
	topics = []string{common.PostbackTopic, common.PostbackPingTopic, common.DelayedPostbackTopic, common.FilteredTopic}

	for _, topic := range topics {
		client, err := pubsub.NewClient(context.Background(), common.ProjectName)
		if err != nil {
			logger.ErrorLogger(err.Error(), "GooglePubSub", "Topic: "+topic+" Initialization Error")
			panic(err.Error())
		}
		topicClient := client.Topic(topic)
		brokers[topic] = topicClient
	}
}
func main() {

	db.InitializeRedisTranxn(100)
	db.InitializeRedisBackup(100)
	db.InitializeMongoBackup()
	db.InitializeMongoSessionPool()
	db.InitialiseBackupES()

	//initialize subscriber options
	router.Initialize(25, 240, common.DelayedPostbackSub, InsertBatch)

	//schedule for 5 secs loop
	go router.ProcessingBatch()

	//create a listner for shutdown
	apiRouter := router.InitRoutes("delayed-sub")
	server := &http.Server{
		Addr:    "0.0.0.0:" + os.Args[1],
		Handler: apiRouter,
	}

	fmt.Println("Listening...")
	err := server.ListenAndServe()
	if err != nil {
		go logger.ErrorLogger(err.Error(), "DelayedPostback", "Server Creation Failed:")
	}
}

func InsertBatch(messages []*pubsub.Message) {

	transactionIDs := map[string]bool{}

	for _, msg := range messages {

		fmt.Println("TransactionID:", msg.Attributes[common.TRANSACTION_ID])
		//add messageID as parameter in request map for checking duplicate in mongo
		msg.Attributes[common.MessageID] = msg.ID
		msg.Attributes[common.SubscriptionTime] = time.Now().UTC().String()

		//check this transactionID is already processed or not within the batch
		_, ok := transactionIDs[msg.Attributes[common.TRANSACTION_ID]]
		retryValid := dao.ValidateRetryTransactionIDOnBackup(msg.Attributes[common.TRANSACTION_ID])
		delayValid := dao.ValidateDelayedTransactionIDOnBackup(msg.Attributes[common.TRANSACTION_ID])
		if !ok && !retryValid && !delayValid {
			//add to trasnactionIDs
			transactionIDs[msg.Attributes[common.TRANSACTION_ID]] = true
			//for mongo backup details
			backupAffiliateID := "999"
			backupOfferID := "1"
			offerID := "1"
			affiliateID := "999"
			isConverted := dao.ValidateConvertedTransactionID(msg.Attributes[common.TRANSACTION_ID])
			offerType := "0"

			//add messageID as parameter in request map for checking duplicate in mongo
			msg.Attributes[common.MessageID] = msg.ID
			msg.Attributes[common.SubscriptionTime] = time.Now().UTC().String()
			msg.Attributes[common.Processor] = common.DELAYED_SUB
			transactionID := msg.Attributes[common.TRANSACTION_ID]
			//validate transaction_id exists or not
			valid := dao.ValidateTransactionID(transactionID)
			backupValid := false
			//check on transaction90 days old
			if !valid {
				var redisBackupModel model.RedisTransactionBackup
				valid, redisBackupModel = dao.SearchRedisKeysFromESBackup(transactionID)
				if valid {
					backupAffiliateID = redisBackupModel.AffiliateID
					backupOfferID = redisBackupModel.OfferID
					backupValid = true
				}
			}
			if !valid {

				fmt.Println("TransactionID is Invalid in postback even after retry:-", transactionID)

				//retry again after increasing count
				retryCount, _ := strconv.Atoi(msg.Attributes[common.RetryCount])
				updatedRetryCount := retryCount + 1

				//check its already present in fitered postback log for more than 5 times
				filters := map[string]interface{}{common.TransactionID: transactionID}
				filteredCount := dao.GetCountFromMongoSession(common.MongoDB, common.FilteredPostBackLog, filters)

				if updatedRetryCount < 4 && filteredCount < 10 && reg.MatchString(transactionID) {
					msg.Attributes[common.RetryCount] = strconv.Itoa(updatedRetryCount)
					msg.Attributes[common.FilteredCount] = strconv.Itoa(filteredCount)
					msg.Attributes[common.Activity] = "51"
					subscribe(common.FilteredTopic, msg.Attributes)
					subscribe(common.DelayedPostbackTopic, msg.Attributes)
				} else {
					fmt.Println("No more retry. Limit is over.")
					msg.Attributes[common.Activity] = "52"
					subscribe(common.FilteredTopic, msg.Attributes)
				}

			} else {

				postbackProcessor := true
				//split for load-balancer adding ip
				ips := strings.Split(msg.Attributes["X-Forwarded-For"], ",")

				//handle backup data diferently
				if backupValid {
					valid := dao.ValidateAdveriserWithOfferID(backupOfferID, ips[0])
					offerType = dao.GetOfferTypeOnTranxn(backupOfferID)
					if !valid {
						//log to filtered postbacks
						msg.Attributes[common.Activity] = "39"
						subscribe(common.FilteredTopic, msg.Attributes)
						postbackProcessor = false
					}
				} else {
					//check received_ip is whitelisted on advertiser or not
					valid, offerID := dao.ValidateAdveriserIPWithTransaction(transactionID, ips[0])
					if !valid {
						//log to filtered postbacks
						msg.Attributes[common.Activity] = "39"
						subscribe(common.FilteredTopic, msg.Attributes)
						postbackProcessor = false
					} else {
						//check transaction id exists, checking for duplicate transaction or postevents
						isConverted, offerType, affiliateID = dao.ValidateTransactionIDForPostback(transactionID, offerID)
						//offer type without transactionID
						if strings.Compare(offerType, "8") == 0 {
							fmt.Println("Wrong Offer type received: Invalid postback.", offerType, msg.Attributes[common.OFFER_ID])
							//log to filtered postbacks
							msg.Attributes[common.Activity] = "43"
							subscribe(common.FilteredTopic, msg.Attributes)
							postbackProcessor = false
						}
					}
				}

				//offer type without postevents and already converted
				if isConverted && (strings.Compare(offerType, "3") == 0 || strings.Compare(offerType, "7") == 0) {
					fmt.Println("TransactionId Exists in converted postbacks and postevents not enabled")
					//log to filtered postbacks
					msg.Attributes[common.Activity] = "37"
					subscribe(common.FilteredTopic, msg.Attributes)
					postbackProcessor = false

				}

				if postbackProcessor {
					//successfully processed for logging
					dao.SaveDelayedPostbackTransaction(msg.Attributes[common.TRANSACTION_ID])

					setpostback := new(handler.PostbackReq)
					setpostback.TransactionId = transactionID
					if backupValid {
						setpostback.AffiliateID = backupAffiliateID
						setpostback.OfferID = backupOfferID
						//for postback ping log
						msg.Attributes[common.OFFER_ID] = backupOfferID
						msg.Attributes[common.AFF_ID] = backupAffiliateID
					} else {
						setpostback.AffiliateID = affiliateID
						setpostback.OfferID = offerID
						//for postback ping log
						msg.Attributes[common.OFFER_ID] = offerID
						msg.Attributes[common.AFF_ID] = affiliateID
					}
					setpostback.OfferType = offerType
					setpostback.IsConverted = isConverted
					setpostback.GoalID = msg.Attributes["goal_id"]
					setpostback.ConversionIP = ips[0]

					response := &handler.PostbackRes{}
					request := client.NewJsonRequest("go.micro.service.v1.postback", "Postback.Setpostback", setpostback)
					if err := client.Call(context.Background(), request, response); err != nil {
						fmt.Println("Client Calling Error In Tracker Set Postback:", err, request, response)
					}

					fmt.Println(response, "response")
					//check redirection is needed or not
					switch response.Activity {

					case 0:
						//process to fraud postback
						msg.Attributes[common.Activity] = "38"
						subscribe(common.FilteredTopic, msg.Attributes)

					case 3:
						//process to sent conversions
						msg.Attributes[common.REDIRECT_URL] = response.Url
						msg.Attributes[common.Activity] = "3"
						subscribe(common.PostbackTopic, msg.Attributes)
						ping(response.Url, msg.Attributes)

					case 4:
						//process this to unsent conversions
						msg.Attributes[common.Activity] = "4"
						subscribe(common.PostbackTopic, msg.Attributes)

					case 5:
						//process to sent post Events
						msg.Attributes[common.REDIRECT_URL] = response.Url
						msg.Attributes[common.Activity] = "5"
						subscribe(common.PostbackTopic, msg.Attributes)
						ping(response.Url, msg.Attributes)

					case 6:
						//process this to unsent post events
						msg.Attributes[common.Activity] = "6"
						subscribe(common.PostbackTopic, msg.Attributes)

					case 7:
						//process to rotated sent conversion
						msg.Attributes[common.REDIRECT_URL] = response.Url
						msg.Attributes[common.Activity] = "7"
						subscribe(common.PostbackTopic, msg.Attributes)
						ping(response.Url, msg.Attributes)

					case 8:
						//process to rotated un-sent conversion
						msg.Attributes[common.Activity] = "8"
						subscribe(common.PostbackTopic, msg.Attributes)

					case 9:
						//process to rotated sent postevents
						msg.Attributes[common.REDIRECT_URL] = response.Url
						msg.Attributes[common.Activity] = "9"
						subscribe(common.PostbackTopic, msg.Attributes)
						ping(response.Url, msg.Attributes)

					case 10:
						//process to rotated un-sent postevents
						msg.Attributes[common.Activity] = "10"
						subscribe(common.PostbackTopic, msg.Attributes)

					case 11:
						//process to sent conversions without transactionID
						msg.Attributes[common.REDIRECT_URL] = response.Url
						msg.Attributes[common.Activity] = "11"
						subscribe(common.PostbackTopic, msg.Attributes)
						ping(response.Url, msg.Attributes)

					case 12:
						//process this to unsent conversions without transactionID
						msg.Attributes[common.Activity] = "12"
						subscribe(common.PostbackTopic, msg.Attributes)

					case 42:
						//process this to unsent conversions without transactionID & without media postback template
						msg.Attributes[common.Activity] = "13"
						subscribe(common.PostbackTopic, msg.Attributes)

					default:
						//no actions taken on postback possible wrong offer types
						msg.Attributes[common.Activity] = "50"
						subscribe(common.FilteredTopic, msg.Attributes)

					}

					if len(response.Url) == 0 && response.Activity != 0 && response.Activity != 12 && response.Activity != 10 &&
						response.Activity != 8 && response.Activity != 6 && response.Activity != 4 {
						//process to filter postback log for template error
						fwdMap := make(map[string]string)
						duplicateMap(fwdMap, msg.Attributes)
						fwdMap[common.Activity] = "42"
						subscribe(common.FilteredTopic, msg.Attributes)
					}
				}
			}
		} else {
			if ok {
				fmt.Println("No need to retry. Already tried.")
				msg.Attributes[common.Activity] = "54"
				subscribe(common.FilteredTopic, msg.Attributes)
			} else if retryValid {
				fmt.Println("No need to retry. Already converted by retry.")
				msg.Attributes[common.Activity] = "55"
				subscribe(common.FilteredTopic, msg.Attributes)
			} else if delayValid {
				fmt.Println("No need to retry. Already converted by delayed.")
				msg.Attributes[common.Activity] = "56"
				subscribe(common.FilteredTopic, msg.Attributes)
			}

		}
	}
}

func ping(url string, msg map[string]string) {

	fmt.Println("Pinging started")
	if len(url) > 0 {
		startTime := time.Now()
		sucess := bson.M{}
		rsp, err := sling.New().Get(url).ReceiveSuccess(sucess)
		endTime := time.Now()
		timeTaken := endTime.Sub(startTime).Seconds()
		fmt.Println("Time Taken For Postback Ping:-", timeTaken)
		//add for Logging
		var rspBytes []byte
		rsp.Body.Read(rspBytes)
		msg[common.RESPONSE_BODY] = string(rspBytes)
		msg[common.RESPONSE_CODE] = strconv.Itoa(rsp.StatusCode)
		if err != nil {
			msg[common.ERROR] = err.Error()
		}
		timeTakenString := strconv.FormatFloat(timeTaken, 'E', -1, 64)
		msg[common.TIME_TAKEN] = timeTakenString
		rsp.Body.Close()
	} else {
		fmt.Println("Received no length URL for Ping")
		msg[common.RESPONSE_BODY] = "0"
		msg[common.ERROR] = "0"
		msg[common.RESPONSE_CODE] = "0"
	}
	//publish for subscriber
	subscribe(common.PostbackPingTopic, msg)
}

func duplicateMap(fwdMap map[string]string, request map[string]string) {
	for k := range request {
		fwdMap[k] = request[k]
	}
}

func subscribe(topic string, msg map[string]string) {
	fmt.Println("Broker Received Message :"+topic, time.Now().UTC())
	fmt.Println("Transaction Received on Subscriber:-", msg[common.TRANSACTION_ID], msg[common.OFFER_ID], msg[common.AFF_ID])
	result := brokers[topic].Publish(context.Background(), &pubsub.Message{Attributes: msg})
	serverID, err := result.Get(context.Background())
	if err != nil {
		fmt.Println("Error in Publishing to Google Pub/Sub", serverID, err, err.Error(), msg)
		msg[common.ErrorMessage] = err.Error()
		dao.InsertToMongoBackup(common.MongoDB, common.ErrorTransaction, msg)
	}
	fmt.Println("Broker End of Received Message :"+topic, time.Now().UTC())
}
