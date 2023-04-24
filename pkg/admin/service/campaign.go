package service

import (
	"context"
	"fmt"
	"time"

	"git.selly.red/Selly-Modules/logger"

	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"

	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"

	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
)

// CampaignInterface ...
type CampaignInterface interface {
	//	CreateWithClientData create new campaign with client data ...
	CreateWithClientData(ctx context.Context, payload requestmodel.CampaignCreate) (campaignID string, err error)

	Update(ctx context.Context, id primitive.ObjectID, payload requestmodel.CampaignUpdate) (campaignID string, err error)

	// All get all campaigns ...
	All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseCampaignAll)

	ChangeStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.CampaignChangeStatus) (result responsemodel.ResponseChangeStatus, err error)

	// Detail get detail campaign by id ...
	Detail(ctx context.Context, id primitive.ObjectID) (result *responsemodel.ResponseCampaignDetail, err error)

	// GetPlatformByCampaign ...
	GetPlatformByCampaign(ctx context.Context, id primitive.ObjectID) (result responsemodel.ResponsePlatformInfos)

	// GetCampaignByIDs ...
	GetCampaignByIDs(ctx context.Context, campaignIds []primitive.ObjectID) []responsemodel.ResponseCampaignShort

	// CreatePlatformWithClientData ...
	CreatePlatformWithClientData(ctx context.Context, campaignID primitive.ObjectID, payload requestmodel.PlatformCreate) (platformID string, err error)

	// UpdateStatisticListCampaign ...
	UpdateStatisticListCampaign(ctx context.Context, campaignIds []primitive.ObjectID)

	// CheckAndUpdateStatusCampaign ..
	CheckAndUpdateStatusCampaign(ctx context.Context, campaignID primitive.ObjectID)

	// GetShortInfoByID ...
	GetShortInfoByID(ctx context.Context, id primitive.ObjectID) responsemodel.ResponseCampaignShortInfo
}

// campaignImplement ...
type campaignImplement struct {
	CurrentStaff    externalauth.User
	PlatformService platformImplement
}

// CheckAndUpdateStatusCampaign ...
func (s campaignImplement) CheckAndUpdateStatusCampaign(ctx context.Context, campaignID primitive.ObjectID) {

	// Check isExistCampaignInactive
	campaign := dao.Campaign().FindOneByCondition(ctx, bson.M{
		"_id": campaignID,
	})
	if campaign.Status == constant.CampaignStatusInActive {
		return
	}

	var (
		platformSvc = platformImplement{}
		cond        = bson.D{
			{"campaignId", campaignID},
			{"status", constant.PlatformStatusActive},
		}
	)
	if count := platformSvc.CountByCondition(ctx, cond); count > 0 {
		return
	}

	var (
		payloadUpdate = bson.M{
			"status":    constant.CampaignStatusInActive,
			"updatedAt": ptime.Now(),
		}
	)

	dao.Campaign().UpdateOneByCondition(ctx, bson.M{"_id": campaignID}, bson.M{"$set": payloadUpdate})

	// audit
	auditSvc := Audit(s.CurrentStaff)
	go auditSvc.Create(
		constant.AuditTargetCampaign,
		campaignID.Hex(),
		payloadUpdate,
		constant.MsgChangeStatusInactiveAffiliateCampaignWhenInactiveAllPlatform,
		constant.AuditActionEdit,
	)

	return
}

// Campaign ...
func Campaign(cs externalauth.User) CampaignInterface {
	return campaignImplement{
		CurrentStaff: cs,
	}
}

