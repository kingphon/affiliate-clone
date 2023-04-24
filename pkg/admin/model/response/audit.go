package responsemodel

import (
	"git.selly.red/Selly-Modules/audit"
)

// ResponseAuditAll ...
type ResponseAuditAll struct {
	List  []audit.Audit `json:"list"`
	Total int64         `json:"total"`
}
