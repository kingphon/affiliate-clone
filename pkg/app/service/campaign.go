package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/request"

	"git.selly.red/Selly-Server/affiliate/pkg/app/locale"

	"git.selly.red/Selly-Server/affiliate/external/utils/parray"

	appconstant "git.selly.red/Selly-Server/affiliate/pkg/app/constant"

	"git.selly.red/Selly-Modules/natsio/client"
	natsmodel "git.selly.red/Selly-Modules/natsio/model"

	"git.selly.red/Selly-Server/affiliate/external/utils/file"

	"git.selly.red/Selly-Server/affiliate/internal/config"

	"git.selly.red/Selly-Modules/mongodb"

	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
	"github.com/friendsofgo/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/utils/pagetoken"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/pkg/app/dao"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/response"
)

// CampaignInterface ...
type CampaignInterface interface {
	// All return campaign with condition ...
	All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseList)

	// Detail ...
	Detail(ctx context.Context, id primitive.ObjectID, sellerID primitive.ObjectID) (result *responsemodel.ResponseDetail, err error)

	// GenerateShareURL ...
	GenerateShareURL(ctx context.Context, id primitive.ObjectID, userID string) (result *responsemodel.ResponseGenerateShareURL, err error)

	// FindCampaignActiveByID ...
	FindCampaignActiveByID(ctx context.Context, id primitive.ObjectID) mgaffiliate.Campaign

	// GetShortInfoByID ...
	GetShortInfoByID(ctx context.Context, id primitive.ObjectID) responsemodel.ResponseCampaignShortInfo

	// GetCampaignFilter ...
	GetCampaignFilter(ctx context.Context) responsemodel.ResponseCampaignFilter

	// GroupAll ...
	GroupAll(ctx context.Context, q mgquery.AppQuery) responsemodel.ResponseListCampaignGroupAll

	// GroupDetail ...
	GroupDetail(ctx context.Context, groupType string) *responsemodel.ResponseCampaignGroupBrief

	// GetSellerCampaignStatistic ...
	GetSellerCampaignStatistic(ctx context.Context, id, sellerID primitive.ObjectID) responsemodel.ResponseGetSellerCampaignStatistic

	// UpdateClick ...
	UpdateClick(ctx context.Context, id primitive.ObjectID, payload requestmodel.ClickUpdateBody) error
}

// Campaign return campaign service
func Campaign() CampaignInterface {
	return campaignImplement{}
}

// campaignImplement ...
type campaignImplement struct{}

// UpdateClick ...
func (s campaignImplement) UpdateClick(ctx context.Context, id primitive.ObjectID, payload requestmodel.ClickUpdateBody) error {
	return dao.Click().UpdateOne(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"finalDetected": mgaffiliate.FinalDetected{
				URL:   payload.FinalURL,
				Click: payload.PartnerClickID,
			},
			"updatedAt": time.Now(),
		},
	})
}

// GetSellerCampaignStatistic ...
func (s campaignImplement) GetSellerCampaignStatistic(ctx context.Context, id, sellerID primitive.ObjectID) responsemodel.ResponseGetSellerCampaignStatistic {
	var (
		d    = dao.SellerCampaignStatistic()
		cond = bson.D{
			{"sellerId", sellerID},
			{"campaignId", id},
		}
	)

	doc := d.FindOneByCondition(ctx, cond)
	return responsemodel.ResponseGetSellerCampaignStatistic{
		Data: responsemodel.SellerCampaignStatisticBrief{
			ID:         doc.ID,
			SellerID:   sellerID,
			CampaignID: id,
			Statistic: responsemodel.SellerCampaignStatistic{
				TransactionTotal: doc.Statistic.TransactionPending +
					doc.Statistic.TransactionApproved +
					doc.Statistic.TransactionCashback,
			},
		},
	}
}

// GroupDetail ...
func (s campaignImplement) GroupDetail(ctx context.Context, groupType string) *responsemodel.ResponseCampaignGroupBrief {
	switch groupType {
	case appconstant.CampaignGroupType.CampaignList.Key:
		var typeGroupCampaignList = appconstant.CampaignGroupType.CampaignList
		return &responsemodel.ResponseCampaignGroupBrief{
			ID:   typeGroupCampaignList.Key,
			Type: typeGroupCampaignList.Key,
			Name: typeGroupCampaignList.Title,
		}
	}
	return nil
}

