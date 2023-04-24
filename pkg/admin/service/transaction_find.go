package service

import (
	"context"
	"git.selly.red/Selly-Modules/mongodb"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	"sync"

	natsmodel "git.selly.red/Selly-Modules/natsio/model"
	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/external/utils/parray"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/model/query"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

// Statistic ...
func (s transactionImplement) Statistic(ctx context.Context, q mgquery.Affiliate) (result responsemodel.ResponseTransactionStatistic) {
	var (
		cond = bson.D{}
	)
	q.AssignCampaignIds(&cond)
	//q.AssignStatus(&cond)
	q.AssignTransactionTime(&cond)
	q.AssignSource(&cond)
	q.AssignSeller(&cond)
	q.AssignStatuses(&cond)
	res := dao.Transaction().AggregateStatisticDashboardByCondition(ctx, cond)
	return responsemodel.ResponseTransactionStatistic{
		TotalTransaction:    res.TransactionTotal,
		TotalCommissionReal: res.CommissionTotal,
		SellerCommission:    res.SellerCommission,
		SellyCommission:     res.SellyCommission,
		TotalSeller:         int64(len(res.Sellers)),
		TotalCampaign:       int64(len(res.Campaigns)),
	}
}

// AggregateStatisticByCondition ...
func (s transactionImplement) AggregateStatisticByCondition(ctx context.Context, cond interface{}) query.TransactionStatistic {
	var d = dao.Transaction()
	return d.AggregateStatisticByCondition(ctx, cond)
}

// All ...
func (s transactionImplement) All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseTransactionAll) {
	var (
		d  = dao.Transaction()
		wg = sync.WaitGroup{}
	)

	// Assign condition
	cond := bson.D{}
	q.Affiliate.AssignKeyword(&cond)
	q.Affiliate.AssignStatuses(&cond)
	q.Affiliate.AssignSource(&cond)
	q.Affiliate.AssignTransactionTime(&cond)
	q.Affiliate.AssignCampaign(&cond)
	q.Affiliate.AssignSeller(&cond)

	wg.Add(2)

	// Find
	go func() {
		defer wg.Done()

		// Find options
		findOpts := q.GetFindOptionsWithPage()
		findOpts.SetProjection(bson.M{
			"_id":                      1,
			"code":                     1,
			"name":                     1,
			"campaignId":               1,
			"platformId":               1,
			"sellerId":                 1,
			"source":                   1,
			"commission":               1,
			"sellerCommissionRate":     1,
			"sellerCommission":         1,
			"sellyCommission":          1,
			"estimateSellerCommission": 1,
			"transactionTime":          1,
			"status":                   1,
			"rejectedReason":           1,
			"estimateCashbackAt":       1,
		})

		docs := d.FindByCondition(ctx, cond, findOpts)
		if len(docs) == 0 {
			result.List = make([]responsemodel.ResponseTransactionBrief, 0)
			return
		}

		data := s.getSellerAndCampaignByListTransaction(ctx, docs)
		result.List = s.getTransactionBriefByList(ctx, docs, data)
	}()

	// Assign total
	go func() {
		defer wg.Done()
		result.Total = d.CountByCondition(ctx, cond)
	}()

	wg.Wait()

	// Assign limit

	result.Limit = q.Limit
	return
}

// Detail ...
func (s transactionImplement) Detail(ctx context.Context, id primitive.ObjectID) (result *responsemodel.ResponseTransactionDetail, err error) {
	var (
		d           = dao.Transaction()
		cond        = bson.M{"_id": id}
		sellerSvc   = sellerImplement{}
		campaignSvc = campaignImplement{}
	)

	// Find
	transaction := d.FindOneByCondition(ctx, cond)
	if transaction.ID.IsZero() {
		return nil, errors.New(errorcode.CampaignNotFound)
	}

	// Get seller and campaign
	sellerIDs := []primitive.ObjectID{transaction.SellerID}
	campaignIDs := []primitive.ObjectID{transaction.CampaignID}

	seller, err := sellerSvc.GetSellerByIDs(ctx, sellerIDs)
	if err != nil || len(seller) == 0 {
		return nil, errors.New(errorcode.SellerNotFound)
	}

	campaign := campaignSvc.GetCampaignByIDs(ctx, campaignIDs)
	if len(campaign) == 0 {
		return nil, errors.New(errorcode.CampaignNotFound)
	}

	result = s.detail(ctx, transaction, seller[0], campaign[0])
	return
}

