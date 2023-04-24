package server

import (
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/affiliate/external/utils/file"
	"git.selly.red/Selly-Server/affiliate/internal/config"
	"git.selly.red/Selly-Server/affiliate/pkg/app/route"
	"git.selly.red/Selly-Server/affiliate/pkg/app/server/initialize"
	"github.com/labstack/echo/v4"
)

// Bootstrap ...
func Bootstrap(e *echo.Echo) {
	logger.Init("selly", "affiliate-app")

	// Init modules
	initialize.Init()

	// file
	cfg := config.GetENV()
	file.SetFileHost(cfg.FileHost)

	// Routes
	route.Init(e)
}
