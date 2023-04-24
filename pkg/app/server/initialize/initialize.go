package initialize

import (
	"git.selly.red/Selly-Server/affiliate/internal/config"
	"git.selly.red/Selly-Server/affiliate/internal/zk"
	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
)

// Init ...
func Init() {
	config.Init()
	errorcode.Init()
	zk.Connect()
	mongoDB()
	nats()
}
