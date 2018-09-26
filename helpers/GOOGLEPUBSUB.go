package v1

import (
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	logger "github.com/adcamie/adserver/logger"
	"golang.org/x/net/context"
)

func CreateSubscription(client *pubsub.Client, topic *pubsub.Topic, name string) {
	fmt.Println("Subscription creation started")
	ctx := context.Background()
	sub := client.Subscription(name)
	ok, err := sub.Exists(ctx)
	if err != nil {
		fmt.Println("Error while checking for existence of subscription :", err)
	}
	if ok {
		fmt.Println("Subscription already exists")
	} else {
		sub, err = client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 10 * time.Second,
		})
		if err != nil {
			fmt.Println("Error while subscription creation :", err)
		}
	}
	fmt.Println("Created subscription :", sub)
}

func PullTopicMessages(shutdownFlag *bool, limit int, name string, client *pubsub.Client) []*pubsub.Message {
	fmt.Println("Pulling topic messages started")
	ctx := context.Background()
	var mu sync.Mutex
	sub := client.Subscription(name)
	cntxt, cancel := context.WithCancel(ctx)
	messageReceived := 0
	var messages []*pubsub.Message
	err := sub.Receive(cntxt, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		if messageReceived >= limit || *shutdownFlag {
			fmt.Println("Inside cancelling Pub/Sub")
			messages = append(messages, msg)
			cancel()
			msg.Nack()
			return
		}
		msg.Ack()
		messages = append(messages, msg)
		messageReceived++
		fmt.Println("Msg-ID:-", msg.ID)
		//return
	})
	fmt.Println("Out and ready to get message")
	if err != nil {
		fmt.Println("Error while pulling topic messages :", err)
		//ping to exception listners
		go logger.ErrorLogger(err.Error(), "GooglePubSub", "Subscription Failed:"+name)
		return nil
	}
	return messages
}
