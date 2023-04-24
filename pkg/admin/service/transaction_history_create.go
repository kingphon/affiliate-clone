package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/locale"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateAndInsertHistoriesByUpdateTransaction ...
func (s transactionHistoryImplement) GenerateAndInsertHistoriesByUpdateTransaction(ctx context.Context, transaction mgaffiliate.Transaction) (err error) {
	history := mgaffiliate.TransactionHistory{
		ID:            primitive.NewObjectID(),
		TransactionID: transaction.ID,
		SellerID:      transaction.SellerID,
		Status:        transaction.Status,
		Desc:          s.generateDescByStatus(transaction.Status),
		CreatedAt:     ptime.Now(),
	}
	if transaction.Status == constant.TransactionStatus.Pending.Key {
		history.CreatedAt = transaction.TransactionTime
	}

	var d = dao.TransactionHistory()
	err = d.InsertOne(ctx, history)
	return
}

// GenerateAndInsertHistoriesByNewTransaction ...
func (s transactionHistoryImplement) GenerateAndInsertHistoriesByNewTransaction(ctx context.Context, transactionNew mgaffiliate.Transaction) (err error) {
	// 1. generate
	click := dao.Click().FindOneByCondition(ctx, bson.M{
		"_id": transactionNew.Click.ClickID,
	})
	var histories = make([]interface{}, 0)
	histories = append(histories, mgaffiliate.TransactionHistory{
		ID:            primitive.NewObjectID(),
		TransactionID: transactionNew.ID,
		SellerID:      transactionNew.SellerID,
		Status:        constant.TransactionHistoryStatusNew,
		Desc:          s.generateDescByStatus(constant.TransactionHistoryStatusNew),
		CreatedAt:     click.CreatedAt,
	})

	// 2. generate by transaction status
	history := mgaffiliate.TransactionHistory{
		ID:            primitive.NewObjectID(),
		TransactionID: transactionNew.ID,
		SellerID:      transactionNew.SellerID,
		Status:        transactionNew.Status,
		Desc:          s.generateDescByStatus(transactionNew.Status),
		CreatedAt:     ptime.Now().Add(time.Millisecond),
	}
	if transactionNew.Status == constant.TransactionStatus.Pending.Key {
		history.CreatedAt = transactionNew.TransactionTime
	}
	histories = append(histories, history)

	var d = dao.TransactionHistory()
	err = d.InsertMany(ctx, histories)
	return
}

//
// PRIVATE METHOD
//

// generateDescByStatus ...
func (transactionHistoryImplement) generateDescByStatus(status string) (desc string) {
	switch status {
	case constant.TransactionHistoryStatusNew:
		desc = locale.TransactionHistoryDescStatusNew
	case constant.TransactionStatus.Pending.Key:
		desc = locale.TransactionHistoryDescStatusPending
	case constant.TransactionStatus.Approved.Key:
		desc = locale.TransactionHistoryDescStatusApproved
	case constant.TransactionStatus.Cashback.Key:
		desc = locale.TransactionHistoryDescStatusCashback
	case constant.TransactionStatus.Rejected.Key:
		desc = locale.TransactionHistoryDescStatusRejected
	}
	return
}
