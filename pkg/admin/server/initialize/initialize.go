package initialize

import (
	"git.selly.red/Selly-Server/affiliate/internal/config"
	"git.selly.red/Selly-Server/affiliate/internal/zk"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
)

// Init ...
func Init() {
	// Config ...
	config.Init()

	// Error code locale ...
	errorcode.Init()

	// Zookeeper connect ...
	zk.Connect()

	// Mongo db connect ...
	mongoDB()

	// Authentication ...
	authentication()

	// Audit ...
	InitAudit()

	// Nats connect ...
	nats()

	// Schedule ...
	initSchedule()
}
