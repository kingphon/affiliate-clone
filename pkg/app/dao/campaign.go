package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/affiliate/pkg/app/database"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CampaignInterface ...
type CampaignInterface interface {
	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Campaign)

	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Campaign)
}

// Campaign return campaign dao
func Campaign() CampaignInterface {
	return campaignImplement{}
}

// campaignImplement ...
type campaignImplement struct{}

// FindOneByCondition ...
func (c campaignImplement) FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Campaign) {
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

// FindByCondition ...
func (c campaignImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Campaign) {
	var (
		col = database.CampaignCol()
	)

	cursor, err := col.Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.Campaign - FindByCondition cursor", logger.LogData{
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
		logger.Error("dao.Campaign - FindByCondition decode", logger.LogData{
			Data: bson.M{
				"cond":  cond,
				"opts":  opts,
				"error": err.Error(),
			},
		})
	}
	return
}
