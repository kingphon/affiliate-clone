package dao

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/model/query"
)

// FindByCondition ...
func (transactionImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Transaction) {
	var (
		col = database.TransactionCol()
	)

	cursor, err := col.Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.Transaction - FindByCondition cursor", logger.LogData{
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
		logger.Error("dao.Transaction - FindByCondition decode", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
	}
	return

}

// CountByCondition ...
func (transactionImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	var (
		col = database.TransactionCol()
	)

	total, _ := col.CountDocuments(ctx, cond)
	return total
}

// FindOneByCondition ...
func (transactionImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Transaction) {
	var col = database.TransactionCol()

	if err := col.FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.Transaction - FindOneByCondition err", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
	}
	return
}

// AggregateStatisticByReconciliationCondition ...
func (d transactionImplement) AggregateStatisticByReconciliationCondition(ctx context.Context, cond interface{}) (result query.TransactionStatisticDashboard) {
	var (
		col = database.TransactionCol()
	)
	match := bson.M{
		"$match": cond,
	}

	group := bson.M{
		"$group": bson.D{
			{"_id", ""},
			{"transactionTotal", bson.D{{"$sum", 1}}},
			{"commissionTotal", bson.D{{"$sum", "$commission"}}},
			{"sellerCommission", bson.D{{"$sum", "$sellerCommission"}}},
			{"sellyCommission", bson.D{{"$sum", "$sellyCommission"}}},
		},
	}

	cursor, err := col.Aggregate(ctx, []bson.M{match, group})

	var data []query.TransactionStatisticDashboard

	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &data); err != nil {
		fmt.Println(err.Error())
		return query.TransactionStatisticDashboard{}
	}

	if len(data) > 0 {
		return data[0]
	}
	return query.TransactionStatisticDashboard{}
}
