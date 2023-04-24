package service

import (
	"context"
	"fmt"
	"sync"

	appconstant "git.selly.red/Selly-Server/affiliate/pkg/app/constant"

	"git.selly.red/Selly-Server/affiliate/pkg/app/locale"

	"git.selly.red/Selly-Server/affiliate/external/utils/format"

	"git.selly.red/Selly-Server/affiliate/external/constant"

	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
	"github.com/friendsofgo/errors"

	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/app/dao"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/external/utils/pagetoken"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/response"
)

// TransactionInterface ...
type TransactionInterface interface {
	// All return transaction with condition ...
	All(ctx context.Context, q mgquery.AppQuery, sellerID primitive.ObjectID) (result responsemodel.ResponseList)

	// Detail ...
	Detail(ctx context.Context, id, sellerID primitive.ObjectID) (result *responsemodel.ResponseDetail, err error)

	// GetFilter ...
	GetFilter(ctx context.Context) responsemodel.ResponseTransactionFilter

	// GetSummary ...
	GetSummary(ctx context.Context, sellerID primitive.ObjectID) responsemodel.ResponseTransactionGetSummary

	// GetStatistic ...
	GetStatistic(ctx context.Context, sellerID primitive.ObjectID) responsemodel.ResponseTransactionStatistic

	// CountByCondition ...
	CountByCondition(ctx context.Context, cond interface{}) int64
}

// Transaction return transaction service
func Transaction() TransactionInterface {
	return transactionImplement{}
}

// transactionImplement ...
type transactionImplement struct{}

// CountByCondition ...
func (s transactionImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	var d = dao.Transaction()
	count := d.CountByCondition(ctx, cond)
	return count
}

// GetStatistic ...
func (s transactionImplement) GetStatistic(ctx context.Context, sellerID primitive.ObjectID) responsemodel.ResponseTransactionStatistic {
	sellerStatistic := dao.SellerStatistic().FindOneByCondition(ctx, bson.M{
		"sellerId": sellerID,
	})
	if sellerStatistic.ID.IsZero() {
		return responsemodel.ResponseTransactionStatistic{}
	}
	return responsemodel.NewResponseTransactionStatistic(sellerStatistic.Statistic)
}

// GetSummary ...
func (s transactionImplement) GetSummary(ctx context.Context, sellerID primitive.ObjectID) responsemodel.ResponseTransactionGetSummary {
	var (
		d    = dao.SellerStatistic()
		cond = bson.D{{"sellerId", sellerID}}
	)

	doc := d.FindOneByCondition(ctx, cond)
	var data = []responsemodel.LabelValue{
		{
			Label: locale.SummaryTitleTransactionTotal,
			Value: format.ConvertInt4ToCurrency(
				doc.Statistic.TransactionPending +
					doc.Statistic.TransactionApproved +
					doc.Statistic.TransactionCashback),
		},
		{
			Label: locale.SummaryTitleTransactionCashback,
			Value: format.ConvertInt4ToCurrency(doc.Statistic.TransactionCashback),
		},
		{
			Label: locale.SummaryTitleCommissionCashback,
			Value: format.ToCurrencyVND(doc.Statistic.CommissionCashback),
		},
		{
			Label: locale.SummaryTitleCommissionApproved,
			Value: format.ToCurrencyVND(doc.Statistic.CommissionApproved),
		},
	}

	return responsemodel.ResponseTransactionGetSummary{
		Data: data,
	}
}

// GetFilter ...
func (transactionImplement) GetFilter(ctx context.Context) responsemodel.ResponseTransactionFilter {
	var transactionStatus = constant.TransactionStatus
	return responsemodel.ResponseTransactionFilter{
		Data: []responsemodel.KeyValue{
			{
				Key:   transactionStatus.All.Key,
				Value: transactionStatus.All.Title,
			},
			{
				Key:   transactionStatus.Pending.Key,
				Value: transactionStatus.Pending.Title,
			},
			{
				Key:   transactionStatus.Approved.Key,
				Value: transactionStatus.Approved.Title,
			},
			{
				Key:   transactionStatus.Cashback.Key,
				Value: transactionStatus.Cashback.Title,
			},
			{
				Key:   transactionStatus.Rejected.Key,
				Value: transactionStatus.Rejected.Title,
			},
		},
	}
}

