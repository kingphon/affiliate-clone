package service

import (
	"context"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/model/query"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

// TransactionInterface ...
type TransactionInterface interface {
	// Statistic ...
	Statistic(ctx context.Context, q mgquery.Affiliate) (result responsemodel.ResponseTransactionStatistic)

	// All ...
	All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseTransactionAll)

	// Detail ...
	Detail(ctx context.Context, id primitive.ObjectID) (result *responsemodel.ResponseTransactionDetail, err error)

	// FindByCode ...
	FindByCode(ctx context.Context, code string) (result mgaffiliate.Transaction)

	// InsertOne ...
	InsertOne(ctx context.Context, doc mgaffiliate.Transaction) (err error)

	// AggregateStatisticByCondition ...
	AggregateStatisticByCondition(ctx context.Context, cond interface{}) query.TransactionStatistic

	// GetHistoriesByTransaction ...
	GetHistoriesByTransaction(ctx context.Context, transactionID primitive.ObjectID) (result responsemodel.ResponseTransactionHistories)

	// GetByCondition get transaction by condition
	GetByCondition(ctx context.Context, q mgquery.AppQuery, doc mgaffiliate.Reconciliation) (result responsemodel.ResponseTransactionAll)

	// AggregateStatisticByReconciliationCondition ...
	AggregateStatisticByReconciliationCondition(ctx context.Context, doc mgaffiliate.Reconciliation, payload requestmodel.ReconciliationPayloadStatistic) responsemodel.ResponseReconciliationStatistic

	// GetByReconciliationCondition ...
	GetByReconciliationCondition(ctx context.Context, cond interface{}) (result []mgaffiliate.Transaction)

	// GetByQuery ...
	GetByQuery(ctx context.Context, q mgquery.AppQuery) (result []mgaffiliate.Transaction)
}

// transactionImplement ...
type transactionImplement struct {
}

// Transaction ...
func Transaction() TransactionInterface {
	return transactionImplement{}
}
