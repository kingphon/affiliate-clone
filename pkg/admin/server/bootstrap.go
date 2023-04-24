package server

import (
	"git.selly.red/Selly-Modules/logger"
	"github.com/labstack/echo/v4"

	"git.selly.red/Selly-Server/affiliate/external/utils/file"
	"git.selly.red/Selly-Server/affiliate/internal/config"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/route"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/server/initialize"
)

// Bootstrap ...
func Bootstrap(e *echo.Echo) {
	logger.Init("selly", "affiliate-admin")

	// Init modules
	initialize.Init()

	// file
	cfg := config.GetENV()
	file.SetFileHost(cfg.FileHost)

	// Routes
	route.Init(e)

}
