package responsemodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseGetSellerCampaignStatistic ...
type ResponseGetSellerCampaignStatistic struct {
	Data SellerCampaignStatisticBrief `json:"data"`
}

// SellerCampaignStatisticBrief ...
type SellerCampaignStatisticBrief struct {
	ID         primitive.ObjectID      `json:"_id"`
	SellerID   primitive.ObjectID      `json:"sellerId"`
	CampaignID primitive.ObjectID      `json:"campaignId"`
	Statistic  SellerCampaignStatistic `json:"statistic"`
}

// SellerCampaignStatistic ...
type SellerCampaignStatistic struct {
	TransactionTotal int64 `bson:"transactionTotal"`
}
