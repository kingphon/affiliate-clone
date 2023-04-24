package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/affiliate/external/auth/permission"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/affiliate/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/admin/route/validation"
)

func reconciliation(e *echo.Group) {
	var (
		g = e.Group("/reconciliations", routeauth.RequiredLogin)
		h = handler.Reconciliation{}
		v = routevalidation.Reconciliation{}
	)

	// Permission
	edit := routeauth.CheckPermission(permission.Affiliate.Reconciliation.Edit)
	view := routeauth.CheckPermission(permission.Affiliate.Reconciliation.View)

	// Create
	g.POST("", h.Create, edit, v.Create)

	// All
	g.GET("", h.All, view, v.All)

	// Detail
	g.GET("/:id", h.Detail, view, v.Detail)

	// GetTransactionsByCondition
	g.GET("/:id/transactions", h.GetTransactionsByCondition, view, v.GetTransactionsByCondition, v.Detail)

	// GetStatistic ...
	g.GET("/:id/statistic", h.GetStatistic, view, v.GetStatistic, v.Detail)

	// ChangeStatus ...
	g.PATCH("/:id/status", h.ChangeStatus, v.ChangeStatus, v.Detail)

}
