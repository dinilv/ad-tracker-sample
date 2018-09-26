package v1

import "time"

//log
type TrackLog struct {

	//Offer details
	OfferID     string `bson:"offerID,omitempty"`
	AffiliateID string `bson:"affiliateID,omitempty"`
	OfferType   string `bson:"offerType,omitempty"`
	Sequence    int    `bson:"sequence,omitempty"`
	Status      string `bson:"status,omitempty"`
	Activity    int    `bson:"activity,omitempty"`

	//Request
	ReceivedOfferID       string     `bson:"receivedOfferID,omitempty"`
	ReceivedAffiliateID   string     `bson:"receivedAffiliateID,omitempty"`
	ReceivedTransactionID string     `bson:"receivedTransactionID,omitempty"`
	CookieID              string     `bson:"cookieID,omitempty"`
	TransactionID         string     `bson:"transactionID,omitempty"`
	Geo                   GeoDetails `bson:"geo,omitempty"`
	ClickGeo              GeoDetails `bson:"clickGeo,omitempty"`
	UserAgent             string     `bson:"userAgent,omitempty"`
	Method                string     `bson:"method,omitempty"`
	ContentType           string     `bson:"contentType,omitempty"`
	RequestBody           string     `bson:"requestBody,omitempty"`

	RequestedPackage string `bson:"requestedPackage,omitempty"`
	Referer          string `bson:"referer,omitempty"`
	IP               string `bson:"ip,omitempty"`
	SessionIP        string `bson:"sessionIP,omitempty"`
	ConversionIP     string `bson:"conversionIP,omitempty"`
	MessageID        string `bson:"messageID,omitempty"`

	//transactions
	ClickDate           time.Time `bson:"clickDate,omitempty"`
	ClickURL            string    `bson:"clickURL,omitempty"`
	ClickRedirectURL    string    `bson:"clickRedirectURL,omitempty"`
	ReceivedPostbackURL string    `bson:"receivedPostbackURL,omitempty"`
	SentPostbackURL     string    `bson:"sentPostbackURL,omitempty"`
	MSISDN              string    `bson:"msisdn,omitempty"`
	Operator            string    `bson:"operator,omitempty"`
	ServiceCode         string    `bson:"serviceCode,omitempty"`
	RequestHeader       string    `bson:"requestH,omitempty"`
	Comment             string    `bson:"comment,omitempty"`
	ServiceID           string    `bson:"serviceID,omitempty"`
	CPCode              string    `bson:"cpCode,omitempty"`

	//for performance monitoring
	ResponseTime     float64 `bson:"responseTime,omitempty"`
	SubscriberName   string  `bson:"subscriberName,omitempty"`
	APITime          string  `bson:"apiTime,omitempty"`
	BrokerTime       string  `bson:"brokerTime,omitempty"`
	GoRoutineTime    string  `bson:"goRoutineTime,omitempty"`
	SubscriptionTime string  `bson:"subscriptionTime,omitempty"`
	RetryCount       string  `bson:"retryCount,omitempty"`
	ErrorMessage     string  `bson:"errorMessage,omitempty"`
	FilteredCount    string  `bson:"filteredCount,omitempty"`
	Processor        string  `bson:"processor,omitempty"`
	Host             string  `bson:"host,omitempty"`

	//time stamping
	Minute  int       `bson:"minute,omitempty"`
	Hour    int       `bson:"hour"`
	Date    time.Time `bson:"date,omitempty"`
	UTCDate time.Time `bson:"utcdate,omitempty"`
	Week    int       `bson:"week,omitempty"`
	Month   string    `bson:"month,omitempty"`
	MonthID int       `bson:"monthID,omitempty"`
	Year    int       `bson:"year,omitempty"`

	//HO parameters
	AdvertiserID    string `bson:"advertiserID,omitempty"`
	AdvertiserRefID string `bson:"advertiserRefID,omitempty"`
	AdvertiserSub   string `bson:"advertiserSub,omitempty"`
	AdvertiserSub1  string `bson:"advertiserSub1,omitempty"`
	AdvertiserSub2  string `bson:"advertiserSub2,omitempty"`
	AdvertiserSub3  string `bson:"advertiserSub3,omitempty"`
	AdvertiserSub4  string `bson:"advertiserSub4,omitempty"`
	AdvertiserSub5  string `bson:"advertiserSub5,omitempty"`
	AffiliateSub    string `bson:"affiliateSub,omitempty"`
	AffiliateSub1   string `bson:"affiliateSub1,omitempty"`
	AffiliateSub2   string `bson:"affiliateSub2,omitempty"`
	AffiliateSub3   string `bson:"affiliateSub3,omitempty"`
	AffiliateSub4   string `bson:"affiliateSub4,omitempty"`
	AffiliateSub5   string `bson:"affiliateSub5,omitempty"`
	AffiliateSub6   string `bson:"affiliateSub6,omitempty"`
	AffiliateSub7   string `bson:"affiliateSub7,omitempty"`
	AffiliateSub8   string `bson:"affiliateSub8,omitempty"`
	AffiliateSub9   string `bson:"affiliateSub9,omitempty"`
	AffiliateSub10  string `bson:"affiliateSub10,omitempty"`
	AffiliateName   string `bson:"affiliateName,omitempty"`
	AffiliateRef    string `bson:"affiliateRef,omitempty"`

	//conversion details
	ConvertedIP       string `bson:"convertedIP,omitempty"`
	ConvertedTime     string `bson:"convertedTime,omitempty"`
	ConvertedDate     string `bson:"convertedDate,omitempty"`
	ConvertedDateTime string `bson:"convertedDateTime,omitempty"`
	Source            string `bson:"source,omitempty"`

	//post events
	EventStatus        string  `bson:"eventStatus,omitempty"`
	Amount             float64 `bson:"amount,omitempty"`
	ConversionUniqueID string  `bson:"conversionUniqueID,omitempty"`

	CreativeFile     string `bson:"creativeFile,omitempty"`
	GoalID           string `bson:"goalID,omitempty"`
	GoalRef          string `bson:"goalRef,omitempty"`
	Currency         string `bson:"currency,omitempty"`
	AffiliatePayout  string `bson:"affiliatePayout,omitempty"`
	AffiliateRevenue string `bson:"affiliateRevenue,omitempty"`
	OfferFileID      string `bson:"offerFileID,omitempty"`
	OfferName        string `bson:"offerName,omitempty"`
	OfferRef         string `bson:"offerRef,omitempty"`
	OfferURLID       string `bson:"offerURLID,omitempty"`
	Ran              string `bson:"ran,omitempty"`
	SaleAmount       string `bson:"saleAmount,omitempty"`

	//Mobile tracking parameters
	GoogleAID       string `bson:"googleAID,omitempty"`
	AndroidID       string `bson:"androidID,omitempty"`
	AndroidIDMD5    string `bson:"androidIDMD5,omitempty"`
	AndroidIDSHA1   string `bson:"androidIDSHA1,omitempty"`
	DeviceBrand     string `bson:"deviceBrand,omitempty"`
	DeviceID        string `bson:"deviceID,omitempty"`
	DeviceIDMD5     string `bson:"deviceIDMD5,omitempty"`
	DeviceIDSHA1    string `bson:"deviceIDSHA1,omitempty"`
	DeviceModel     string `bson:"deviceModel,omitempty"`
	DeviceOS        string `bson:"deviceOS,omitempty"`
	DeviceOSVersion string `bson:"deviceOSVersion,omitempty"`
	IOSIfa          string `bson:"iosIFA,omitempty"`
	IOSIfaMD5       string `bson:"iosIFAMD5,omitempty"`
	IOSIfaSHA1      string `bson:"iosIFASHA1,omitempty"`
	IOSIfv          string `bson:"iosIFV,omitempty"`
	MacAddress      string `bson:"macAddress,omitempty"`
	MacAddressMD5   string `bson:"macAddressMD5,omitempty"`
	MacAddressSHA1  string `bson:"macAddressSHA1,omitempty"`
	WindowsAID      string `bson:"windowsAID,omitempty"`
	WindowsAIDMD5   string `bson:"windowsAIDMD5,omitempty"`
	WindowsSHA1     string `bson:"windowsAIDSHA1,omitempty"`
	ODIN            string `bson:"odin,omitempty"`
	OpenUDID        string `bson:"openUDID,omitempty"`
	UNID            string `bson:"unid,omitempty"`
	AppUserID       string `bson:"appUserID,omitempty"`
}

