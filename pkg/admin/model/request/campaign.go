package requestmomdel

import (
	"math"

	"git.selly.red/Selly-Modules/mongodb"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/file"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
)

// CampaignCreate ...
type CampaignCreate struct {
	Name                 string                   `json:"name"`
	Order                int                      `json:"order"`
	Desc                 string                   `json:"desc"`
	Logo                 *file.FilePhotoRequest   `json:"logo"`
	Covers               []*file.FilePhotoRequest `json:"covers"`
	Commission           CampaignCommission       `json:"commission"`
	From                 string                   `json:"from"`
	To                   string                   `json:"to"`
	Platforms            []PlatformCreate         `json:"platforms"`
	EstimateCashback     CampaignEstimateCashback `json:"estimateCashback"`
	ShareDesc            string                   `json:"shareDesc"`
	AllowShowShareAction bool                     `json:"allowShowShareAction"`
}

// Validate ...
func (m CampaignCreate) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name,
			validation.Required.Error(errorcode.CampaignIsRequiredName)),
		validation.Field(&m.Commission),
		validation.Field(&m.EstimateCashback),
		validation.Field(&m.Platforms),
	)
}

// CampaignCommission ...
type CampaignCommission struct {
	Real          float64 `json:"real"`
	SellerPercent float64 `json:"sellerPercent"`
}

// Validate ...
func (m CampaignCommission) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(
			&m.Real,
			validation.Required.Error(errorcode.CampaignInvalidReal),
			validation.Min(float64(1)).Error(errorcode.CampaignInvalidReal),
		),
		validation.Field(
			&m.SellerPercent,
			validation.Required.Error(errorcode.CampaignInValidSellerPercent),
			validation.Min(float64(1)).Error(errorcode.CampaignInValidSellerPercent),
			validation.Max(float64(100)).Error(errorcode.CampaignInValidSellerPercent),
		),
	)
}

// CampaignEstimateCashback ...
type CampaignEstimateCashback struct {
	Day       int `json:"day"`
	NextMonth int `json:"nextMonth"`
}

func (m CampaignEstimateCashback) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Day, validation.Required.Error(errorcode.CampaignInvalidDay),
			validation.Min(1).Error(errorcode.CampaignInvalidDay),
			validation.Max(31).Error(errorcode.CampaignInvalidDay),
		),
		validation.Field(&m.NextMonth, validation.Required.Error(errorcode.CampaignInvalidNextMonth)),
	)
}

// ConvertToBSON ...
func (m CampaignCreate) ConvertToBSON() mgaffiliate.Campaign {
	// 1. Calculate
	sellerReceive := math.Round(m.Commission.Real * m.Commission.SellerPercent / 100)
	sellyReceive := m.Commission.Real - sellerReceive

	// 2. logo
	logo := m.Logo.ConvertToFilePhoto()
	logoID, _ := mongodb.NewIDFromString(logo.ID)

	// 3. Covers
	covers := make([]*mgaffiliate.FilePhoto, 0)
	for _, cover := range m.Covers {
		if cover == nil {
			continue
		}
		coverData := cover.ConvertToFilePhoto()
		coverID, _ := mongodb.NewIDFromString(cover.ID)
		covers = append(covers, &mgaffiliate.FilePhoto{
			ID:   coverID,
			Name: coverData.Name,
			Dimensions: &mgaffiliate.FileDimensions{
				Small: &mgaffiliate.FileSize{
					Width:  coverData.Dimensions.Small.Width,
					Height: coverData.Dimensions.Small.Height,
				},
				Medium: &mgaffiliate.FileSize{
					Width:  coverData.Dimensions.Medium.Width,
					Height: coverData.Dimensions.Medium.Height,
				},
			},
		})
	}

	return mgaffiliate.Campaign{
		ID:           mongodb.NewObjectID(),
		Name:         m.Name,
		SearchString: mongodb.NonAccentVietnamese(m.Name),
		Status:       constant.PlatformStatusInActive,
		Logo: &mgaffiliate.FilePhoto{
			ID:   logoID,
			Name: logo.Name,
			Dimensions: &mgaffiliate.FileDimensions{
				Small: &mgaffiliate.FileSize{
					Width:  logo.Dimensions.Small.Width,
					Height: logo.Dimensions.Small.Height,
				},
				Medium: &mgaffiliate.FileSize{
					Width:  logo.Dimensions.Medium.Width,
					Height: logo.Dimensions.Medium.Height,
				},
			},
		},
		Covers: covers,
		Order:  m.Order,
		Desc:   m.Desc,
		Commission: mgaffiliate.CampaignCommission{
			Real:          m.Commission.Real,
			SellerPercent: m.Commission.SellerPercent,
			Selly:         sellyReceive,
			Seller:        sellerReceive,
		},
		EstimateCashback: mgaffiliate.CampaignEstimateCashback{
			Day:       m.EstimateCashback.Day,
			NextMonth: m.EstimateCashback.NextMonth,
		},
		ShareDesc:            m.ShareDesc,
		From:                 ptime.TimeParseISODate(m.From),
		To:                   ptime.TimeParseISODate(m.To),
		CreatedAt:            ptime.Now(),
		UpdatedAt:            ptime.Now(),
		Platforms:            make([]string, 0),
		AllowShowShareAction: m.AllowShowShareAction,
	}
}

