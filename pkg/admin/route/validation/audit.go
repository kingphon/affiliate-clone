package routevalidation

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
)

// Audit
type Audit struct{}

// GetAll ...
func (Audit) GetAll(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			query requestmodel.AllQuery
		)

		if err := c.Bind(&query); err != nil {
			return response.R400(c, nil, "")
		}

		if err := query.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetQuery(c, query)
		return next(c)
	}
}