//Entities
type Offer struct {
	OfferID    string `bson:"offerID,omitempty"`
	OfferRefID string `bson:"offerRefID,omitempty"`
	OfferName  string `bson:"offerName,omitempty,omitempty"`
	ServiceID  string `bson:"serviceID,omitempty"`
	Template   string `bson:"template,omitempty"`
	GoalID     string `bson:"goalID,omitempty"`
	OfferType  string `bson:"offerType,omitempty"`
}

type Affiliate struct {
	AffiliateID    string  `bson:"affiliateID,omitempty"`
	AffiliateName  string  `bson:"affiliateName,omitempty"`
	Mqf            float64 `bson:"mqf,omitempty"`
	AffiliateRefID string  `bson:"affiliateRefID,omitempty"`
	MediaTemplate  string  `bson:"mediaTemplate,omitempty"`
}

type TrackerModel struct {
	Offer
	Affiliate
	TrackLog
}

type ID struct {
	Date        string `bson:"date,omitempty"`
	AffiliateID string `bson:"affiliateID,omitempty"`
	OfferID     string `bson:"offerID,omitempty"`
	Hour        string `bson:"hour,omitempty"`
}
type Report struct {
	Ids ID `bson:"_id,omitempty"`
}

type APIReport struct {
	UtcDate           time.Time `json:"utcdate" bson:"utcdate"`
	Hour              int       `json:"hour" bson:"hour"`
	Date              string    `json:"date,omitempty" bson:"date,omitempty"`
	AffiliateID       string    `json:"affiliate_id,omitempty" bson:"affiliate_id,omitempty"`
	AffiliateName     string    `json:"affiliate_name,omitempty" bson:"affiliate_name,omitempty"`
	AffiliateRefID    string    `json:"affiliate_reference_id,omitempty" bson:"affiliate_reference_id,omitempty"`
	OfferID           string    `json:"offer_id,omitempty" bson:"offer_id,omitempty"`
	FwdOfferID        string    `json:"fwd_offer_id,omitempty" bson:"fwd_offer_id,omitempty"`
	FwdAffiliateID    string    `json:"fwd_affiliate_id,omitempty" bson:"fwd_affiliate_id,omitempty"`
	OfferRefID        string    `json:"offer_reference_id,omitempty" bson:"offer_reference_id,omitempty"`
	OfferName         string    `json:"offer_name,omitempty" bson:"offer_name,omitempty"`
	RecvOfferID       string    `json:"recv_offer_id,omitempty" bson:"recv_offer_id,omitempty"`
	RecvOfferName     string    `json:"recv_offer_name,omitempty" bson:"recv_offer_name,omitempty"`
	RecvAffiliateID   string    `json:"recv_affiliate_id,omitempty" bson:"recv_affiliate_id,omitempty"`
	RecvAffiliateName string    `json:"recv_affiliate_name,omitempty" bson:"recv_affiliate_name,omitempty"`

	Impressions                 int               `json:"impressions" bson:"impressions,omitempty"`
	DuplicateImpressions        int               `json:"duplicate_impressions,omitempty" bson:"duplicate_impressions,omitempty"`
	UniqueImpressions           int               `json:"unique_impressions,omitempty" bson:"unique_impressions,omitempty"`
	Clicks                      int               `json:"clicks" bson:"clicks,omitempty"`
	DuplicateClicks             int               `json:"duplicate_clicks,omitempty" bson:"duplicate_clicks,omitempty"`
	UniqueClicks                int               `json:"unique_clicks,omitempty" bson:"unique_clicks,omitempty"`
	RotatedClicksFwd            int               `json:"rotated_clicks_fwd" bson:"rotated_clicks_fwd,omitempty"`
	Conversions                 int               `json:"conversions" bson:"conversions,omitempty"`
	SentConversions             int               `json:"sent_conversions" bson:"sent_conversions,omitempty"`
	UnSentConversions           int               `json:"unsent_conversions" bson:"unsent_conversions,omitempty"`
	SentRotatedConversionsFwd   int               `json:"sent_rotated_conversions_fwd" bson:"sent_rotated_conversions_fwd,omitempty"`
	UnSentRotatedConversionsFwd int               `json:"unsent_rotated_conversions_fwd" bson:"unsent_rotated_conversions_fwd,omitempty"`
	RotatedConversionsFwd       int               `json:"rotated_conversions_fwd" bson:"rotated_conversions_fwd,omitempty"`
	DuplicateConversions        int               `json:"duplicate_conversions,omitempty" bson:"duplicate_conversions,omitempty"`
	Events                      int               `json:"events,omitempty" bson:"events,omitempty"`
	Frauds                      int               `json:"frauds,omitempty" bson:"frauds,omitempty"`
	Filtered                    int               `json:"filtered,omitempty" bson:"filtered,omitempty"`
	Country                     string            `json:"geo,omitempty" bson:"geo,omitempty"`
	ClickCountry                string            `json:"clickGeo,omitempty" bson:"clickGeo,omitempty"`
	AffiliateSub3               string            `json:"affiliateSub3,omitempty" bson:"affiliateSub3,omitempty"`
	Rotations                   []RotationDetails `json:"rotations,omitempty" bson:"rotations,omitempty"`
	ReceivedRotations           []RotationDetails `json:"recv_rotations,omitempty" bson:"recv_rotations,omitempty"`
}

