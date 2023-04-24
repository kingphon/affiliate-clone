package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
)

// CreateWithClientData ...
func (s platformImplement) CreateWithClientData(ctx context.Context, payload []requestmodel.PlatformCreate, campaignID primitive.ObjectID) (err error) {
	var (
		d          = dao.Platform()
		createBody = make([]interface{}, 0)
		auditSvc   = Audit(s.CurrentStaff)
	)

	for _, platform := range payload {
		doc := platform.ConvertToBSON(campaignID)
		createBody = append(createBody, doc)

		// audit
		auditSvc.Create(
			constant.AuditTargetPlatform,
			doc.ID.Hex(),
			payload,
			constant.MsgCreateAffiliatePlatform,
			constant.AuditActionCreate,
		)
	}

	err = d.InsertMany(ctx, createBody)

	return
}

// CreateOnePlatform ...
func (s platformImplement) CreateOnePlatform(ctx context.Context, campaignID primitive.ObjectID, payload requestmodel.PlatformCreate) (platformId string, err error) {
	var (
		d        = dao.Platform()
		doc      = payload.ConvertToBSON(campaignID)
		auditSvc = Audit(s.CurrentStaff)
	)

	if err = d.InsertOne(ctx, doc); err != nil {
		return
	}

	platformId = doc.ID.Hex()

	// audit
	auditSvc.Create(
		constant.AuditTargetPlatform,
		platformId,
		payload,
		constant.MsgCreateAffiliatePlatform,
		constant.AuditActionCreate,
	)

	return
}
