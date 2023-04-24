package auth

import (
	"github.com/friendsofgo/errors"
	"time"

	"git.selly.red/Selly-Modules/authentication"
	"git.selly.red/Selly-Modules/natsio"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/auth/permission"
	"git.selly.red/Selly-Server/affiliate/internal/config"
)

var client *authentication.Client

// InitAuthentication ...
func InitAuthentication(apiKey string, nats config.NatsConfig) {
	setSecretKey, err := externalauth.Init(apiKey, permission.AffiliateSource, natsio.Config{
		URL:            nats.URL,
		User:           nats.Username,
		Password:       nats.Password,
		RequestTimeout: 3 * time.Minute,
	})

	if err != nil {
		panic(err)
	}

	envVars := config.GetENV()
	envVars.SecretKey = setSecretKey
	client = externalauth.GetClient()
}

// VerifyCode ...
func VerifyCode(payload authentication.StaffVerifyCodeBody) (*authentication.Response, error) {
	if client == nil {
		return nil, errors.New("failed init client authentication")
	}

	return client.Request.VerifyCode(payload)
}
