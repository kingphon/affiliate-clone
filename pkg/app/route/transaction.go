package route

import (
	"git.selly.red/Selly-Server/affiliate/pkg/app/handler"
	routeauth "git.selly.red/Selly-Server/affiliate/pkg/app/route/auth"
	routevalidation "git.selly.red/Selly-Server/affiliate/pkg/app/route/validation"
	"github.com/labstack/echo/v4"
)

// Transaction ...
type Transaction struct{}

// transaction ...
func transaction(e *echo.Group) {
	var (
		g = e.Group("/transactions", routeauth.RequiredLogin)
		h = handler.Transaction{}
		v = routevalidation.Transaction{}
	)

	// All
	g.GET("", h.All, v.All)

	// Detail
	g.GET("/:id", h.Detail, v.Detail)

	// Get history
	g.GET("/:id/histories", h.GetHistories, v.Detail)

	// Filter
	g.GET("/filter", h.TransactionFilter)

	// Summary
	g.GET("/summary", h.GetSummary)

	// Statistic
	g.GET("/statistic", h.GetStatistic)
}
