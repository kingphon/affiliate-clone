package requestmodel

import (
	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// GenerateAffiliateLinkBody ...
type GenerateAffiliateLinkBody struct {
	Code       string `json:"code"`
	PlatformID string `json:"platformId"`
	Checksum   string `json:"checksum"`
	From       string `json:"from"`
}

// Validate ...
func (m GenerateAffiliateLinkBody) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Code,
			validation.Required.Error(errorcode.CodeInvalid)),
		validation.Field(&m.PlatformID,
			validation.Required.Error(errorcode.PlatformIDInvalid)),
		validation.Field(&m.Checksum,
			validation.Required.Error(errorcode.ChecksumInvalid)),
	)
}
