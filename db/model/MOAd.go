package v1

import (
	"time"
)

// Mobile Ad is a structure used for serializing/deserializing ad-request data in Elasticsearch.
type MOAd struct {

	//Campaign
	CampaignID   int     `json:"campaignID,omitempty"`
	CampaignName string  `json:"campaignName,omitempty"`
	Bid          float64 `json:"bid,omitempty"`
	CategoryType int     `json:"categoryType,omitempty"`
	Priority     int     `json:"priority,omitempty"`
	Status       int     `json:"status,omitempty"`
	BidType      int     `json:"bidType,omitempty"`
	//Creative
	CreativeID   int    `json:"creativeID,omitempty"`
	CreativeLink string `json:"creativeLink,omitempty"`
	CreativeType int    `json:"creativeType,omitempty"`
	ClickURL     string `json:"clickURL,omitempty"`
	PixelURL     string `json:"pixelURL,omitempty"`
	Height       int    `json:"height,omitempty"`
	Width        int    `json:"width,omitempty"`
	//targetting
	ConnectionTypes []int    `json:"connectionTypes,omitempty"`
	CountryIds      []int    `json:"countryIds,omitempty"`
	CookieIds       []string `json:"cookieIds,omitempty"`
	//For scoring
	Impressions int     `json:"impressions,omitempty"`
	Clicks      int     `json:"clicks,omitempty"`
	ECPM        float64 `json:"ECPM,omitempty"`
	CTR         float64 `json:"CTR,omitempty"`
}

type MOAdTracking struct {

	//Db related
	ID           string `json:"rowID,omitempty"`
	PriorityType string `json:"pType,omitempty"`

	//Campaign
	CampaignID   int     `json:"campaignID,omitempty"`
	CampaignName string  `json:"campaignName,omitempty"`
	Bid          float64 `json:"bid,omitempty"`
	CategoryType int     `json:"categoryType,omitempty"`

	//Creative
	CreativeID   int    `json:"creativeID,omitempty"`
	CreativeType int    `json:"creativeType,omitempty"`
	ClickURL     string `json:"clickURL,omitempty"`
	PixelURL     string `json:"pixelURL,omitempty"`
	Type         int    `json:"type,omitempty"`
	Height       int    `json:"height,omitempty"`
	Width        int    `json:"width,omitempty"`
	Status       int    `json:"status,omitempty"`
	Priority     int    `json:"priority,omitempty"`

	//Operator
	Operator string `json:"operator,omitempty"`
	MSISDN   string `json:"msisdn,omitempty"`

	//Request
	RequestID     string     `json:"requestID,omitempty"`
	UserAgent     string     `json:"userAgent,omitempty"`
	Geo           GeoDetails `json:"geo,omitempty"`
	RequestHeader string     `json:"requestH,omitempty"`

	//Publisher
	Publisher          string `json:"publisher,omitempty"`
	ContentCategory    int    `json:"contentCategory,omitempty"`
	ContentSubCategory int    `json:"contentSubCategory,omitempty"`
	ContentPageID      int    `json:"contentPageID,omitempty"`
	AppID              string `json:"appID,omitempty"`
	Version            string `json:"version,omitempty"`

	//Resposne
	CookieId       string `json:"cookieId,omitempty"`
	ResponseStatus string `json:"responseStatus,omitempty"`
	TrackingID     string `json:"trackingID,omitempty"`
	ResponseTime   int    `json:"responseTime,omitempty"`
	Comment        string `json:"comment,omitempty"`

	//calculated fields
	ClickSessionTime int `json:"clickSessionTime,omitempty"`

	//time stamping
	Hour    int       `json:"hour,omitempty"`
	Date    time.Time `json:"date,omitempty"`
	UTCDate time.Time `json:"utcdate,omitempty"`
	Week    int       `json:"week,omitempty"`
	Month   string    `json:"month,omitempty"`
	MonthID int       `json:"monthID,omitempty"`
	Year    int       `json:"year,omitempty"`
}

type GeoDetails struct {
	IP          string  `json:"ip"`
	CountryName string  `json:"country_name"`
	CountryCode string  `json:"country_code"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
}
