package service

import (
	"context"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"

	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetHistoriesByTransactionID ...
func (s transactionHistoryImplement) GetHistoriesByTransactionID(ctx context.Context, transactionID primitive.ObjectID) []responsemodel.ResponseTransactionHistory {
	var (
		d    = dao.TransactionHistory()
		cond = bson.D{
			{"transactionId", transactionID},
		}
		opts   = &options.FindOptions{Sort: bson.D{{"createdAt", -1}}}
		result = make([]responsemodel.ResponseTransactionHistory, 0)
	)

	histories := d.FindByCondition(ctx, cond, opts)
	for _, history := range histories {
		historyInfo := s.brief(ctx, history)
		result = append(result, historyInfo)
	}

	return result
}

//
// PRIVATE METHOD
//

// brief ...
func (transactionHistoryImplement) brief(ctx context.Context, history mgaffiliate.TransactionHistory) responsemodel.ResponseTransactionHistory {
	return responsemodel.ResponseTransactionHistory{
		ID:        history.ID,
		Status:    history.Status,
		Desc:      history.Desc,
		CreatedAt: ptime.TimeResponseInit(history.CreatedAt),
	}
}
