package v1

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var smsHttpClient http.Client
var messages map[int]string

func init() {
	smsHttpClient = http.Client{}
	messages = map[int]string{0: "Hi Adcamiens, Redis server is down . Please monitor its health.",
		1: "Hi Adcamiens, Mongo Server is down. Please Take action as required.",
		2: "Hi Adcamiens, Tracker Click Server is down. Please Take action as required.",
		3: "Hi Adcamiens, Tracker Postback Server is down. Please Take action as required.",
		4: "Hi Adcamiens, No conversion for the last hour today. Check adcamie-tracker-app.",
		5: "Hi Adcamiens, Subscription To API time is more than 5 Minutes. Increase subscriber Batch.",
		6: "Hi Adcamiens, Tracker Mongo Backup Server is down. Please Take action as required.",
		7: "Hi Adcamiens, Tracker Redis Backup Server is down . Please monitor its health.",
		8: "Hi Adcamiens, Tracker Redis Transaction Server is down. Please Take action as required.",
		9: "Hi Adcamiens, Tracker Elastic Search Server is down. Please Take action as required.",
	}
}

type SMS struct {
	mobile    map[string]bool
	messageID int
}

func NewSMS(id int, mob map[string]bool) *SMS {
	return &SMS{
		mobile:    mob,
		messageID: id,
	}
}

func (s *SMS) Send() {

	//smsGatewayURL, _ := url.Parse("http://login.bulksmsgateway.in/sendmessage.php")
	smsGatewayURL, _ := url.Parse("http://cloud.smsindiahub.in/vendorsms/pushsms.aspx?user=kunalag&password=kunvib@129&")

	var params = url.Values{}
	params.Add("sid", "SCAMIE")
	params.Add("msg", messages[s.messageID])
	params.Add("fl", "0")
	params.Add("gwid", "2")

	for mob := range s.mobile {

		log.Println(mob)
		params.Add("msisdn", mob)
		req, _ := http.NewRequest("GET", smsGatewayURL.String()+params.Encode(), nil)
		log.Println(req)
		res, _ := smsHttpClient.Do(req)
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(res)
		fmt.Println(string(body))
		defer res.Body.Close()
	}

}
