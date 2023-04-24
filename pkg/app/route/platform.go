package route

import (
	"git.selly.red/Selly-Server/affiliate/pkg/app/handler"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/app/route/validation"
	"github.com/labstack/echo/v4"
)

// Platform ...
type Platform struct{}

// platform ...
func platform(e *echo.Group) {
	var (
		g = e.Group("/campaign-platforms")
		h = handler.Platform{}
		v = routevalidation.Platform{}
	)

	// generate affiliate link
	g.PATCH("/generate-affiliate-link", h.GenerateAffiliateLink, v.GenerateAffiliateLink)
}
