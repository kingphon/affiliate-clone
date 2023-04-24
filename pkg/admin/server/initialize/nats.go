package initialize

import (
	"git.selly.red/Selly-Modules/natsio"
	"git.selly.red/Selly-Modules/natsio/subject"
	internalnats "git.selly.red/Selly-Server/affiliate/internal/nats"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/handler"
)

func nats() {
	// Init nats
	var natsServer = internalnats.ServerNats{}
	c := natsServer.Connect()

	// Subs ...
	initQueueSubscribeForAdmin(c)

}

// initQueueSubscribeForAdmin ...
func initQueueSubscribeForAdmin(c *natsio.JSONEncoder) {
	if c == nil {
		panic("Errors NewJSONEncodedConn")
	}

	h := handler.Nats{EncodedConn: c}
	subj := subject.Affiliate

	c.QueueSubscribe(subj.GetTransactions, subj.GetTransactions, h.GetTransactions)

}