// Detail ...
func (s transactionImplement) Detail(ctx context.Context, id, sellerID primitive.ObjectID) (result *responsemodel.ResponseDetail, err error) {
	var cond = bson.D{
		{"_id", id},
		{"sellerId", sellerID},
	}

	transaction := dao.Transaction().FindOneByCondition(ctx, cond)
	if transaction.ID.IsZero() {
		err = errors.New(errorcode.TransactionNotFound)
		return
	}

	data := s.detail(ctx, transaction)
	result = &responsemodel.ResponseDetail{
		Data: data,
	}
	return
}

// All ...
func (s transactionImplement) All(ctx context.Context, q mgquery.AppQuery, sellerID primitive.ObjectID) (result responsemodel.ResponseList) {
	// 1. Init value
	var list = make([]responsemodel.ResponseTransactionBrief, 0)

	// 2. Assign condition
	var cond = bson.D{
		{"sellerId", sellerID},
	}

	q.AssignStatus(&cond)

	// 3. Find
	docs := dao.Transaction().FindByCondition(ctx, cond, q.GetFindOptionsWithPage())

	// 4. Convert response
	for _, doc := range docs {
		list = append(list, s.brief(ctx, doc))
	}

	// Page token
	endData := len(list) < int(q.Limit)
	var nextPageToken = ""
	if len(list) == int(q.Limit) {
		nextPageToken = pagetoken.PageTokenUsingPage(int(q.Page) + 1)
	}

	// Response
	result = responsemodel.ResponseList{
		List:          list,
		EndData:       endData,
		NextPageToken: nextPageToken,
	}
	return
}

//
// PRIVATE METHODS
//

// getStatusInfoByStatus ...
func (transactionImplement) getStatusInfoByStatus(status string) responsemodel.ResponseTransactionStatus {
	var (
		text, color    string
		statusConstant = constant.TransactionStatus
	)

	switch status {
	case statusConstant.Pending.Key:
		text = statusConstant.Pending.Title
		color = statusConstant.Pending.Color
	case statusConstant.Approved.Key:
		text = statusConstant.Approved.Title
		color = statusConstant.Approved.Color
	case statusConstant.Cashback.Key:
		text = statusConstant.Cashback.Title
		color = statusConstant.Cashback.Color
	case statusConstant.Rejected.Key:
		text = statusConstant.Rejected.Title
		color = statusConstant.Rejected.Color
	}
	return responsemodel.ResponseTransactionStatus{
		Type:  status,
		Text:  text,
		Color: color,
	}
}

// getProcessIconsByStatus ...
func (transactionImplement) getProcessIconsByStatus(status string) (result []appconstant.TransactionProcessIcon) {
	var transactionIconList = appconstant.TransactionProcessIcons
	switch status {
	case constant.TransactionStatus.Pending.Key:
		result = transactionIconList.Pending
	case constant.TransactionStatus.Approved.Key:
		result = transactionIconList.Approved
	case constant.TransactionStatus.Cashback.Key:
		result = transactionIconList.Cashback
	case constant.TransactionStatus.Rejected.Key:
		result = transactionIconList.Rejected
	}

	for _, icon := range result {
		icon.Icon = icon.Icon.GetResponseData()
	}
	return
}

