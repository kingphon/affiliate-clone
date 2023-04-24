package responsemodel

import (
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseTransactionHistories ...
type ResponseTransactionHistories struct {
	Data []ResponseTransactionHistory `json:"data"`
}

// ResponseTransactionHistory ...
type ResponseTransactionHistory struct {
	ID        primitive.ObjectID  `json:"_id"`
	Status    string              `json:"status"`
	Desc      string              `json:"desc"`
	CreatedAt *ptime.TimeResponse `json:"createdAt"`
}
