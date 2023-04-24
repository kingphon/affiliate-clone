package service

import (
	"context"

	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateStatusByID ...
func (s clickImplement) UpdateStatusByID(ctx context.Context, id primitive.ObjectID, status string) error {
	var (
		d       = dao.Click()
		payload = bson.D{
			{"$set", bson.D{
				{"status", status},
				{"updatedAt", ptime.Now()},
			}},
		}
		cond = bson.D{{"_id", id}}
	)

	err := d.UpdateOneByCondition(ctx, cond, payload)
	return err
}
