package middleware

import (
	"git.selly.red/Selly-Modules/authentication"
	"git.selly.red/Selly-Server/affiliate/internal/config"
)

// VerifyOtpCode ...
func VerifyOtpCode(body authentication.StaffVerifyCodeBody) error {
	if config.IsEnableAuthenticationService() {
		return VerifyCodeWithAuthentication(body)
	}
	return nil
}
