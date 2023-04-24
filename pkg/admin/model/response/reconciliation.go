package responsemodel

import (
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
)

// ResponseReconciliationAll ...
type ResponseReconciliationAll struct {
	List  []ResponseReconciliationBrief `json:"list"`
	Total int64                         `json:"total"`
	Limit int64                         `json:"limit"`
}

// ResponseReconciliationBrief ...
type ResponseReconciliationBrief struct {
	ID        string                           `json:"_id"`
	Name      string                           `json:"name"`
	Type      string                           `json:"type"`
	Status    string                           `json:"status"`
	Condition *ResponseReconciliationCondition `json:"condition,omitempty"`
	CreatedAt *ptime.TimeResponse              `json:"createdAt"`
	ActionBy  *ResponseReconciliationActionBy  `json:"actionBy"`
}

// ResponseReconciliationCondition ...
type ResponseReconciliationCondition struct {
	Campaign ResponseCampaignShortInfo `json:"campaign"`
	Source   string                    `json:"source"`
	FromAt   *ptime.TimeResponse       `json:"fromAt"`
	ToAt     *ptime.TimeResponse       `json:"toAt"`
}

// ResponseReconciliationActionBy ...
type ResponseReconciliationActionBy struct {
	CreateBy  ResponseStaff  `json:"createBy"`
	ApproveBy *ResponseStaff `json:"approveBy,omitempty"`
}

// ResponseReconciliationDetail ...
type ResponseReconciliationDetail struct {
	ID           string                              `json:"_id"`
	Name         string                              `json:"name"`
	Type         string                              `json:"type"`
	Status       string                              `json:"status"`
	Condition    *ResponseReconciliationCondition    `json:"condition,omitempty"`
	CreatedAt    *ptime.TimeResponse                 `json:"createdAt"`
	ActionBy     *ResponseReconciliationActionBy     `json:"actionBy"`
	TrackingTime *ResponseReconciliationTrackingTime `json:"trackingTime"`
}

type ResponseReconciliationTrackingTime struct {
	ChangeStatusApprovedAt *ptime.TimeResponse `json:"changeStatusApprovedAt,omitempty"`
}

// ResponseReconciliationStatistic ...
type ResponseReconciliationStatistic struct {
	TotalTransaction    int64   `json:"totalTransaction"`
	TotalCommissionReal float64 `json:"totalCommissionReal"`
	SellerCommission    float64 `json:"sellerCommission"`
	SellyCommission     float64 `json:"sellyCommission"`
}
