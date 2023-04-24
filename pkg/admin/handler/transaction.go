package handler

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/service"
)

// Transaction ...
type Transaction struct {
}

// Statistic godoc
// @tags Transaction
// @summary All
// @id admin-transaction-statistic
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.TransactionStatistic true "Query"
// @success 200 {object} responsemodel.ResponseTransactionStatistic
// @router /transactions/statistic [get]
func (Transaction) Statistic(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		qParams = echocontext.GetQuery(c).(requestmodel.TransactionStatistic)
		s       = service.Transaction()
		q       = mgquery.Affiliate{
			CampaignIds: qParams.CampaignIds,
			Status:      qParams.Status,
			FromAt:      ptime.TimeParseISODate(qParams.FromAt),
			ToAt:        ptime.TimeParseISODate(qParams.ToAt),
			Source:      qParams.Source,
			Seller:      qParams.Seller,
		}
	)
	result := s.Statistic(ctx, q)
	return response.R200(c, result, "")
}

// All godoc
// @tags Transaction
// @summary All
// @id admin-transaction-all
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.TransactionAll true "Query"
// @success 200 {object} responsemodel.ResponseTransactionAll
// @router /transactions [get]
func (Transaction) All(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		qParams = echocontext.GetQuery(c).(requestmodel.TransactionAll)
		s       = service.Transaction()
		q       = mgquery.AppQuery{
			Page:          qParams.Page,
			Limit:         qParams.Limit,
			SortInterface: bson.M{"transactionTime": -1},
			Affiliate: mgquery.Affiliate{
				Campaign: qParams.Campaign,
				Status:   qParams.Status,
				Keyword:  qParams.Keyword,
				FromAt:   ptime.TimeParseISODate(qParams.FromAt),
				ToAt:     ptime.TimeParseISODate(qParams.ToAt),
				Source:   qParams.Source,
				Seller:   qParams.Seller,
			},
		}
	)

	result := s.All(ctx, q)
	return response.R200(c, result, "")
}

// Detail godoc
// @tags Transaction
// @summary Detail
// @id admin-transaction-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Transaction id"
// @success 200 {object} responsemodel.ResponseTransactionDetail
// @router /transactions/{id} [get]
func (Transaction) Detail(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.Transaction()
		id  = echocontext.GetParam(c, "id").(primitive.ObjectID)
	)

	result, err := s.Detail(ctx, id)
	if err != nil {
		response.R400(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}

// GetHistoryByTransaction godoc
// @tags Transaction
// @summary Get transaction histories
// @id admin-transaction-histories
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Transaction id"
// @success 200 {object} responsemodel.ResponseTransactionHistories
// @router /transactions/{id}/histories [get]
func (Transaction) GetHistoryByTransaction(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.Transaction()
		id  = echocontext.GetParam(c, "id").(primitive.ObjectID)
	)

	result := s.GetHistoriesByTransaction(ctx, id)
	return response.R200(c, result, "")
}

// Crawl godoc
// @tags Transaction
// @summary All
// @id admin-transaction-crawl
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string false "Device id"
// @router /transactions/crawl [get]
func (Transaction) Crawl(c echo.Context) error {
	var (
		s = service.TransactionCrawl()
	)

	go s.CrawlTransactionDataFromAdmin()
	return response.R200(c, nil, "")
}
