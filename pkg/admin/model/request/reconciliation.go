package requestmomdel

import (
	"git.selly.red/Selly-Modules/mongodb"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/constant"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
)

// ReconciliationCreate ...
type ReconciliationCreate struct {
	Name      string                        `json:"name"`
	Type      string                        `json:"type"`
	Condition ReconciliationConditionCreate `json:"condition"`
}

func (m ReconciliationCreate) Validate() error {
	// Validation condition when type == "filter"
	if funk.Contains(constant.ReconciliationType.Filter.Key, m.Type) {
		if err := m.Condition.Validate(); err != nil {
			return err
		}
	}

	var status = []interface{}{
		constant.ReconciliationType.Filter.Key,
		constant.ReconciliationType.Import.Key,
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required.Error(errorcode.ReconciliationIsRequiredName)),
		validation.Field(&m.Type, validation.Required.Error(errorcode.ReconciliationIsRequiredType),
			validation.In(status...).Error(errorcode.ReconciliationInvalidType),
		),
	)
}

// ReconciliationConditionCreate ...
type ReconciliationConditionCreate struct {
	Campaign string `json:"campaign"`
	Source   string `json:"source"`
	FromAt   string `json:"fromAt"`
	ToAt     string `json:"toAt"`
}

func (m ReconciliationConditionCreate) Validate() error {
	var sources = []interface{}{
		constant.ReconciliationSourceOther,
		constant.ReconciliationSourceAccessTrade,
	}
	return validation.ValidateStruct(&m,
		validation.Field(&m.Campaign, is.MongoID.Error(errorcode.ReconciliationConditionInvalidCampaignId),
			validation.Required.Error(errorcode.ReconciliationConditionIsRequiredCampaignId)),
		validation.Field(&m.Source, validation.Required.Error(errorcode.ReconciliationConditionIsRequiredSource),
			validation.In(sources...).Error(errorcode.ReconciliationConditionInvalidSource)),
		validation.Field(&m.FromAt, validation.Required.Error(errorcode.ReconciliationConditionIsRequiredFromAt),
			validation.Date(ptime.DateLayoutFull).Error(errorcode.ReconciliationConditionInvalidFromAt)),
		validation.Field(&m.ToAt, validation.Required.Error(errorcode.ReconciliationConditionIsRequiredToAt),
			validation.Date(ptime.DateLayoutFull).Error(errorcode.ReconciliationConditionInvalidToAt)),
	)
}

func (m ReconciliationCreate) ConvertToBSON(staff externalauth.User) mgaffiliate.Reconciliation {
	var (
		condition *mgaffiliate.ReconciliationCondition
	)

	// Check type payload
	if m.Type == constant.ReconciliationType.Filter.Key {
		condition = m.Condition.ConvertToBSON()
	}

	// Staff
	var createdBy primitive.ObjectID
	if staff.ID != "" {
		staff, _ := mongodb.NewIDFromString(staff.ID)
		createdBy = staff
	}

	return mgaffiliate.Reconciliation{
		ID:           mongodb.NewObjectID(),
		Name:         m.Name,
		Type:         m.Type,
		Condition:    condition,
		Status:       constant.ReconciliationStatus.New.Key,
		SearchString: mongodb.NonAccentVietnamese(m.Name),
		CreatedAt:    ptime.Now(),
		UpdatedAt:    ptime.Now(),
		ActionBy: &mgaffiliate.ReconciliationActionBy{
			CreatedBy: createdBy,
		},
	}
}

func (m ReconciliationConditionCreate) ConvertToBSON() *mgaffiliate.ReconciliationCondition {
	res := &mgaffiliate.ReconciliationCondition{}

	if m.Campaign != "" {
		campaignId, _ := mongodb.NewIDFromString(m.Campaign)
		res.CampaignId = campaignId
	}

	return &mgaffiliate.ReconciliationCondition{
		CampaignId: res.CampaignId,
		Source:     m.Source,
		FromAt:     ptime.TimeParseISODate(m.FromAt),
		ToAt:       ptime.TimeParseISODate(m.ToAt),
	}
}

