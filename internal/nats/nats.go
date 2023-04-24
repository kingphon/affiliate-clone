package nats

import (
	"fmt"
	"time"

	"git.selly.red/Selly-Modules/natsio"
	"git.selly.red/Selly-Server/affiliate/internal/config"
)

// ServerNats ...
type ServerNats struct {
	c *natsio.JSONEncoder
}

// Connect ...
func (ns ServerNats) Connect() *natsio.JSONEncoder {
	cfg := config.GetENV().Nats
	err := natsio.Connect(natsio.Config{
		URL:            cfg.URL,
		User:           cfg.Username,
		Password:       cfg.Password,
		RequestTimeout: 2 * time.Minute,
	})
	if err != nil {
		panic(err)
	}

	s := natsio.GetServer()
	c, err := s.NewJSONEncodedConn()
	if err != nil {
		panic(err)
	}

	fmt.Println(c)
	return c
}