// FindByCode ...
func (s transactionImplement) FindByCode(ctx context.Context, code string) (result mgaffiliate.Transaction) {
	var (
		d    = dao.Transaction()
		cond = bson.D{{"code", code}}
	)

	result = d.FindOneByCondition(ctx, cond)
	return
}

// GetHistoriesByTransaction ...
func (s transactionImplement) GetHistoriesByTransaction(ctx context.Context, transactionID primitive.ObjectID) (result responsemodel.ResponseTransactionHistories) {
	var historyService = transactionHistoryImplement{}
	list := historyService.GetHistoriesByTransactionID(ctx, transactionID)

	result = responsemodel.ResponseTransactionHistories{Data: list}
	return
}

// GetByCondition ...
func (s transactionImplement) GetByCondition(ctx context.Context, q mgquery.AppQuery, doc mgaffiliate.Reconciliation) (result responsemodel.ResponseTransactionAll) {
	var (
		d  = dao.Transaction()
		wg = sync.WaitGroup{}
	)

	// Assign condition
	cond := bson.D{}
	if doc.Status == constant.ReconciliationStatus.Completed.Key {
		cond = bson.D{
			{"campaignId", doc.Condition.CampaignId},
			{"source", doc.Condition.Source},
			{"transactionTime", bson.M{
				"$gte": doc.Condition.FromAt,
				"$lte": doc.Condition.ToAt,
			}},
			{"$or", []interface{}{
				bson.D{{"reconciliationId", doc.ID},
					{"status", constant.TransactionStatus.Cashback.Key},
				},
				bson.D{{"status", constant.TransactionStatus.Approved.Key}},
			}},
		}
	} else {
		cond = bson.D{
			{"campaignId", doc.Condition.CampaignId},
			{"source", doc.Condition.Source},
			{"transactionTime", bson.M{
				"$gte": doc.Condition.FromAt,
				"$lte": doc.Condition.ToAt,
			}},
			{"status", constant.TransactionStatus.Approved.Key},
		}
	}

	q.Affiliate.AssignKeyword(&cond)
	q.Affiliate.AssignSeller(&cond)
	q.Affiliate.AssignStatuses(&cond)

	wg.Add(2)

	// Find
	go func() {
		defer wg.Done()

		// Find options
		findOpts := q.GetFindOptionsWithPage()

		docs := d.FindByCondition(ctx, cond, findOpts)

		data := s.getSellerAndCampaignByListTransaction(ctx, docs)
		result.List = s.getTransactionBriefByList(ctx, docs, data)

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

func (s transactionImplement) GetByReconciliationCondition(ctx context.Context, cond interface{}) (result []mgaffiliate.Transaction) {
	return dao.Transaction().FindByCondition(ctx, cond)
}

func (s transactionImplement) GetByQuery(ctx context.Context, q mgquery.AppQuery) (result []mgaffiliate.Transaction) {
	// 1. Init
	result = make([]mgaffiliate.Transaction, 0)

	var (
		d  = dao.Transaction()
		wg = sync.WaitGroup{}
	)

	// 2. Assign condition
	cond := bson.D{}
	q.Affiliate.AssignSource(&cond)
	q.Affiliate.AssignTransactionTime(&cond)
	q.Affiliate.AssignCampaign(&cond)
	q.Affiliate.AssignStatuses(&cond)
	q.Affiliate.AssignTimeCashback(&cond)

	wg.Add(1)

	// 3. Find
	go func() {
		defer wg.Done()

		// Find options
		findOpts := q.GetFindOptionsWithPage()

		docs := d.FindByCondition(ctx, cond, findOpts)
		if len(docs) == 0 {
			return
		}

		for _, doc := range docs {
			result = append(result, doc)
		}
	}()

	wg.Wait()

	return
}

//
// PRIVATE METHOD
//

// detail ...
func (transactionImplement) detail(ctx context.Context, doc mgaffiliate.Transaction, seller natsmodel.ResponseSellerInfo, campaign responsemodel.ResponseCampaignShort) *responsemodel.ResponseTransactionDetail {
	return &responsemodel.ResponseTransactionDetail{
		ID:   doc.ID.Hex(),
		Code: doc.Code,
		Seller: responsemodel.ResponseSellerShort{
			ID:   seller.ID,
			Name: seller.Name,
		},
		Campaign: responsemodel.ResponseCampaignShort{
			ID:   campaign.ID,
			Name: campaign.Name,
			Logo: campaign.Logo,
		},
		TransactionTime: ptime.TimeResponseInit(doc.TransactionTime),
		Commission: responsemodel.ResponseCampaignCommission{
			Real:          doc.Commission,
			SellerPercent: doc.SellerCommissionRate,
			Selly:         doc.SellyCommission,
			Seller:        doc.SellerCommission,
		},
		EstimateSellerCommission: doc.EstimateSellerCommission,
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
		Status:             doc.Status,
		EstimateCashbackAt: ptime.TimeResponseInit(doc.EstimateCashbackAt),
	}
}

// brief ...
func (transactionImplement) brief(ctx context.Context, doc mgaffiliate.Transaction, campaign responsemodel.ResponseCampaignShort, seller natsmodel.ResponseSellerInfo) responsemodel.ResponseTransactionBrief {
	return responsemodel.ResponseTransactionBrief{
		ID:   doc.ID.Hex(),
		Code: doc.Code,
		Campaign: responsemodel.ResponseCampaignShort{
			ID:   campaign.ID,
			Name: campaign.Name,
			Logo: campaign.Logo,
		},
		Seller: responsemodel.ResponseSellerShort{
			ID:   seller.ID,
			Name: seller.Name,
		},
		Source: doc.Source,
		Commission: responsemodel.ResponseCampaignCommission{
			Real:          doc.Commission,
			SellerPercent: doc.SellerCommissionRate,
			Selly:         doc.SellyCommission,
			Seller:        doc.SellerCommission,
		},
		EstimateSellerCommission: doc.EstimateSellerCommission,
		TransactionTime:          ptime.TimeResponseInit(doc.TransactionTime),
		Status:                   doc.Status,
		RejectedReason:           doc.RejectedReason,
		EstimateCashbackAt:       ptime.TimeResponseInit(doc.EstimateCashbackAt),
	}
}

// getSellerIDsByTransactionList ...
func (transactionImplement) getSellerIDsByTransactionList(ctx context.Context, docs []mgaffiliate.Transaction) []primitive.ObjectID {
	sellerIDs := make([]primitive.ObjectID, 0)
	for _, doc := range docs {
		sellerIDs = append(sellerIDs, doc.SellerID)
	}

	sellerIDsUniq := mongodb.UniqObjectIds(sellerIDs)
	return sellerIDsUniq
}

// getCampaignIDsByTransactionList ...
func (transactionImplement) getCampaignIDsByTransactionList(ctx context.Context, docs []mgaffiliate.Transaction) []primitive.ObjectID {
	campaignIDs := make([]primitive.ObjectID, 0)
	for _, doc := range docs {
		campaignIDs = append(campaignIDs, doc.CampaignID)
	}

	campaignIDsUniq := mongodb.UniqObjectIds(campaignIDs)
	return campaignIDsUniq
}

// getSellerAndCampaignByListTransaction ...
func (s transactionImplement) getSellerAndCampaignByListTransaction(ctx context.Context, docs []mgaffiliate.Transaction) (data responsemodel.DataSellerCampaign) {
	var (
		wg = sync.WaitGroup{}
	)

	wg.Add(2)

	// Get sellers
	go func() {
		defer wg.Done()
		sellerSvc := sellerImplement{}
		sellerIDs := s.getSellerIDsByTransactionList(ctx, docs)
		data.Sellers, _ = sellerSvc.GetSellerByIDs(ctx, sellerIDs)
	}()

	// Get campaigns
	go func() {
		defer wg.Done()
		campaignSvc := campaignImplement{}
		campaignIDs := s.getCampaignIDsByTransactionList(ctx, docs)
		data.Campaigns = campaignSvc.GetCampaignByIDs(ctx, campaignIDs)
	}()

	wg.Wait()
	return
}

// getTransactionBriefByList ...
func (s transactionImplement) getTransactionBriefByList(ctx context.Context, docs []mgaffiliate.Transaction, data responsemodel.DataSellerCampaign) (result []responsemodel.ResponseTransactionBrief) {
	total := len(docs)
	result = make([]responsemodel.ResponseTransactionBrief, total)

	wg := sync.WaitGroup{}

	wg.Add(total)

	for i, doc := range docs {
		go func(i int, doc mgaffiliate.Transaction) {
			defer wg.Done()
			result[i] = s.getInfoBriefByTransaction(ctx, doc, data)
		}(i, doc)
	}

	wg.Wait()
	return
}

// getInfoBriefByTransaction ...
func (s transactionImplement) getInfoBriefByTransaction(ctx context.Context, doc mgaffiliate.Transaction, preData responsemodel.DataSellerCampaign) responsemodel.ResponseTransactionBrief {
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

// AggregateStatisticByReconciliationCondition ...
func (transactionImplement) AggregateStatisticByReconciliationCondition(ctx context.Context, doc mgaffiliate.Reconciliation, payload requestmodel.ReconciliationPayloadStatistic) (result responsemodel.ResponseReconciliationStatistic) {
	var cond = bson.D{
		{"campaignId", doc.Condition.CampaignId},
		{"source", doc.Condition.Source},
		{"transactionTime", bson.M{
			"$gte": doc.Condition.FromAt,
			"$lte": doc.Condition.ToAt,
		}},
	}

	switch payload.Status {
	case constant.TransactionStatus.All.Key, "":
		appendCond := bson.E{
			"$or", []interface{}{
				bson.D{{"status", constant.TransactionStatus.Approved.Key}},
				bson.D{
					{"status", constant.TransactionStatus.Cashback.Key},
					{"reconciliationId", doc.ID},
				},
			},
		}
		cond = append(cond, appendCond)
	case constant.TransactionStatus.Approved.Key:
		appendCond := bson.E{"status", constant.TransactionStatus.Approved.Key}
		cond = append(cond, appendCond)
	case constant.TransactionStatus.Cashback.Key:
		appendCond := bson.D{
			{"status", constant.TransactionStatus.Cashback.Key},
			{"reconciliationId", doc.ID},
		}
		cond = append(cond, appendCond...)
	}

	if payload.Seller != "" {
		id := mongodb.ConvertStringToObjectID(payload.Seller)
		appendCond := bson.E{"sellerId", id}
		cond = append(cond, appendCond)
	}

	if payload.SearchCode != "" {
		appendCond := bson.E{"code", payload.SearchCode}
		cond = append(cond, appendCond)
	}

	res := dao.Transaction().AggregateStatisticByReconciliationCondition(ctx, cond)

	return responsemodel.ResponseReconciliationStatistic{
		TotalTransaction:    res.TransactionTotal,
		TotalCommissionReal: res.CommissionTotal,
		SellerCommission:    res.SellerCommission,
		SellyCommission:     res.SellyCommission,
	}
}
