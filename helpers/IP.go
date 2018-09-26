package v1

import (
	"log"
	"math/rand"
	"time"

	"github.com/rdegges/go-ipify"
)

var freeGeoIP = []string{"5000", "5001", "5002"}
var freeGeoIPClick = []string{"6000", "6001", "6002"}

var localHost = "http://localhost:"
var JSON = "/json/"

func GetIP() string {
	ip, _ := ipify.GetIp()
	return ip
}

//freegeoip pooling
func GetFreeGeoIP(ip string) string {
	rand.Seed(time.Now().Unix())
	i := rand.Intn(2)
	url := localHost + freeGeoIP[i] + JSON + ip
	log.Println("freeGeoIP", url)
	return url
}
func GetFreeGeoIPClick(ip string) string {
	rand.Seed(time.Now().Unix())
	i := rand.Intn(2)
	url := localHost + freeGeoIPClick[i] + JSON + ip
	log.Println("freeGeoIPClick", url)
	return url
}
