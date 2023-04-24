package dao

import (
	"context"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

func (reconciliationImplement) InsertOne(ctx context.Context, payload interface{}) (err error) {
	if _, err := database.ReconciliationCol().InsertOne(ctx, payload); err != nil {
		logger.Error("dao.Reconciliation - InsertOne", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return
}
