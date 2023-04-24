package service

import (
	"context"

	"github.com/friendsofgo/errors"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
)

func (s reconciliationImplement) CreateWithClientData(ctx context.Context, payload requestmodel.ReconciliationCreate) (reconciliationId string, err error) {
	var (
		d   = dao.Reconciliation()
		doc = payload.ConvertToBSON(s.currentStaff)
	)

	// InsertOne
	if err = d.InsertOne(ctx, doc); err != nil {
		err = errors.New(errorcode.ReconciliationErrorWhenCreated)
		return
	}

	// Audit
	auditSvc := Audit(s.currentStaff)
	go auditSvc.Create(
		constant.AuditTargetReconciliation,
		doc.ID.Hex(),
		payload,
		constant.MsgCreateAffiliateReconciliation,
		constant.AuditActionCreate,
	)

	// Response
	reconciliationId = doc.ID.Hex()
	return
}
