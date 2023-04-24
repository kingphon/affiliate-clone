package handler

import (
	"context"
	"git.selly.red/Selly-Modules/natsio"
	"git.selly.red/Selly-Modules/natsio/model"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/service"
)

// Nats ...
type Nats struct {
	EncodedConn *natsio.JSONEncoder
}

// GetTransactions ...
func (h *Nats) GetTransactions(subject, reply string, req *model.GetTransactionsRequest) {
	var (
		fileResSvc = service.FileResponse()

		ctx = context.Background()
		q   = mgquery.AppQuery{
			Page:  req.Page,
			Limit: req.Limit,
			Affiliate: mgquery.Affiliate{
				Status:   req.Status,
				Keyword:  req.Keyword,
				Source:   req.Source,
				Campaign: req.Campaign,
				Seller:   req.Seller,
				FromAt:   req.FromAt,
				ToAt:     req.ToAt,
			},
		}
	)

	s := service.Transaction()
	data := s.All(ctx, q)
	res := model.GetTransactionsResponse{
		Total: data.Total,
		Limit: data.Limit,
		List:  make([]model.TransactionInfo, len(data.List)),
	}

	for i, item := range data.List {
		res.List[i] = model.TransactionInfo{
			ID:   item.ID,
			Code: item.Code,
			Campaign: model.ResponseCampaignShort{
				ID:   item.Campaign.ID,
				Name: item.Campaign.Name,
				Logo: fileResSvc.ConvertResponseFilePhotoNats(item.Campaign.Logo),
			},
			Seller: model.ResponseSellerInfo{
				ID:   item.Seller.ID,
				Name: item.Seller.Name,
			},
			Source: item.Source,
			Commission: model.ResponseCampaignCommission{
				Real:          item.Commission.Real,
				SellerPercent: item.Commission.SellerPercent,
				Selly:         item.Commission.Selly,
				Seller:        item.Commission.Seller,
			},
			EstimateSellerCommission: item.EstimateSellerCommission,
			TransactionTime:          item.TransactionTime.FormatISODate(),
			Status:                   item.Status,
			RejectedReason:           item.RejectedReason,
			EstimateCashbackAt:       item.EstimateCashbackAt.FormatISODate(),
		}

	}

	h.response(reply, res, nil)
}

// response ...
func (h *Nats) response(reply string, data interface{}, err error) {
	res := map[string]interface{}{"data": data}
	if err != nil {
		res["error"] = err.Error()
	}
	h.EncodedConn.Publish(reply, res)
}
