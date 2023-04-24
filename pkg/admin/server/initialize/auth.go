package initialize

import (
	"git.selly.red/Selly-Server/affiliate/internal/config"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/auth"
)

// authentication
func authentication() {
	env := config.GetENV()
	auth.InitAuthentication(env.Nats.APIKey, env.Nats)
}
