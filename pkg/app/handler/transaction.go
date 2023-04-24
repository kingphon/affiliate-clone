package handler

import (
	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/request"
	"git.selly.red/Selly-Server/affiliate/pkg/app/service"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction ...
type Transaction struct{}

// GetStatistic godoc
// @tags Transaction
// @summary GetStatistic
// @id app-transaction-statistic-by-me
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @success 200 {object} responsemodel.ResponseTransactionStatistic
// @router /transactions/statistic [get]
func (Transaction) GetStatistic(c echo.Context) error {
	var (
		ctx         = echocontext.GetContext(c)
		s           = service.Transaction()
		sellerID, _ = mongodb.NewIDFromString(echocontext.GetCurrentUserID(c))
	)
	statistic := s.GetStatistic(ctx, sellerID)
	return response.R200(c, statistic, "")
}

// All godoc
// @tags Transaction
// @summary All
// @id app-transaction-list
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.TransactionAll false "Query"
// @success 200 {object} responsemodel.ResponseList
// @router /transactions [get]
func (Transaction) All(c echo.Context) error {
	var (
		ctx       = echocontext.GetContext(c)
		s         = service.Transaction()
		qParams   = echocontext.GetQuery(c).(requestmodel.TransactionAll)
		pageToken = echocontext.PageTokenDecode(qParams.PageToken)
		q         = mgquery.AppQuery{
			Page:          pageToken.Page,
			Limit:         int64(constant.Limit20),
			Keyword:       qParams.Keyword,
			Status:        qParams.Status,
			SortInterface: bson.D{{"transactionTime", -1}},
		}
		sellerID, _ = mongodb.NewIDFromString(echocontext.GetCurrentUserID(c))
	)

	result := s.All(ctx, q, sellerID)
	return response.R200(c, result, "")
}

// Detail godoc
// @tags Transaction
// @summary Detail
// @id app-transaction-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Transaction id"
// @success 200 {object} responsemodel.ResponseDetail
// @router /transactions/{id} [get]
func (Transaction) Detail(c echo.Context) error {
	var (
		ctx         = echocontext.GetContext(c)
		s           = service.Transaction()
		id          = echocontext.GetParam(c, "id").(primitive.ObjectID)
		sellerID, _ = mongodb.NewIDFromString(echocontext.GetCurrentUserID(c))
	)

	result, err := s.Detail(ctx, id, sellerID)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}

// GetHistories godoc
// @tags Transaction
// @summary GetHistories
// @id app-transaction-histories
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Transaction id"
// @success 200 {object} responsemodel.ResponseTransactionHistories
// @router /transactions/{id}/histories [get]
func (Transaction) GetHistories(c echo.Context) error {
	var (
		ctx           = echocontext.GetContext(c)
		s             = service.TransactionHistory()
		transactionID = echocontext.GetParam(c, "id").(primitive.ObjectID)
		sellerID, _   = mongodb.NewIDFromString(echocontext.GetCurrentUserID(c))
	)

	result, err := s.GetHistories(ctx, transactionID, sellerID)
	if err != nil {
		return response.R400(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}

// TransactionFilter godoc
// @tags Transaction
// @summary Filter
// @id app-transaction-filter
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @success 200 {object} responsemodel.ResponseTransactionFilter
// @router /transactions/filter [get]
func (Transaction) TransactionFilter(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.Transaction()
	)

	result := s.GetFilter(ctx)
	return response.R200(c, result, "")

}

// GetSummary godoc
// @tags Transaction
// @summary Summary
// @id app-transaction-summary
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @success 200 {object} responsemodel.ResponseTransactionGetSummary
// @router /transactions/summary [get]
func (Transaction) GetSummary(c echo.Context) error {
	var (
		ctx         = echocontext.GetContext(c)
		s           = service.Transaction()
		sellerID, _ = mongodb.NewIDFromString(echocontext.GetCurrentUserID(c))
	)

	result := s.GetSummary(ctx, sellerID)
	return response.R200(c, result, "")
}
