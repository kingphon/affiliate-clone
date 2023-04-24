package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// PlatformInterface ...
type PlatformInterface interface {
	// InsertMany ...
	InsertMany(ctx context.Context, payload []interface{}) (err error)

	UpdateOneByCondition(ctx context.Context, cond, payload interface{}) (err error)

	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Platform)

	DistinctByCondition(ctx context.Context, fieldName string, cond interface{}) (result []interface{}, err error)

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Platform)

	// InsertOne ...
	InsertOne(ctx context.Context, payload interface{}) (err error)

	// CountByCondition ...
	CountByCondition(ctx context.Context, cond interface{}) int64
}

// platformImplement ...
type platformImplement struct{}

// Platform ...
func Platform() PlatformInterface {
	return platformImplement{}
}

// FindOneByCondition ...
func (d platformImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Platform) {
	var col = database.PlatformCol()

	if err := col.FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.Platform - FindOneByCondition err", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
	}
	return
}

// UpdateOneByCondition ...
func (d platformImplement) UpdateOneByCondition(ctx context.Context, cond, payload interface{}) (err error) {
	var col = database.PlatformCol()

	_, err = col.UpdateOne(ctx, cond, payload)
	return
}

// DistinctByCondition ...
func (d platformImplement) DistinctByCondition(ctx context.Context, fieldName string, cond interface{}) (result []interface{}, err error) {
	var col = database.PlatformCol()

	result, err = col.Distinct(ctx, fieldName, cond)
	return
}
