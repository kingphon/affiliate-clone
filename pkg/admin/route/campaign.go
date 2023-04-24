package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/affiliate/external/auth/permission"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/affiliate/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/admin/route/validation"
)

// campaign ...
func campaign(e *echo.Group) {
	var (
		g         = e.Group("/campaigns", routeauth.RequiredLogin)
		h         = handler.Campaign{}
		v         = routevalidation.Campaign{}
		vPlatform = routevalidation.Platform{}
	)

	// Permission
	edit := routeauth.CheckPermission(permission.Affiliate.Campaign.Edit)
	view := routeauth.CheckPermission(permission.Affiliate.Campaign.View)

	// Create
	g.POST("", h.Create, edit, v.Create)

	// Update
	g.PUT("/:id", h.Update, edit, v.Update, v.Detail)

	// All
	g.GET("", h.All, view, v.All)

	// Change status
	g.PATCH("/:id/status", h.ChangeStatus, edit, v.ChangeStatus, v.Detail)

	// Detail
	g.GET("/:id", h.Detail, view, v.Detail)

	// GetPlatform
	g.GET("/:id/platforms", h.GetPlatformByCampaign, view, v.Detail)

	// Create platform
	g.POST("/:id/platforms", h.CreatePlatform, edit, vPlatform.Create, v.Detail)
}
