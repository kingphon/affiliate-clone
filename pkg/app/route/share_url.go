package route

import (
	"git.selly.red/Selly-Server/affiliate/pkg/app/handler"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/app/route/validation"
	"github.com/labstack/echo/v4"
)

// ShareURL ...
type ShareURL struct{}

// shareURL ...
func shareURL(e *echo.Group) {
	var (
		g = e.Group("/campaign-share-url")
		h = handler.ShareURL{}
		v = routevalidation.ShareURL{}
	)

	// Get info share URL
	g.GET("/:code", h.GetInfoShareURL, v.Code)
}
