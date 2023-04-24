package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

const (
	ReconciliationErrorWhenCreated                = "reconciliation_error_when_created"
	ReconciliationInvalidStatus                   = "reconciliation_invalid_status"
	ReconciliationInvalidID                       = "reconciliation_invalid_id"
	ReconciliationNotFound                        = "reconciliation_not_found"
	ReconciliationErrorWhenCensored               = "reconciliation_error_when_censored"
	ReconciliationErrorWhenDeleted                = "reconciliation_error_when_deleted"
	ReconciliationIsRequiredName                  = "reconciliation_is_required_name"
	ReconciliationIsRequiredType                  = "reconciliation_is_required_type"
	ReconciliationConditionInvalidCampaignId      = "reconciliation_condition_invalid_campaign_id"
	ReconciliationConditionIsRequiredCampaignId   = "reconciliation_condition_is_required_campaign_id"
	ReconciliationConditionIsRequiredSource       = "reconciliation_condition_is_required_source"
	ReconciliationConditionIsRequiredFromAt       = "reconciliation_condition_is_required_from_at"
	ReconciliationConditionIsRequiredToAt         = "reconciliation_condition_is_required_to_at"
	ReconciliationConditionInvalidSource          = "reconciliation_condition_is_invalid_source"
	StaffCreateReconciliationUncensored           = "staff_create_reconciliation_uncensored"
	OnlyStaffCreatorCanDeleteReconciliation       = "only_staff_creator_can_delete_reconciliation"
	ReconciliationConditionInvalidFromAt          = "reconciliation_condition_invalid_from_at"
	ReconciliationConditionInvalidToAt            = "reconciliation_condition_invalid_to_at"
	ReconciliationInvalidNewStatusInDatabase      = "reconciliation_invalid_new_status_in_database"
	ReconciliationCanNotUpdateStatus              = "reconciliation_can_not_update_status"
	ReconciliationInvalidType                     = "reconciliation_invalid_type"
	ReconciliationStatusIsNewAndApproved          = "reconciliation_status_is_type_and_approved"
	ReconciliationIsRequiredStatus                = "reconciliation_is_required_status"
	ReconciliationIsRequiredCodeAuthGoogle        = "reconciliation_is_required_auth_google"
	ReconciliationInvalidApprovedStatusInDatabase = "reconciliation_invalid_approved_status_in_database"
	ReconciliationErrorUnauthorized               = "reconciliation_error_unauthorized"
)

var reconciliation = []response.Code{
	{
		Key:     ReconciliationErrorWhenCreated,
		Message: "lỗi khi tạo mới yêu cầu đối soát",
		Code:    700,
	},
	{
		Key:     ReconciliationInvalidStatus,
		Message: "trạng thái yêu cầu đối soát không hợp lệ",
		Code:    701,
	},
	{
		Key:     ReconciliationInvalidID,
		Message: "id yêu cầu đối soát không hợp lệ",
		Code:    702,
	},
	{
		Key:     ReconciliationNotFound,
		Message: "đối soát không tồn tại",
		Code:    703,
	},
	{
		Key:     ReconciliationErrorWhenCensored,
		Message: "lỗi khi kiểm duyệt yêu cầu đối soát",
		Code:    704,
	},
	{
		Key:     ReconciliationErrorWhenDeleted,
		Message: "lỗi khi xóa yêu cầu đối soát",
		Code:    705,
	},
	{
		Key:     ReconciliationIsRequiredName,
		Message: "tên yêu cầu đối soát không được trống",
		Code:    706,
	},
	{
		Key:     ReconciliationIsRequiredType,
		Message: "loại yêu cầu đối soát không được trống",
		Code:    707,
	},

	//
	{
		Key:     ReconciliationConditionInvalidCampaignId,
		Message: "campaignID condition yêu cầu đối soát không không hợp lệ",
		Code:    708,
	},
	{
		Key:     ReconciliationConditionIsRequiredCampaignId,
		Message: "campaignID condition yêu cầu đối soát không được trống",
		Code:    709,
	},
	{
		Key:     ReconciliationConditionIsRequiredSource,
		Message: "source condition yêu cầu đối soát không được trống",
		Code:    710,
	},
	{
		Key:     ReconciliationConditionIsRequiredFromAt,
		Message: "fromAt condition yêu cầu đối soát không được trống",
		Code:    711,
	},
	{
		Key:     ReconciliationConditionIsRequiredToAt,
		Message: "toAt condition yêu cầu đối soát không được trống",
		Code:    712,
	},
	{
		Key:     ReconciliationConditionInvalidSource,
		Message: "source condition yêu cầu đối soát không hợp lệ",
		Code:    713,
	},
	{
		Key:     StaffCreateReconciliationUncensored,
		Message: "Người tạo yêu cầu đối soát không được duyệt",
		Code:    714,
	},
	{
		Key:     OnlyStaffCreatorCanDeleteReconciliation,
		Message: "Chỉ người tạo mới có thể xóa yêu cầu đối soát",
		Code:    715,
	},
	{
		Key:     ReconciliationConditionInvalidFromAt,
		Message: "fromAt condition yêu cầu đối soát không hợp lệ",
		Code:    716,
	},
	{
		Key:     ReconciliationConditionInvalidToAt,
		Message: "toAt condition yêu cầu đối soát không hợp lệ",
		Code:    717,
	},
	{
		Key:     ReconciliationInvalidNewStatusInDatabase,
		Message: "trạng thái yêu cầu đối soát trong database phải là New mới có thể update sang approved",
		Code:    718,
	},
	{
		Key:     ReconciliationCanNotUpdateStatus,
		Message: "không thể update trạng thái yêu cầu đối soát",
		Code:    719,
	},
	{
		Key:     ReconciliationInvalidType,
		Message: "loại yêu cầu đối soát không hợp lệ",
		Code:    720,
	},
	{
		Key:     ReconciliationStatusIsNewAndApproved,
		Message: "trạng thái yêu cầu đối soát phải là new hoặc approve mới có thể update",
		Code:    720,
	},
	{
		Key:     ReconciliationIsRequiredStatus,
		Message: "trạng thái yêu cầu đối soát phải không được trống",
		Code:    721,
	},
	{
		Key:     ReconciliationIsRequiredCodeAuthGoogle,
		Message: "code authentication google không được trống",
		Code:    722,
	},
	{
		Key:     ReconciliationInvalidApprovedStatusInDatabase,
		Message: "trạng thái yêu cầu đối soát phải là new hoặc approve mới có thể update sang completed",
		Code:    723,
	},
	{
		Key:     ReconciliationErrorUnauthorized,
		Message: "bạn không có quyền thực hiện hành động này",
		Code:    724,
	},
}
