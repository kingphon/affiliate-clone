package route

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/external/utils/routemiddleware"
	"git.selly.red/Selly-Server/affiliate/internal/config"
	routeauth "git.selly.red/Selly-Server/affiliate/pkg/app/route/auth"
)

// Init ...
func Init(e *echo.Echo) {
	var (
		envVars = config.GetENV()
	)

	// Middlewares ...
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(envVars.SecretKey),
		Skipper: func(c echo.Context) bool {
			token := echocontext.GetToken(c)
			return token == ""
		},
	}))

	e.Use(routeauth.Jwt())

	e.Use(routemiddleware.CORSConfig())
	e.Use(routemiddleware.Locale)

	r := e.Group("/app/affiliate")

	// Components
	common(r)
	campaign(r)
	platform(r)
	transaction(r)
	shareURL(r)
}
