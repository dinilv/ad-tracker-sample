package main

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

type Transactions struct {
	TransactionID string    `bson:"transactionID,omitempty"`
	UTCDate       time.Time `bson:"utcdate,omitempty"`
	Date          time.Time `bson:"date,omitempty"`
}

func main() {

	files := []string{"nohup.out"}
	//transactionIDMap := map[string]bool{"gg1515999738750476995124117233536363521613": true}

	for _, fileName := range files {
		//file, err := os.Open("logs/" + fileName) /Users/dinilv/Desktop/micro

		file, err := os.Open("/Users/dinilv/Desktop/micro/" + fileName)

		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		log.Println("On start of scanner")
		scanner := bufio.NewScanner(file)
		//log.Println("error", scanner.Err(), len(scanner.Bytes()))
		const maxCapacity = 102400 * 1024
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)

		var i = 0

		for scanner.Scan() {

			line := scanner.Text()
			log.Println("on scanner", line)

			splitted := strings.Split(line, " ")
			//log.Println("splitted", splitted[0], splitted[6], splitted[8])
			log.Println(len(splitted), "length")
			if len(splitted) > 8 {
				log.Println(splitted, "SPLITTED")
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
						//log.Println("transactionID", params[constants.TRANSACTION_ID])
						//transactionID := params[constants.TRANSACTION_ID]

					}
					if err := scanner.Err(); err != nil {
						log.Println(err)
					}
					log.Println("number of lines:-", i)
					i++
				}
			}

		}
	}
}
