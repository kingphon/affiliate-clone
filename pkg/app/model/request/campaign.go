package requestmodel

import (
	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CampaignAll ...
type CampaignAll struct {
	Keyword   string `json:"keyword" query:"keyword"`
	PageToken string `json:"pageToken" query:"pageToken"`
	Sort      string `json:"sort" query:"sort"`
}

// Validate ...
func (m CampaignAll) Validate() error {
	return validation.ValidateStruct(
		&m,
	)
}

// CampaignGroupAll ...
type CampaignGroupAll struct {
	Screen string `json:"screen" query:"screen"`
}

// Validate ...
func (m CampaignGroupAll) Validate() error {
	return validation.ValidateStruct(
		&m,
	)
}

// ClickUpdateBody ...
type ClickUpdateBody struct {
	PartnerClickID string `json:"partnerClickID"`
	FinalURL       string `json:"finalURL"`
}

// Validate ...
func (m ClickUpdateBody) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FinalURL,
			validation.Required.Error(errorcode.CampaignFinalURLIsRequired)),
	)
}