// brief ...
func (s transactionImplement) brief(ctx context.Context, doc mgaffiliate.Transaction) responsemodel.ResponseTransactionBrief {
	campaignInfo := Campaign().GetShortInfoByID(ctx, doc.CampaignID)
	campaignInfo.Commission = doc.SellerCommission
	return responsemodel.ResponseTransactionBrief{
		ID:              doc.ID,
		Code:            doc.Code,
		Commission:      s.getSellerCommission(doc),
		Source:          doc.Source,
		Status:          s.getStatusInfoByStatus(doc.Status),
		TransactionTime: ptime.TimeResponseInit(doc.TransactionTime),
		Campaign:        campaignInfo,
		RejectedReason:  doc.RejectedReason,
	}
}

// getSellerCommission ...
func (s transactionImplement) getSellerCommission(doc mgaffiliate.Transaction) (commission float64) {
	commission = doc.SellerCommission
	if doc.IsValidStatus(constant.TransactionStatus.Pending.Key) {
		commission = doc.EstimateSellerCommission
	}
	return
}

// detail ...
func (s transactionImplement) detail(ctx context.Context, transaction mgaffiliate.Transaction) (result responsemodel.ResponseTransactionDetail) {
	var (
		wg          = sync.WaitGroup{}
		campaign    responsemodel.ResponseCampaignShortInfo
		lastHistory responsemodel.ResponseTransactionHistory
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		campaign = campaignImplement{}.GetShortInfoByID(ctx, transaction.CampaignID)
		campaign.Commission = transaction.SellerCommission
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		lastHistory = transactionHistoryImplement{}.GetLastDocBriefByTransactionID(ctx, transaction.ID)
	}()
	wg.Wait()

	result = responsemodel.ResponseTransactionDetail{
		ID:                 transaction.ID,
		Code:               transaction.Code,
		TransactionTime:    ptime.TimeResponseInit(transaction.TransactionTime),
		Source:             transaction.Source,
		Commission:         s.getSellerCommission(transaction),
		EstimateCashbackAt: ptime.TimeResponseInit(transaction.EstimateCashbackAt),
		Status:             s.getStatusInfoByStatus(transaction.Status),
		Device: responsemodel.ResponseTransactionDevice{
			Model:          transaction.Device.Model,
			UserAgent:      transaction.Device.UserAgent,
			Manufacturer:   transaction.Device.Manufacturer,
			OSName:         transaction.Device.OSName,
			OSVersion:      transaction.Device.OSVersion,
			BrowserVersion: transaction.Device.BrowserVersion,
			BrowserName:    transaction.Device.BrowserName,
			DeviceType:     transaction.Device.DeviceType,
			DeviceID:       transaction.Device.DeviceID,
		},
		DeviceText:     s.generateDeviceText(transaction.Device),
		Campaign:       campaign,
		ProcessIcons:   s.getProcessIconsByStatus(transaction.Status),
		LastHistory:    lastHistory,
		RejectedReason: transaction.RejectedReason,
	}

	result.SupportChannel = result.GetSupportChannel()
	return
}

// generateDeviceText ...
func (s transactionImplement) generateDeviceText(device mgaffiliate.Device) string {
	var text = ""
	if device.Manufacturer != "" {
		text = fmt.Sprintf("%s: %s\n", locale.DeviceTitleManufacturer, device.Manufacturer)
	}

	if device.DeviceType != "" || device.Model != "" {
		deviceType := device.Model
		if device.DeviceType != "" {
			deviceType = device.DeviceType + " - " + device.Model
		}
		text += fmt.Sprintf("%s: %s\n", locale.DeviceTitleDeviceType, deviceType)
	}

	if device.OSName != "" && device.OSVersion != "" {
		text += fmt.Sprintf("%s: %s - %s\n", locale.DeviceTitleOsName, device.OSName, device.OSVersion)
	}

	if device.BrowserName != "" && device.BrowserVersion != "" {
		text += fmt.Sprintf("%s: %s - %s", locale.DeviceTitleBrowserName, device.BrowserName, device.BrowserVersion)
	}
	return text
}
