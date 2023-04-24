package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// FindOneByCondition ...
func (c clickImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Click) {
	var col = database.ClickCol()

	if err := col.FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.Click - FindOneByCondition err", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
	}
	return
}

// FindByCondition ...
func (clickImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Click) {
	var (
		col = database.ClickCol()
	)

	cursor, err := col.Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.Click - FindByCondition cursor", logger.LogData{
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
		logger.Error("dao.Click - FindByCondition decode", logger.LogData{
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
func (clickImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	var (
		col = database.ClickCol()
	)

	total, _ := col.CountDocuments(ctx, cond)
	return total
}
