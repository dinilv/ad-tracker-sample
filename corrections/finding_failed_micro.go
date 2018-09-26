package main

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	"github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	"github.com/dghubble/sling"
	"gopkg.in/mgo.v2/bson"
)

type Transactions struct {
	TransactionID string    `bson:"transactionID,omitempty"`
	UTCDate       time.Time `bson:"utcdate,omitempty"`
	Date          time.Time `bson:"date,omitempty"`
}

func main() {
	log.Println("on start")
	config.InitializeMongo()
	config.InitializeRedisTranxn(100)
	today := time.Now().UTC()
	rounded := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()).UTC()
	//r, _ := regexp.Compile("[aA-zZ]{3}\\s[0-9]{2}\\s[aA-zZ]{3}\\s[0-9]{4}")
	files := []string{"nohup.out_1", "nohup.out_2", "nohup.out_3", "nohup.out_4", "nohup.out_5", "nohup.out_6", "nohup.out_7", "nohup.out_8", "nohup.out_9", "nohup.out_10", "nohup.out_11", "nohup.out_12", "nohup.out_13", "nohup.out_14", "nohup.out_15", "nohup.out_16", "nohup.out_17", "nohup.out_18", "nohup.out_19", "nohup.out_20", "nohup.out_21", "nohup.out_22", "nohup.out_23"}
	offerIDs := map[string]bool{"1797": true, "1756": true}
	transactionIDMap := map[string]bool{}

	for _, fileName := range files {
		//file, err := os.Open("logs/" + fileName) /Users/dinilv/Desktop/micro
		file, err := os.Open("logs")

		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		log.Println("on start of scanner")
		scanner := bufio.NewScanner(file)
		//log.Println("error", scanner.Err(), len(scanner.Bytes()))
		const maxCapacity = 102400 * 1024
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)

		var i = 0

		for scanner.Scan() {

			line := scanner.Text()
			//log.Println("on scanner", line)

			splitted := strings.Split(line, " ")
			//log.Println("splitted", splitted[0], splitted[6], splitted[8])
			log.Println(len(splitted), "length")
			if len(splitted) > 8 {
				if strings.Compare(splitted[8], "500") == 0 {
					log.Println(splitted)
				}

				//check url is not null
				if len(splitted[6]) > 0 {
					receivedUrl := "http://track.adcamie.com" + splitted[6] + "&ip=" + splitted[0]
					tempUrl, _ := url.Parse(receivedUrl)
					log.Println("Temp URL", tempUrl, splitted[6])
					if tempUrl != nil {
						params := make(map[string]string)

						for key, value := range tempUrl.Query() {
							params[key] = value[0]
						}
						log.Println("transactionID", params[constants.TRANSACTION_ID])
						transactionID := params[constants.TRANSACTION_ID]
						//validate transactionID
						if dao.ValidateTransactionID(transactionID) {
							if !dao.ValidateConvertedTransactionID(transactionID) {
								//check offerIDs of postbacks
								offerID, _ := dao.GetTransaction(transactionID)
								if !transactionIDMap[transactionID] {
									transactionIDMap[transactionID] = true
									sucess := bson.M{}
									_, err := sling.New().Get(receivedUrl).ReceiveSuccess(sucess)
									log.Println("Response", receivedUrl, sucess)
									if err != nil {
										log.Println("Error:", err.Error())
									}
								} else {
									log.Println("Duplicate transactionID within retry")
								}

							} else {
								log.Println("Converted Transaction")
							}
						} else {
							log.Println("Not a Valid transactionID")
							filtersForDelete := map[string]interface{}{"transactionID": transactionID}
							dao.DeleteFromMongo(constants.MongoDB, "FailedTransactions", filtersForDelete)
							dao.InsertToMongo(constants.MongoDB, "FailedTransactions", &Transactions{
								TransactionID: transactionID,
								UTCDate:       time.Now().UTC(),
								Date:          rounded})
						}

					}
				}
			}
			if err := scanner.Err(); err != nil {
				log.Println(err)
			}
			log.Println("number of lines:-", i)
			i++
		}
	}
}