// UpdateStatisticListCampaign ...
func (s campaignImplement) UpdateStatisticListCampaign(ctx context.Context, campaignIds []primitive.ObjectID) {
	if len(campaignIds) == 0 {
		return
	}
	cond := bson.M{
		"campaignId": bson.M{
			"$in": campaignIds,
		},
	}
	fmt.Println("UpdateStatisticListCampaign ------------------------")
	results := dao.Transaction().AggregateCampaignStatistic(ctx, cond)
	fmt.Println("Result : ", results)
	if len(results) == 0 {
		return
	}

	var (
		wModel []mongo.WriteModel
	)
	for _, result := range results {
		update := bson.M{
			"$set": bson.M{
				"statistic.rewardTotal": result.TotalTransactionApproved + result.TotalTransactionCashback + result.TotalTransactionPending,
				"updatedAt":             time.Now(),
			},
		}
		wModel = append(wModel, mongo.NewUpdateOneModel().SetFilter(bson.M{
			"_id": result.CampaignID,
		}).SetUpdate(update))
	}

	if err := dao.Campaign().BulkWrite(ctx, wModel); err != nil {
		logger.Error("dao.Transaction().BulkWrite", logger.LogData{
			Data: bson.M{
				"error": err.Error(),
			},
		})
	}
}

// FindByID ...
func (campaignImplement) FindByID(ctx context.Context, id primitive.ObjectID) (result mgaffiliate.Campaign, err error) {
	var (
		d    = dao.Campaign()
		cond = bson.D{{"_id", id}}
	)

	campaign := d.FindOneByCondition(ctx, cond)
	if campaign.ID.IsZero() {
		err = errors.New(errorcode.CampaignNotFound)
		return
	}

	result = campaign
	return
}

// Update ...
func (s campaignImplement) Update(ctx context.Context, id primitive.ObjectID, payload requestmodel.CampaignUpdate) (campaignID string, err error) {
	campaign, err := s.FindByID(ctx, id)
	if err != nil {
		return
	}

	var (
		data          = payload.ConvertToBSON()
		payloadUpdate = bson.M{
			"name":                 data.Name,
			"searchString":         data.SearchString,
			"logo":                 data.Logo,
			"covers":               data.Covers,
			"desc":                 data.Desc,
			"order":                data.Order,
			"updatedAt":            data.UpdatedAt,
			"commission":           data.Commission,
			"from":                 data.From,
			"to":                   data.To,
			"estimateCashback":     data.EstimateCashback,
			"shareDesc":            data.ShareDesc,
			"allowShowShareAction": data.AllowShowShareAction,
		}
		cond     = bson.D{{"_id", campaign.ID}}
		auditSvc = Audit(s.CurrentStaff)
	)

	if err = dao.Campaign().UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		return
	}

	campaignID = campaign.ID.Hex()

	// audit
	go auditSvc.Create(
		constant.AuditTargetCampaign,
		campaignID,
		payloadUpdate,
		constant.MsgEditAffiliateCampaign,
		constant.AuditActionEdit,
	)

	return
}

// ChangeStatus ...
func (s campaignImplement) ChangeStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.CampaignChangeStatus) (result responsemodel.ResponseChangeStatus, err error) {
	campaign, err := s.FindByID(ctx, id)
	if err != nil {
		return
	}

	// Pre condition: campaign
	if payload.Status == constant.CampaignStatusActive {
		// Check platforms
		var platformCond = bson.D{
			{"campaignId", campaign.ID},
			{"status", constant.PlatformStatusActive},
		}

		if count := s.PlatformService.CountByCondition(ctx, platformCond); count == 0 {
			err = errors.New(errorcode.CampaignErrorWhenUpdateStatus)
			return
		}
	}

	var (
		payloadUpdate = bson.M{
			"status":    payload.Status,
			"updatedAt": ptime.Now(),
		}
		cond     = bson.D{{"_id", campaign.ID}}
		auditSvc = Audit(s.CurrentStaff)
	)

	if err = dao.Campaign().UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		return
	}

	// audit
	go auditSvc.Create(
		constant.AuditTargetCampaign,
		id.Hex(),
		payloadUpdate,
		constant.MsgChangeStatusAffiliateCampaign,
		constant.AuditActionEdit,
	)

	return
}
