package v1

import (
	"time"

	"cloud.google.com/go/bigquery"
)

//log
type BigQueryTrackLog struct {

	//Offer details
	OfferID     string `bigquery:"OfferID,nullable"`
	AffiliateID string `bigquery:"AffiliateID,nullable"`
	Status      string `bigquery:"Status,nullable"`
	Activity    int    `bigquery:"Activity,nullable"`

	//Request
	ReceivedOfferID       string `bigquery:"ReceivedOfferID,nullable"`
	ReceivedAffiliateID   string `bigquery:"ReceivedAffiliateID,nullable"`
	ReceivedTransactionID string `bigquery:"ReceivedTransactionID,nullable"`
	CookieID              string `bigquery:"CookieID,nullable"`
	TransactionID         string `bigquery:"TransactionID,nullable"`
	//Geo                   GeoDetails `bigquery:"Geo,nullable"`
	//ClickGeo              GeoDetails `bigquery:"ClickGeo,nullable"`
	UserAgent string `bigquery:"UserAgent,nullable"`

	RequestedPackage string `bigquery:"RequestedPackage,nullable"`
	Referer          string `bigquery:"Referer,nullable"`
	IP               string `bigquery:"IP,nullable"`
	SessionIP        string `bigquery:"SessionIP,nullable"`
	ConversionIP     string `bigquery:"ConversionIP,nullable"`
	MessageID        string `bigquery:"MessageID,nullable"`

	//transactions
	ClickDate           time.Time `bigquery:"ClickDate,nullable"`
	ClickURL            string    `bigquery:"ClickURL,nullable"`
	ClickRedirectURL    string    `bigquery:"ClickRedirectURL,nullable"`
	RecievedPostbackURL string    `bigquery:"RecievedPostbackURL,nullable"`
	SentPostbackURL     string    `bigquery:"SentPostbackURL,nullable"`
	MSISDN              string    `bigquery:"MSISDN,nullable"`
	Operator            string    `bigquery:"Operator,nullable"`
	ServiceCode         string    `bigquery:"ServiceCode,nullable" `
	RequestHeader       string    `bigquery:"RequestHeader,nullable"`
	Comment             string    `bigquery:"Comment,nullable"`
	ServiceID           string    `bigquery:"ServiceID,nullable"`
	CPCode              string    `bigquery:"CPCode,nullable"`

	//time stamping
	Minute  int       `bigquery:"Minute,nullable"`
	Hour    int       `bigquery:"Hour,nullable"`
	Date    time.Time `bigquery:"Date,nullable"`
	UTCDate time.Time `bigquery:"UTCDate,nullable"`
	Week    int       `bigquery:"Week,nullable"`
	Month   string    `bigquery:"Month,nullable"`
	MonthID int       `bigquery:"MonthID,nullable"`
	Year    int       `bigquery:"Year,nullable"`

	//HO parameters
	AdvertiserID    string `bigquery:"AdvertiserID,nullable"`
	AdvertiserRefID string `bigquery:"AdvertiserRefID,nullable"`
	AdvertiserSub   string `bigquery:"AdvertiserSub,nullable"`
	AdvertiserSub1  string `bigquery:"AdvertiserSub1,nullable"`
	AdvertiserSub2  string `bigquery:"AdvertiserSub2,nullable"`
	AdvertiserSub3  string `bigquery:"AdvertiserSub3,nullable"`
	AdvertiserSub4  string `bigquery:"AdvertiserSub4,nullable"`
	AdvertiserSub5  string `bigquery:"AdvertiserSub5,nullable"`
	AffiliateSub    string `bigquery:"AffiliateSub,nullable"`
	AffiliateSub1   string `bigquery:"AffiliateSub1,nullable"`
	AffiliateSub2   string `bigquery:"AffiliateSub2,nullable"`
	AffiliateSub3   string `bigquery:"AffiliateSub3,nullable"`
	AffiliateSub4   string `bigquery:"AffiliateSub4,nullable"`
	AffiliateSub5   string `bigquery:"AffiliateSub5,nullable"`
	AffiliateSub6   string `bigquery:"AffiliateSub6,nullable"`
	AffiliateSub7   string `bigquery:"AffiliateSub7,nullable"`
	AffiliateSub8   string `bigquery:"AffiliateSub8,nullable"`
	AffiliateSub9   string `bigquery:"AffiliateSub9,nullable"`
	AffiliateSub10  string `bigquery:"AffiliateSub10,nullable"`
	AffiliateName   string `bigquery:"AffiliateName,nullable"`
	AffiliateRef    string `bigquery:"AffiliateRef,nullable"`

	//conversion details
	ConvertedIP       string `bigquery:"ConvertedIP,nullable"`
	ConvertedTime     string `bigquery:"ConvertedTime,nullable"`
	ConvertedDate     string `bigquery:"ConvertedDate,nullable"`
	ConvertedDateTime string `bigquery:"ConvertedDateTime,nullable"`
	Source            string `bigquery:"Source,nullable"`

	//post events
	EventStatus        string  `bigquery:"EventStatus,nullable"`
	Amount             float64 `bigquery:"Amount,nullable"`
	ConversionUniqueID string  `bigquery:"ConversionUniqueID,nullable"`

	CreativeFile     string `bigquery:"CreativeFile,nullable"`
	GoalID           string `bigquery:"GoalID,nullable"`
	GoalRef          string `bigquery:"GoalRef,nullable"`
	Currency         string `bigquery:"Currency,nullable"`
	AffiliatePayout  string `bigquery:"AffiliatePayout,nullable"`
	AffiliateRevenue string `bigquery:"AffiliateRevenue,nullable"`
	OfferFileID      string `bigquery:"OfferFileID,nullable"`
	OfferName        string `bigquery:"OfferName,nullable"`
	OfferRef         string `bigquery:"OfferRef,nullable"`
	OfferURLID       string `bigquery:"OfferURLID,nullable"`
	Ran              string `bigquery:"Ran,nullable"`
	SaleAmount       string `bigquery:"SaleAmount,nullable"`

	//Mobile tracking parameters
	//GoogleAID       string `bigquery:"GoogleAID,nullable"`
	AndroidID       string `bigquery:"AndroidID,nullable"`
	AndroidIDMD5    string `bigquery:"AndroidIDMD5,nullable"`
	AndroidIDSHA1   string `bigquery:"AndroidIDSHA1,nullable"`
	DeviceBrand     string `bigquery:"DeviceBrand,nullable"`
	DeviceID        string `bigquery:"DeviceID,nullable"`
	DeviceIDMD5     string `bigquery:"DeviceIDMD5,nullable"`
	DeviceIDSHA1    string `bigquery:"DeviceIDSHA1,nullable"`
	DeviceModel     string `bigquery:"DeviceModel,nullable"`
	DeviceOS        string `bigquery:"DeviceOS,nullable"`
	DeviceOSVersion string `bigquery:"DeviceOSVersion,nullable"`
	IOSIfa          string `bigquery:"IOSIfa,nullable"`
	IOSIfaMD5       string `bigquery:"IOSIfaMD5,nullable"`
	IOSIfaSHA1      string `bigquery:"IOSIfaSHA1,nullable"`
	IOSIfv          string `bigquery:"IOSIfv,nullable"`
	MacAddress      string `bigquery:"MacAddress,nullable"`
	MacAddressMD5   string `bigquery:"MacAddressMD5,nullable"`
	MacAddressSHA1  string `bigquery:"MacAddressSHA1,nullable"`
	WindowsAID      string `bigquery:"WindowsAID,nullable"`
	WindowsAIDMD5   string `bigquery:"WindowsAIDMD5,nullable"`
	WindowsSHA1     string `bigquery:"WindowsSHA1,nullable"`
	ODIN            string `bigquery:"ODIN,nullable"`
	OpenUDID        string `bigquery:"OpenUDID,nullable"`
	UNID            string `bigquery:"UNID,nullable"`
	AppUserID       string `bigquery:"AppUserID,nullable"`
}

