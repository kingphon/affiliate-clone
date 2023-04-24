package service

import (
	"context"

	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SellerCampaignStatisticInterface ...
type SellerCampaignStatisticInterface interface {
	UpsertStatisticBySellerAndCampaignID(ctx context.Context, sellerID, campaignID primitive.ObjectID) (err error)
}

// sellerCampaignStatisticImplement ...
type sellerCampaignStatisticImplement struct {
	TransactionService transactionImplement
}

// SellerCampaignStatistic ...
func SellerCampaignStatistic() SellerCampaignStatisticInterface {
	return &sellerCampaignStatisticImplement{}
}

// UpsertStatisticBySellerAndCampaignID ...
func (s sellerCampaignStatisticImplement) UpsertStatisticBySellerAndCampaignID(ctx context.Context, sellerID, campaignID primitive.ObjectID) (err error) {
	var (
		cond = bson.D{
			{"sellerId", sellerID},
			{"campaignId", campaignID},
		}
	)

	statistic := s.TransactionService.AggregateStatisticByCondition(ctx, cond)
	var (
		d       = dao.SellerCampaignStatistic()
		opts    = options.Update().SetUpsert(true)
		payload = bson.M{
			"$set": bson.M{
				"statistic.transactionTotal":    statistic.TransactionTotal,
				"statistic.transactionCashback": statistic.TransactionCashback,
				"statistic.transactionPending":  statistic.TransactionPending,
				"statistic.transactionApproved": statistic.TransactionApproved,
				"statistic.transactionRejected": statistic.TransactionRejected,
				"statistic.commissionTotal":     statistic.CommissionTotal,
				"statistic.commissionCashback":  statistic.CommissionCashback,
				"statistic.commissionPending":   statistic.CommissionPending,
				"statistic.commissionApproved":  statistic.CommissionApproved,
				"statistic.commissionRejected":  statistic.CommissionRejected,
				"updatedAt":                     ptime.Now(),
			},
		}
	)

	err = d.UpdateByCondition(ctx, cond, payload, opts)
	return
}
