package service

import (
	"context"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	"sync"

	"git.selly.red/Selly-Modules/mongodb"
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
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

//
// PUBLIC METHODS
//

// All ...
func (s reconciliationImplement) All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseReconciliationAll) {
	var (
		d  = dao.Reconciliation()
		wg = sync.WaitGroup{}
	)

	// Assign conditions
	cond := bson.D{}
	q.Affiliate.AssignKeyword(&cond)
	q.Affiliate.AssignReconciliationSource(&cond)
	q.Affiliate.AssignReconciliationCampaign(&cond)
	q.Affiliate.AssignStatus(&cond)
	q.Affiliate.AssignCreatedAt(&cond)

	// Find
	wg.Add(2)

	go func() {
		defer wg.Done()

		// Prepare data
		result.List = make([]responsemodel.ResponseReconciliationBrief, 0)

		// Find options
		findOpts := q.GetFindOptionsWithPage()
		findOpts.SetProjection(bson.M{
			"_id":       1,
			"name":      1,
			"type":      1,
			"status":    1,
			"condition": 1,
			"createdAt": 1,
			"actionBy":  1,
		})

		docs := d.FindByCondition(ctx, cond, findOpts)

		// Get list info staff by reconciliation
		data := s.getActionByListReconciliation(ctx, docs)

		result.List = s.getReconciliationBriefByList(ctx, docs, data)

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

// Detail ...
func (s reconciliationImplement) Detail(ctx context.Context, id primitive.ObjectID) (result *responsemodel.ResponseReconciliationDetail, err error) {

	var (
		d    = dao.Reconciliation()
		cond = bson.M{"_id": id}
	)

	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		return nil, errors.New(errorcode.ReconciliationNotFound)
	}

	// Get  info staff actionBy reconciliation
	docs := []mgaffiliate.Reconciliation{doc}
	data := s.getActionByListReconciliation(ctx, docs)

	result = s.getInfoDetailByReconciliation(ctx, doc, data)

	return
}

// FindByID ...
func (reconciliationImplement) FindByID(ctx context.Context, id primitive.ObjectID) (result mgaffiliate.Reconciliation, err error) {
	var (
		d    = dao.Reconciliation()
		cond = bson.M{"_id": id}
	)

	reconciliation := d.FindOneByCondition(ctx, cond)
	if reconciliation.ID.IsZero() {
		err = errors.New(errorcode.ReconciliationNotFound)
		return
	}

	result = reconciliation
	return
}

// GetTransactionsByCondition ...
func (reconciliationImplement) GetTransactionsByCondition(ctx context.Context, q mgquery.AppQuery, id primitive.ObjectID) (result responsemodel.ResponseTransactionAll) {
	result.List = make([]responsemodel.ResponseTransactionBrief, 0)
	result.Total = 0

	// Find reconciliation
	var (
		d    = dao.Reconciliation()
		cond = bson.M{"_id": id}
	)
	reconciliation := d.FindOneByCondition(ctx, cond)

	if reconciliation.Condition == nil {
		return
	}

	if reconciliation.Condition.CampaignId.IsZero() {
		return
	}

	transactionSvc := transactionImplement{}
	result = transactionSvc.GetByCondition(ctx, q, reconciliation)

	return
}

// GetStatistic ...
func (reconciliationImplement) GetStatistic(ctx context.Context, id primitive.ObjectID, payload requestmodel.ReconciliationPayloadStatistic) (result responsemodel.ResponseReconciliationStatistic, err error) {
	result = responsemodel.ResponseReconciliationStatistic{
		TotalTransaction:    0,
		TotalCommissionReal: 0,
		SellerCommission:    0,
		SellyCommission:     0,
	}

	var (
		transactionSvc = transactionImplement{}
		d              = dao.Reconciliation()
		cond           = bson.M{"_id": id}
	)

	reconciliation := d.FindOneByCondition(ctx, cond)

	if reconciliation.Condition == nil {
		return
	}

	if reconciliation.Condition.CampaignId.IsZero() {
		return
	}

	result = transactionSvc.AggregateStatisticByReconciliationCondition(ctx, reconciliation, payload)

	return
}

//
// PRIVATE METHODS
//

// brief ...
func (s reconciliationImplement) brief(ctx context.Context, doc mgaffiliate.Reconciliation, actionByInfo responsemodel.ResponseReconciliationActionBy) responsemodel.ResponseReconciliationBrief {
	result := responsemodel.ResponseReconciliationBrief{
		ID:        doc.ID.Hex(),
		Name:      doc.Name,
		Type:      doc.Type,
		Status:    doc.Status,
		CreatedAt: ptime.TimeResponseInit(doc.CreatedAt),
		ActionBy: &responsemodel.ResponseReconciliationActionBy{
			CreateBy: responsemodel.ResponseStaff{
				ID:   actionByInfo.CreateBy.ID,
				Name: actionByInfo.CreateBy.Name,
			},
		},
	}

	switch doc.Type {
	case constant.ReconciliationType.Filter.Key:
		result.Condition = s.convertReconciliationCondition(ctx, doc.Condition)
	case constant.ReconciliationType.Import.Key:
	}

	if doc.ActionBy != nil {
		if !doc.ActionBy.ApproveBy.IsZero() {
			result.ActionBy.ApproveBy = &responsemodel.ResponseStaff{
				ID:   actionByInfo.ApproveBy.ID,
				Name: actionByInfo.ApproveBy.Name,
			}
		}
	}

	return result
}

// convertReconciliationCondition ...
func (reconciliationImplement) convertReconciliationCondition(ctx context.Context, doc *mgaffiliate.ReconciliationCondition) *responsemodel.ResponseReconciliationCondition {
	// Campaign
	campaignSvc := campaignImplement{}
	campaign := campaignSvc.GetShortInfoByID(ctx, doc.CampaignId)

	return &responsemodel.ResponseReconciliationCondition{
		Campaign: responsemodel.ResponseCampaignShortInfo{
			ID:   campaign.ID,
			Name: campaign.Name,
		},
		Source: doc.Source,
		FromAt: ptime.TimeResponseInit(doc.FromAt),
		ToAt:   ptime.TimeResponseInit(doc.ToAt),
	}
}

// getActionByListReconciliation ...
func (s reconciliationImplement) getActionByListReconciliation(ctx context.Context, docs []mgaffiliate.Reconciliation) (data []responsemodel.ResponseStaff) {
	// uniq list id staff
	listIDsStaffUniq := s.uniqListIDsStaff(ctx, docs)

	// get list staff by ids
	data = s.getListStaffBysIDs(ctx, listIDsStaffUniq)

	return
}

// getListStaffBysIDs ...
func (reconciliationImplement) getListStaffBysIDs(ctx context.Context, ids []primitive.ObjectID) []responsemodel.ResponseStaff {
	var (
		staffSvc = staffImplement{}
	)

	return staffSvc.GetListStaffInfoByIDs(ids)
}

// getReconciliationBriefByList ...
func (s reconciliationImplement) getReconciliationBriefByList(ctx context.Context, docs []mgaffiliate.Reconciliation, idsStaff []responsemodel.ResponseStaff) (result []responsemodel.ResponseReconciliationBrief) {

	total := len(docs)
	result = make([]responsemodel.ResponseReconciliationBrief, total)

	wg := sync.WaitGroup{}
	wg.Add(total)

	for i, doc := range docs {
		go func(i int, doc mgaffiliate.Reconciliation) {
			defer wg.Done()
			result[i] = s.getInfoBriefByReconciliation(ctx, doc, idsStaff)
		}(i, doc)
	}

	wg.Wait()
	return
}

// getInfoBriefByReconciliation ...
func (s reconciliationImplement) getInfoBriefByReconciliation(ctx context.Context, doc mgaffiliate.Reconciliation, preData []responsemodel.ResponseStaff) responsemodel.ResponseReconciliationBrief {
	var (
		createdBy  responsemodel.ResponseStaff
		approvedBy responsemodel.ResponseStaff
		wg         = sync.WaitGroup{}
	)
	wg.Add(1)

	go func() {
		defer wg.Done()

		foundR := parray.Find(preData, func(item responsemodel.ResponseStaff) bool {
			return item.ID == doc.ActionBy.CreatedBy.Hex()
		})
		if foundR != nil {
			createdBy = foundR.(responsemodel.ResponseStaff)
		}
	}()

	if !doc.ActionBy.ApproveBy.IsZero() {
		wg.Add(1)
		go func() {
			defer wg.Done()

			foundR := parray.Find(preData, func(item responsemodel.ResponseStaff) bool {
				return item.ID == doc.ActionBy.ApproveBy.Hex()
			})
			if foundR != nil {
				approvedBy = foundR.(responsemodel.ResponseStaff)
			}
		}()

	}

	wg.Wait()

	data := responsemodel.ResponseReconciliationActionBy{
		CreateBy:  createdBy,
		ApproveBy: &approvedBy,
	}

	return s.brief(ctx, doc, data)
}

// uniqListIDsStaff ...
func (s reconciliationImplement) uniqListIDsStaff(ctx context.Context, docs []mgaffiliate.Reconciliation) (result []primitive.ObjectID) {
	// init
	listStringIDsStaff := make([]string, 0)

	// append data
	for _, r := range docs {
		if !r.ActionBy.CreatedBy.IsZero() {
			listStringIDsStaff = append(listStringIDsStaff, r.ActionBy.CreatedBy.Hex())
		}
		if !r.ActionBy.ApproveBy.IsZero() {
			listStringIDsStaff = append(listStringIDsStaff, r.ActionBy.ApproveBy.Hex())
		}
	}

	uniqIDs := parray.UniqueArrayStrings(listStringIDsStaff)

	result = make([]primitive.ObjectID, 0)
	for _, id := range uniqIDs {
		objID, _ := mongodb.NewIDFromString(id)
		result = append(result, objID)
	}

	return
}

// detail ...
func (s reconciliationImplement) detail(ctx context.Context, doc mgaffiliate.Reconciliation, actionByInfo responsemodel.ResponseReconciliationActionBy) *responsemodel.ResponseReconciliationDetail {
	result := &responsemodel.ResponseReconciliationDetail{
		ID:        doc.ID.Hex(),
		Name:      doc.Name,
		Type:      doc.Type,
		Status:    doc.Status,
		Condition: nil,
		CreatedAt: ptime.TimeResponseInit(doc.CreatedAt),
		ActionBy: &responsemodel.ResponseReconciliationActionBy{
			CreateBy: responsemodel.ResponseStaff{
				ID:   actionByInfo.CreateBy.ID,
				Name: actionByInfo.CreateBy.Name,
			},
		},
	}

	switch doc.Type {
	case constant.ReconciliationType.Filter.Key:
		result.Condition = s.convertReconciliationCondition(ctx, doc.Condition)
	case constant.ReconciliationType.Import.Key:
	}

	if doc.ActionBy != nil {
		if !doc.ActionBy.ApproveBy.IsZero() {
			result.ActionBy.ApproveBy = &responsemodel.ResponseStaff{
				ID:   actionByInfo.ApproveBy.ID,
				Name: actionByInfo.ApproveBy.Name,
			}
		}
	}

	if doc.TrackingTime != nil {
		if !doc.TrackingTime.ChangeStatusApprovedAt.IsZero() {
			if result.TrackingTime == nil {
				result.TrackingTime = &responsemodel.ResponseReconciliationTrackingTime{}
			}
			result.TrackingTime.ChangeStatusApprovedAt = ptime.TimeResponseInit(doc.TrackingTime.ChangeStatusApprovedAt)
		}
	}

	return result
}

// getInfoDetailByReconciliation ...
func (s reconciliationImplement) getInfoDetailByReconciliation(ctx context.Context, doc mgaffiliate.Reconciliation, preData []responsemodel.ResponseStaff) *responsemodel.ResponseReconciliationDetail {
	var (
		createdBy  responsemodel.ResponseStaff
		approvedBy responsemodel.ResponseStaff
		wg         = sync.WaitGroup{}
	)
	wg.Add(1)

	go func() {
		defer wg.Done()

		foundR := parray.Find(preData, func(item responsemodel.ResponseStaff) bool {
			return item.ID == doc.ActionBy.CreatedBy.Hex()
		})
		if foundR != nil {
			createdBy = foundR.(responsemodel.ResponseStaff)
		}
	}()

	if !doc.ActionBy.ApproveBy.IsZero() {
		wg.Add(1)
		go func() {
			defer wg.Done()

			foundR := parray.Find(preData, func(item responsemodel.ResponseStaff) bool {
				return item.ID == doc.ActionBy.ApproveBy.Hex()
			})
			if foundR != nil {
				approvedBy = foundR.(responsemodel.ResponseStaff)
			}
		}()

	}

	wg.Wait()

	data := responsemodel.ResponseReconciliationActionBy{
		CreateBy:  createdBy,
		ApproveBy: &approvedBy,
	}

	return s.detail(ctx, doc, data)
}
