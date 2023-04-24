package dao

import (
	"context"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SellerStatisticInterface ...
type SellerStatisticInterface interface {
	// InsertOne ...
	InsertOne(ctx context.Context, payload interface{}) (err error)

	// UpdateByCondition ...
	UpdateByCondition(ctx context.Context, cond, payload interface{}, opts ...*options.UpdateOptions) (err error)
}

// sellerStatisticImplement ...
type sellerStatisticImplement struct{}

// SellerStatistic ...
func SellerStatistic() SellerStatisticInterface {
	return sellerStatisticImplement{}
}

// UpdateByCondition ...
func (s sellerStatisticImplement) UpdateByCondition(ctx context.Context, cond, payload interface{}, opts ...*options.UpdateOptions) (err error) {
	var d = database.SellerStatisticCol()

	_, err = d.UpdateOne(ctx, cond, payload, opts...)
	return
}
