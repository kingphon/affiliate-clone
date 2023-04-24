package service

import (
	"context"

	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"

	"go.mongodb.org/mongo-driver/bson/primitive"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
)

// TransactionHistoryInterface ...
type TransactionHistoryInterface interface {
	// GenerateAndInsertHistoriesByNewTransaction ...
	GenerateAndInsertHistoriesByNewTransaction(ctx context.Context, transactionNew mgaffiliate.Transaction) (err error)

	// GenerateAndInsertHistoriesByUpdateTransaction ...
	GenerateAndInsertHistoriesByUpdateTransaction(ctx context.Context, transaction mgaffiliate.Transaction) (err error)

	// GetHistoriesByTransactionID ...
	GetHistoriesByTransactionID(ctx context.Context, transactionID primitive.ObjectID) []responsemodel.ResponseTransactionHistory
}

// transactionHistoryImplement ...
type transactionHistoryImplement struct {
}

// TransactionHistory ...
func TransactionHistory() TransactionHistoryInterface {
	return transactionHistoryImplement{}
}
