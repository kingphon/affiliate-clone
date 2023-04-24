package service

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"

	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/pkg/app/dao"

	"go.mongodb.org/mongo-driver/bson"

	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransactionHistoryInterface ...
type TransactionHistoryInterface interface {
	GetHistories(ctx context.Context, transactionID, sellerID primitive.ObjectID) (responsemodel.ResponseTransactionHistories, error)

	GetLastDocBriefByTransactionID(ctx context.Context, transactionID primitive.ObjectID) responsemodel.ResponseTransactionHistory
}

// TransactionHistory return transaction history service
func TransactionHistory() TransactionHistoryInterface {
	return transactionHistoryImplement{}
}

// transactionHistoryImplement ...
type transactionHistoryImplement struct{}

// GetHistories ...
func (s transactionHistoryImplement) GetHistories(ctx context.Context, transactionID, sellerID primitive.ObjectID) (data responsemodel.ResponseTransactionHistories, err error) {
	// 1. Init value
	var list = make([]responsemodel.ResponseTransactionHistory, 0)

	var (
		cond = bson.D{
			{"transactionId", transactionID},
			{"sellerId", sellerID},
		}
		q = mgquery.AppQuery{
			SortInterface: bson.D{{"createdAt", -1}},
		}
		d = dao.TransactionHistory()
	)

	docs := d.FindByCondition(ctx, cond, q.GetFindOptionsWithSort())
	for _, doc := range docs {
		list = append(list, s.brief(ctx, doc))
	}

	data.Data = list
	return
}

// GetLastDocBriefByTransactionID ...
func (s transactionHistoryImplement) GetLastDocBriefByTransactionID(ctx context.Context, transactionID primitive.ObjectID) responsemodel.ResponseTransactionHistory {
	var (
		d    = dao.TransactionHistory()
		cond = bson.D{{"transactionId", transactionID}}
		opts = &options.FindOneOptions{
			Sort: bson.D{{"createdAt", -1}},
		}
	)

	doc := d.FindOneByCondition(ctx, cond, opts)
	return s.brief(ctx, doc)
}

//
// PRIVATE METHOD
//

// brief ...
func (transactionHistoryImplement) brief(ctx context.Context, doc mgaffiliate.TransactionHistory) responsemodel.ResponseTransactionHistory {
	return responsemodel.ResponseTransactionHistory{
		ID:        doc.ID,
		Status:    doc.Status,
		Desc:      doc.Desc,
		CreatedAt: ptime.TimeResponseInit(doc.CreatedAt),
	}
}
