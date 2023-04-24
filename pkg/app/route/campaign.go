package route

import (
	"git.selly.red/Selly-Server/affiliate/pkg/app/handler"
	routeauth "git.selly.red/Selly-Server/affiliate/pkg/app/route/auth"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/app/route/validation"
	"github.com/labstack/echo/v4"
)

// Campaign ...
type Campaign struct{}

// campaign ...
func campaign(e *echo.Group) {
	var (
		g = e.Group("/campaigns")

		h = handler.Campaign{}
		v = routevalidation.Campaign{}
	)

	// Groups
	g.GET("/groups", h.GroupAll, v.GroupAll)

	// Detail
	g.GET("/groups/:id", h.GroupDetail, v.GroupDetail)

	// Get list campaign by group id
	g.GET("/groups/:id/items", h.GroupItems, v.All)

	// List
	g.GET("", h.All, v.All)

	// Filter
	g.GET("/filter", h.Filter)

	// Detail
	g.GET("/:id", h.Detail, v.Detail)

	// generate share url
	g.PATCH("/:id/share-url", h.GenerateShareURL, v.Detail, routeauth.RequiredLogin)

	// update click id
	g.PUT("/click/:clickID", h.UpdateClick, v.ClickDetail, v.UpdateClickBody)
}
