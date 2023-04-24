package route

import (
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/affiliate/external/auth/permission"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/handler"
	routeauth "git.selly.red/Selly-Server/affiliate/pkg/admin/route/auth"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/admin/route/validation"
)

// platform ...
func platform(e *echo.Group) {
	var (
		g = e.Group("/campaign-platforms", routeauth.RequiredLogin)
		h = handler.Platform{}
		v = routevalidation.Platform{}
	)

	// Permission
	edit := routeauth.CheckPermission(permission.Affiliate.Campaign.Edit)

	// Update
	g.PUT("/:id", h.Update, edit, v.Update, v.Detail)

	// Change status
	g.PATCH("/:id/status", h.ChangeStatus, edit, v.ChangeStatus, v.Detail)
}
