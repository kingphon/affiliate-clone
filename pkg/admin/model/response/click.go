package responsemodel

import "git.selly.red/Selly-Server/affiliate/external/utils/ptime"

// ResponseClickAll ...
type ResponseClickAll struct {
	List  []ResponseClickBrief `json:"list"`
	Total int64                `json:"total"`
	Limit int64                `json:"limit"`
}

// ResponseClickBrief ...
type ResponseClickBrief struct {
	ID            string                     `json:"_id"`
	Campaign      ResponseCampaignShort      `json:"campaign"`
	Seller        ResponseSellerShort        `json:"seller"`
	PartnerSource string                     `json:"partnerSource"`
	AffiliateURL  string                     `json:"affiliateURL"`
	ShareURL      string                     `json:"shareURL"`
	CampaignURL   string                     `json:"campaignURL"`
	Status        string                     `json:"status"`
	Device        ResponseTransactionDevice  `json:"device"`
	CreatedAt     *ptime.TimeResponse        `json:"createdAt"`
	From          string                     `json:"from"`
	FinalDetected ResponseClickFinalDetected `json:"finalDetected"`
}

// ResponseClickFinalDetected ...
type ResponseClickFinalDetected struct {
	URL   string `json:"URL"`
	Click string `json:"Click"`
}

// ResponseClickStatistic ...
type ResponseClickStatistic struct {
	Total          int64 `json:"total"`
	TotalPending   int64 `json:"totalPending"`
	TotalCompleted int64 `json:"totalCompleted"`
}
