package dao

import (
	"context"
	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// UpdateOneByCondition ...
func (reconciliationImplement) UpdateOneByCondition(ctx context.Context, cond, payload interface{}) (err error) {
	_, err = database.ReconciliationCol().UpdateOne(ctx, cond, payload)
	if err != nil {
		logger.Error("dao.ReconciliationCol - UpdateOneByCondition err", logger.LogData{
			Data: bson.M{
				"cond":    cond,
				"payload": payload,
				"error":   err.Error(),
			},
		})
	}
	return
}
