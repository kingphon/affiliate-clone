package service

import (
	"context"

	jsmodel "git.selly.red/Selly-Modules/natsio/js/model"

	clientjstream "git.selly.red/Selly-Server/affiliate/internal/nats/jstream"

	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SellerStatisticInterface ...
type SellerStatisticInterface interface {
	UpsertStatisticBySellerID(ctx context.Context, sellerID primitive.ObjectID) (err error)
}

// sellerStatisticImplement ...
type sellerStatisticImplement struct {
	TransactionService transactionImplement
}

// SellerStatistic ...
func SellerStatistic() SellerStatisticInterface {
	return &sellerStatisticImplement{}
}

// UpsertStatisticBySellerID ...
func (s sellerStatisticImplement) UpsertStatisticBySellerID(ctx context.Context, sellerID primitive.ObjectID) (err error) {
	var cond = bson.D{{"sellerId", sellerID}}

	statistic := s.TransactionService.AggregateStatisticByCondition(ctx, cond)
	var (
		d       = dao.SellerStatistic()
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

	// Nats publish seller update statistic
	var (
		jsPayload = jsmodel.PayloadUpdateSellerAffiliateStatistic{
			SellerID: sellerID.Hex(),
			Statistic: jsmodel.SellerAffiliateStatistic{
				TransactionTotal:              int(statistic.TransactionTotal),
				TransactionCashback:           int(statistic.TransactionCashback),
				TransactionPending:            int(statistic.TransactionPending),
				TransactionApproved:           int(statistic.TransactionApproved),
				TransactionRejected:           int(statistic.TransactionRejected),
				CommissionTransactionTotal:    statistic.CommissionTotal,
				CommissionTransactionCashback: statistic.CommissionCashback,
				CommissionTransactionApproved: statistic.CommissionApproved,
				CommissionTransactionPending:  statistic.CommissionPending,
				CommissionTransactionRejected: statistic.CommissionRejected,
			},
		}
		clientJS = clientjstream.ClientJestStreamPull{}
	)

	clientJS.UpdateSellerAffiliateStatistic(jsPayload)
	return
}
