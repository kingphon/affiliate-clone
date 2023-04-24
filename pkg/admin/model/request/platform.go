package requestmomdel

import (
	"regexp"
	"strings"

	"git.selly.red/Selly-Modules/mongodb"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/utils/prandom"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
)

// PlatformCreate ...
type PlatformCreate struct {
	Partner      PlatformPartner `json:"partner"`
	URL          string          `json:"url"`
	PlatformType string          `json:"platform"`
}

// Validate ...
func (m PlatformCreate) Validate() error {
	platforms := []interface{}{
		constant.PlatformAll,
		constant.PlatformIOS,
		constant.PlatformAndroid,
		constant.PlatformWeb,
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.URL, validation.Required.Error(errorcode.PlatformIsRequiredUrl)),
		validation.Field(&m.PlatformType,
			validation.Required.Error(errorcode.PlatformIsRequiredPlatform),
			validation.In(platforms...).Error(errorcode.PlatformInvalidPlatform)),
		validation.Field(&m.Partner),
	)
}

// PlatformPartner ...
type PlatformPartner struct {
	Source     string `json:"source"`
	CampaignID string `json:"campaignID"`
}

// Validate ...
func (m PlatformPartner) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Source, validation.Required.Error(errorcode.PlatformIsRequiredSource)),
		validation.Field(&m.CampaignID,
			validation.When(m.Source != constant.PlatformPartnerSourceOther,
				validation.Required.Error(errorcode.PlatformIsInvalidCampaignID),
				validation.Match(regexp.MustCompile("^([a-z0-9A-Z])*$")).Error(errorcode.PlatformIsInvalidCampaignID))),
	)
}

// ConvertToBSON ...
func (m PlatformCreate) ConvertToBSON(campaignID primitive.ObjectID) mgaffiliate.Platform {

	// remove space characters campaignID
	campaignId := strings.ReplaceAll(m.Partner.CampaignID, " ", "")

	res := mgaffiliate.Platform{
		ID:         mongodb.NewObjectID(),
		Code:       GenerateCodePlatform(),
		CampaignID: campaignID,
		Status:     constant.CampaignStatusInActive,
		Partner: mgaffiliate.PlatformPartner{
			Source: m.Partner.Source,
		},
		URL:          m.URL,
		PlatformType: m.PlatformType,
		CreatedAt:    ptime.Now(),
		UpdatedAt:    ptime.Now(),
	}
	if res.Partner.Source != constant.PlatformPartnerSourceOther {
		res.Partner.CampaignID = campaignId
	}

	return res
}

// GenerateCodePlatform ...
func GenerateCodePlatform() string {
	return prandom.RandomStringWithLength(4)
}

// PlatformChangeStatus ...
type PlatformChangeStatus struct {
	Status string `json:"status"`
}

// Validate ...
func (m PlatformChangeStatus) Validate() error {
	var statuses = []interface{}{
		constant.PlatformStatusInActive,
		constant.PlatformStatusActive,
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.Status,
			validation.In(statuses...).Error(errorcode.PlatformStatusInvalid)),
	)
}
