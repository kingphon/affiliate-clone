package responsemodel

import (
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
)

// ResponseTransactionGetSummary ...
type ResponseTransactionGetSummary struct {
	Data []LabelValue `json:"data"`
}

// ResponseTransactionStatistic ...
type ResponseTransactionStatistic struct {
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

// NewResponseTransactionStatistic ...
func NewResponseTransactionStatistic(statistic mgaffiliate.Statistic) ResponseTransactionStatistic {
	return ResponseTransactionStatistic{
		TransactionTotal:    statistic.TransactionTotal,
		TransactionCashback: statistic.TransactionCashback,
		TransactionPending:  statistic.TransactionPending,
		TransactionApproved: statistic.TransactionApproved,
		TransactionRejected: statistic.TransactionRejected,
		CommissionTotal:     statistic.CommissionTotal,
		CommissionCashback:  statistic.CommissionCashback,
		CommissionPending:   statistic.CommissionPending,
		CommissionApproved:  statistic.CommissionApproved,
		CommissionRejected:  statistic.CommissionRejected,
	}
}