// GroupAll ...
func (s campaignImplement) GroupAll(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseListCampaignGroupAll) {
	result.Data = make([]responsemodel.ResponseCampaignGroup, 0)

	// Check screen
	if isAllow := parray.ContainsStr(appconstant.ScreenListAllowCampaignGroup, q.Screen); !isAllow {
		return
	}

	// 2. Get items
	data := s.All(ctx, mgquery.AppQuery{Page: 0, Limit: 6, SortInterface: s.getSortCampaign("")})
	if data.List != nil {
		list := data.List.([]responsemodel.ResponseCampaignBrief)
		if len(list) == 0 {
			return
		}

		// 1. Init group campaign list
		var typeGroupCampaignList = appconstant.CampaignGroupType.CampaignList
		campaignGroup := responsemodel.ResponseCampaignGroup{
			ID:    typeGroupCampaignList.Key,
			Type:  typeGroupCampaignList.Key,
			Name:  typeGroupCampaignList.Title,
			Items: list,
		}

		// response
		result.Data = append(result.Data, campaignGroup)
	}

	return
}

// GetShortInfoByID ...
func (s campaignImplement) GetShortInfoByID(ctx context.Context, id primitive.ObjectID) responsemodel.ResponseCampaignShortInfo {
	var cond = bson.D{{"_id", id}}
	campaign := dao.Campaign().FindOneByCondition(ctx, cond)

	return responsemodel.ResponseCampaignShortInfo{
		ID:         campaign.ID,
		Name:       campaign.Name,
		Commission: campaign.Commission.Seller,
		Logo:       s.convertResponseFilePhoto(campaign.Logo).GetResponseData(),
	}
}

// FindCampaignActiveByID ...
func (s campaignImplement) FindCampaignActiveByID(ctx context.Context, id primitive.ObjectID) mgaffiliate.Campaign {
	var (
		d    = dao.Campaign()
		cond = bson.M{"_id": id, "status": constant.CampaignStatusActive}
	)

	campaign := d.FindOneByCondition(ctx, cond)
	return campaign
}

// GenerateShareURL ...
func (s campaignImplement) GenerateShareURL(ctx context.Context, id primitive.ObjectID, userIDStr string) (result *responsemodel.ResponseGenerateShareURL, err error) {
	// 1. check campaign id
	campaign := s.FindCampaignActiveByID(ctx, id)
	if campaign.ID.IsZero() || !campaign.IsAvailable() {
		err = errors.New(errorcode.CampaignNotFound)
		return
	}

	// 2. check platform id
	var (
		platformDao  = dao.Platform()
		condPlatform = bson.M{"campaignId": campaign.ID, "status": constant.PlatformStatusActive}
	)

	platforms := platformDao.FindByCondition(ctx, condPlatform)
	if len(platforms) == 0 {
		err = errors.New(errorcode.PlatformNotFound)
		return
	}

	// 3. check existed doc share url
	var (
		sellerShareURLDao  = dao.SellerShareURL()
		userID, _          = mongodb.NewIDFromString(userIDStr)
		condSellerShareURL = bson.D{
			{"sellerId", userID},
			{"campaignId", campaign.ID},
		}
	)

	shareDocs := sellerShareURLDao.FindByCondition(ctx, condSellerShareURL)
	if len(shareDocs) == 0 {
		// 4. generate new share url
		shareDocs, err = s.insertNewSellerShareURlByData(ctx, campaign, platforms, userID)
		if err != nil {
			return
		}
	}

	// Response
	var responsePlatformShareURL = make([]responsemodel.PlatformGenerateShareURL, 0)
	for _, doc := range shareDocs {
		responsePlatformShareURL = append(responsePlatformShareURL, responsemodel.PlatformGenerateShareURL{
			ID:        doc.PlatformID,
			Platform:  doc.Platform,
			ShareURL:  s.generateShareURL(ctx, doc),
			Title:     s.getTitleShareURLByPlatform(ctx, doc.Platform),
			ShareCode: doc.Code,
			OpenVia:   appconstant.OpenViaInApp,
		})
	}

	result = &responsemodel.ResponseGenerateShareURL{
		CampaignID: campaign.ID,
		Platforms:  responsePlatformShareURL,
	}
	return
}

