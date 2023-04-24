package dao

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/external/constant"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/model/query"

	"go.mongodb.org/mongo-driver/mongo/options"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
)

// TransactionInterface ...
type TransactionInterface interface {
	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (docs []mgaffiliate.Transaction)

	// CountByCondition ..
	CountByCondition(ctx context.Context, cond interface{}) int64

	// FindOneByCondition ...
	FindOneByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgaffiliate.Transaction)

	// InsertOne ...
	InsertOne(ctx context.Context, payload interface{}) (err error)

	// AggregateStatisticByCondition ...
	AggregateStatisticByCondition(ctx context.Context, cond interface{}) query.TransactionStatistic

	// AggregateStatisticDashboardByCondition ...
	AggregateStatisticDashboardByCondition(ctx context.Context, cond interface{}) query.TransactionStatisticDashboard

	// AggregateCampaignStatistic ...
	AggregateCampaignStatistic(ctx context.Context, cond interface{}) []mgaffiliate.TransactionWithStatistic

	// UpdateOneByCondition ...
	UpdateOneByCondition(ctx context.Context, cond, payload interface{}) error

	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error

	// UpdateManyByCondition ...
	UpdateManyByCondition(ctx context.Context, cond, payload interface{}) error

	// AggregateStatisticByReconciliationCondition ...
	AggregateStatisticByReconciliationCondition(ctx context.Context, cond interface{}) (result query.TransactionStatisticDashboard)

	// UpsertOneByCondition ...
	UpsertOneByCondition(ctx context.Context, cond, payload interface{}) error
}

// transactionImplement ...
type transactionImplement struct{}

// Transaction ...
func Transaction() TransactionInterface {
	return transactionImplement{}
}

// BulkWrite ...
func (d transactionImplement) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) error {
	_, err := database.TransactionCol().BulkWrite(ctx, models, opts...)
	return err
}

// AggregateStatisticDashboardByCondition ...
func (d transactionImplement) AggregateStatisticDashboardByCondition(ctx context.Context, cond interface{}) query.TransactionStatisticDashboard {
	var (
		col = database.TransactionCol()
	)
	match := bson.M{
		"$match": cond,
	}

	group := bson.M{
		"$group": bson.M{
			"_id": "",
			"transactionTotal": bson.M{
				"$sum": 1,
			},
			"commissionTotal": bson.M{
				"$sum": "$commission",
			},
			"sellers": bson.M{
				"$addToSet": "$sellerId",
			},
			"campaigns": bson.M{
				"$addToSet": "$campaignId",
			},
			"sellerCommission": bson.M{
				"$sum": "$sellerCommission",
			},
			"sellyCommission": bson.M{
				"$sum": "$sellyCommission",
			},
		},
	}

	cursor, err := col.Aggregate(ctx, []bson.M{match, group})

	var data []query.TransactionStatisticDashboard

	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &data); err != nil {
		fmt.Println(err.Error())
		return query.TransactionStatisticDashboard{}
	}

	if len(data) > 0 {
		return data[0]
	}
	return query.TransactionStatisticDashboard{}
}

// AggregateCampaignStatistic ...
func (d transactionImplement) AggregateCampaignStatistic(ctx context.Context, cond interface{}) []mgaffiliate.TransactionWithStatistic {
	var (
		col = database.TransactionCol()
	)
	match := bson.M{
		"$match": cond,
	}

	group := bson.M{
		"$group": bson.M{
			"_id": "$campaignId",
			"totalTransaction": bson.M{
				"$sum": 1,
			},
			"totalTransactionPending": bson.M{
				"$sum": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []string{"$status", constant.TransactionStatus.Pending.Key}},
						1,
						0,
					},
				},
			},
			"totalTransactionApproved": bson.M{
				"$sum": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []string{"$status", constant.TransactionStatus.Approved.Key}},
						1,
						0,
					},
				},
			},
			"totalTransactionRejected": bson.M{
				"$sum": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []string{"$status", constant.TransactionStatus.Rejected.Key}},
						1,
						0,
					},
				},
			},
			"totalTransactionCashback": bson.M{
				"$sum": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []string{"$status", constant.TransactionStatus.Cashback.Key}},
						1,
						0,
					},
				},
			},
		},
	}

	cursor, err := col.Aggregate(ctx, []bson.M{match, group})

	var data []mgaffiliate.TransactionWithStatistic

	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &data); err != nil {
		fmt.Println(err.Error())
	}
	return data
}

// AggregateStatisticByCondition ...
func (d transactionImplement) AggregateStatisticByCondition(ctx context.Context, cond interface{}) (result query.TransactionStatistic) {
	var dao = database.TransactionCol()

	cursor, err := dao.Aggregate(ctx, bson.A{
		bson.D{{"$match", cond}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", ""},
					{"transactionTotal", bson.D{{"$sum", 1}}},
					{"transactionCashback",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											bson.D{
												{"$eq",
													bson.A{
														"$status",
														constant.TransactionStatus.Cashback.Key,
													},
												},
											},
											1,
											0,
										},
									},
								},
							},
						},
					},
					{"transactionPending",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											bson.D{
												{"$eq",
													bson.A{
														"$status",
														constant.TransactionStatus.Pending.Key,
													},
												},
											},
											1,
											0,
										},
									},
								},
							},
						},
					},
					{"transactionApproved",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											bson.D{
												{"$eq",
													bson.A{
														"$status",
														constant.TransactionStatus.Approved.Key,
													},
												},
											},
											1,
											0,
										},
									},
								},
							},
						},
					},
					{"transactionRejected",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											bson.D{
												{"$eq",
													bson.A{
														"$status",
														constant.TransactionStatus.Rejected.Key,
													},
												},
											},
											1,
											0,
										},
									},
								},
							},
						},
					},
					{"commissionTotal", bson.D{{"$sum", "$sellerCommission"}}},
					{"commissionCashback",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											bson.D{
												{"$eq",
													bson.A{
														"$status",
														constant.TransactionStatus.Cashback.Key,
													},
												},
											},
											"$sellerCommission",
											0,
										},
									},
								},
							},
						},
					},
					{"commissionPending",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											bson.D{
												{"$eq",
													bson.A{
														"$status",
														constant.TransactionStatus.Pending.Key,
													},
												},
											},
											"$sellerCommission",
											0,
										},
									},
								},
							},
						},
					},
					{"commissionApproved",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											bson.D{
												{"$eq",
													bson.A{
														"$status",
														constant.TransactionStatus.Approved.Key,
													},
												},
											},
											"$sellerCommission",
											0,
										},
									},
								},
							},
						},
					},
					{"commissionRejected",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.A{
											bson.D{
												{"$eq",
													bson.A{
														"$status",
														constant.TransactionStatus.Rejected.Key,
													},
												},
											},
											"$sellerCommission",
											0,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	var data []query.TransactionStatistic

	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &data); err != nil {
		return
	}

	if len(data) > 0 {
		result = data[0]
	}
	return
}