type RotationDetails struct {
	OfferID           string `json:"offer_id,omitempty" bson:"offer_id,omitempty"`
	OfferRefID        string `json:"offer_reference_id,omitempty" bson:"offer_reference_id,omitempty"`
	AffiliateID       string `json:"affiliate_id,omitempty" bson:"affiliate_id,omitempty"`
	AffiliateRefID    string `json:"affiliate_reference_id,omitempty" bson:"affiliate_reference_id,omitempty"`
	Clicks            int    `json:"clicks,omitempty" bson:"clicks,omitempty"`
	Conversions       int    `json:"conversions,omitempty" bson:"conversions,omitempty"`
	SentConversions   int    `json:"sent_conversions,omitempty" bson:"sent_conversions,omitempty"`
	UnSentConversions int    `json:"unsent_conversions,omitempty" bson:"unsent_conversions,omitempty"`
}

type OfferID struct {
	OfferID   string    `json:"offerID,omitempty" bson:"offerID,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
type AffiliateMQF struct {
	Mqf float64 `json:"mqf,omitempty" bson:"mqf,omitempty"`
}
type ExhaustedOffers struct {
	OfferID     string `json:"offerID,omitempty" bson:"offerID,omitempty"`
	AffiliateID string `json:"affiliateID,omitempty" bson:"affiliateID,omitempty"`
	Status      string `json:"status,omitempty" bson:"status,omitempty"`
}

type MSISDNDetails struct {
	MSISDN     string    `json:"msisdn,omitempty" bson:"msisdn,omitempty"`
	Operator   string    `json:"operator,omitempty" bson:"operator,omitempty"`
	Blocked    string    `json:"blocked,omitempty" bson:"blocked,omitempty"`
	ServiceIDs []string  `json:"serviceIDs,omitempty" bson:"serviceIDs,omitempty"`
	OfferIDs   []string  `json:"offerIDs,omitempty" bson:"offerIDs,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type OfferStack struct {
	OfferID     string    `json:"offerID,omitempty" bson:"offerID,omitempty"`
	OfferType   string    `json:"offerType,omitempty" bson:"offerType,omitempty"`
	LastEventID string    `json:"lastEventID,omitempty" bson:"lastEventID,omitempty"`
	LastEvent   string    `json:"lastEvent,omitempty" bson:"lastEventID,omitempty"`
	Bid         float64   `json:"bid,omitempty" bson:"clicks,omitempty"`
	Status      string    `json:"status,omitempty" bson:"status,omitempty"`
	AvoidStatus string    `json:"avoidStatus,omitempty" bson:"avoidStatus,omitempty"`
	Group       string    `json:"group,omitempty" bson:"group,omitempty"`
	ECPC        float64   `json:"ecpc,omitempty" bson:"ecpc,omitempty"`
	Clicks      int       `json:"clicks,omitempty" bson:"clicks,omitempty"`
	Conversions int       `json:"conversions,omitempty" bson:"conversions,omitempty"`
	CountryIDs  []string  `json:"countryIDs,omitempty" bson:"countryIDs,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
type AdcamieEvent struct {
	OfferID       string `json:"offerID,omitempty" bson:"offerID,omitempty"`
	OfferType     string `json:"offerType,omitempty" bson:"offerType,omitempty"`
	EventID       string `json:"eventID,omitempty" bson:"eventID,omitempty"`
	EventName     string `json:"eventName,omitempty" bson:"eventName,omitempty"`
	Comment       string `json:"comment,omitempty" bson:"comment,omitempty"`
	AdcamieStatus string `json:"adcamieStatus,omitempty" bson:"adcamieStatus,omitempty"`
	AvoidStatus   string `json:"avoidStatus,omitempty" bson:"avoidStatus,omitempty"`
	//Campaign Details for update
	CampaignID  string    `json:"campaignID,omitempty" bson:"campaignID,omitempty"`
	Bid         float64   `json:"bid,omitempty" bson:"bid,omitempty"`
	Status      string    `json:"status,omitempty" bson:"status,omitempty"`
	Group       string    `json:"group,omitempty" bson:"group,omitempty"`
	ECPC        float64   `json:"ecpc,omitempty" bson:"ecpc,omitempty"`
	Clicks      int       `json:"clicks,omitempty" bson:"clicks,omitempty"`
	Conversions int       `json:"conversions,omitempty" bson:"conversions,omitempty"`
	CountryIDs  []string  `json:"countryIDs,omitempty" bson:"countryIDs,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UTCDate     time.Time `bson:"utcdate,omitempty"`
	IP          string    `json:"ip,omitempty" bson:"ip,omitempty"`
}