// CampaignUpdate ...
type CampaignUpdate struct {
	Name                 string                   `json:"name"`
	Order                int                      `json:"order"`
	Desc                 string                   `json:"desc"`
	Logo                 *file.FilePhotoRequest   `json:"logo"`
	Covers               []*file.FilePhotoRequest `json:"covers"`
	Commission           CampaignCommission       `json:"commission"`
	From                 string                   `json:"from"`
	To                   string                   `json:"to"`
	EstimateCashback     CampaignEstimateCashback `json:"estimateCashback"`
	ShareDesc            string                   `json:"shareDesc"`
	AllowShowShareAction bool                     `json:"allowShowShareAction"`
}

// Validate ...
func (m CampaignUpdate) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name,
			validation.Required.Error(errorcode.CampaignIsRequiredName)),
		validation.Field(&m.Commission),
	)
}

// ConvertToBSON ...
func (m CampaignUpdate) ConvertToBSON() mgaffiliate.Campaign {
	// 1. Calculate
	sellerReceive := math.Round(m.Commission.Real * m.Commission.SellerPercent / 100)
	sellyReceive := m.Commission.Real - sellerReceive

	// 2. logo
	logo := m.Logo.ConvertToFilePhoto()
	logoID, _ := mongodb.NewIDFromString(logo.ID)

	// 3. Covers
	covers := make([]*mgaffiliate.FilePhoto, 0)
	for _, cover := range m.Covers {
		if cover == nil {
			continue
		}

		coverData := cover.ConvertToFilePhoto()
		coverID, _ := mongodb.NewIDFromString(cover.ID)
		covers = append(covers, &mgaffiliate.FilePhoto{
			ID:   coverID,
			Name: coverData.Name,
			Dimensions: &mgaffiliate.FileDimensions{
				Small: &mgaffiliate.FileSize{
					Width:  coverData.Dimensions.Small.Width,
					Height: coverData.Dimensions.Small.Height,
				},
				Medium: &mgaffiliate.FileSize{
					Width:  coverData.Dimensions.Medium.Width,
					Height: coverData.Dimensions.Medium.Height,
				},
			},
		})
	}

	return mgaffiliate.Campaign{
		Name:         m.Name,
		SearchString: mongodb.NonAccentVietnamese(m.Name),
		Logo: &mgaffiliate.FilePhoto{
			ID:   logoID,
			Name: logo.Name,
			Dimensions: &mgaffiliate.FileDimensions{
				Small: &mgaffiliate.FileSize{
					Width:  logo.Dimensions.Small.Width,
					Height: logo.Dimensions.Small.Height,
				},
				Medium: &mgaffiliate.FileSize{
					Width:  logo.Dimensions.Medium.Width,
					Height: logo.Dimensions.Medium.Height,
				},
			},
		},
		Covers: covers,
		Order:  m.Order,
		Desc:   m.Desc,
		Commission: mgaffiliate.CampaignCommission{
			Real:          m.Commission.Real,
			SellerPercent: m.Commission.SellerPercent,
			Selly:         sellyReceive,
			Seller:        sellerReceive,
		},
		EstimateCashback: mgaffiliate.CampaignEstimateCashback{
			Day:       m.EstimateCashback.Day,
			NextMonth: m.EstimateCashback.NextMonth,
		},
		From:                 ptime.TimeParseISODate(m.From),
		To:                   ptime.TimeParseISODate(m.To),
		UpdatedAt:            ptime.Now(),
		ShareDesc:            m.ShareDesc,
		AllowShowShareAction: m.AllowShowShareAction,
	}
}

// CampaignAll ...
type CampaignAll struct {
	Page    int64  `query:"page"`
	Limit   int64  `query:"limit"`
	Keyword string `query:"keyword"`
	Status  string `query:"status"`
	FromAt  string `query:"fromAt"`
	ToAt    string `query:"toAt"`
	Sort    string `query:"sort"`
}

func (m CampaignAll) Validate() error {
	var status = []interface{}{
		constant.CampaignStatusActive,
		constant.CampaignStatusInActive,
	}

	return validation.ValidateStruct(
		&m,
		validation.Field(&m.Status,
			validation.In(status...).Error(errorcode.CampaignIsRequiredFrom)),
	)
}

// CampaignChangeStatus ...
type CampaignChangeStatus struct {
	Status string `json:"status"`
}

// Validate ...
func (m CampaignChangeStatus) Validate() error {
	var statuses = []interface{}{
		constant.CampaignStatusInActive,
		constant.CampaignStatusActive,
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.Status,
			validation.In(statuses...).Error(errorcode.PlatformStatusInvalid)),
	)
}
