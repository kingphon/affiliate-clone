package handler

import (
	"git.selly.red/Selly-Server/affiliate/external/constant"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/service"
)

// Click ...
type Click struct{}

// Statistic godoc
// @tags Click
// @summary Statistic
// @id admin-click-statistic
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.ClickStatistic true "Query"
// @success 200 {object} responsemodel.ResponseClickStatistic
// @router /clicks/statistic [get]
func (Click) Statistic(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		qParams = echocontext.GetQuery(c).(requestmodel.ClickStatistic)
		s       = service.Click()
		q       = mgquery.Affiliate{
			Source:      qParams.Source,
			FromAt:      ptime.TimeParseISODate(qParams.FromAt),
			ToAt:        ptime.TimeParseISODate(qParams.ToAt),
			CampaignIds: qParams.CampaignIds,
			Seller:      qParams.Seller,
		}
	)
	result := s.Statistic(ctx, q)
	return response.R200(c, result, "")
}

// All godoc
// @tags Click
// @summary All
// @id admin-click-all
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.ClickAll true "Query"
// @success 200 {object} responsemodel.ResponseClickAll
// @router /clicks [get]
func (Click) All(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		qParams = echocontext.GetQuery(c).(requestmodel.ClickAll)
		s       = service.Click()
		q       = mgquery.AppQuery{
			Page:          qParams.Page,
			Limit:         qParams.Limit,
			SortInterface: bson.M{"createdAt": -1},
			Affiliate: mgquery.Affiliate{
				Status:   qParams.Status,
				Keyword:  qParams.Keyword,
				Source:   qParams.Source,
				FromAt:   ptime.TimeParseISODate(qParams.FromAt),
				ToAt:     ptime.TimeParseISODate(qParams.ToAt),
				Campaign: qParams.Campaign,
				Seller:   qParams.Seller,
				From:     qParams.From,
			},
		}
	)

	result := s.All(ctx, q)
	return response.R200(c, result, "")
}

// MigrationSearchString godoc
// @tags Click
// @summary MigrationSearchString
// @id admin-click-migration-search-string
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param key query string true "Key"
// @success 200 {object} nil
// @router /clicks/migration-search-string [get]
func (h Click) MigrationSearchString(c echo.Context) error {
	var (
		s   = service.Click()
		key = c.QueryParam("key")
	)
	if key == constant.KeyMigration {
		go s.MigrationSearchString()
	}
	return response.R200(c, nil, "")
}
