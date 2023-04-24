package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
)

// ReconciliationInterface ...
type ReconciliationInterface interface {
	// InsertOne ...
	InsertOne(ctx context.Context, payload interface{}) (err error)

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Reconciliation)

	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Reconciliation)

	// CountByCondition ...
	CountByCondition(ctx context.Context, cond interface{}) int64

	// UpdateOneByCondition ...
	UpdateOneByCondition(ctx context.Context, cond, payload interface{}) (err error)
}

// reconciliationImplement ...
type reconciliationImplement struct {
}

// Reconciliation ...
func Reconciliation() ReconciliationInterface {
	return &reconciliationImplement{}
}
