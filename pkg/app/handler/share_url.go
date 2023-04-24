package handler

import (
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/pkg/app/service"
	"github.com/labstack/echo/v4"
)

// ShareURL ...
type ShareURL struct{}

// GetInfoShareURL godoc
// @tags ShareURL
// @summary Get info share URL
// @id app-share-url
// @security ApiKeyAuth
// @accept json
// @produce json
// @Param  code path string true "Seller share url code"
// @success 200 {object} responsemodel.ResponseShareURLInfo
// @router /campaign-share-url/{code} [get]
func (ShareURL) GetInfoShareURL(c echo.Context) error {
	var (
		ctx  = echocontext.GetContext(c)
		s    = service.ShareURL()
		code = echocontext.GetParam(c, "code").(string)
	)

	result, err := s.GetInfoShareURL(ctx, code)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}
