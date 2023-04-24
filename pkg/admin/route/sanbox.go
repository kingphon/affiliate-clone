package route

import (
	"git.selly.red/Selly-Server/affiliate/pkg/admin/handler"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/admin/route/validation"
	"github.com/labstack/echo/v4"
)

// sandbox ...
func sandbox(e *echo.Group) {
	var (
		g = e.Group("/sandbox")
		h = handler.Sandbox{}
		v = routevalidation.Sandbox{}
	)

	g.GET("/transactions", h.GenerateTransaction, v.GenerateTransaction)

	g.POST("/crawler", h.Crawlers, v.CrawlerTransaction)
}
