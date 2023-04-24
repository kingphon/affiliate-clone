package routeauth

import (
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"github.com/labstack/echo/v4"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/response"
)

// RequiredLogin ...
func RequiredLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// check invalid token
		user := externalauth.GetCurrentUserByToken(c.Get("user"))
		if user == nil || user.ID == "" {
			return response.R403(c, echo.Map{}, response.CommonForbidden)
		}

		c.Set(echocontext.KeyCurrentUserID, user.ID)
		return next(c)
	}
}
