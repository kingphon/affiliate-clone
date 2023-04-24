package responsemodel

// ResponseListCampaignGroupAll ...
type ResponseListCampaignGroupAll struct {
	Data []ResponseCampaignGroup `json:"data"`
}

// ResponseCampaignGroup ...
type ResponseCampaignGroup struct {
	ID    string                  `json:"_id"`
	Type  string                  `json:"type"`
	Name  string                  `json:"name"`
	Items []ResponseCampaignBrief `json:"items"`
}

// ResponseCampaignGroupBrief ...
type ResponseCampaignGroupBrief struct {
	ID   string `json:"_id"`
	Type string `json:"type"`
	Name string `json:"name"`
}