// Detail ...
func (s campaignImplement) Detail(ctx context.Context, id, sellerID primitive.ObjectID) (result *responsemodel.ResponseDetail, err error) {
	var (
		d    = dao.Campaign()
		cond = bson.D{
			{"_id", id},
			{"status", constant.CampaignStatusActive},
		}
	)

	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() || !doc.IsAvailable() {
		err = errors.New(errorcode.CampaignNotFound)
		return
	}

	result = &responsemodel.ResponseDetail{Data: s.detail(ctx, doc, sellerID)}
	return
}

// getSortCampaign ...
func (s campaignImplement) getSortCampaign(sortStr string) interface{} {
	switch sortStr {
	case constant.SortKeyCreatedAt:
		return bson.D{{"createdAt", -1}}
	case constant.SortKeyCommission:
		return bson.D{
			{"commission.seller", -1},
			{"createdAt", -1},
		}
	case constant.SortKeyRewardTotal:
		return bson.D{
			{"statistic.rewardTotal", -1},
			{"createdAt", -1},
		}
	default:
		return bson.D{
			{"order", -1},
			{"createdAt", -1},
		}
	}
}

// All ...
func (s campaignImplement) All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseList) {
	// 1. Init value
	var list = make([]responsemodel.ResponseCampaignBrief, 0)

	// 2. Assign condition
	var (
		now  = ptime.Now()
		cond = bson.D{
			{"status", constant.CampaignStatusActive},
			{"$or", []interface{}{
				bson.M{
					"from": bson.M{
						"$lte": now,
					},
					"to": bson.M{
						"$gte": now,
					},
				},
				bson.M{
					"$and": []bson.M{
						{"from": bson.M{"$exists": false}},
						{"to": bson.M{"$exists": false}},
					},
				},
				bson.M{
					"$and": []bson.M{
						{"from": time.Time{}},
						{"from": time.Time{}},
					},
				},
			}},
		}
	)

	// Assign sort
	q.SortInterface = s.getSortCampaign(q.SortStr)

	// 3. Find
	var d = dao.Campaign()
	findOpts := q.GetFindOptionsWithPage()

	docs := d.FindByCondition(ctx, cond, findOpts)

	// 4. Convert response
	for _, doc := range docs {
		list = append(list, s.brief(ctx, doc))
	}

	// Page token
	endData := len(list) < int(q.Limit)
	var nextPageToken = ""
	if len(list) == int(q.Limit) {
		nextPageToken = pagetoken.PageTokenUsingPage(int(q.Page) + 1)
	}

	// Response
	result = responsemodel.ResponseList{
		List:          list,
		EndData:       endData,
		NextPageToken: nextPageToken,
	}
	return
}

//
// PRIVATE METHODS
//

// brief ...
func (s campaignImplement) brief(ctx context.Context, doc mgaffiliate.Campaign) responsemodel.ResponseCampaignBrief {
	return responsemodel.ResponseCampaignBrief{
		ID:         doc.ID,
		Name:       doc.Name,
		Commission: doc.Commission.Seller,
		Statistic: responsemodel.CampaignStatisticBrief{
			RewardTotal: doc.Statistic.RewardTotal,
		},
		Logo:                 s.convertResponseFilePhoto(doc.Logo).GetResponseData(),
		From:                 ptime.TimeResponseInit(doc.From),
		To:                   ptime.TimeResponseInit(doc.To),
		AllowShowShareAction: doc.AllowShowShareAction,
	}
}

// detail ...
func (s campaignImplement) detail(ctx context.Context, doc mgaffiliate.Campaign, sellerID primitive.ObjectID) *responsemodel.ResponseCampaignDetail {
	var (
		wg                = sync.WaitGroup{}
		platforms         = make([]responsemodel.PlatformBrief, 0)
		isCanJoinCampaign = true
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		platforms = Platform().GetPlatformBriefByCampaignID(ctx, doc.ID)
	}()

	if !sellerID.IsZero() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var condCount = bson.D{
				{"sellerId", sellerID},
				{"campaignId", doc.ID},
				{"from", constant.FromSeller},
				{"status", bson.M{
					"$ne": constant.TransactionStatus.Rejected.Key,
				}},
			}
			if count := Transaction().CountByCondition(ctx, condCount); count > 0 {
				isCanJoinCampaign = false
			}
		}()
	}

	wg.Wait()
	return &responsemodel.ResponseCampaignDetail{
		ID:        doc.ID,
		Name:      doc.Name,
		Logo:      fileResponseImplement{}.ConvertResponseFilePhoto(doc.Logo),
		Covers:    fileResponseImplement{}.ConvertResponseListFilePhoto(doc.Covers),
		Desc:      doc.Desc,
		ShareDesc: doc.ShareDesc,
		Statistic: responsemodel.CampaignStatisticBrief{
			RewardTotal: doc.Statistic.RewardTotal,
		},
		Commission:           doc.Commission.Seller,
		From:                 ptime.TimeResponseInit(doc.From),
		To:                   ptime.TimeResponseInit(doc.To),
		Platforms:            platforms,
		IsCanJoinCampaign:    isCanJoinCampaign,
		AllowShowShareAction: doc.AllowShowShareAction,
	}
}

