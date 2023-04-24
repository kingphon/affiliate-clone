package responsemodel

import (
	"time"

	"git.selly.red/Selly-Server/affiliate/external/constant"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseTransactionCrawlData ...
type ResponseTransactionCrawlData struct {
	Code int `json:"code"`
	Data struct {
		Data []ResponseTransactionCrawl `json:"data"`
	} `json:"data"`
	Message string `json:"message"`
}

// ResponseTransactionCrawl ...
type ResponseTransactionCrawl struct {
	ID               primitive.ObjectID                `bson:"_id" json:"_id"`
	Campaign         primitive.ObjectID                `bson:"campaign" json:"campaign"`
	ClickID          primitive.ObjectID                `bson:"clickId" json:"clickId"`
	Source           string                            `bson:"source" json:"source"`
	Status           string                            `bson:"status" json:"status"`
	Commission       float64                           `bson:"commission" json:"commission"`
	TransactionID    string                            `bson:"transactionId" json:"transactionId"`
	TransactionTime  time.Time                         `bson:"transactionTime" json:"transactionTime"`
	Category         string                            `bson:"category" json:"category"`
	ClickURL         string                            `bson:"clickURL" json:"clickURL"`
	RejectedReason   string                            `bson:"rejectedReason" json:"rejectedReason"`
	Device           ResponseTransactionCrawlDevice    `bson:"device" json:"device"`
	Products         []ResponseTransactionCrawlProduct `bson:"products" json:"products"`
	CreatedAt        time.Time                         `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time                         `bson:"updatedAt" json:"updatedAt"`
	UpdatedHash      string                            `bson:"updatedHash" json:"updatedHash"`
	PartnerAffId     string                            `bson:"partnerAffId" json:"partnerAffId"`
	Options          ResponseTransactionCrawlOptions   `bson:"options" json:"options"`
	IsFraud          bool                              `bson:"isFraud" json:"isFraud"`
	SellerID         primitive.ObjectID                `bson:"sellerId" json:"-"`
	ReconciliationID primitive.ObjectID                `bson:"reconciliationId" json:"-"`
}

// GetSourceTransaction ...
func (m ResponseTransactionCrawl) GetSourceTransaction() (source string) {
	switch m.Source {
	case constant.CampaignTransactionCrawlSourceAccessTradeSelly:
		source = constant.CampaignTransactionSourceAccessTrade
	case constant.CampaignTransactionSourceOtherSelly:
		source = constant.CampaignTransactionSourceOther
	}
	return
}

// ResponseTransactionCrawlOptions ...
type ResponseTransactionCrawlOptions struct {
	ATIsConfirmed bool `bson:"atIsConfirmed" json:"atIsConfirmed"`
}

// ResponseTransactionCrawlProduct ...
type ResponseTransactionCrawlProduct struct {
	ID         string  `bson:"id" json:"id"`
	Name       string  `bson:"name" json:"name"`
	Quantity   int     `bson:"quantity" json:"quantity"`
	Status     string  `bson:"status" json:"status"`
	Commission float64 `bson:"commission" json:"commission"`
}

// ResponseTransactionCrawlDevice ...
type ResponseTransactionCrawlDevice struct {
	Model     string `bson:"model" json:"model"`
	Type      string `bson:"type" json:"type"`
	Device    string `bson:"device" json:"device"`
	OS        string `bson:"os" json:"os"`
	Browser   string `bson:"browser" json:"browser"`
	Payload   string `bson:"payload" json:"payload"`
	UserAgent string `bson:"userAgent" json:"userAgent"`
	Brand     string `bson:"brand" json:"brand"`
	OsVersion string `bson:"osVersion" json:"osVersion"`
}

// AggregateTransactionTempGroupSeller ...
type AggregateTransactionTempGroupSeller struct {
	ID    primitive.ObjectID         `json:"_id" bson:"_id"`
	Temps []ResponseTransactionCrawl `json:"temps"`
}
