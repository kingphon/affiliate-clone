package service

import (
	"context"

	"github.com/friendsofgo/errors"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
)

// CreateWithClientData ...
func (s campaignImplement) CreateWithClientData(ctx context.Context, payload requestmodel.CampaignCreate) (campaignID string, err error) {
	var (
		d        = dao.Campaign()
		doc      = payload.ConvertToBSON()
		auditSvc = Audit(s.CurrentStaff)
	)

	if err = d.InsertOne(ctx, doc); err != nil {
		return
	}

	// Insert platforms
	servicePlatform := Platform(s.CurrentStaff)
	if err = servicePlatform.CreateWithClientData(ctx, payload.Platforms, doc.ID); err != nil {
		return
	}

	// TODO: audit

	// Response
	campaignID = doc.ID.Hex()

	// audit
	go auditSvc.Create(
		constant.AuditTargetCampaign,
		doc.ID.Hex(),
		payload,
		constant.MsgCreateAffiliateCampaign,
		constant.AuditActionCreate,
	)

	return
}

// CreatePlatformWithClientData ...
func (s campaignImplement) CreatePlatformWithClientData(ctx context.Context, campaignID primitive.ObjectID, payload requestmodel.PlatformCreate) (platformID string, err error) {
	// 1. Check existed platform
	var (
		//d = dao.Campaign()

		condCheckExisted = bson.D{
			{"campaignId", campaignID},
			{"platform", payload.PlatformType},
		}
	)

	var platformService = platformImplement{}

	// Check existed platform
	if count := platformService.CountByCondition(ctx, condCheckExisted); count > 0 {
		err = errors.New(errorcode.PlatformIsExisted)
		return
	}

	platformID, err = platformService.CreateOnePlatform(ctx, campaignID, payload)
	return
}
