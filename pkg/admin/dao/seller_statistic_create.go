package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// InsertOne ...
func (s sellerStatisticImplement) InsertOne(ctx context.Context, payload interface{}) (err error) {
	var col = database.SellerStatisticCol()

	if _, err = col.InsertOne(ctx, payload); err != nil {
		logger.Error("dao.SellerStatistic - InsertOne", logger.LogData{
			Data: bson.M{
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return
}
