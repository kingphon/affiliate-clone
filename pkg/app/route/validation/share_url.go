package routevalidation

import (
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
	"github.com/labstack/echo/v4"
)

// ShareURL ...
type ShareURL struct{}

// Code ...
func (ShareURL) Code(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var code = c.Param("code")
		if code == "" {
			return response.R404(c, nil, errorcode.CampaignNotFound)
		}
		echocontext.SetParam(c, "code", code)
		return next(c)
	}
}
