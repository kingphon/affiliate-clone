package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
)

// UpdateByID ...
func (s transactionImplement) UpdateByID(ctx context.Context, id primitive.ObjectID, data bson.M) (err error) {
	var (
		d       = dao.Transaction()
		cond    = bson.D{{"_id", id}}
		payload = bson.M{
			"$set": data,
		}
	)
	err = d.UpdateOneByCondition(ctx, cond, payload)
	return
}

// UpsertByID ...
func (s transactionImplement) UpsertByID(ctx context.Context, id primitive.ObjectID, data interface{}) (err error) {
	var (
		d       = dao.Transaction()
		cond    = bson.M{"_id": id}
		payload = bson.M{
			"$set": data,
		}
	)

	err = d.UpsertOneByCondition(ctx, cond, payload)
	return
}
