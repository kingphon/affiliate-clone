package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/app/database"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SellerShareURLInterface ...
type SellerShareURLInterface interface {
	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.SellerShareURL)

	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.SellerShareURL)

	// InsertOne ...
	InsertOne(ctx context.Context, payload mgaffiliate.SellerShareURL) error

	// CountByCondition ...
	CountByCondition(ctx context.Context, cond interface{}) int64

	// InsertMany ...
	InsertMany(ctx context.Context, payload []interface{}) error
}

// SellerShareURL return seller share url dao
func SellerShareURL() SellerShareURLInterface {
	return sellerShareURLImplement{}
}

// sellerShareURLImplement ...
type sellerShareURLImplement struct{}

// InsertMany ...
func (p sellerShareURLImplement) InsertMany(ctx context.Context, payload []interface{}) error {
	var col = database.SellerShareURLCol()

	_, err := col.InsertMany(ctx, payload)
	return err
}

// CountByCondition ...
func (p sellerShareURLImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	var col = database.SellerShareURLCol()

	total, _ := col.CountDocuments(ctx, cond)
	return total
}

// InsertOne ...
func (p sellerShareURLImplement) InsertOne(ctx context.Context, payload mgaffiliate.SellerShareURL) error {
	var (
		col = database.SellerShareURLCol()
	)

	_, err := col.InsertOne(ctx, payload)
	if err != nil {
		logger.Error("dao.SellerShareURL - InsertOne", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return err
}

// FindByCondition ...
func (p sellerShareURLImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.SellerShareURL) {
	var (
		col = database.SellerShareURLCol()
	)

	cursor, err := col.Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.SellerShareURL - FindByCondition cursor", logger.LogData{
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
		logger.Error("dao.SellerShareURL - FindByCondition decode", logger.LogData{
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
func (p sellerShareURLImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.SellerShareURL) {
	var col = database.SellerShareURLCol()

	if err := col.FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.SellerShareURL - FindOneByCondition err", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
	}
	return
}
