package service

import (
	"context"

	"git.selly.red/Selly-Modules/audit"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/external/utils/format"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

// AuditInterface ...
type AuditInterface interface {
	// Create ...
	Create(target string, targetId string, data interface{}, msg string, action string)

	// All ...
	All(ctx context.Context, q mgquery.AppQuery) responsemodel.ResponseAuditAll
}

// auditImplement ...
type auditImplement struct {
	CurrentStaff externalauth.User
}

// Audit ...
func Audit(cs externalauth.User) AuditInterface {
	return &auditImplement{
		CurrentStaff: cs,
	}
}

// Create ...
func (s auditImplement) Create(target string, targetId string, data interface{}, msg string, action string) {
	p := audit.CreatePayload{
		Target:   target,
		TargetID: targetId,
		Action:   action,
		Data:     format.ToString(data),
		Message:  msg,
		Author: audit.CreatePayloadAuthor{
			ID:   s.CurrentStaff.ID,
			Name: s.CurrentStaff.Name,
			Type: constant.TypeStaffAudit,
		},
	}
	audit.GetInstance().Create(p)
}

// GetAll ...
func (s auditImplement) GetAll(q audit.AllQuery) {
	audit.GetInstance().All(q)
}

// All return audit
func (auditImplement) All(ctx context.Context, q mgquery.AppQuery) responsemodel.ResponseAuditAll {
	query := audit.AllQuery{
		Target:   q.Audit.Target,
		TargetID: q.Audit.TargetID,
		Limit:    q.Limit,
		Page:     q.Page,
		Sort:     q.SortInterface,
	}

	result, total := audit.GetInstance().All(query)
	return responsemodel.ResponseAuditAll{
		List:  result,
		Total: total,
	}

}
