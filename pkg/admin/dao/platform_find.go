package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// FindByCondition ...
func (platformImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Platform) {
	var (
		col = database.PlatformCol()
	)

	cursor, err := col.Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.Platform - FindByCondition cursor", logger.LogData{
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
		logger.Error("dao.Platform - FindByCondition decode", logger.LogData{
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
func (platformImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	var (
		col = database.PlatformCol()
	)

	total, _ := col.CountDocuments(ctx, cond)
	return total
}
