package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

const (
	CampaignInvalidID          = "campaign_invalid_id"
	CampaignNotFound           = "campaign_not_found"
	CampaignGroupNotFound      = "campaign_group_not_found"
	CampaignFinalURLIsRequired = "campaign_final_url_is_required"
)

var campaign = []response.Code{
	{
		Key:     CampaignInvalidID,
		Message: "id campaign không hợp lệ",
		Code:    201,
	},
	{
		Key:     CampaignNotFound,
		Message: "campaign không tồn tại",
		Code:    202,
	},
	{
		Key:     CampaignGroupNotFound,
		Message: "campaign group không tồn tại",
		Code:    203,
	},
	{
		Key:     CampaignFinalURLIsRequired,
		Message: "FinalURL không thể trống",
		Code:    204,
	},
}
