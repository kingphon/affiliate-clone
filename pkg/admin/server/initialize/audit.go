package initialize

import (
	"git.selly.red/Selly-Modules/audit"
	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/internal/config"
)

// InitAudit ...
func InitAudit() {
	var cfg = config.GetENV().MongoAudit

	// Init
	if err := audit.NewInstance(audit.Config{
		Targets: []string{
			constant.AuditTargetAffiliate,
		},
		MongoDB: cfg.GetConnectOptions(),
	}); err != nil {
		panic(err)
	}
}
