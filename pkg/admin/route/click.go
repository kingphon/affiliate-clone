package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/affiliate/external/auth/permission"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/affiliate/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/admin/route/validation"
)

func click(e *echo.Group) {
	var (
		g = e.Group("/clicks", routeauth.RequiredLogin)
		h = handler.Click{}
		v = routevalidation.Click{}
	)

	// Permission
	view := routeauth.CheckPermission(permission.Affiliate.Click.View)

	// All
	g.GET("", h.All, view, v.All)

	g.GET("/statistic", h.Statistic, view, v.Statistic)

	g.GET("/migration-search-string", h.MigrationSearchString)
}
