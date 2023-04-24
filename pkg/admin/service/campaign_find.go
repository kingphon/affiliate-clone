package service

import (
	"context"
	"strings"
	"sync"

	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

//
// PUBLIC METHOD
//

// All ...
func (s campaignImplement) All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseCampaignAll) {
	var (
		d  = dao.Campaign()
		wg = sync.WaitGroup{}
	)

	// Assign conditions
	cond := bson.D{}
	q.Affiliate.AssignKeyword(&cond)
	q.Affiliate.AssignStatus(&cond)
	q.Affiliate.AssignCreatedAt(&cond)

	// Find
	wg.Add(2)

	go func() {
		defer wg.Done()

		// prepare data
		result.List = make([]responsemodel.ResponseCampaignBrief, 0)

		// Find options
		q = s.sortString(ctx, q)
		findOpts := q.GetFindOptionsWithPage()
		findOpts.SetProjection(bson.M{
			"_id":              1,
			"name":             1,
			"logo":             1,
			"desc":             1,
			"commission":       1,
			"estimateCashback": 1,
			"platforms":        1,
			"order":            1,
			"status":           1,
			"createdAt":        1,
		})

		docs := d.FindByCondition(ctx, cond, findOpts)
		for _, doc := range docs {
			result.List = append(result.List, s.brief(ctx, doc))
		}
	}()

	// Assign total
	go func() {
		defer wg.Done()

		result.Total = d.CountByCondition(ctx, cond)
	}()

	wg.Wait()

	//Assign limit
	result.Limit = q.Limit

	return
}

func (s campaignImplement) sortString(ctx context.Context, q mgquery.AppQuery) mgquery.AppQuery {
	if q.SortStr != "" {
		var (
			keySort   string
			valueSort int64
		)

		// Get value sort in queryParam
		qSort := strings.Trim(q.SortStr, "-")
		switch qSort {
		case "createdAt":
			keySort = "createdAt"
		case "real":
			keySort = "commission.real"
		default:
			keySort = "createdAt"
		}

		// Get first character in string
		valueSort = -1
		if !strings.HasPrefix(q.SortStr, "-") && q.SortStr != "" {
			valueSort = 1
		}

		q.SortInterface = bson.D{
			{keySort, valueSort},
			{"_id", -1},
		}
	}

	return q
}

// Detail ...
func (s campaignImplement) Detail(ctx context.Context, id primitive.ObjectID) (result *responsemodel.ResponseCampaignDetail, err error) {
	var (
		d    = dao.Campaign()
		cond = bson.M{"_id": id}
	)

	campaign := d.FindOneByCondition(ctx, cond)
	if campaign.ID.IsZero() {
		return nil, errors.New(errorcode.CampaignNotFound)
	}

	result = s.detail(ctx, campaign)

	return
}

// GetPlatformByCampaign ...
func (s campaignImplement) GetPlatformByCampaign(ctx context.Context, id primitive.ObjectID) (result responsemodel.ResponsePlatformInfos) {
	var (
		platformSvc = Platform(externalauth.User{})
	)

	result.Data = platformSvc.GetByCampaign(ctx, id)
	return
}

// GetCampaignByIDs ...
func (s campaignImplement) GetCampaignByIDs(ctx context.Context, campaignIds []primitive.ObjectID) []responsemodel.ResponseCampaignShort {
	var (
		d    = dao.Campaign()
		cond = bson.M{
			"_id": bson.M{"$in": campaignIds},
		}
	)

	// List campaign
	campaigns := d.FindByCondition(ctx, cond)

	// List campaign short
	campaignsShort := make([]responsemodel.ResponseCampaignShort, 0)
	for _, doc := range campaigns {
		cs := s.short(ctx, doc)
		campaignsShort = append(campaignsShort, cs)
	}

	return campaignsShort
}

// GetShortInfoByID ...
func (s campaignImplement) GetShortInfoByID(ctx context.Context, id primitive.ObjectID) responsemodel.ResponseCampaignShortInfo {
	doc := dao.Campaign().FindOneByCondition(ctx, bson.M{"_id": id})

	return responsemodel.ResponseCampaignShortInfo{
		ID:   doc.ID.Hex(),
		Name: doc.Name,
	}
}

//
// PRIVATE METHOD
//

// brief ...
func (s campaignImplement) brief(ctx context.Context, doc mgaffiliate.Campaign) responsemodel.ResponseCampaignBrief {
	var (
		fileResSvc = fileResponseImplement{}
	)

	return responsemodel.ResponseCampaignBrief{
		ID:   doc.ID.Hex(),
		Name: doc.Name,
		Logo: fileResSvc.ConvertResponseFilePhoto(doc.Logo).GetResponseData(),
		Desc: doc.Desc,
		Commission: responsemodel.ResponseCampaignCommission{
			Real:          doc.Commission.Real,
			SellerPercent: doc.Commission.SellerPercent,
			Selly:         doc.Commission.Selly,
			Seller:        doc.Commission.Seller,
		},
		EstimateCashback: responsemodel.ResponseCampaignEstimateCashback{
			Day:       doc.EstimateCashback.Day,
			NextMonth: doc.EstimateCashback.NextMonth,
		},
		Platforms:            doc.Platforms,
		Order:                doc.Order,
		Status:               doc.Status,
		AllowShowShareAction: doc.AllowShowShareAction,
		CreatedAt:            ptime.TimeResponseInit(doc.CreatedAt),
	}
}

// detail
func (s campaignImplement) detail(ctx context.Context, doc mgaffiliate.Campaign) *responsemodel.ResponseCampaignDetail {
	var (
		fileResSvc = fileResponseImplement{}
	)

	return &responsemodel.ResponseCampaignDetail{
		ID:     doc.ID.Hex(),
		Name:   doc.Name,
		Logo:   fileResSvc.ConvertResponseFilePhoto(doc.Logo).GetResponseData(),
		Covers: fileResSvc.ConvertResponseListFilePhoto(doc.Covers),
		Desc:   doc.Desc,
		Commission: responsemodel.ResponseCampaignCommission{
			Real:          doc.Commission.Real,
			SellerPercent: doc.Commission.SellerPercent,
			Selly:         doc.Commission.Selly,
			Seller:        doc.Commission.Seller,
		},
		From:   ptime.TimeResponseInit(doc.From),
		To:     ptime.TimeResponseInit(doc.To),
		Status: doc.Status,
		EstimateCashback: responsemodel.ResponseCampaignEstimateCashback{
			Day:       doc.EstimateCashback.Day,
			NextMonth: doc.EstimateCashback.NextMonth,
		},
		CreatedAt:            ptime.TimeResponseInit(doc.CreatedAt),
		ShareDesc:            doc.ShareDesc,
		Order:                doc.Order,
		Platforms:            doc.Platforms,
		AllowShowShareAction: doc.AllowShowShareAction,
	}
}

// short ...
func (campaignImplement) short(ctx context.Context, doc mgaffiliate.Campaign) responsemodel.ResponseCampaignShort {
	var (
		fileResSvc = fileResponseImplement{}
	)
	return responsemodel.ResponseCampaignShort{
		ID:   doc.ID.Hex(),
		Name: doc.Name,
		Logo: fileResSvc.ConvertResponseFilePhoto(doc.Logo).GetResponseData(),
	}
}
