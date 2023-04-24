package query

import "go.mongodb.org/mongo-driver/bson/primitive"

// TransactionStatistic ...
type TransactionStatistic struct {
	TransactionTotal    int64   `json:"transactionTotal"`
	TransactionCashback int64   `json:"transactionCashback"`
	TransactionPending  int64   `json:"transactionPending"`
	TransactionApproved int64   `json:"transactionApproved"`
	TransactionRejected int64   `json:"transactionRejected"`
	CommissionTotal     float64 `json:"commissionTotal"`
	CommissionCashback  float64 `json:"commissionCashback"`
	CommissionPending   float64 `json:"commissionPending"`
	CommissionApproved  float64 `json:"commissionApproved"`
	CommissionRejected  float64 `json:"commissionRejected"`
}

// TransactionStatisticDashboard ...
type TransactionStatisticDashboard struct {
	TransactionTotal int64                `bson:"transactionTotal"`
	CommissionTotal  float64              `bson:"commissionTotal"`
	SellerCommission float64              `bson:"sellerCommission"`
	SellyCommission  float64              `bson:"sellyCommission"`
	Sellers          []primitive.ObjectID `bson:"sellers"`
	Campaigns        []primitive.ObjectID `bson:"campaigns"`
}
