package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

// ReconciliationInterface ...
type ReconciliationInterface interface {
	// CreateWithClientData ...
	CreateWithClientData(ctx context.Context, payload requestmodel.ReconciliationCreate) (ReconciliationId string, err error)

	// All get all reconciliation
	All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseReconciliationAll)

	// Detail get detail reconciliation
	Detail(ctx context.Context, id primitive.ObjectID) (result *responsemodel.ResponseReconciliationDetail, err error)

	// FindByID get reconciliation by id
	FindByID(ctx context.Context, id primitive.ObjectID) (result mgaffiliate.Reconciliation, err error)

	// GetTransactionsByCondition  Get list transaction by condition ...
	GetTransactionsByCondition(ctx context.Context, q mgquery.AppQuery, id primitive.ObjectID) (result responsemodel.ResponseTransactionAll)

	// GetStatistic get statistics reconciliation
	GetStatistic(ctx context.Context, id primitive.ObjectID, payload requestmodel.ReconciliationPayloadStatistic) (result responsemodel.ResponseReconciliationStatistic, err error)

	// ChangeStatus ...
	ChangeStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.ReconciliationPayloadChangeStatus) (result responsemodel.ResponseChangeStatus, err error)
}

// ReconciliationImplement ...
type reconciliationImplement struct {
	currentStaff externalauth.User
}

// Reconciliation ...
func Reconciliation(cs externalauth.User) ReconciliationInterface {
	return &reconciliationImplement{
		currentStaff: cs,
	}
}
