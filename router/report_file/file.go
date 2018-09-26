package main

import (
	"log"
	"net/http"

	listner "github.com/adcamie/adserver/listners/v1/router"
	"github.com/rs/cors"
)

func main() {
	startServer()
}

func startServer() {

	router := InitRoutes()
	server := &http.Server{
		Addr:    "0.0.0.0:8081",
		Handler: router,
	}
	log.Println("Listening...")
	server.ListenAndServe()
	log.Println("not Listening........")

}

func InitRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/report/csv/download", listner.ReportCSVDownload)
	mux.HandleFunc("/track/mocampaign/blocked_msisdn", listner.ReportCSVDownload)
	mux.HandleFunc("/track/mocampaign/subscribed_msisdn", listner.UploadCSVSubscribed)
	mux.HandleFunc("/report/csv/ip", listner.UploadCSVIP)
	//not decided on flow
	mux.HandleFunc("/track/mocampaign/update/service_id", listner.ReportCSVDownload)
	mux.HandleFunc("/track/mocampaign/upload/blocked_msisdn", listner.UploadCSVBlocked)
	mux.HandleFunc("/track/mocampaign/upload/subscribed_msisdn", listner.UploadCSVSubscribed)
	mux.HandleFunc("/track/upload/transaction", listner.UploadCSVTransactions)

	router := cors.Default().Handler(mux)

	return router
}
