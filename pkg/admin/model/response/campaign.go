package responsemodel

import (
	"git.selly.red/Selly-Server/affiliate/external/utils/file"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
)

// ResponseCampaignAll ...
type ResponseCampaignAll struct {
	List  []ResponseCampaignBrief `json:"list"`
	Total int64                   `json:"total"`
	Limit int64                   `json:"limit"`
}

// ResponseCampaignBrief ...
type ResponseCampaignBrief struct {
	ID                   string                           `json:"_id"`
	Name                 string                           `json:"name"`
	Logo                 *file.FilePhoto                  `json:"logo"`
	Desc                 string                           `json:"desc"`
	Commission           ResponseCampaignCommission       `json:"commission"`
	EstimateCashback     ResponseCampaignEstimateCashback `json:"estimateCashback"`
	Platforms            []string                         `json:"platforms"`
	Order                int                              `json:"order"`
	Status               string                           `json:"status"`
	AllowShowShareAction bool                             `json:"allowShowShareAction"`
	CreatedAt            *ptime.TimeResponse              `json:"createdAt"`
}

// ResponseCampaignCommission ...
type ResponseCampaignCommission struct {
	Real          float64 `json:"real"`
	SellerPercent float64 `json:"sellerPercent"`
	Selly         float64 `json:"selly"`
	Seller        float64 `json:"seller"`
}

// ResponseCampaignEstimateCashback ...
type ResponseCampaignEstimateCashback struct {
	Day       int `json:"day"`
	NextMonth int `json:"nextMonth"`
}

// ResponseCampaignDetail ...
type ResponseCampaignDetail struct {
	ID                   string                           `json:"_id"`
	Name                 string                           `json:"name"`
	Logo                 *file.FilePhoto                  `json:"logo"`
	Covers               []*file.FilePhoto                `json:"covers"`
	Desc                 string                           `json:"desc"`
	Commission           ResponseCampaignCommission       `json:"commission"`
	From                 *ptime.TimeResponse              `json:"from"`
	To                   *ptime.TimeResponse              `json:"to"`
	Status               string                           `json:"status"`
	EstimateCashback     ResponseCampaignEstimateCashback `json:"estimateCashback"`
	ShareDesc            string                           `json:"shareDesc"`
	CreatedAt            *ptime.TimeResponse              `json:"createdAt"`
	Order                int                              `json:"order"`
	Platforms            []string                         `json:"platforms"`
	AllowShowShareAction bool                             `json:"allowShowShareAction"`
}

// ResponseCampaignShortInfo ...
type ResponseCampaignShortInfo struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}
