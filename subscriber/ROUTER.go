package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
	helper "github.com/adcamie/adserver/helpers/v1"
	logger "github.com/adcamie/adserver/logger"
	"github.com/rs/cors"
)

var client *pubsub.Client
var batchLimit *int
var shutDownFlag *bool
var duration int
var subscriptionProcesorPointer func(messages []*pubsub.Message)
var subscriber string

func Initialize(batchSize int, durationReceived int, subscriberReceived string, functionPointer func(messages []*pubsub.Message)) {

	//default values
	duration = durationReceived
	subscriptionProcesorPointer = functionPointer
	subscriber = subscriberReceived

	//variables for subscriber flexibility
	i := batchSize
	flag := false
	batchLimit = &i
	shutDownFlag = &flag

	var err error
	client, err = pubsub.NewClient(context.Background(), common.ProjectName)
	if err != nil {
		fmt.Println("Could not create pubsub client:", err)
		logger.ErrorLogger(err.Error(), "GooglePubSub", "Client Creation for Subscription Failed :"+subscriber)
	}
}

func InitRoutes(module string) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/subscriber/"+module, func(w http.ResponseWriter, r *http.Request) {
		//change flag on google pub/sub
		*shutDownFlag = true
		time.Sleep(5 * time.Second)
		fmt.Println("Shutdown Flag is set to true now:" + module)
	})
	mux.HandleFunc("/subscriber/"+module+"/restart", func(w http.ResponseWriter, r *http.Request) {
		//change flag on google pub/sub
		*shutDownFlag = false
		go ProcessingBatch()
		fmt.Println("Shutdown Flag is set to false now.....RESTARTING...:" + module)
	})

	mux.HandleFunc("/subscriber/"+module+"/shutdown", func(w http.ResponseWriter, r *http.Request) {
		//change flag on google pub/sub
		*shutDownFlag = true
		time.Sleep(25 * time.Second)
		//close sessions
		db.ShutdownMongo()
		db.ShutdownMongoBackup()
		db.ShutdownMongoSessionPool()

		fmt.Println("Shutdown Flag is set to true now:" + module)
	})

	mux.HandleFunc("/subscriber/"+module+"/change/batch-limit", func(w http.ResponseWriter, r *http.Request) {
		newBatchLimit, _ := strconv.Atoi(r.FormValue(common.BatchLimit))
		*batchLimit = *batchLimit + newBatchLimit
		fmt.Println("Changed Batch Limit:" + module)
		w.Write([]byte(strconv.Itoa(*batchLimit)))
	})
	mux.HandleFunc("/subscriber/"+module+"/change/duration", func(w http.ResponseWriter, r *http.Request) {
		newDuration, _ := strconv.Atoi(r.FormValue(common.Duration))
		duration = duration + newDuration
		fmt.Println("Changed Batch Limit:" + module)
		w.Write([]byte(strconv.Itoa(duration)))
	})
	mux.HandleFunc("/subscriber/"+module+"/batch-limit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Current Batch Limit:" + module)
		w.Write([]byte(strconv.Itoa(*batchLimit)))
	})
	mux.HandleFunc("/subscriber/"+module+"/db", func(w http.ResponseWriter, r *http.Request) {
		InitilaiseSessions()
		fmt.Println("Reconfigured db connections:" + module)
		w.Write([]byte(strconv.Itoa(*batchLimit)))
	})

	router := cors.Default().Handler(mux)
	return router
}

func ProcessingBatch() {
	for {
		messages := helper.PullTopicMessages(shutDownFlag, *batchLimit, subscriber, client)
		subscriptionProcesorPointer(messages)
		time.Sleep(time.Second * time.Duration(duration))
		if *shutDownFlag {
			fmt.Println("Exiting due to Shutdown Flag is set to true.")
			break
		}

	}
}

func InitilaiseSessions() {
	//elastic
	db.InitialiseBackupES()
	db.InitialiseMasterES()
	//redis
	db.InitializeRedisTranxn(100)
	db.InitializeRedisBackup(100)
	//mongo
	db.InitializeMongo()
	db.InitializeMongoBackup()
	db.InitializeMongoSessionPool()
}

func CloseSessions() {
	//close sessions
	db.ShutdownMongo()
	db.ShutdownMongoSessionPool()
	db.ShutdownMongoBackup()
}
