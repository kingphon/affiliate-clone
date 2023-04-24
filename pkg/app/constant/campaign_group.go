package appconstant

import "git.selly.red/Selly-Server/affiliate/external/constant"

// campaignGroupType ...
type campaignGroupType struct {
	Key   string
	Title string
}

// CampaignGroupType ...
var CampaignGroupType = struct {
	CampaignList campaignGroupType
}{
	CampaignList: campaignGroupType{
		Key:   "campaign_list",
		Title: "5 phút có thưởng",
	},
}

// ScreenListAllowCampaignGroup ...
var ScreenListAllowCampaignGroup = []string{
	constant.Screen.Campaign,
}
