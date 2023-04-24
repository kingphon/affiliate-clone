package service

import (
	"context"
	"fmt"
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Modules/mongodb"
	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	"math/rand"
	"sync"
	"time"

	"git.selly.red/Selly-Server/affiliate/external/utils/prandom"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"

	"github.com/friendsofgo/errors"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

//
// PUBLIC METHODS
//

// ChangeStatus ...
func (s reconciliationImplement) ChangeStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.ReconciliationPayloadChangeStatus) (result responsemodel.ResponseChangeStatus, err error) {
	// 1. Check isExited reconciliation
	var d = dao.Reconciliation()
	reconciliation := d.FindOneByCondition(ctx, bson.M{"_id": id})
	if reconciliation.ID.IsZero() {
		return result, errors.New(errorcode.ReconciliationNotFound)
	}

	// 2. Check status in db
	statuses := []string{constant.ReconciliationStatus.New.Key, constant.ReconciliationStatus.Approved.Key}
	if !funk.Contains(statuses, reconciliation.Status) {
		return result, errors.New(errorcode.ReconciliationStatusIsNewAndApproved)
	}

	// 3. Check payload status
	switch payload.Status {
	case constant.ReconciliationStatus.Deleted.Key:
		return s.changeStatusDeleted(ctx, reconciliation)
	case constant.ReconciliationStatus.Rejected.Key:
		return s.changeStatusRejected(ctx, reconciliation)
	case constant.ReconciliationStatus.Approved.Key:
		return s.changeStatusApproved(ctx, reconciliation)
	case constant.ReconciliationStatus.Completed.Key:
		return s.changeStatusCompleted(ctx, reconciliation)
	default:
		return result, errors.New(errorcode.ReconciliationInvalidStatus)
	}

}

//
// PRIVATE METHODS
//

// changeStatusDeleted ...
func (s reconciliationImplement) changeStatusDeleted(ctx context.Context, doc mgaffiliate.Reconciliation) (result responsemodel.ResponseChangeStatus, err error) {
	// Check permission
	if s.currentStaff.ID != doc.ActionBy.CreatedBy.Hex() {
		return result, errors.New(errorcode.OnlyStaffCreatorCanDeleteReconciliation)

	}

	// Check status db
	if doc.Status != constant.ReconciliationStatus.New.Key {
		return result, errors.New(errorcode.ReconciliationInvalidNewStatusInDatabase)
	}

	// Update
	var (
		cond          = bson.M{"_id": doc.ID}
		payloadUpdate = bson.D{
			{"status", constant.ReconciliationStatus.Deleted.Key},
			{"updateAt", ptime.Now()},
			{"trackingTime.changeStatusDeletedAt", ptime.Now()},
			{"actionBy.deletedBy", s.currentStaff.ID},
		}
	)
	if err = dao.Reconciliation().UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		return responsemodel.ResponseChangeStatus{}, errors.New(errorcode.ReconciliationCanNotUpdateStatus)
	}

	// Audit
	auditSvc := Audit(s.currentStaff)
	go auditSvc.Create(
		constant.AuditTargetReconciliation,
		doc.ID.Hex(),
		payloadUpdate,
		constant.MsgChangeStatusDeletedAffiliateReconciliation,
		constant.AuditActionEdit,
	)

	// Response
	result.ID = doc.ID.Hex()
	result.Status = constant.ReconciliationStatus.Deleted.Key
	return
}

// changeStatusRejected ...
func (s reconciliationImplement) changeStatusRejected(ctx context.Context, doc mgaffiliate.Reconciliation) (result responsemodel.ResponseChangeStatus, err error) {

	// Check permission
	if s.currentStaff.ID == doc.ActionBy.CreatedBy.Hex() {
		return result, errors.New(errorcode.StaffCreateReconciliationUncensored)
	}

	// Check status db
	if doc.Status != constant.ReconciliationStatus.New.Key {
		return result, errors.New(errorcode.ReconciliationInvalidNewStatusInDatabase)
	}

	// Update
	var (
		cond          = bson.M{"_id": doc.ID}
		payloadUpdate = bson.D{
			{"status", constant.ReconciliationStatus.Rejected.Key},
			{"updateAt", ptime.Now()},
			{"trackingTime.changeStatusRejectedAt", ptime.Now()},
			{"actionBy.rejectedBy", s.currentStaff.ID},
		}
	)
	if err = dao.Reconciliation().UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		return responsemodel.ResponseChangeStatus{}, errors.New(errorcode.ReconciliationCanNotUpdateStatus)
	}

	// Audit
	auditSvc := Audit(s.currentStaff)
	go auditSvc.Create(
		constant.AuditTargetReconciliation,
		doc.ID.Hex(),
		payloadUpdate,
		constant.MsgChangeStatusRejectedAffiliateReconciliation,
		constant.AuditActionEdit,
	)

	// Response
	result.ID = doc.ID.Hex()
	result.Status = constant.ReconciliationStatus.Rejected.Key
	return
}

