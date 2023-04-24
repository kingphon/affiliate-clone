package requestmomdel

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// TransactionAll ...
type TransactionAll struct {
	Page     int64  `query:"page"`
	Limit    int64  `query:"limit"`
	Keyword  string `query:"keyword"`
	Status   string `query:"status"`
	FromAt   string `query:"fromAt"`
	ToAt     string `query:"toAt"`
	Source   string `query:"source"`
	Campaign string `query:"campaign"`
	Seller   string `query:"seller"`
}

// Validate ...
func (m TransactionAll) Validate() error {
	return validation.ValidateStruct(&m)
}

// TransactionStatistic ...
type TransactionStatistic struct {
	CampaignIds string `query:"campaignIds"`
	Source      string `query:"source"`
	FromAt      string `query:"fromAt"`
	ToAt        string `query:"toAt"`
	Status      string `query:"status"`
	Seller      string `query:"seller"`
}

// Validate ...
func (m TransactionStatistic) Validate() error {
	return validation.ValidateStruct(&m)
}

// AdminTransactionCrawl ...
type AdminTransactionCrawl struct {
}
