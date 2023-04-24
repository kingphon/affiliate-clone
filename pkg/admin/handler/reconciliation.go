package handler

import (
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/auth/permission"
	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/service"
)

// Reconciliation ...
type Reconciliation struct{}

// Create godoc
// @tags Reconciliation
// @summary Create
// @id app-reconciliation-create
// @security ApiKeyAuth
// @accept json
// @produce json
// @param payload body requestmodel.ReconciliationCreate true "Payload"
// @success 200 {object} responsemodel.ResponseCreate
// @router /reconciliations [post]
func (Reconciliation) Create(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		payload = echocontext.GetPayload(c).(requestmodel.ReconciliationCreate)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Reconciliation(cs)
	)

	result, err := s.CreateWithClientData(ctx, payload)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}

	return response.R200(c, responsemodel.ResponseCreate{ID: result}, "")
}

// All godoc
// @tags Reconciliation
// @summary All
// @id admin-reconciliation-all
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.ReconciliationAll true "Query"
// @success 200 {object} responsemodel.ResponseReconciliationAll
// @router /reconciliations [get]
func (Reconciliation) All(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		qParams = echocontext.GetQuery(c).(requestmodel.ReconciliationAll)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Reconciliation(cs)
		q       = mgquery.AppQuery{
			Page:          qParams.Page,
			Limit:         qParams.Limit,
			SortInterface: bson.M{"createdAt": -1},
			Affiliate: mgquery.Affiliate{
				Status:   qParams.Status,
				Keyword:  qParams.Keyword,
				FromAt:   ptime.TimeParseISODate(qParams.FromAt),
				ToAt:     ptime.TimeParseISODate(qParams.ToAt),
				Campaign: qParams.Campaign,
				Source:   qParams.Source,
			},
		}
	)

	result := s.All(ctx, q)
	return response.R200(c, result, "")
}

// Detail godoc
// @tags Reconciliation
// @summary Detail
// @id admin-reconciliation-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Reconciliation id"
// @success 200 {object} responsemodel.ResponseReconciliationDetail
// @router /reconciliations/{id} [get]
func (Reconciliation) Detail(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		id  = echocontext.GetParam(c, "id").(primitive.ObjectID)
		cs  = echocontext.GetStaff(c).(externalauth.User)
		s   = service.Reconciliation(cs)
	)

	result, err := s.Detail(ctx, id)
	if err != nil {
		return response.R400(c, echo.Map{}, err.Error())
	}

	return response.R200(c, result, "")
}

// GetTransactionsByCondition  godoc
// @tags Reconciliation
// @summary GetTransactionsByConditions
// @id admin-reconciliation-getTransactionsByConditions
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Reconciliation id"
// @param payload query requestmodel.ReconciliationTransactionAll true "Query"
// @success 200 {object}  responsemodel.ResponseReconciliationAll
// @router /reconciliations/{id}/transactions [get]
func (Reconciliation) GetTransactionsByCondition(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Reconciliation(cs)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		qParams = echocontext.GetQuery(c).(requestmodel.ReconciliationTransactionAll)
		q       = mgquery.AppQuery{
			Page:  qParams.Page,
			Limit: qParams.Limit,
			Affiliate: mgquery.Affiliate{
				Keyword: qParams.Keyword,
				Status:  qParams.Status,
				Seller:  qParams.Seller,
			},
		}
	)

	result := s.GetTransactionsByCondition(ctx, q, id)
	return response.R200(c, result, "")
}

// GetStatistic  godoc
// @tags Reconciliation
// @summary GetStatistic
// @id admin-reconciliation-getStatistic
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Reconciliation id"
// @param payload query requestmodel.ReconciliationPayloadStatistic true "Query"
// @success 200 {object}  responsemodel.ResponseReconciliationStatistic
// @router /reconciliations/{id}/statistic [get]
func (Reconciliation) GetStatistic(c echo.Context) error {
	var (
		ctx   = echocontext.GetContext(c)
		cs    = echocontext.GetStaff(c).(externalauth.User)
		s     = service.Reconciliation(cs)
		id    = echocontext.GetParam(c, "id").(primitive.ObjectID)
		query = echocontext.GetQuery(c).(requestmodel.ReconciliationPayloadStatistic)
	)

	result, err := s.GetStatistic(ctx, id, query)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}

	return response.R200(c, result, "")
}

// ChangeStatus godoc
// @tags Reconciliation
// @summary ChangeStatus
// @id app-reconciliation-Change status
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Reconciliation id"
// @param payload body requestmodel.ReconciliationPayloadChangeStatus true "Payload status"
// @success 200 {object} responsemodel.ResponseChangeStatus
// @router /reconciliations/{id}/status [Patch]
func (h Reconciliation) ChangeStatus(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Reconciliation(cs)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.ReconciliationPayloadChangeStatus)
	)

	if err := h.checkPermissionReconciliation(c); err != nil {
		return response.R403(c, nil, err.Error())
	}

	result, err := s.ChangeStatus(ctx, id, payload)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}

	return response.R200(c, result, "")
}

// checkPermissionReconciliation ...
func (Reconciliation) checkPermissionReconciliation(c echo.Context) error {
	var (
		payload = echocontext.GetPayload(c).(requestmodel.ReconciliationPayloadChangeStatus)
		scopes  = make([]string, 0)
		cs      = echocontext.GetStaff(c).(externalauth.User)
	)

	switch payload.Status {
	case constant.ReconciliationStatus.Approved.Key:
		scopes = permission.Affiliate.Reconciliation.Approve
	case constant.ReconciliationStatus.Rejected.Key:
		scopes = permission.Affiliate.Reconciliation.Approve
	case constant.ReconciliationStatus.Deleted.Key:
		scopes = permission.Affiliate.Reconciliation.Edit
	case constant.ReconciliationStatus.Completed.Key:
		scopes = permission.Affiliate.Reconciliation.Edit
	default:
		return errors.New(errorcode.ReconciliationInvalidStatus)
	}

	err := externalauth.CheckPermission(scopes, externalauth.StaffCheckPermissionBody{
		StaffID:  cs.ID,
		DeviceID: echocontext.GetDeviceID(c),
		Token:    echocontext.GetToken(c),
	})

	return err
}
