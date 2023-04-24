package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/affiliate/pkg/app/database"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SellerStatisticInterface ...
type SellerStatisticInterface interface {
	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.SellerStatistic)

	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.SellerStatistic)
}

// SellerStatistic return seller statistic dao
func SellerStatistic() SellerStatisticInterface {
	return sellerStatisticImplement{}
}

// sellerStatisticImplement ...
type sellerStatisticImplement struct{}

// FindByCondition ...
func (s sellerStatisticImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.SellerStatistic) {
	var (
		col = database.SellerStatisticCol()
	)

	cursor, err := col.Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.SellerStatistic - FindByCondition cursor", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
		return
	}

	// Close cursor
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &docs); err != nil {
		logger.Error("dao.SellerStatistic - FindByCondition decode", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
	}
	return
}

// FindOneByCondition ...
func (s sellerStatisticImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.SellerStatistic) {
	var col = database.SellerStatisticCol()

	if err := col.FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.SellerStatistic - FindOneByCondition err", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
	}
	return
}
