package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"
)

// UpdateOneByCondition ...
func (d transactionImplement) UpdateOneByCondition(ctx context.Context, cond, payload interface{}) error {
	var dao = database.TransactionCol()
	_, err := dao.UpdateOne(ctx, cond, payload)
	return err
}

func (d transactionImplement) UpdateManyByCondition(ctx context.Context, cond, payload interface{}) error {
	var dao = database.TransactionCol()
	_, err := dao.UpdateMany(ctx, cond, payload)
	return err
}

// UpsertOneByCondition ...
func (d transactionImplement) UpsertOneByCondition(ctx context.Context, cond, payload interface{}) error {
	var dao = database.TransactionCol()
	fmt.Println("cond in dao: ", cond)
	fmt.Println("payload in dao: ", payload)
	_, err := dao.UpdateOne(ctx, cond, payload, options.Update().SetUpsert(true))
	return err
}
