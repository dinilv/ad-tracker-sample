package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	constants "github.com/adcamie/adserver/common/v1"
	logger "github.com/adcamie/adserver/logger"
)

var topics []string
var host = "http://track.adcamie.com"

var brokers = map[string]*pubsub.Topic{}

func init() {

	topics = []string{constants.ImpressionTopic, constants.ClickTopic, constants.RotatedTopic, constants.PostbackTopic, constants.PostbackPingTopic, constants.DelayedPostbackTopic, constants.FilteredTopic}

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

func subscribe(topic string, msg map[string]string) {
	fmt.Println("Transaction Received on Subscriber:-", msg[constants.TRANSACTION_ID], msg[constants.OFFER_ID], msg[constants.AFF_ID])
	result := brokers[topic].Publish(context.Background(), &pubsub.Message{Attributes: msg})
	serverID, err := result.Get(context.Background())
	if err != nil {
		fmt.Println("Error in Publishing to Google Pub/Sub", serverID, err)
		fmt.Println(msg)
	}
}

func main() {
	var message = map[string]string{"api_time": "2018-02-15 02:35:05.843472121 +0000 UTC", "offer_id": "1919", "url": "", "X-Forwarded-For": "13.113.52.3, 35.201.111.41", "aff_id": "151382", "aff_sub3": "7100_1447_01152d1a", "aff_sub": "b2ZmZXJpZD00NDM2OSZ1c2VyaWQ9NzEwMCZjbGlja3RpbWU9MjAxODAyMTUwMjM1MDUmY2hhbm5lbD0xNDQ3XzAxMTUyZDFhJmdhaWQ9JmFuZGlkPXthbmRpZH0maWRmYT0mYWZmX3N1Yj0yMjBhNTNmMy02MzgyLTQ4YTctYjZjZC04OTgyZjgzNzlmODRfX3BzcG0mc3ViMT17c3ViMX0mc3ViMj17c3ViMn0mZ2VvPVVT",
		"activity": "2", "session_ip": "13.113.52.3", "adsauce_id": "uid_1518662105843439267", "click_red_url": "https://cld-m.tlnk.io/serve?site_id=94040&agency_id=1371&ref_id=gg15186621051919151382843331311&sub_site=7100_1447_01152d1a&sub_campaign=1919&action=click&publisher_id=359098", "Host": "track.adcamie.com", "transaction_id": "gg15186621051919151382843331311", "click_url": "http://track.adcamie.com/aff_c?aff_sub=b2ZmZXJpZD00NDM2OSZ1c2VyaWQ9NzEwMCZjbGlja3RpbWU9MjAxODAyMTUwMjM1MDUmY2hhbm5lbD0xNDQ3XzAxMTUyZDFhJmdhaWQ9JmFuZGlkPXthbmRpZH0maWRmYT0mYWZmX3N1Yj0yMjBhNTNmMy02MzgyLTQ4YTctYjZjZC04OTgyZjgzNzlmODRfX3BzcG0mc3ViMT17c3ViMX0mc3ViMj17c3ViMn0mZ2VvPVVT&offer_id=1919&aff_id=151382&aff_sub3=7100_1447_01152d1a", "User-Agent": "ôS 2.1.0 rv:2.1.0.48547 (iPhone; iOS 10.3.3; zh-Hans_JP)", "Content-Type": "application/x-www-form-urlencoded", "time_taken": "0.000717"}

	fmt.Println(message)
	//subscribe(constants.ClickTopic, message)

}