type AdcamieEvents struct {
	Events []AdcamieEvent `json:"events,omitempty" bson:"events,omitempty"`
	IP     string         `json:"ip,omitempty" bson:"ip,omitempty"`
}

type OfferGeo struct {
	OfferID    string   `json:"offerID,omitempty" bson:"offerID,omitempty"`
	CountryIDs []string `json:"countryIDs,omitempty" bson:"countryIDs,omitempty"`
}

type PostBackReq struct {
	AffiliateID   string `json:"aff_id,omitempty" bson:"aff_id,omitempty"`
	OfferID       string `json:"offer_id,omitempty" bson:"offer_id,omitempty"`
	TransactionID string `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
}

type RedisTransaction struct {
	Transaction string `json:"transaction,omitempty" bson:"transaction,omitempty"`
	Value       string `json:"value,omitempty" bson:"value,omitempty"`
	Page        int    `json:"page,omitempty" bson:"page,omitempty"`
}

type Message struct {
	ID  string            `json:"id,omitempty" bson:"id,omitempty"`
	Msg map[string]string `json:"attributes,omitempty" bson:"attributes,omitempty"`
}

type APIMetaDataReport struct {
	UTCDate                     time.Time `bson:"utcdate,omitempty"`
	Hour                        int       `json:"hour" bson:"hour"`
	Date                        time.Time `json:"date,omitempty" bson:"date,omitempty"`
	AffiliateID                 string    `json:"affiliate_id,omitempty" bson:"affiliate_id,omitempty"`
	OfferID                     string    `json:"offer_id,omitempty" bson:"offer_id,omitempty"`
	FwdOfferID                  string    `json:"fwd_offer_id,omitempty" bson:"fwd_offer_id,omitempty"`
	FwdAffiliateID              string    `json:"fwd_affiliate_id,omitempty" bson:"fwd_affiliate_id,omitempty"`
	OfferRefID                  string    `json:"offer_reference_id,omitempty" bson:"offer_reference_id,omitempty"`
	RecvOfferID                 string    `json:"recv_offer_id,omitempty" bson:"recv_offer_id,omitempty"`
	RecvAffiliateID             string    `json:"recv_affiliate_id,omitempty" bson:"recv_affiliate_id,omitempty"`
	Impressions                 int       `json:"impressions" bson:"impressions,omitempty"`
	Clicks                      int       `json:"clicks" bson:"clicks,omitempty"`
	UniqueClicks                int       `json:"unique_clicks" bson:"unique_clicks,omitempty"`
	DuplicateClicks             int       `json:"duplicate_clicks" bson:"duplicate_clicks,omitempty"`
	RotatedClicksFwd            int       `json:"rotated_clicks_fwd" bson:"rotated_clicks_fwd,omitempty"`
	Conversions                 int       `json:"conversions" bson:"conversions,omitempty"`
	SentConversions             int       `json:"sent_conversions" bson:"sent_conversions,omitempty"`
	UnSentConversions           int       `json:"unsent_conversions" bson:"unsent_conversions,omitempty"`
	SentRotatedConversionsFwd   int       `json:"sent_rotated_conversions_fwd" bson:"sent_rotated_conversions_fwd,omitempty"`
	UnSentRotatedConversionsFwd int       `json:"unsent_rotated_conversions_fwd" bson:"unsent_rotated_conversions_fwd,omitempty"`
	RotatedConversionsFwd       int       `json:"rotated_conversions_fwd" bson:"rotated_conversions_fwd,omitempty"`
	DuplicateConversions        int       `json:"duplicate_conversions,omitempty" bson:"duplicate_conversions,omitempty"`
	PostEvents                  int       `json:"events,omitempty" bson:"events,omitempty"`
	Collection                  string    `json:"collection,omitempty" bson:"collection,omitempty"`
	CookieCount                 int       `json:"cookie_count,omitempty" bson:"cookie_count,omitempty"`
	CookieID                    string    `json:"cookie_id,omitempty" bson:"cookie_id,omitempty"`
}

type PostBackPingLog struct {
	Date            time.Time `json:"date,omitempty" bson:"date,omitempty"`
	APITime         string    `bson:"apiTime,omitempty"`
	OfferID         string    `bson:"offerID,omitempty"`
	AffiliateID     string    `bson:"affiliateID,omitempty"`
	TransactionID   string    `bson:"transactionID,omitempty"`
	UTCDate         time.Time `bson:"utcdate,omitempty"`
	ResponseCode    string    `bson:"responseCode,omitempty"`
	Response        string    `bson:"response,omitempty"`
	TimeTaken       string    `bson:"timeTaken,omitempty"`
	ErrorMessage    string    `bson:"errorMessage,omitempty"`
	SentPostbackURL string    `bson:"sentPostbackURL,omitempty"`
}
type RedisTransactionBackup struct {
	UTCDate       time.Time `json:"utcdate,omitempty" bson:"utcdate,omitempty"`
	TransactionID string    `json:"transactionID,omitempty" bson:"transactionID,omitempty"`
	OfferID       string    `json:"offerID,omitempty" bson:"offerID,omitempty"`
	AffiliateID   string    `json:"affiliateID,omitempty" bson:"affiliateID,omitempty"`
	AffiliateSub  string    `json:"affiliateSub,omitempty" bson:"affiliateSub,omitempty"`
	AffiliateSub2 string    `json:"affiliateSub2,omitempty" bson:"affiliateSub2,omitempty"`
}

type Transactions struct {
	TransactionID string    `bson:"transactionID,omitempty"`
	UTCDate       time.Time `bson:"utcdate,omitempty"`
	Date          time.Time `bson:"date,omitempty"`
}

type Rotation struct {
	Key   string `bson:"key,omitempty"`
	Value string `bson:"value,omitempty"`
}
