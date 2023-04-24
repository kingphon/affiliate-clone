package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/affiliate/external/utils/routemiddleware"
	routeauth "git.selly.red/Selly-Server/affiliate/pkg/admin/route/auth"
)

// Init ...
func Init(e *echo.Echo) {

	// Middlewares ...
	e.Use(routeauth.Jwt())

	e.Use(routemiddleware.CORSConfig())
	e.Use(routemiddleware.Locale)

	r := e.Group("/admin/affiliate")

	// Components
	common(r)
	campaign(r)
	platform(r)
	transaction(r)
	sandbox(r)
	click(r)
	reconciliation(r)
	audit(r)
}
