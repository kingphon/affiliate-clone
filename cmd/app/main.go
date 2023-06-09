package main

import (
	"os"

	"git.selly.red/Selly-Server/affiliate/docs/app"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.elastic.co/apm/module/apmechov4"

	"git.selly.red/Selly-Server/affiliate/internal/config"
	"git.selly.red/Selly-Server/affiliate/pkg/app/server"
)

// @title Selly Affiliate - APP API
// @version 1.0
// @description All APIs for Affiliate app.
// @description
// @description ******************************
// @description - Add description
// @description ******************************
// @description
// @termsOfService https://selly.vn
// @contact.name Dev team
// @contact.url https://selly.vn
// @contact.email dev@selly.vn
// @basePath /app/affiliate

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Echo instance
	e := echo.New()

	e.Use(apmechov4.Middleware())

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${remote_ip} | ${method} ${uri} - ${status} - ${latency_human}\n",
	}))
	e.Use(middleware.Gzip())
	if os.Getenv("ENV") == "release" {
		e.Use(middleware.Recover())
	}

	// Bootstrap things
	server.Bootstrap(e)

	// Swagger
	if config.IsEnvDevelop() {
		domain := os.Getenv("DOMAIN_AFFILIATE_APP")
		app.SwaggerInfo.Host = domain
		e.GET(app.SwaggerInfo.BasePath+"/swagger/*", echoSwagger.WrapHandler)
	}

	// Start server
	e.Logger.Fatal(e.Start(config.GetENV().App.Port))
}
