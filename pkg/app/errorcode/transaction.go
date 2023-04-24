package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

const (
	TransactionNotFound = "transaction_not_found"
)

var transaction = []response.Code{
	{
		Key:     TransactionNotFound,
		Message: "đơn hàng không hợp lệ",
		Code:    401,
	},
}
