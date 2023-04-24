package dao

import (
	"context"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionHistoryInterface ...
type TransactionHistoryInterface interface {
	// InsertOne ...
	InsertOne(ctx context.Context, payload interface{}) (err error)

	// InsertMany ...
	InsertMany(ctx context.Context, payload []interface{}) (err error)

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.TransactionHistory)
}

// transactionHistoryImplement ...
type transactionHistoryImplement struct{}

// TransactionHistory ...
func TransactionHistory() TransactionHistoryInterface {
	return transactionHistoryImplement{}
}
