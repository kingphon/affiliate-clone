package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/affiliate/external/auth/permission"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/affiliate/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/admin/route/validation"
)

// transaction ...
func transaction(e *echo.Group) {
	var (
		g = e.Group("/transactions", routeauth.RequiredLogin)
		h = handler.Transaction{}
		v = routevalidation.Transaction{}
	)

	// Permission
	view := routeauth.CheckPermission(permission.Affiliate.Transaction.View)
	edit := routeauth.CheckPermission(permission.Affiliate.Transaction.Edit)

	// All
	g.GET("", h.All, view, v.All)

	// Detail
	g.GET("/:id", h.Detail, view, v.Detail)

	// Histories
	g.GET("/:id/histories", h.GetHistoryByTransaction, view, v.Detail)

	// Statistic
	g.GET("/statistic", h.Statistic, view, v.Statistic)

	// Crawl
	g.GET("/crawl", h.Crawl, edit)
}
