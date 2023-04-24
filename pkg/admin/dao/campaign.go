package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// CampaignInterface ...
type CampaignInterface interface {
	// InsertOne ...
	InsertOne(ctx context.Context, payload interface{}) (err error)

	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Campaign)

	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Campaign)

	UpdateOneByCondition(ctx context.Context, cond, payload interface{}) (err error)

	// CountByCondition ..
	CountByCondition(ctx context.Context, cond interface{}) int64
	// BulkWrite ...
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error
}

// campaignImplement ...
type campaignImplement struct{}

// Campaign ...
func Campaign() CampaignInterface {
	return campaignImplement{}
}

// BulkWrite ...
func (d campaignImplement) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error {
	_, err := database.CampaignCol().BulkWrite(ctx, models, opts...)
	return err
}

// UpdateOneByCondition ...
func (d campaignImplement) UpdateOneByCondition(ctx context.Context, cond, payload interface{}) (err error) {
	var col = database.CampaignCol()

	_, err = col.UpdateOne(ctx, cond, payload)
	return
}

// FindOneByCondition ...
func (d campaignImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Campaign) {
	var col = database.CampaignCol()

	if err := col.FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.Campaign - FindOneByCondition err", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
	}
	return
}
