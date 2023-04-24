package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/affiliate/pkg/app/database"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionInterface ...
type TransactionInterface interface {
	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Transaction)

	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Transaction)

	// CountByCondition ...
	CountByCondition(ctx context.Context, cond interface{}) int64
}

// Transaction return transaction dao
func Transaction() TransactionInterface {
	return transactionImplement{}
}

// transactionImplement ...
type transactionImplement struct{}

// CountByCondition ...
func (t transactionImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	var col = database.TransactionCol()

	total, _ := col.CountDocuments(ctx, cond)
	return total
}

// FindByCondition ...
func (t transactionImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Transaction) {
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

// FindOneByCondition ...
func (t transactionImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Transaction) {
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