// changeStatusApproved ...
func (s reconciliationImplement) changeStatusApproved(ctx context.Context, doc mgaffiliate.Reconciliation) (result responsemodel.ResponseChangeStatus, err error) {
	// Check permission
	if s.currentStaff.ID == doc.ActionBy.CreatedBy.Hex() {
		return result, errors.New(errorcode.StaffCreateReconciliationUncensored)
	}

	// Check status db
	if doc.Status != constant.ReconciliationStatus.New.Key {
		return result, errors.New(errorcode.ReconciliationInvalidNewStatusInDatabase)
	}

	// Update
	var (
		cond          = bson.M{"_id": doc.ID}
		payloadUpdate = bson.D{
			{"status", constant.ReconciliationStatus.Approved.Key},
			{"updateAt", ptime.Now()},
			{"trackingTime.changeStatusApprovedAt", ptime.Now()},
			{"actionBy.approvedBy", s.currentStaff.ID},
		}
	)
	if err = dao.Reconciliation().UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		return responsemodel.ResponseChangeStatus{}, errors.New(errorcode.ReconciliationCanNotUpdateStatus)
	}

	// Audit
	auditSvc := Audit(s.currentStaff)
	go auditSvc.Create(
		constant.AuditTargetReconciliation,
		doc.ID.Hex(),
		payloadUpdate,
		constant.MsgChangeStatusApprovedAffiliateReconciliation,
		constant.AuditActionEdit,
	)

	// Response
	result.ID = doc.ID.Hex()
	result.Status = constant.ReconciliationStatus.Running.Key
	return
}

// changeStatusCompleted ...
func (s reconciliationImplement) changeStatusCompleted(ctx context.Context, doc mgaffiliate.Reconciliation) (result responsemodel.ResponseChangeStatus, err error) {
	// Check status db
	if doc.Status != constant.ReconciliationStatus.Approved.Key {
		return result, errors.New(errorcode.ReconciliationInvalidApprovedStatusInDatabase)
	}

	// Update
	var (
		cond          = bson.M{"_id": doc.ID}
		payloadUpdate = bson.D{
			{"status", constant.ReconciliationStatus.Running.Key},
			{"updateAt", ptime.Now()},
			{"trackingTime.changeStatusRunningAt", ptime.Now()},
		}
	)
	if err = dao.Reconciliation().UpdateOneByCondition(ctx, cond, bson.M{"$set": payloadUpdate}); err != nil {
		return responsemodel.ResponseChangeStatus{}, errors.New(errorcode.ReconciliationCanNotUpdateStatus)
	}

	// Response
	result.ID = doc.ID.Hex()
	result.Status = constant.ReconciliationStatus.Running.Key

	var (
		auditSvc = Audit(s.currentStaff)
		ctxBg    = context.Background()
	)

	// Audit
	go auditSvc.Create(
		constant.AuditTargetReconciliation,
		doc.ID.Hex(),
		payloadUpdate,
		constant.MsgChangeStatusRunningAffiliateReconciliation,
		constant.AuditActionEdit,
	)

	// Cashback
	go s.runCashback(ctxBg, doc.ID)
	return
}

