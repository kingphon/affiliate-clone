package handler

import (
	"context"

	"git.selly.red/Selly-Server/affiliate/internal/config"

	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/service"
	"github.com/labstack/echo/v4"
)

// Sandbox ...
type Sandbox struct {
}

// GenerateTransaction godoc
// @tags Sandbox
// @summary Generate transaction
// @id admin-sandbox-generate-transaction
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.SandboxGenerateTransaction true "Query"
// @router /sandbox/transactions [get]
func (Sandbox) GenerateTransaction(c echo.Context) error {
	var (
		ctx           = context.Background()
		qParams       = echocontext.GetQuery(c).(requestmodel.SandboxGenerateTransaction)
		s             = service.MockData()
		campaignID, _ = mongodb.NewIDFromString(qParams.CampaignID)
		sellerID, _   = mongodb.NewIDFromString(qParams.SellerID)
	)

	go s.CreateMockData(ctx, campaignID, sellerID)
	return response.R200(c, nil, "")
}

// Crawlers godoc
// @tags Sandbox
// @summary Crawlers Sandbox
// @id sandbox-crawler
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload body requestmodel.SandboxCrawlerTransaction true "Payload"
// @router /sandbox/crawler [post]
func (Sandbox) Crawlers(c echo.Context) error {
	var (
		qParams = echocontext.GetPayload(c)
		s       = service.MockData()
	)
	if qParams == nil {
		return response.R400(c, nil, response.CommonBadRequest)
	}
	if config.IsEnvDevelop() {
		go s.SandboxCrawl(qParams.(requestmodel.SandboxCrawlerTransaction))
	}
	return response.R200(c, nil, "")
}
