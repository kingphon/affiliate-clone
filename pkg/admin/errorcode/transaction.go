package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

const (
	TransactionNotFound      = "transaction_not_found"
	TransactionInvalidID     = "transaction_invalid_id"
	TransactionInvalidStatus = "transaction_invalid_status"
	TransactionInvalidTime   = "transaction_invalid_time"
)

var transaction = []response.Code{
	{
		Key:     TransactionNotFound,
		Message: "lượt thưởng không tồn tại",
		Code:    500,
	},
	{
		Key:     TransactionInvalidID,
		Message: "id lượt thưởng không hợp lệ",
		Code:    501,
	},
	{
		Key:     TransactionInvalidStatus,
		Message: "trạng thái lượt thưởng không hợp lệ",
		Code:    502,
	},
	{
		Key:     TransactionInvalidTime,
		Message: "thời gian phát sinh không hợp lệ",
		Code:    503,
	},
}