// runCashback ...
func (s reconciliationImplement) runCashback(ctx context.Context, reconciliationId primitive.ObjectID) (err error) {
	// 1. Get transactions
	reconciliation, err := s.FindByID(ctx, reconciliationId)
	if err != nil {
		return errors.New(errorcode.ReconciliationNotFound)
	}

	// No condition
	if reconciliation.Condition == nil {
		fmt.Println("Không có condition")
		return
	}

	var (
		transactionSvc            = transactionImplement{}
		transactionCrawlSvc       = transactionCrawlImplement{}
		limit               int64 = 1000
		page                int64 = 0
		q                         = mgquery.AppQuery{
			Affiliate: mgquery.Affiliate{
				Campaign:     reconciliation.Condition.CampaignId.Hex(),
				FromAt:       reconciliation.Condition.FromAt,
				ToAt:         reconciliation.Condition.ToAt,
				Source:       reconciliation.Condition.Source,
				Status:       constant.TransactionStatus.Approved.Key,
				TimeCashback: ptime.Now(),
			},
			Page:          page,
			Limit:         limit,
			SortInterface: bson.D{{"createdAt", -1}},
		}
	)

	// 2. Generate name temp
	nameTemp := fmt.Sprintf("%s_%d", prandom.RandomStringWithLength(5), time.Now().Unix())
	daoTemp := database.ReconcilicationTempByNameCol(nameTemp)

	// 3. Insert transaction temp
	for {
		q.Page = page
		transactions := transactionSvc.GetByQuery(ctx, q)
		if len(transactions) == 0 {
			break
		}

		// Convert temps
		transactionCrawls := s.convertToTransactionCraw(transactions, reconciliation.ID)

		// insert many temps
		if _, err := daoTemp.InsertMany(ctx, transactionCrawls); err != nil {
			fmt.Println("err insert many temps", err.Error())
			break
		}

		if len(transactions) != int(q.Limit) {
			break
		}

		page++
	}

	// 4. Check and update transaction from transaction-temps
	if err = transactionCrawlSvc.checkAndUpdateTransactionFromTemp(daoTemp); err != nil {
		logger.Error("Error-checkAndUpdateTransactionFromTemp:", logger.LogData{
			Data: bson.M{
				"error": err.Error(),
			},
		})
	}

	// 5. Update campaign statistic
	campaignIds := transactionCrawlSvc.distinctCampaignID(ctx, daoTemp)
	go Campaign(externalauth.User{}).UpdateStatisticListCampaign(ctx, campaignIds)

	// 6. Update status success
	var (
		cond          = bson.M{"_id": reconciliationId}
		payloadUpdate = bson.M{"$set": bson.D{
			{"status", constant.ReconciliationStatus.Completed.Key},
			{"updateAt", ptime.Now()},
			{"trackingTime.changeStatusCompletedAt", ptime.Now()},
			{"actionBy.cashbackBy", s.currentStaff.ID},
		}}
	)
	dao.Reconciliation().UpdateOneByCondition(ctx, cond, payloadUpdate)

	// 7. Drop collection
	daoTemp.Drop(ctx)

	// 8. Audit
	auditSvc := Audit(s.currentStaff)
	go auditSvc.Create(
		constant.AuditTargetReconciliation,
		reconciliationId.Hex(),
		payloadUpdate,
		constant.MsgChangeStatusCompletedAffiliateReconciliation,
		constant.AuditActionEdit,
	)

	return
}

// ConvertToTransactionCraw ...
func (s reconciliationImplement) convertToTransactionCraw(docs []mgaffiliate.Transaction, reconciliationID primitive.ObjectID) (result []interface{}) {
	// 1. Init
	result = make([]interface{}, len(docs))

	var wg = sync.WaitGroup{}
	wg.Add(len(docs))

	// 2. ConvertToTransactionCraw
	for i, trans := range docs {
		go func(index int, doc mgaffiliate.Transaction) {
			defer wg.Done()
			data := responsemodel.ResponseTransactionCrawl{
				ID:              mongodb.NewObjectID(),
				Campaign:        doc.CampaignID,
				ClickID:         doc.Click.ClickID,
				Source:          doc.Source,
				Status:          constant.TransactionStatus.Cashback.Key,
				Commission:      doc.Commission,
				TransactionID:   doc.Code,
				TransactionTime: doc.TransactionTime,
				Category:        doc.Category,
				ClickURL:        doc.Click.AffiliateURL,
				RejectedReason:  doc.RejectedReason,
				Device: responsemodel.ResponseTransactionCrawlDevice{
					Model:     doc.Device.Model,
					Type:      doc.Device.DeviceType,
					Device:    doc.Device.DeviceID,
					OS:        doc.Device.OSName,
					Browser:   doc.Device.BrowserVersion,
					UserAgent: doc.Device.UserAgent,
					OsVersion: doc.Device.OSVersion,
				},

				CreatedAt:        doc.CreatedAt,
				UpdatedAt:        ptime.Now(),
				UpdatedHash:      RandStringBytes(),
				SellerID:         doc.SellerID,
				ReconciliationID: reconciliationID,
			}
			result[index] = data
		}(i, trans)
	}

	wg.Wait()
	return
}

// RandStringBytes ...
func RandStringBytes() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
