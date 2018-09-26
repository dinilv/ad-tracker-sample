package v1

type PostbackReq struct {
	IsConverted   bool
	AffiliateID   string
	OfferID       string
	OfferType     string
	TransactionID string
	ConversionIP  string
	SessionIP     string
	GoalID        string
	ClickURL      string
}

type PostbackRes struct {
	RequestParams       map[string]string
	Status              string
	ReceivedOfferID     string
	ReceivedAffiliateID string
	URL                 string
	ClickURL            string
	Activity            int
	IsPingRequired      bool
	SessionIP           string
}
