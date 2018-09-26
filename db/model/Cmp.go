package v1

import "time"

type CMPAPIReport struct {
	UtcDate     time.Time `json:"utcdate" bson:"utcdate"`
	Minute      int       `json:"minute" bson:"minute"`
	Hour        int       `json:"hour" bson:"hour"`
	Date        string    `json:"date,omitempty" bson:"date,omitempty"`
	AffiliateID string    `json:"affiliate_id,omitempty" bson:"affiliate_id,omitempty"`
	OfferID     string    `json:"offer_id,omitempty" bson:"offer_id,omitempty"`

	//standard stats
	Impressions          int `json:"impressions" bson:"impressions,omitempty"`
	DuplicateImpressions int `json:"duplicate_impressions,omitempty" bson:"duplicate_impressions,omitempty"`
	UniqueImpressions    int `json:"unique_impressions,omitempty" bson:"unique_impressions,omitempty"`
	InvalidImpressions   int `json:"invalid_impressions,omitempty" bson:"invalid_impressions,omitempty"`

	BannerClicks          int `json:"banner_clicks" bson:"banner_clicks,omitempty"`
	DuplicateBannerClicks int `json:"duplicate_banner_clicks,omitempty" bson:"duplicate_banner_clicks,omitempty"`
	UniqueBannerClicks    int `json:"unique_banner_clicks,omitempty" bson:"unique_banner_clicks,omitempty"`
	InvalidBannerClicks   int `json:"invalid_banner_clicks,omitempty" bson:"invalid_banner_clicks,omitempty"`

	LandingPageViews          int `json:"landing_page_views" bson:"landing_page_views,omitempty"`
	DuplicateLandingPageViews int `json:"duplicate_landing_page_views,omitempty" bson:"duplicate_landing_page_views,omitempty"`
	UniqueLandingPageViews    int `json:"unique_landing_page_views,omitempty" bson:"unique_landing_page_views,omitempty"`
	InvalidLandingPageViews   int `json:"invalid_landing_page_views,omitempty" bson:"invalid_landing_page_views,omitempty"`

	Engagments          int `json:"engagments" bson:"engagments,omitempty"`
	DuplicateEngagments int `json:"duplicate_engagments,omitempty" bson:"duplicate_engagments,omitempty"`
	UniqueEngagments    int `json:"unique_engagments,omitempty" bson:"unique_engagments,omitempty"`
	InvalidEngagments   int `json:"invalid_engagments,omitempty" bson:"invalid_engagments,omitempty"`

	Conversions               int `json:"conversions" bson:"conversions,omitempty"`
	SentConversions           int `json:"sent_conversions" bson:"sent_conversions,omitempty"`
	UnSentConversions         int `json:"unsent_conversions" bson:"unsent_conversions,omitempty"`
	SentRotatedConversionsFwd int `json:"sent_rotated_conversions_fwd" bson:"sent_rotated_conversions_fwd,omitempty"`
	InvalidConversions        int `json:"invalid_conversions,omitempty" bson:"invalid_conversions,omitempty"`

	ContentViews int `json:"content_views" bson:"content_views,omitempty"`

	//rotation details
	RotatedClicksFwd            int                  `json:"rotated_clicks_fwd" bson:"rotated_clicks_fwd,omitempty"`
	UnSentRotatedConversionsFwd int                  `json:"unsent_rotated_conversions_fwd" bson:"unsent_rotated_conversions_fwd,omitempty"`
	RotatedConversionsFwd       int                  `json:"rotated_conversions_fwd" bson:"rotated_conversions_fwd,omitempty"`
	Rotations                   []CMPRotationDetails `json:"rotations,omitempty" bson:"rotations,omitempty"`
	ReceivedRotations           []CMPRotationDetails `json:"recv_rotations,omitempty" bson:"recv_rotations,omitempty"`

	//operator details
	OperatorReportDetails []CMPOperatorAPIReport `json:"operator_report_details,omitempty" bson:"operator_report_details,omitempty"`
}

type CMPOperatorAPIReport struct {
	OperatorID string `json:"operator_id,omitempty" bson:"operator_id,omitempty"`
	//subscription details
	UnSubscriptions []UnSubscriptionDetails `json:"un_subscriptions" bson:"un_subscriptions,omitempty"`
	Subscriptions   int                     `json:"subscriptions" bson:"subscriptions,omitempty"`
	MTSent          int                     `json:"mt_sent" bson:"mt_sent,omitempty"`
	MTSucess        int                     `json:"mt_success" bson:"mt_success,omitempty"`
	MTFail          int                     `json:"mt_fail" bson:"mt_fail,omitempty"`
	MTUnknown       int                     `json:"mt_unknown" bson:"mt_unknown,omitempty"`
}

type CMPRotationDetails struct {
	OfferID           string `json:"offer_id,omitempty" bson:"offer_id,omitempty"`
	AffiliateID       string `json:"affiliate_id,omitempty" bson:"affiliate_id,omitempty"`
	Clicks            int    `json:"clicks,omitempty" bson:"clicks,omitempty"`
	Conversions       int    `json:"conversions,omitempty" bson:"conversions,omitempty"`
	SentConversions   int    `json:"sent_conversions,omitempty" bson:"sent_conversions,omitempty"`
	UnSentConversions int    `json:"unsent_conversions,omitempty" bson:"unsent_conversions,omitempty"`
}

type UnSubscriptionDetails struct {
	Day   string `json:"day" bson:"day,omitempty"`
	Count int    `json:"count" bson:"count,omitempty"`
}
