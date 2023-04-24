package middleware

import (
	"encoding/json"
	"fmt"
	"git.selly.red/Selly-Modules/authentication"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/auth"
	"github.com/friendsofgo/errors"
)

// VerifyCodeWithAuthentication ...
func VerifyCodeWithAuthentication(body authentication.StaffVerifyCodeBody) error {
	// Check permission
	res, err := auth.VerifyCode(body)
	if err != nil {
		return errors.New(fmt.Sprintf("Xác thực thất bại! (%s)", err.Error()))
	}

	var result authentication.StaffVerifyCodeResponse
	_ = json.Unmarshal(res.Data, &result)
	if !result.Success {
		return errors.New(fmt.Sprintf("Xác thực thất bại!"))
	}

	return nil
}
