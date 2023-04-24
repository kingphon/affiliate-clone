package service

import (
	"context"

	"git.selly.red/Selly-Server/affiliate/external/constant"

	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"

	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"

	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"

	"go.mongodb.org/mongo-driver/bson/primitive"

	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
)

// PlatformInterface ...
type PlatformInterface interface {
	// CreateWithClientData create new platform with client data ...
	CreateWithClientData(ctx context.Context, payload []requestmodel.PlatformCreate, campaignID primitive.ObjectID) (err error)

	// Update ...
	Update(ctx context.Context, id primitive.ObjectID, payload requestmodel.PlatformCreate) (platformID string, err error)

	// ChangeStatus ...
	ChangeStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.PlatformChangeStatus) (result responsemodel.ResponseChangeStatus, err error)

	// GetByCampaign ...
	GetByCampaign(ctx context.Context, id primitive.ObjectID) (result []responsemodel.ResponsePlatformInfo)

	// CreateOnePlatform ...
	CreateOnePlatform(ctx context.Context, campaignID primitive.ObjectID, payload requestmodel.PlatformCreate) (platformId string, err error)

	// CountByCondition ...
	CountByCondition(ctx context.Context, cond interface{}) int64
}

// platformImplement ...
type platformImplement struct {
	CurrentStaff externalauth.User
}

// Platform ...
func Platform(cs externalauth.User) PlatformInterface {
	return platformImplement{CurrentStaff: cs}
}

// FindByID ...
func (platformImplement) FindByID(ctx context.Context, id primitive.ObjectID) (result mgaffiliate.Platform, err error) {
	var (
		d    = dao.Platform()
		cond = bson.D{{"_id", id}}
	)

	platform := d.FindOneByCondition(ctx, cond)
	if platform.ID.IsZero() {
		err = errors.New(errorcode.PlatformNotFound)
		return
	}

	result = platform
	return
}

// Update ...
func (s platformImplement) Update(ctx context.Context, id primitive.ObjectID, payload requestmodel.PlatformCreate) (platformID string, err error) {
	platform, err := s.FindByID(ctx, id)
	if err != nil {
		return
	}

	var (
		data          = payload.ConvertToBSON(platform.CampaignID)
		payloadUpdate = bson.M{
			"partner":   data.Partner,
			"url":       data.URL,
			"platform":  data.PlatformType,
			"updatedAt": data.UpdatedAt,
		}
		cond = bson.D{{"_id", platform.ID}}

		auditSvc = Audit(s.CurrentStaff)
	)

	if err = dao.Platform().UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		return
	}

	platformID = platform.ID.Hex()
	go s.updatePlatformForCampaign(platform.CampaignID)

	go auditSvc.Create(
		constant.AuditTargetPlatform,
		id.Hex(),
		payloadUpdate,
		constant.MsgEditAffiliatePlatform,
		constant.AuditActionEdit,
	)

	return
}

// ChangeStatus ...
func (s platformImplement) ChangeStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.PlatformChangeStatus) (result responsemodel.ResponseChangeStatus, err error) {
	platform, err := s.FindByID(ctx, id)
	if err != nil {
		return
	}

	var (
		payloadUpdate = bson.M{
			"status":    payload.Status,
			"updatedAt": ptime.Now(),
		}
		cond     = bson.D{{"_id", platform.ID}}
		auditSvc = Audit(s.CurrentStaff)
	)

	if err = dao.Platform().UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		return
	}

	// Check status platforms by campaignId and update status campaign
	if payload.Status == constant.PlatformStatusInActive {
		go campaignImplement{CurrentStaff: s.CurrentStaff}.CheckAndUpdateStatusCampaign(context.Background(), platform.CampaignID)
	}

	// Update field platforms in campaign document
	go s.updatePlatformForCampaign(platform.CampaignID)

	// audit
	go auditSvc.Create(
		constant.AuditTargetPlatform,
		id.Hex(),
		payloadUpdate,
		constant.MsgChangeStatusAffiliatePlatform,
		constant.AuditActionEdit,
	)
	return
}

//  updatePlatformForCampaign ...
func (s platformImplement) updatePlatformForCampaign(campaignID primitive.ObjectID) {
	var (
		ctx  = context.Background()
		d    = dao.Platform()
		cond = bson.D{
			{"campaignId", campaignID},
			{"status", constant.PlatformStatusActive},
		}
	)

	platforms, err := d.DistinctByCondition(ctx, "platform", cond)
	if err != nil {
		return
	}

	var (
		condCampaign = bson.D{{"_id", campaignID}}
		payload      = bson.M{
			"$set": bson.M{
				"platforms": platforms,
				"updatedAt": ptime.Now(),
			},
		}
	)
	dao.Campaign().UpdateOneByCondition(ctx, condCampaign, payload)

	// audit
	auditSvc := Audit(s.CurrentStaff)
	go auditSvc.Create(
		constant.AuditTargetPlatform,
		campaignID.Hex(),
		payload,
		constant.MsgEditAffiliateCampaign,
		constant.AuditActionEdit,
	)
}

// CountByCondition ...
func (platformImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	d := dao.Platform()
	total := d.CountByCondition(ctx, cond)
	return total
}
