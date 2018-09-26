package v1

import "time"

type RotationStack struct {
	OfferID     string    `bson:"offerID,omitempty"`
	AffiliateID string    `bson:"affiliateID,omitempty"`
	AddedDate   time.Time `bson:"addedDate,omitempty"`
	Event       string    `bson:"event,omitempty"`
}

type RotationGroupStack struct {
	RotatedFromOffer         string    `bson:"rotatedFromOffer,omitempty"`
	Country                  string    `bson:"country,omitempty"`
	Group                    string    `bson:"group,omitempty"`
	SelectedOfferForRotation string    `bson:"selectedOfferForRotation,omitempty"`
	AddedDate                time.Time `bson:"addedDate,omitempty"`
}
