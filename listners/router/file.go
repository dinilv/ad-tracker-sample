package router

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	constants "github.com/adcamie/adserver/common/v1"
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
	helper "github.com/adcamie/adserver/helpers/v1"
	"github.com/dghubble/sling"
	"github.com/olivere/elastic"
	"gopkg.in/mgo.v2/bson"
)

func ReportCSVDownload(w http.ResponseWriter, r *http.Request) {

	log.Print("Listener got the request to download report :")

	ts := strconv.Itoa(time.Now().Nanosecond())
	filename := "TrackerLog_" + ts + ".csv"
	absPath, _ := filepath.Abs(constants.Router_Path + filename)
	_, err := os.Create(absPath)
	log.Println(err)
	//redirect url with generated filename
	url := "http://localhost:8000/v1/log/csv?filename=" + filename + "&"

	//format received url
	receivedURL := strings.Replace(r.URL.String(), "/report/csv/download?", "", -1)
	receivedURL = strings.Replace(receivedURL, "/report/csv/download/?", "", -1)
	receivedURL = strings.Replace(receivedURL, "%2F", "/", -1)
	receivedURL = strings.Replace(receivedURL, "%2C", ",", -1)

	redirectedURL := url + receivedURL
	fmt.Println(redirectedURL)
	sucess := bson.M{}
	rsp, error := sling.New().Get(redirectedURL).ReceiveSuccess(sucess)
	fmt.Println(rsp, error)

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Set("Content-type", "application/csv")
	w.Header().Set("Access-Control-Expose-Headers", "File-Name")
	w.Header().Set("File-Name", filename)

	filebytes, err := ioutil.ReadFile(absPath)

	if err != nil {
		fmt.Println(err)

	}

	b := bytes.NewBuffer(filebytes)

	if _, err := b.WriteTo(w); err != nil { // <----- here!
		fmt.Fprintf(w, "%s", err)
	}

	w.Write(b.Bytes())

}

func UploadCSVBlocked(w http.ResponseWriter, r *http.Request) {

	log.Print("Listener got the request to upload blocked csv report :")
	log.Print("Request", r)
	log.Print("Request Data:")
	//parse req to file
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("File error:", err)
		return
	}

	operator := r.FormValue("operator")
	blocked := r.FormValue("block")
	overwrite, _ := strconv.ParseBool(r.FormValue("overwrite"))
	if overwrite {
		//delete other records, which are already blocked
		query := elastic.NewTermQuery("blocked", "true")
		dao.DeleteFromES(constants.MOTracker, constants.MSISDN, query)
	}
	//parse file to csv
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	csvData, _ := reader.ReadAll()

	bulk := db.ESMasterClient.Bulk().Index(constants.MOTracker).Type(constants.MSISDN)
	i := 0
	script := elastic.NewScript("ctx._source.blocked=params.blocked;ctx._source.operator=params.operator;").Params(map[string]interface{}{"blocked": blocked, "operator": operator}).Lang("painless")
	for _, each := range csvData {
		i = i + 1

		if i > 10000 {
			commitToES(bulk)
			bulk = db.ESMasterClient.Bulk().Index(constants.MOTracker).Type(constants.MSISDN)
			i = 0
		}
		log.Println("MSISDN", each[0])
		//parse each row
		obj := &model.MSISDNDetails{
			MSISDN:   each[0],
			Blocked:  blocked,
			Operator: operator,
			OfferIDs: []string{"000"},
		}

		bulk.Add(elastic.NewBulkUpdateRequest().Id(each[0]).Script(script).ScriptedUpsert(false).Upsert(obj))
	}

	log.Println("Outside loop", bulk.NumberOfActions())
	defer file.Close()
	commitToES(bulk)

}

