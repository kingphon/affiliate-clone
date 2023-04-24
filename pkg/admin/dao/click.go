package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// ClickInterface ...
type ClickInterface interface {
	// InsertOne ...
	InsertOne(ctx context.Context, payload interface{}) (err error)

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Click)

	// CountByCondition ...
	CountByCondition(ctx context.Context, cond interface{}) int64

	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Click)

	// UpdateOneByCondition ...
	UpdateOneByCondition(ctx context.Context, cond, payload interface{}) error

	// AggregateStatistic ...
	AggregateStatistic(ctx context.Context, cond interface{}) mgaffiliate.ClickStatistic

	// BulkWrite ...
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error
}

// clickImplement ...
type clickImplement struct{}

// Click ...
func Click() ClickInterface {
	return clickImplement{}
}

// AggregateStatistic ...
func (c clickImplement) AggregateStatistic(ctx context.Context, cond interface{}) mgaffiliate.ClickStatistic {
	var (
		data = make([]mgaffiliate.ClickStatistic, 0)
		col  = database.ClickCol()
	)
	match := bson.M{
		"$match": cond,
	}
	group := bson.M{
		"$group": bson.M{
			"_id": "",
			"total": bson.M{
				"$sum": 1,
			},
			"totalPending": bson.M{
				"$sum": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []string{"$status", constant.CampaignTransactionClickStatusPending}},
						1,
						0,
					},
				},
			},
			"totalCompleted": bson.M{
				"$sum": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []string{"$status", constant.CampaignTransactionClickStatusCompleted}},
						1,
						0,
					},
				},
			},
		},
	}
	cursor, err := col.Aggregate(ctx, []bson.M{match, group})

	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &data); err != nil {
		fmt.Println(err.Error())
		return mgaffiliate.ClickStatistic{}
	}
	if len(data) > 0 {
		return data[0]
	}
	return mgaffiliate.ClickStatistic{}
}

// InsertOne ...
func (c clickImplement) InsertOne(ctx context.Context, payload interface{}) (err error) {
	var col = database.ClickCol()

	if _, err = col.InsertOne(ctx, payload); err != nil {
		logger.Error("dao.Click - InsertOne", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return
}

// UpdateOneByCondition ...
func (c clickImplement) UpdateOneByCondition(ctx context.Context, cond, payload interface{}) (err error) {
	var col = database.ClickCol()

	if _, err = col.UpdateOne(ctx, cond, payload); err != nil {
		logger.Error("dao.Click - UpdateOneByCondition", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return
}

// BulkWrite ...
func (c clickImplement) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error {
	var col = database.ClickCol()

	res, err := col.BulkWrite(ctx, models, opts...)
	if err != nil {
		logger.Error("SocialPostDAO.BulkWrite result: ", logger.LogData{
			Data: bson.M{
				"res":   res,
				"error": err.Error(),
			},
		})
	}
	return err
}
