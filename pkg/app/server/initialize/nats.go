package initialize

import (
	internalnats "git.selly.red/Selly-Server/affiliate/internal/nats"
)

func nats() {
	var natsServer = internalnats.ServerNats{}
	natsServer.Connect()
}
