package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

const (
	PlatformNotFound  = "platform_not_found"
	PlatformIsExisted = "platform_is_existed"
	PlatformIDInvalid = "platform_id_invalid"
	ChecksumInvalid   = "checksum_invalid"
	CodeInvalid       = "code_invalid"
)

var platform = []response.Code{
	{
		Key:     PlatformNotFound,
		Message: "platform không hợp lệ",
		Code:    301,
	},
	{
		Key:     PlatformIsExisted,
		Message: "platform đã tồn tại",
		Code:    302,
	},
	{
		Key:     PlatformIDInvalid,
		Message: "platform không hợp lệ",
		Code:    303,
	},
	{
		Key:     ChecksumInvalid,
		Message: "Checksum không hợp lệ",
		Code:    304,
	},
	{
		Key:     CodeInvalid,
		Message: "Mã không hợp lệ",
		Code:    305,
	},
}
