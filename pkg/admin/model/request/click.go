package requestmomdel

import validation "github.com/go-ozzo/ozzo-validation/v4"

// ClickAll ...
type ClickAll struct {
	Page     int64  `query:"page"`
	Limit    int64  `query:"limit"`
	Keyword  string `query:"keyword"`
	Status   string `query:"status"`
	Source   string `query:"source"`
	FromAt   string `query:"fromAt"`
	ToAt     string `query:"toAt"`
	Campaign string `query:"campaign"`
	Seller   string `query:"seller"`
	From     string `query:"from"`
}

func (m ClickAll) Validate() error {
	return validation.ValidateStruct(&m)
}

// ClickStatistic ...
type ClickStatistic struct {
	CampaignIds string `query:"campaignIds"`
	Source      string `query:"source"`
	FromAt      string `query:"fromAt"`
	ToAt        string `query:"toAt"`
	Seller      string `query:"seller"`
}

func (m ClickStatistic) Validate() error {
	return validation.ValidateStruct(&m)
}
