package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

const (
	CampaignIsRequiredName          = "campaign_is_required_name"
	CampaignIsRequiredOrder         = "campaign_is_required_order"
	CampaignIsRequiredFrom          = "campaign_is_required_from"
	CampaignIsRequiredTo            = "campaign_is_required_to"
	CampaignIsRequiredReal          = "campaign_is_required_real"
	CampaignIsRequiredSellerPercent = "campaign_is_required_seller_percent"
	CampaignInvalidReal             = "campaign_invalid_Real"
	CampaignInValidSellerPercent    = "campaign_invalid_seller_percent"
	CampaignErrorWhenCreated        = "campaign_error_when_created"
	CampaignNotFound                = "campaign_not_found"
	CampaignInvalidDay              = "campaign_invalid_day"
	CampaignInvalidNextMonth        = "campaign_invalid_next_month"
	CampaignErrorWhenUpdateStatus   = "campaign_error_when_update_status"
)

var campaign = []response.Code{
	{
		Key:     CampaignIsRequiredName,
		Message: "tên sản phẩm dịch vụ không được trống",
		Code:    101,
	},
	{
		Key:     CampaignIsRequiredOrder,
		Message: "số thứ tự sắp xếp không được trống",
		Code:    102,
	},
	{
		Key:     CampaignIsRequiredFrom,
		Message: "thời gian phát hành không được trống",
		Code:    103,
	},
	{
		Key:     CampaignIsRequiredTo,
		Message: "thời gian phát hành không được trống",
		Code:    104,
	},
	{
		Key:     CampaignIsRequiredReal,
		Message: "tiền nhận từ đối tác  không được trống",
		Code:    105,
	},
	{
		Key:     CampaignIsRequiredSellerPercent,
		Message: "phần trăm cho seller không được trống",
		Code:    106,
	},
	{
		Key:     CampaignInvalidReal,
		Message: "tiền nhận từ đối tác phải lớn hơn 0",
		Code:    107,
	},
	{
		Key:     CampaignInValidSellerPercent,
		Message: "phần trăm cho seller chỉ từ 1 đến 100",
		Code:    108,
	},
	{
		Key:     CampaignErrorWhenCreated,
		Message: "Lỗi khi tạo sản phẩm dịch vụ",
		Code:    109,
	},
	{
		Key:     CampaignNotFound,
		Message: "không tìm thấy sản phẩm dịch vụ",
		Code:    110,
	},
	{
		Key:     CampaignInvalidDay,
		Message: "ngày đối soát chỉ từ 1 đến 31",
		Code:    111,
	}, {
		Key:     CampaignInvalidNextMonth,
		Message: "số tháng đối soát phải lớn hơn 0",
		Code:    112,
	},
	{
		Key:     CampaignErrorWhenUpdateStatus,
		Message: "không thể thay đổi trạng thái do không có platform nào đang hoạt động",
		Code:    113,
	},
}