// convertResponseFilePhoto ...
func (campaignImplement) convertResponseFilePhoto(f *mgaffiliate.FilePhoto) *file.FilePhoto {
	if f == nil {
		return nil
	}

	return &file.FilePhoto{
		ID:   f.ID.Hex(),
		Name: f.Name,
		Dimensions: &file.FileDimensions{
			Small: &file.FileSize{
				Width:  f.Dimensions.Small.Width,
				Height: f.Dimensions.Small.Height,
			},
			Medium: &file.FileSize{
				Width:  f.Dimensions.Medium.Width,
				Height: f.Dimensions.Medium.Height,
			},
		},
	}
}

// convertResponseListFilePhoto ...
func (s campaignImplement) convertResponseListFilePhoto(f []*mgaffiliate.FilePhoto) []*file.FilePhoto {
	var result = make([]*file.FilePhoto, 0)
	if len(f) == 0 {
		return result
	}

	for _, photo := range f {
		fs := s.convertResponseFilePhoto(photo)
		if fs != nil {
			result = append(result, fs)
		}
	}
	return result
}

// insertNewSellerShareURlByData ...
func (s campaignImplement) insertNewSellerShareURlByData(ctx context.Context, campaign mgaffiliate.Campaign, platforms []mgaffiliate.Platform, userID primitive.ObjectID) (result []mgaffiliate.SellerShareURL, err error) {
	// 1. init value
	var payload = make([]interface{}, 0)

	// 2. Get user info
	seller, err := client.GetSeller().GetSellerInfoByID(natsmodel.GetSellerByIDRequest{SellerID: userID})
	if err != nil {
		return
	}

	for _, platform := range platforms {
		raw := mgaffiliate.SellerShareURL{
			ID:         primitive.NewObjectID(),
			SellerID:   userID,
			CampaignID: campaign.ID,
			PlatformID: platform.ID,
			Code:       fmt.Sprintf("%s%s", seller.Code, platform.Code),
			CreatedAt:  time.Now(),
			Platform:   platform.PlatformType,
		}
		payload = append(payload, raw)
		result = append(result, raw)
	}

	// 3. insert
	var d = dao.SellerShareURL()
	err = d.InsertMany(ctx, payload)
	return
}

// generateShareURL ...
func (s campaignImplement) generateShareURL(ctx context.Context, doc mgaffiliate.SellerShareURL) string {
	var host = config.GetENV().HostShareURL
	return fmt.Sprintf("%s/aff/%s", host, doc.Code)
}

// GetCampaignFilter ...
func (s campaignImplement) GetCampaignFilter(ctx context.Context) responsemodel.ResponseCampaignFilter {
	data := []responsemodel.CampaignFilter{
		{
			Key:   constant.KeySort,
			Title: "Sắp xếp theo",
			List: []responsemodel.KeyValue{
				{
					Key:   constant.SortKeyCreatedAt,
					Value: "Mới nhất",
				},
				{
					Key:   constant.SortKeyCommission,
					Value: "Thưởng cao nhất",
				},
				{
					Key:   constant.SortKeyRewardTotal,
					Value: "Lượt thưởng cao nhất",
				},
			},
		},
	}
	return responsemodel.ResponseCampaignFilter{
		Data: data,
	}
}

// getTitleShareURLByPlatform ...
func (s campaignImplement) getTitleShareURLByPlatform(ctx context.Context, platform string) (title string) {
	switch platform {
	case constant.PlatformIOS:
		title = locale.TitlePlatformIOS
	case constant.PlatformAndroid:
		title = locale.TitlePlatformAndroid
	case constant.PlatformAll:
		title = locale.TitlePlatformAll
	case constant.PlatformWeb:
		title = locale.TitlePlatformWeb
	}

	return
}
