package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

const (
	PlatformIsRequiredName      = "platform_is_required_name"
	PlatformIsRequiredPartner   = "platform_is_required_partner"
	PlatformIsRequiredUrl       = "platform_is_required_url"
	PlatformIsRequiredPlatform  = "platform_is_required_platform"
	PlatformIsRequiredSource    = "platform_is_required_source"
	PlatformInvalidSource       = "platform_invalid_source"
	PlatformInvalidPlatform     = "platform_invalid_platform"
	PlatformErrorWhenCreated    = "platform_error_when_created"
	PlatformNotFound            = "platform_not_found"
	PlatformStatusInvalid       = "platform_invalid_status"
	PlatformIsInvalidCampaignID = "platform_invalid_campaign_id"
	PlatformIsExisted           = "platform_is_existed"
)

var platform = []response.Code{
	{
		Key:     PlatformIsRequiredName,
		Message: "tên platform không được trống",
		Code:    300,
	},
	{
		Key:     PlatformIsRequiredPartner,
		Message: "đối tác không được trống",
		Code:    301,
	},
	{
		Key:     PlatformIsRequiredUrl,
		Message: "url không được trống",
		Code:    302,
	},
	{
		Key:     PlatformIsRequiredPlatform,
		Message: "platform không được trống",
		Code:    303,
	},
	{
		Key:     PlatformIsRequiredSource,
		Message: "source không được trống",
		Code:    304,
	},
	{
		Key:     PlatformInvalidPlatform,
		Message: "platform không hợp lệ",
		Code:    305,
	},
	{
		Key:     PlatformErrorWhenCreated,
		Message: "platform gặp lỗi khi tạo",
		Code:    306,
	},
	{
		Key:     PlatformNotFound,
		Message: "platform không tìm thấy",
		Code:    307,
	},
	{
		Key:     PlatformStatusInvalid,
		Message: "trạng thái không hợp lệ",
		Code:    308,
	},
	{
		Key:     PlatformIsInvalidCampaignID,
		Message: "campaign id của đối tác không hợp lệ",
		Code:    309,
	},
	{
		Key:     PlatformIsExisted,
		Message: "platform đã tồn tại",
		Code:    310,
	},
	{
		Key:     PlatformInvalidSource,
		Message: "source không hợp lệ",
		Code:    221,
	},
}
