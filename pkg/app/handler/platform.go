package handler

import (
	"fmt"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/external/utils/routemiddleware"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/request"
	"git.selly.red/Selly-Server/affiliate/pkg/app/service"
	"github.com/labstack/echo/v4"
)

// Platform ...
type Platform struct{}

// GenerateAffiliateLink godoc
// @tags Campaign platform
// @summary GenerateAffiliateLink
// @id app-campaign-platform-generate-aff-link
// @security ApiKeyAuth
// @accept json
// @produce json
// @param DeviceModel header string false "Model"
// @param User-Agent header string false "user agent"
// @param OS-NAME header string false "OS name"
// @param OS-VERSION header string false "OS version"
// @param BROWSER-VERSION header string false "browser version"
// @param BROWSER-NAME header string false "browser name"
// @param DeviceType header string false "device type"
// @param DeviceId header string false "device id"
// @param payload body requestmodel.GenerateAffiliateLinkBody true "Payload"
// @success 200 {object} responsemodel.ResponseAffiliateLink
// @router /campaign-platforms/generate-affiliate-link [patch]
func (Platform) GenerateAffiliateLink(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		s       = service.Platform()
		payload = echocontext.GetPayload(c).(requestmodel.GenerateAffiliateLinkBody)
		device  = mgaffiliate.Device{
			Model:          echocontext.GetHeaders(c).Get(routemiddleware.HeaderModel),
			UserAgent:      echocontext.GetHeaders(c).Get(routemiddleware.HeaderUserAgent),
			OSName:         echocontext.GetHeaders(c).Get(routemiddleware.HeaderOSName),
			OSVersion:      echocontext.GetHeaders(c).Get(routemiddleware.HeaderOSVersion),
			BrowserVersion: echocontext.GetHeaders(c).Get(routemiddleware.HeaderBrowserVersion),
			BrowserName:    echocontext.GetHeaders(c).Get(routemiddleware.HeaderBrowserName),
			DeviceType:     echocontext.GetHeaders(c).Get(routemiddleware.HeaderDeviceType),
			DeviceID:       echocontext.GetHeaders(c).Get(routemiddleware.HeaderDeviceID),
			Manufacturer:   echocontext.GetHeaders(c).Get(routemiddleware.HeaderDeviceManufacturer),
		}
	)
	fmt.Println("Headers : ", c.Request().Header)
	// GenerateAffiliateLink ...
	result, err := s.GenerateAffiliateLink(ctx, device, payload)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}
