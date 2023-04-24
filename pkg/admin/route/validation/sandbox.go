package routevalidation

import (
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	"github.com/labstack/echo/v4"
)

// Sandbox ...
type Sandbox struct{}

// GenerateTransaction ...
func (Sandbox) GenerateTransaction(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var query requestmodel.SandboxGenerateTransaction

		if err := c.Bind(&query); err != nil {
			return response.R400(c, nil, "")
		}

		if err := query.Validate(); err != nil {
			return err
		}

		echocontext.SetQuery(c, query)
		return next(c)
	}
}

// CrawlerTransaction ...
func (Sandbox) CrawlerTransaction(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var query requestmodel.SandboxCrawlerTransaction

		if err := c.Bind(&query); err != nil {
			return response.R400(c, nil, "")
		}

		if err := query.Validate(); err != nil {
			return err
		}

		echocontext.SetPayload(c, query)
		return next(c)
	}
}