type BigQueryTrackLogBackUp struct {

	//Offer details
	OfferID     bigquery.NullString `json:"offerID,omitempty" bson:"offerID,omitempty"`
	AffiliateID bigquery.NullString `json:"affiliateID,omitempty" bson:"affiliateID,omitempty"`
	Status      bigquery.NullString `json:"status,omitempty" bson:"status,omitempty"`
	Activity    bigquery.NullInt64  `json:"activity,omitempty" bson:"activity,omitempty"`

	//Request
	ReceivedOfferID       bigquery.NullString `json:"receivedOfferID,omitempty" bson:"receivedOfferID,omitempty"`
	ReceivedAffiliateID   bigquery.NullString `json:"receivedAffiliateID,omitempty" bson:"receivedAffiliateID,omitempty"`
	ReceivedTransactionID bigquery.NullString `json:"receivedTransactionID,omitempty" bson:"receivedTransactionID,omitempty"`
	CookieID              bigquery.NullString `json:"cookieID,omitempty" bson:"cookieID,omitempty"`
	TransactionID         bigquery.NullString `json:"transactionID,omitempty" bson:"transactionID,omitempty"`
	Geo                   *BigQueryGeoDetails `bigquery:"Geo,nullable" json:"geo,omitempty" bson:"geo,omitempty"`
	ClickGeo              *BigQueryGeoDetails `bigquery:"ClickGeo,nullable" json:"clickGeo,omitempty" bson:"clickGeo,omitempty"`
	UserAgent             bigquery.NullString `json:"userAgent,omitempty" bson:"userAgent,omitempty"`

	RequestedPackage bigquery.NullString `json:"requestedPackage,omitempty" bson:"requestedPackage,omitempty"`
	Referer          bigquery.NullString `json:"referer,omitempty" bson:"referer,omitempty"`
	IP               bigquery.NullString `json:"ip,omitempty" bson:"ip,omitempty"`
	SessionIP        bigquery.NullString `json:"sessionIP,omitempty" bson:"sessionIP,omitempty"`
	ConversionIP     bigquery.NullString `json:"conversionIP,omitempty" bson:"conversionIP,omitempty"`
	MessageID        bigquery.NullString `json:"messageID,omitempty" bson:"messageID,omitempty"`

	//transactions
	ClickDate           bigquery.NullTimestamp `json:"clickDate,omitempty" bson:"clickDate,omitempty"`
	ClickURL            bigquery.NullString    `json:"clickURL,omitempty" bson:"clickURL,omitempty"`
	ClickRedirectURL    bigquery.NullString    `json:"clickRedirectURL,omitempty" bson:"clickRedirectURL,omitempty"`
	ReceivedPostbackURL bigquery.NullString    `json:"receivedPostbackURL,omitempty" bson:"receivedPostbackURL,omitempty"`
	SentPostbackURL     bigquery.NullString    `json:"sentPostbackURL,omitempty" bson:"sentPostbackURL,omitempty"`
	MSISDN              bigquery.NullString    `json:"msisdn,omitempty" bson:"msisdn,omitempty"`
	Operator            bigquery.NullString    `json:"operator,omitempty" bson:"operator,omitempty"`
	ServiceCode         bigquery.NullString    `json:"serviceCode,omitempty" bson:"serviceCode,omitempty"`
	RequestHeader       bigquery.NullString    `json:"requestedHeader,omitempty" bson:"requestedHeader,omitempty"`
	Comment             bigquery.NullString    `json:"comment,omitempty" bson:"comment,omitempty"`
	ServiceID           bigquery.NullString    `json:"serviceID,omitempty" bson:"serviceID,omitempty"`
	CPCode              bigquery.NullString    `json:"cpCode,omitempty" bson:"cpCode,omitempty"`

	//time stamping
	Minute  bigquery.NullInt64     `json:"minute,omitempty" bson:"minute,omitempty"`
	Hour    bigquery.NullInt64     `json:"hour,omitempty" bson:"hour,omitempty"`
	Date    bigquery.NullTimestamp `json:"date,omitempty" bson:"date,omitempty"`
	UTCDate bigquery.NullTimestamp `json:"utcdate,omitempty" bson:"utcdate,omitempty"`
	Week    bigquery.NullInt64     `json:"week,omitempty" bson:"week,omitempty"`
	Month   bigquery.NullString    `json:"month,omitempty" bson:"month,omitempty"`
	MonthID bigquery.NullInt64     `json:"monthID,omitempty" bson:"monthID,omitempty"`
	Year    bigquery.NullInt64     `json:"year,omitempty" bson:"year,omitempty"`

	//HO parameters
	AdvertiserID    bigquery.NullString `json:"advertiserID,omitempty" bson:"advertiserID,omitempty"`
	AdvertiserRefID bigquery.NullString `json:"advertiserRefID,omitempty" bson:"advertiserRefID,omitempty"`
	AdvertiserSub   bigquery.NullString `json:"advertiserSub,omitempty" bson:"advertiserSub,omitempty"`
	AdvertiserSub1  bigquery.NullString `json:"advertiserSub1,omitempty" bson:"advertiserSub1,omitempty"`
	AdvertiserSub2  bigquery.NullString `json:"advertiserSub2,omitempty" bson:"advertiserSub2,omitempty"`
	AdvertiserSub3  bigquery.NullString `json:"advertiserSub3,omitempty" bson:"advertiserSub3,omitempty"`
	AdvertiserSub4  bigquery.NullString `json:"advertiserSub4,omitempty" bson:"advertiserSub4,omitempty"`
	AdvertiserSub5  bigquery.NullString `json:"advertiserSub5,omitempty" bson:"advertiserSub5,omitempty"`
	AffiliateSub    bigquery.NullString `json:"affiliateSub,omitempty" bson:"affiliateSub,omitempty"`
	AffiliateSub1   bigquery.NullString `json:"affiliateSub1,omitempty" bson:"affiliateSub1,omitempty"`
	AffiliateSub2   bigquery.NullString `json:"affiliateSub2,omitempty" bson:"affiliateSub2,omitempty"`
	AffiliateSub3   bigquery.NullString `json:"affiliateSub3,omitempty" bson:"affiliateSub3,omitempty"`
	AffiliateSub4   bigquery.NullString `json:"affiliateSub4,omitempty" bson:"affiliateSub4,omitempty"`
	AffiliateSub5   bigquery.NullString `json:"affiliateSub5,omitempty" bson:"affiliateSub5,omitempty"`
	AffiliateSub6   bigquery.NullString `json:"affiliateSub6,omitempty" bson:"affiliateSub6,omitempty"`
	AffiliateSub7   bigquery.NullString `json:"affiliateSub7,omitempty" bson:"affiliateSub7,omitempty"`
	AffiliateSub8   bigquery.NullString `json:"affiliateSub8,omitempty" bson:"affiliateSub8,omitempty"`
	AffiliateSub9   bigquery.NullString `json:"affiliateSub9,omitempty" bson:"affiliateSub9,omitempty"`
	AffiliateSub10  bigquery.NullString `json:"affiliateSub10,omitempty" bson:"affiliateSub10,omitempty"`
	AffiliateName   bigquery.NullString `json:"affiliateName,omitempty" bson:"affiliateName,omitempty"`
	AffiliateRef    bigquery.NullString `json:"affiliateRef,omitempty" bson:"affiliateRef,omitempty"`

	//conversion details
	ConvertedIP       bigquery.NullString `json:"convertedIP,omitempty" bson:"convertedIP,omitempty"`
	ConvertedTime     bigquery.NullString `json:"convertedTime,omitempty" bson:"convertedTime,omitempty"`
	ConvertedDate     bigquery.NullString `json:"convertedDate,omitempty" bson:"convertedDate,omitempty"`
	ConvertedDateTime bigquery.NullString `json:"convertedDateTime,omitempty" bson:"convertedDateTime,omitempty"`
	Source            bigquery.NullString `json:"source,omitempty" bson:"source,omitempty"`

	//post events
	EventStatus        bigquery.NullString  `json:"eventStatus,omitempty" bson:"eventStatus,omitempty"`
	Amount             bigquery.NullFloat64 `json:"amount,omitempty" bson:"amount,omitempty"`
	ConversionUniqueID bigquery.NullString  `json:"conversionUniqueID,omitempty" bson:"conversionUniqueID,omitempty"`

	CreativeFile     bigquery.NullString `json:"creativeFile,omitempty" bson:"creativeFile,omitempty"`
	GoalID           bigquery.NullString `json:"goalID,omitempty" bson:"goalID,omitempty"`
	GoalRef          bigquery.NullString `json:"goalRef,omitempty" bson:"goalRef,omitempty"`
	Currency         bigquery.NullString `json:"currency,omitempty" bson:"currency,omitempty"`
	AffiliatePayout  bigquery.NullString `json:"affiliatePayout,omitempty" bson:"affiliatePayout,omitempty"`
	AffiliateRevenue bigquery.NullString `json:"affiliateRevenue,omitempty" bson:"affiliateRevenue,omitempty"`
	OfferFileID      bigquery.NullString `json:"offerFileID,omitempty" bson:"offerFileID,omitempty"`
	OfferName        bigquery.NullString `json:"offerName,omitempty" bson:"offerName,omitempty"`
	OfferRef         bigquery.NullString `json:"offerRef,omitempty" bson:"offerRef,omitempty"`
	OfferURLID       bigquery.NullString `json:"offerURLID,omitempty" bson:"offerURLID,omitempty"`
	Ran              bigquery.NullString `json:"ran,omitempty" bson:"ran,omitempty"`
	SaleAmount       bigquery.NullString `json:"saleAmount,omitempty" bson:"saleAmount,omitempty"`

	//Mobile tracking parameters
	GoogleAID       bigquery.NullString `json:"googleAID,omitempty" bson:"googleAID,omitempty"`
	AndroidID       bigquery.NullString `json:"androidID,omitempty" bson:"androidID,omitempty"`
	AndroidIDMD5    bigquery.NullString `json:"androidIDMD5,omitempty" bson:"androidIDMD5,omitempty"`
	AndroidIDSHA1   bigquery.NullString `json:"androidIDSHA1,omitempty" bson:"androidIDSHA1,omitempty"`
	DeviceBrand     bigquery.NullString `json:"deviceBrand,omitempty" bson:"deviceBrand,omitempty"`
	DeviceID        bigquery.NullString `json:"deviceID,omitempty" bson:"deviceID,omitempty"`
	DeviceIDMD5     bigquery.NullString `json:"deviceIDMD5,omitempty" bson:"deviceIDMD5,omitempty"`
	DeviceIDSHA1    bigquery.NullString `json:"deviceIDSHA1,omitempty" bson:"deviceIDSHA1,omitempty"`
	DeviceModel     bigquery.NullString `json:"deviceModel,omitempty" bson:"deviceModel,omitempty"`
	DeviceOS        bigquery.NullString `json:"deviceOS,omitempty" bson:"deviceOS,omitempty"`
	DeviceOSVersion bigquery.NullString `json:"deviceOSVersion,omitempty" bson:"deviceOSVersion,omitempty"`
	IOSIfa          bigquery.NullString `json:"iosIFA,omitempty" bson:"iosIFA,omitempty"`
	IOSIfaMD5       bigquery.NullString `json:"iosIFAMD5,omitempty" bson:"iosIFAMD5,omitempty"`
	IOSIfaSHA1      bigquery.NullString `json:"iosIFASHA1,omitempty" bson:"iosIFASHA1,omitempty"`
	IOSIfv          bigquery.NullString `json:"iosIFV,omitempty" bson:"iosIFV,omitempty"`
	MacAddress      bigquery.NullString `json:"macAddress,omitempty" bson:"macAddress,omitempty"`
	MacAddressMD5   bigquery.NullString `json:"macAddressMD5,omitempty" bson:"macAddressMD5,omitempty"`
	MacAddressSHA1  bigquery.NullString `json:"macAddressSHA1,omitempty" bson:"macAddressSHA1,omitempty"`
	WindowsAID      bigquery.NullString `json:"windowsAID,omitempty" bson:"windowsAID,omitempty"`
	WindowsAIDMD5   bigquery.NullString `json:"windowsAIDMD5,omitempty" bson:"windowsAIDMD5,omitempty"`
	WindowsSHA1     bigquery.NullString `json:"windowsAIDSHA1,omitempty" bson:"windowsAIDSHA1,omitempty"`
	ODIN            bigquery.NullString `json:"odin,omitempty" bson:"odin,omitempty"`
	OpenUDID        bigquery.NullString `json:"openUDID,omitempty" bson:"openUDID,omitempty"`
	UNID            bigquery.NullString `json:"unid,omitempty" bson:"unid,omitempty"`
	AppUserID       bigquery.NullString `json:"appUserID,omitempty" bson:"appUserID,omitempty"`
}

type BigQueryGeoDetails struct {
	IP          bigquery.NullString  `json:"ip,omitempty" bson:"ip,omitempty"`
	CountryName bigquery.NullString  `json:"countryname,omitempty" bson:"countryname,omitempty"`
	CountryCode bigquery.NullString  `json:"countrycode,omitempty" bson:"countrycode,omitempty"`
	RegionCode  bigquery.NullString  `json:"regioncode,omitempty" bson:"regioncode,omitempty"`
	RegionName  bigquery.NullString  `json:"regionname,omitempty" bson:"regionname,omitempty"`
	City        bigquery.NullString  `json:"city,omitempty" bson:"city,omitempty"`
	ZipCode     bigquery.NullString  `json:"zipcode,omitempty" bson:"zipcode,omitempty"`
	TimeZone    bigquery.NullString  `json:"timezone,omitempty" bson:"timezone,omitempty"`
	Latitude    bigquery.NullFloat64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   bigquery.NullFloat64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
	MetroCode   bigquery.NullInt64   `json:"metrocode,omitempty" bson:"metrocode,omitempty"`
}
