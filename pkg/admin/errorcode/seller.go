package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

const (
	SellerNotFound  = "seller_not_found"
	SellerInvalidID = "seller_invalid_id"
)

var seller = []response.Code{
	{
		Key:     SellerNotFound,
		Message: "seller không tồn tại",
		Code:    900,
	}, {
		Key:     SellerInvalidID,
		Message: "id seller không hợp lệ",
		Code:    901,
	},
}
