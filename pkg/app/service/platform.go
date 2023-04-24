package service

import (
	"context"

	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"

	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/request"

	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/response"

	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
	"github.com/friendsofgo/errors"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/pkg/app/dao"
	"go.mongodb.org/mongo-driver/bson"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlatformInterface ...
type PlatformInterface interface {
	// GenerateAffiliateLink ...
	GenerateAffiliateLink(ctx context.Context, device mgaffiliate.Device, payload requestmodel.GenerateAffiliateLinkBody) (result responsemodel.ResponseAffiliateLink, err error)

	// GetPlatformBriefByCampaignID ...
	GetPlatformBriefByCampaignID(ctx context.Context, campaignID primitive.ObjectID) (platforms []responsemodel.PlatformBrief)
}

// Platform return platform service
func Platform() PlatformInterface {
	return platformImplement{}
}

// platformImplement ...
type platformImplement struct{}

// GenerateAffiliateLink ...
func (s platformImplement) GenerateAffiliateLink(ctx context.Context, device mgaffiliate.Device, payload requestmodel.GenerateAffiliateLinkBody) (result responsemodel.ResponseAffiliateLink, err error) {
	// Check code
	var shareURLService = shareURLImplement{}
	shareURLDoc, err := shareURLService.GetSellerShareURLByCode(ctx, payload.Code)
	if err != nil {
		return
	}

	if shareURLDoc == nil {
		return
	}

	if shareURLDoc.PlatformID.Hex() != payload.PlatformID {
		err = errors.New(errorcode.PlatformNotFound)
		return
	}

	// 1. Check platform
	platform := s.FindPlatformActiveByID(ctx, shareURLDoc.PlatformID)
	if platform.ID.IsZero() {
		err = errors.New(errorcode.PlatformNotFound)
		return
	}

	// 2. Check campaign active
	var cpService = Campaign()
	campaign := cpService.FindCampaignActiveByID(ctx, platform.CampaignID)
	if campaign.ID.IsZero() {
		err = errors.New(errorcode.CampaignNotFound)
		return
	}

	// 3. Insert doc click
	var (
		clickService = Click()
		payloadClick = requestmodel.PayloadInsertClick{
			SellerID:           shareURLDoc.SellerID,
			CampaignID:         campaign.ID,
			Device:             device,
			Platform:           platform,
			CampaignCommission: campaign.Commission,
			From:               payload.From,
			ShareURL:           campaignImplement{}.generateShareURL(ctx, *shareURLDoc),
		}
	)

	clickID, err := clickService.GenerateDocAndInsert(ctx, payloadClick)
	if err != nil {
		return
	}

	// 4. Generate affiliate link
	affiliateLink := s.generateAffLink(ctx, platform, clickID)

	// 5. update link affiliate link  => doc click
	go func() {
		clickService.UpdateByID(context.Background(), clickID,
			bson.M{"$set": bson.M{
				"affiliateURL": affiliateLink,
				"updatedAt":    ptime.Now(),
			}})
	}()
	result.URL = affiliateLink
	result.ClickID = clickID.Hex()
	return
}

// FindPlatformActiveByID ...
func (platformImplement) FindPlatformActiveByID(ctx context.Context, platformID primitive.ObjectID) (result mgaffiliate.Platform) {
	var (
		d    = dao.Platform()
		cond = bson.D{
			{"_id", platformID},
			{"status", constant.PlatformStatusActive},
		}
	)

	result = d.FindOneByCondition(ctx, cond)
	return
}

// GetPlatformBriefByCampaignID ...
func (platformImplement) GetPlatformBriefByCampaignID(ctx context.Context, campaignID primitive.ObjectID) (platforms []responsemodel.PlatformBrief) {
	platforms = make([]responsemodel.PlatformBrief, 0)

	var (
		dao  = dao.Platform()
		cond = bson.D{
			{"campaignId", campaignID},
			{"status", constant.PlatformStatusActive}}
	)

	docs := dao.FindByCondition(ctx, cond)
	for _, doc := range docs {
		platforms = append(platforms, responsemodel.PlatformBrief{
			ID:       doc.ID,
			Code:     doc.Code,
			Platform: doc.PlatformType,
		})
	}

	return
}