func UploadCSVSubscribed(w http.ResponseWriter, r *http.Request) {

	log.Print("Listener got the request to upload subscription csv report for different services:", r)
	//get offer details for serviceID mapping
	offers := dao.QueryAllOffers()
	offerRefMap := make(map[string]string)
	for _, v := range offers {
		offerRefMap[v.ServiceID] = v.OfferID
	}

	//parse req to file
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}

	operator := r.FormValue("operator")
	overwrite, _ := strconv.ParseBool(r.FormValue("overwrite"))
	log.Println(r.FormValue("overwrite"), "overwrite")

	//parse file to csv
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	csvData, _ := reader.ReadAll()

	bulk := db.ESMasterClient.Bulk().Index(constants.MOTracker).Type(constants.MSISDN)
	i := 0
	if overwrite {
		//delete uploaded existing msisdns from existing DB
		for _, each := range csvData {
			i = i + 1

			if i > 10000 {
				commitToES(bulk)
				bulk = db.ESMasterClient.Bulk().Index(constants.MOTracker).Type(constants.MSISDN)
				i = 0
				break
			}
			bulk.Add(elastic.NewBulkDeleteRequest().Id(each[1]))
		}
		commitToES(bulk)
	}
	i = 0
	for _, each := range csvData {
		i = i + 1

		if i > 10000 {
			commitToES(bulk)
			bulk = db.ESMasterClient.Bulk().Index(constants.MOTracker).Type(constants.MSISDN)
			i = 0
			break
		}
		log.Println("MSISDN", each[1], "ServiceID", each[0])
		//parse each row
		obj := &model.MSISDNDetails{
			MSISDN:    each[1],
			Blocked:   "false",
			OfferIDs:  []string{offerRefMap[each[0]]},
			Operator:  operator,
			UpdatedAt: time.Now().UTC(),
		}

		script := elastic.NewScript("ctx._source.offerIDs.contains(params.offerid) ? (ctx.op = none) : ctx._source.offerIDs.add( params.count);ctx._source.blocked=params.blocked;ctx._source.updatedAt").
			Params(map[string]interface{}{"offerid": offerRefMap[each[0]], "blocked": "false", "updatedAt": time.Now().UTC()}).Lang("painless")
		bulk.Add(elastic.NewBulkUpdateRequest().Id(each[1]).Script(script).ScriptedUpsert(false).Upsert(obj))
	}

	defer file.Close()
	commitToES(bulk)

}

func UploadCSVIP(w http.ResponseWriter, r *http.Request) {

	log.Print("Listener got the request to upload ips csv :")

	//parse req to file
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}

	//parse file to csv
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	csvData, _ := reader.ReadAll()

	//Location and geography details
	geo := new(model.GeoDetails)
	ts := strconv.Itoa(time.Now().Nanosecond())
	filename := "Geo_" + ts + ".csv"
	helper.InitializeWriters(filename, "", []string{"IP", "Country"})
	for _, each := range csvData {
		sling.New().Get(helper.GetFreeGeoIP(ip)).ReceiveSuccess(geo)
		//write to csv
		helper.SaveRow(each[0] + "," + geo.CountryName)
	}
	file.Close()
	helper.CloseWriters()

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Set("Content-type", "application/csv")
	w.Header().Set("Access-Control-Expose-Headers", "File-Name")
	w.Header().Set("File-Name", filename)

	filebytes, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)

	}

	b := bytes.NewBuffer(filebytes)

	if _, err := b.WriteTo(w); err != nil { // <----- here!
		fmt.Fprintf(w, "%s", err)
	}

	w.Write(b.Bytes())

}

func UploadCSVTransactions(w http.ResponseWriter, r *http.Request) {

	log.Print("Listener got the request to upload Transactions csv report for different services:", r)

	//parse req to file
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}

	//parse file to csv
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	csvData, _ := reader.ReadAll()

	for _, each := range csvData {
		//link := "http://track2.adcamie.com/aff_lsr?transaction_id=" + each[1]
		link := "localhost:8000/v1/postback/track?transaction_id=" + each[1]
		ping(link)
	}

	defer file.Close()

}

func commitToES(bulk *elastic.BulkService) {
	log.Println("Before bulk")
	res, err := bulk.Do(context.Background())
	if err != nil {
		log.Println(err, bulk.NumberOfActions())
	}
	if res.Errors {
		// Look up the failed documents with res.Failed(), and e.g. recommit
		log.Println(res.Errors)
	}
	log.Println("After  Bulk")
}

func ping(url string) {
	log.Println("Pinging started")
	startTime := time.Now()
	sucess := bson.M{}
	rsp, error := sling.New().Get(url).ReceiveSuccess(sucess)
	fmt.Println(rsp, error)
	endTime := time.Now()
	log.Println("Time Taken For Postback Ping:-", endTime.Sub(startTime).Seconds())
}
