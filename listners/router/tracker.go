package router

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	helper "github.com/adcamie/adserver/helpers/v1"
)

var ip string
var version string
var port string
var versionURL = "http://35.194.251.197:8080/deployment/version?server_type=tracker-app"
var separator = ":"
var HTTP = "http://"
var impSubURL = "/impression/track?"
var clickSubURL = "/click/track?"
var convSubURL = "/postback/track?"
var portMap = map[string]string{"7000": "8000", "7001": "8001", "7002": "8002",
	"7003": "8003", "7004": "8004", "7005": "8005", "7006": "8006", "7007": "8007",
	"7008": "8008", "7009": "8009", "7010": "8010", "7011": "8011", "7012": "8012",
	"7013": "8013", "7014": "8014", "7015": "8015", "7016": "8016", "7017": "8017",
	"7018": "8018", "7019": "8019"}

func ChangePort(p string) {
	port = portMap[p]
}
func init() {

	//get version latest
	client := &http.Client{}
	req, _ := http.NewRequest("GET", versionURL, nil)
	resp, _ := client.Do(req)
	versionBytes, _ := ioutil.ReadAll(resp.Body)
	version = "/" + "v" + string(versionBytes)
	defer resp.Body.Close()
	//get ip of existing system
	ip = helper.GetIP()
	//ip = "0.0.0.0"
}

func ImpressionAPI(w http.ResponseWriter, r *http.Request) {
	//format received url
	log.Println(r.URL.String(), "url")
	receivedURL := strings.Replace(r.URL.String(), "/aff_i?", "", -1)
	receivedURL = strings.Replace(receivedURL, "/aff_i/?", "", -1)
	redirectedURL := HTTP + ip + separator + port + version + impSubURL + receivedURL
	http.Redirect(w, r, redirectedURL, 307)
	log.Println("Impression API", redirectedURL)
}

func ClickAPI(w http.ResponseWriter, r *http.Request) {
	//format received url
	receivedURL := strings.Replace(r.URL.String(), "/aff_c?", "", -1)
	receivedURL = strings.Replace(receivedURL, "/aff_c/?", "", -1)
	redirectedURL := HTTP + ip + separator + port + version + clickSubURL + receivedURL
	http.Redirect(w, r, redirectedURL, 307)
	log.Println("Click API", redirectedURL)
}

func PostbackAPI(w http.ResponseWriter, r *http.Request) {
	//format received url
	receivedURL := strings.Replace(r.URL.String(), "/aff_lsr?", "", -1)
	receivedURL = strings.Replace(receivedURL, "/aff_lsr/?", "", -1)
	redirectedURL := HTTP + ip + separator + port + version + convSubURL + receivedURL
	http.Redirect(w, r, redirectedURL, 307)
	log.Println("Postback API", redirectedURL)
}
