package requestmomdel

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SandboxGenerateTransaction ...
type SandboxGenerateTransaction struct {
	CampaignID string `json:"campaignId" query:"campaignId"`
	SellerID   string `json:"sellerId" query:"sellerId"`
}

// Validate ...
func (m SandboxGenerateTransaction) Validate() error {
	return validation.ValidateStruct(&m)
}

// SandboxCrawlerTransaction ...
type SandboxCrawlerTransaction struct {
	Campaign        string  `json:"campaign"`
	ClickId         string  `json:"clickId"`
	Source          string  `json:"source"`
	Status          string  `json:"status"`
	Commission      float64 `json:"commission"`
	TransactionID   string  `json:"transactionId"`
	TransactionTime string  `json:"transactionTime"`
	RejectedReason  string  `json:"rejectedReason"`
	UpdatedAt       string  `json:"updatedAt"`
	UpdatedHash     string  `json:"updatedHash"`
}

// Validate ...
func (m SandboxCrawlerTransaction) Validate() error {
	return validation.ValidateStruct(&m)
}
