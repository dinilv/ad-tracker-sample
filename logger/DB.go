package v1

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/adcamie/adserver/common"
)

func ErrorLogger(details string, module string, event string) {
	hc := http.Client{}
	form := url.Values{}
	host, _ := os.Hostname()
	form.Add("host", host)
	form.Add("module", module)
	form.Add("details", details)
	form.Add("event", event)
	req, _ := http.NewRequest("POST", common.LoggingUrl, strings.NewReader(form.Encode()))
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	log.Println(resp, err)
}
