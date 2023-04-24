package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/app/database"
)

// ClickInterface ...
type ClickInterface interface {
	// InsertOne ...
	InsertOne(ctx context.Context, payload mgaffiliate.Click) error

	// UpdateOne ...
	UpdateOne(ctx context.Context, cond, payload interface{}) error
}

// Click return click dao
func Click() ClickInterface {
	return clickImplement{}
}

// clickImplement ...
type clickImplement struct{}

// UpdateOne ...
func (c clickImplement) UpdateOne(ctx context.Context, cond, payload interface{}) error {
	var (
		col = database.ClickCol()
	)

	_, err := col.UpdateOne(ctx, cond, payload)
	if err != nil {
		logger.Error("dao.Click - UpdateOne", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return err
}

func (c clickImplement) InsertOne(ctx context.Context, payload mgaffiliate.Click) error {
	var (
		col = database.ClickCol()
	)

	_, err := col.InsertOne(ctx, payload)
	if err != nil {
		logger.Error("dao.Click - InsertOne", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return err
}