// ReconciliationAll ...
type ReconciliationAll struct {
	Page     int64  `query:"page"`
	Limit    int64  `query:"limit"`
	Keyword  string `query:"keyword"`
	FromAt   string `query:"fromAt"`
	ToAt     string `query:"toAt"`
	Campaign string `query:"campaign"`
	Status   string `query:"status"`
	Source   string `query:"source"`
}

func (m ReconciliationAll) Validate() error {
	var statuses = []interface{}{
		constant.ReconciliationStatus.New.Key,
		constant.ReconciliationStatus.Rejected.Key,
		constant.ReconciliationStatus.Approved.Key,
		constant.ReconciliationStatus.Running.Key,
		constant.ReconciliationStatus.Completed.Key,
		constant.ReconciliationStatus.Deleted.Key,
		constant.ReconciliationStatus.Fail.Key,
	}

	var sources = []interface{}{
		constant.ReconciliationSourceAll,
		constant.ReconciliationSourceAccessTrade,
		constant.ReconciliationSourceOther,
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.Status, validation.In(statuses...).Error(errorcode.ReconciliationInvalidStatus)),
		validation.Field(&m.Campaign, is.MongoID),
		validation.Field(&m.Source, validation.In(sources...).Error(errorcode.ReconciliationConditionInvalidSource)),
	)
}

// ReconciliationPayloadChangeStatus ...
type ReconciliationPayloadChangeStatus struct {
	Status     string `json:"status"`
	CodeAuthGG string `json:"codeAuthGG"`
}

func (m ReconciliationPayloadChangeStatus) Validate() error {
	var statuses = []interface{}{
		constant.ReconciliationStatus.Approved.Key,
		constant.ReconciliationStatus.Rejected.Key,
		constant.ReconciliationStatus.Deleted.Key,
		constant.ReconciliationStatus.Completed.Key,
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.Status, validation.Required.Error(errorcode.ReconciliationIsRequiredStatus), validation.In(statuses...).Error(errorcode.ReconciliationInvalidStatus)),
		validation.Field(&m.CodeAuthGG, validation.Required.Error(errorcode.ReconciliationIsRequiredCodeAuthGoogle), validation.Match(regexp.MustCompile("[0-9]{5}"))),
	)
}

// ReconciliationTransactionAll ...
type ReconciliationTransactionAll struct {
	Page    int64  `query:"page"`
	Limit   int64  `query:"limit"`
	Keyword string `query:"keyword"`
	Status  string `query:"status"`
	Seller  string `query:"seller"`
}

func (m ReconciliationTransactionAll) Validate() error {
	statuses := []interface{}{
		constant.TransactionStatus.New.Key,
		constant.TransactionStatus.Approved.Key,
		constant.TransactionStatus.Rejected.Key,
		constant.TransactionStatus.All.Key,
		constant.TransactionStatus.Cashback.Key,
		constant.TransactionStatus.Pending.Key,
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.Status, validation.In(statuses...).Error(errorcode.TransactionInvalidStatus)),
		validation.Field(&m.Seller, is.MongoID.Error(errorcode.TransactionInvalidID)),
	)
}

// ReconciliationPayloadStatistic ...
type ReconciliationPayloadStatistic struct {
	Status     string `query:"status"`
	Seller     string `query:"seller"`
	SearchCode string `query:"searchCode"`
}

func (m ReconciliationPayloadStatistic) Validate() error {
	statuses := []interface{}{
		constant.TransactionStatus.Approved.Key,
		constant.TransactionStatus.All.Key,
		constant.TransactionStatus.Cashback.Key,
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.Status, validation.In(statuses...).Error(errorcode.TransactionInvalidStatus)),
		validation.Field(&m.Seller, is.MongoID.Error(errorcode.SellerInvalidID)),
	)
}
