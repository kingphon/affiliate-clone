package routeauth

import (
	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/internal/config"
)

var envVars = config.GetENV()

// Jwt ...
func Jwt() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(envVars.SecretKey),
		Skipper: func(c echo.Context) bool {
			token := echocontext.GetToken(c)
			user := externalauth.GetCurrentUserByToken(c.Get("user"))
			var userID = ""
			if user != nil && user.ID != "" {
				userID = user.ID
			}

			c.Set(echocontext.KeyCurrentUserID, userID)
			return token == ""
		},
	})
}
