package responsemodel

import (
	natsmodel "git.selly.red/Selly-Modules/natsio/model"
	"git.selly.red/Selly-Server/affiliate/external/utils/file"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
)

// ResponseTransactionAll ...
type ResponseTransactionAll struct {
	List  []ResponseTransactionBrief `json:"list"`
	Total int64                      `json:"total"`
	Limit int64                      `json:"limit"`
}

// ResponseTransactionBrief ...
type ResponseTransactionBrief struct {
	ID                       string                     `json:"_id"`
	Code                     string                     `json:"code"`
	Campaign                 ResponseCampaignShort      `json:"campaign"`
	Seller                   ResponseSellerShort        `json:"seller"`
	Source                   string                     `json:"source"`
	Commission               ResponseCampaignCommission `json:"commission"`
	EstimateSellerCommission float64                    `json:"estimateSellerCommission"`
	TransactionTime          *ptime.TimeResponse        `json:"transactionTime"`
	Status                   string                     `json:"status"`
	RejectedReason           string                     `json:"rejectedReason"`
	EstimateCashbackAt       *ptime.TimeResponse        `json:"estimateCashbackAt"`
}

// ResponseCampaignShort ...
type ResponseCampaignShort struct {
	ID   string          `json:"_id"`
	Name string          `json:"name"`
	Logo *file.FilePhoto `json:"logo"`
}

// DataSellerCampaign ...
type DataSellerCampaign struct {
	Sellers   []natsmodel.ResponseSellerInfo
	Campaigns []ResponseCampaignShort
}

// ResponseTransactionDetail ...
type ResponseTransactionDetail struct {
	ID                       string                     `json:"_id"`
	Code                     string                     `json:"code"`
	Seller                   ResponseSellerShort        `json:"seller"`
	Campaign                 ResponseCampaignShort      `json:"campaign"`
	TransactionTime          *ptime.TimeResponse        `json:"transactionTime"`
	Commission               ResponseCampaignCommission `json:"commission"`
	Device                   ResponseTransactionDevice  `json:"device"`
	Status                   string                     `json:"status"`
	EstimateCashbackAt       *ptime.TimeResponse        `json:"estimateCashbackAt"`
	EstimateSellerCommission float64                    `json:"estimateSellerCommission"`
}

// ResponseTransactionDevice ...
type ResponseTransactionDevice struct {
	Model          string `bson:"model"`
	UserAgent      string `json:"userAgent"`
	OSName         string `json:"osName"`
	OSVersion      string `json:"osVersion"`
	BrowserVersion string `json:"browserVersion"`
	BrowserName    string `json:"browserName"`
	DeviceType     string `json:"deviceType"`
	Manufacturer   string `bson:"manufacturer"`
	DeviceID       string `json:"deviceId"`
}

// ResponseTransactionStatistic ...
type ResponseTransactionStatistic struct {
	TotalTransaction    int64   `json:"totalTransaction"`
	TotalCommissionReal float64 `json:"totalCommissionReal"`
	SellerCommission    float64 `json:"sellerCommission"`
	SellyCommission     float64 `json:"sellyCommission"`
	TotalSeller         int64   `json:"totalSeller"`
	TotalCampaign       int64   `json:"totalCampaign"`
}
