package service

import (
	"context"
	"git.selly.red/Selly-Modules/mongodb"
	"sync"

	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"

	natsmodel "git.selly.red/Selly-Modules/natsio/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/external/utils/parray"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

// Statistic ...
func (s clickImplement) Statistic(ctx context.Context, q mgquery.Affiliate) (result responsemodel.ResponseClickStatistic) {
	var (
		cond = bson.D{}
	)
	q.AssignCreatedAt(&cond)
	q.AssignSource(&cond)
	q.AssignCampaignIds(&cond)
	q.AssignSeller(&cond)

	data := dao.Click().AggregateStatistic(ctx, cond)
	return responsemodel.ResponseClickStatistic{
		Total:          data.Total,
		TotalPending:   data.TotalPending,
		TotalCompleted: data.TotalCompleted,
	}
}

// FindByID ...
func (s clickImplement) FindByID(ctx context.Context, id primitive.ObjectID) (result mgaffiliate.Click) {
	var (
		d    = dao.Click()
		cond = bson.D{{"_id", id}}
	)

	result = d.FindOneByCondition(ctx, cond)
	return
}

// All ...
func (s clickImplement) All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseClickAll) {
	var (
		d  = dao.Click()
		wg = sync.WaitGroup{}
	)

	// Assign condition
	cond := bson.D{}
	q.Affiliate.AssignKeyword(&cond)
	q.Affiliate.AssignStatus(&cond)
	q.Affiliate.AssignSourceClick(&cond)
	q.Affiliate.AssignCampaign(&cond)
	q.Affiliate.AssignCreatedAt(&cond)
	q.Affiliate.AssignSeller(&cond)
	q.Affiliate.AssignFromClick(&cond)

	wg.Add(2)

	// Find
	go func() {
		defer wg.Done()

		// Find options
		findOpts := q.GetFindOptionsWithPage()
		docs := d.FindByCondition(ctx, cond, findOpts)
		if len(docs) == 0 {
			result.List = make([]responsemodel.ResponseClickBrief, 0)
			return
		}

		data := s.getSellerAndCampaignByListClicks(ctx, docs)
		result.List = s.getClickBriefByList(ctx, docs, data)

	}()

	// Total
	go func() {
		defer wg.Done()

		result.Total = d.CountByCondition(ctx, cond)
	}()

	wg.Wait()

	// Limit
	result.Limit = q.Limit
	return
}

//
// PRIVATE METHODS
//

// getSellerAndCampaignByListClicks ...
func (s clickImplement) getSellerAndCampaignByListClicks(ctx context.Context, docs []mgaffiliate.Click) (data responsemodel.DataSellerCampaign) {
	var (
		wg = sync.WaitGroup{}
	)

	wg.Add(2)

	// getSeller
	go func() {
		defer wg.Done()

		sellerSvc := sellerImplement{}
		sellerIDs := s.getSellerIDsByClickList(ctx, docs)
		data.Sellers, _ = sellerSvc.GetSellerByIDs(ctx, sellerIDs)
	}()

	// get campaign
	go func() {
		defer wg.Done()

		campaignSvc := campaignImplement{}
		campaignIDs := s.getCampaignIDsByClickList(ctx, docs)
		data.Campaigns = campaignSvc.GetCampaignByIDs(ctx, campaignIDs)
	}()

	wg.Wait()

	return
}

// getSellerIDsByClickList ...
func (clickImplement) getSellerIDsByClickList(ctx context.Context, docs []mgaffiliate.Click) []primitive.ObjectID {
	sellerIDs := make([]primitive.ObjectID, 0)
	for _, doc := range docs {
		sellerIDs = append(sellerIDs, doc.SellerID)
	}

	sellerIDsUniq := mongodb.UniqObjectIds(sellerIDs)
	return sellerIDsUniq
}

// getCampaignIDsByClickList ...
func (clickImplement) getCampaignIDsByClickList(ctx context.Context, docs []mgaffiliate.Click) []primitive.ObjectID {
	campaignIDs := make([]primitive.ObjectID, 0)
	for _, doc := range docs {
		campaignIDs = append(campaignIDs, doc.CampaignID)
	}

	campaignIDsUniq := mongodb.UniqObjectIds(campaignIDs)
	return campaignIDsUniq

}

// getClickBriefByList ...
func (s clickImplement) getClickBriefByList(ctx context.Context, docs []mgaffiliate.Click, data responsemodel.DataSellerCampaign) (result []responsemodel.ResponseClickBrief) {
	total := len(docs)
	result = make([]responsemodel.ResponseClickBrief, total)

	wg := sync.WaitGroup{}

	wg.Add(total)

	for i, doc := range docs {
		go func(i int, doc mgaffiliate.Click) {
			defer wg.Done()
			result[i] = s.getInfoBriefByClick(ctx, doc, data)
		}(i, doc)

	}

	wg.Wait()
	return
}

// getInfoBriefByClick ...
func (s clickImplement) getInfoBriefByClick(ctx context.Context, doc mgaffiliate.Click, preData responsemodel.DataSellerCampaign) responsemodel.ResponseClickBrief {
	var (
		campaign responsemodel.ResponseCampaignShort
		seller   natsmodel.ResponseSellerInfo
		wg       = sync.WaitGroup{}
	)

	wg.Add(2)

	// campaign
	go func() {
		defer wg.Done()

		foundC := parray.Find(preData.Campaigns, func(item responsemodel.ResponseCampaignShort) bool {
			return item.ID == doc.CampaignID.Hex()
		})
		if foundC != nil {
			campaign = foundC.(responsemodel.ResponseCampaignShort)
		}
	}()

	// seller
	go func() {
		defer wg.Done()

		foundS := parray.Find(preData.Sellers, func(item natsmodel.ResponseSellerInfo) bool {
			return item.ID == doc.SellerID.Hex()
		})
		if foundS != nil {
			seller = foundS.(natsmodel.ResponseSellerInfo)
		}
	}()

	wg.Wait()

	return s.brief(ctx, doc, campaign, seller)
}

// brief ...
func (clickImplement) brief(ctx context.Context, doc mgaffiliate.Click, campaign responsemodel.ResponseCampaignShort, seller natsmodel.ResponseSellerInfo) responsemodel.ResponseClickBrief {
	return responsemodel.ResponseClickBrief{
		ID: doc.ID.Hex(),
		Campaign: responsemodel.ResponseCampaignShort{
			ID:   campaign.ID,
			Name: campaign.Name,
			Logo: campaign.Logo,
		},
		Seller: responsemodel.ResponseSellerShort{
			ID:   seller.ID,
			Name: seller.Name,
		},
		PartnerSource: doc.PartnerSource,
		AffiliateURL:  doc.AffiliateURL,
		CampaignURL:   doc.CampaignURL,
		ShareURL:      doc.ShareURL,
		Status:        doc.Status,
		Device: responsemodel.ResponseTransactionDevice{
			Model:          doc.Device.Model,
			UserAgent:      doc.Device.UserAgent,
			OSName:         doc.Device.OSName,
			OSVersion:      doc.Device.OSVersion,
			BrowserVersion: doc.Device.BrowserVersion,
			BrowserName:    doc.Device.BrowserName,
			DeviceType:     doc.Device.DeviceType,
			Manufacturer:   doc.Device.Manufacturer,
			DeviceID:       doc.Device.DeviceID,
		},
		From: doc.From,
		FinalDetected: responsemodel.ResponseClickFinalDetected{
			URL:   doc.FinalDetected.URL,
			Click: doc.FinalDetected.Click,
		},
		CreatedAt: ptime.TimeResponseInit(doc.CreatedAt),
	}
}
