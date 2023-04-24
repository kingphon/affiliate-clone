package dao

import (
	"context"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// SellerCampaignStatisticInterface ...
type SellerCampaignStatisticInterface interface {
	// InsertOne ...
	InsertOne(ctx context.Context, payload interface{}) (err error)

	// UpdateByCondition ...
	UpdateByCondition(ctx context.Context, cond, payload interface{}, opts ...*options.UpdateOptions) (err error)
}

// sellerCampaignStatisticImplement ...
type sellerCampaignStatisticImplement struct{}

// SellerCampaignStatistic ...
func SellerCampaignStatistic() SellerCampaignStatisticInterface {
	return sellerCampaignStatisticImplement{}
}

// UpdateByCondition ...
func (s sellerCampaignStatisticImplement) UpdateByCondition(ctx context.Context, cond, payload interface{}, opts ...*options.UpdateOptions) (err error) {
	var col = database.SellerCampaignStatisticCol()
	_, err = col.UpdateOne(ctx, cond, payload, opts...)
	return
}
