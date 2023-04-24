package responsemodel

// ResponsePlatformInfos ...
type ResponsePlatformInfos struct {
	Data []ResponsePlatformInfo `json:"data"`
}
type ResponsePlatformInfo struct {
	ID       string                  `json:"_id"`
	Code     string                  `json:"code"`
	Status   string                  `json:"status"`
	URL      string                  `json:"url"`
	Partner  ResponsePlatformPartner `json:"partner"`
	Platform string                  `json:"platform"`
}

// ResponsePlatformPartner ...
type ResponsePlatformPartner struct {
	Source     string `json:"source"`
	CampaignID string `json:"campaignId"`
}
