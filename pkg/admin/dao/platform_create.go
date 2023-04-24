package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// InsertMany ...
func (d platformImplement) InsertMany(ctx context.Context, payload []interface{}) (err error) {
	var col = database.PlatformCol()

	if _, err := col.InsertMany(ctx, payload); err != nil {
		logger.Error("dao.Platform - InsertMany", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return
}

// InsertOne ...
func (d platformImplement) InsertOne(ctx context.Context, payload interface{}) (err error) {
	var (
		col = database.PlatformCol()
	)

	if _, err := col.InsertOne(ctx, payload); err != nil {
		logger.Error("dao.Platform - InsertOne", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return
}
