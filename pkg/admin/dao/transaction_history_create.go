package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// InsertOne ...
func (t transactionHistoryImplement) InsertOne(ctx context.Context, payload interface{}) (err error) {
	var col = database.TransactionHistoryCol()

	if _, err = col.InsertOne(ctx, payload); err != nil {
		logger.Error("dao.TransactionHistory - InsertOne", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return
}

// InsertMany ...
func (t transactionHistoryImplement) InsertMany(ctx context.Context, payload []interface{}) (err error) {
	var col = database.TransactionHistoryCol()

	if _, err = col.InsertMany(ctx, payload); err != nil {
		logger.Error("dao.TransactionHistory - InsertMany", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return
}
