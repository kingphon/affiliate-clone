package routevalidation

import (
	"git.selly.red/Selly-Modules/authentication"
	"git.selly.red/Selly-Modules/mongodb"
	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/auth/permission"
	"git.selly.red/Selly-Server/affiliate/internal/middleware"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
)

// Reconciliation ...
type Reconciliation struct{}

// Create ...
func (Reconciliation) Create(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			payload requestmodel.ReconciliationCreate
		)

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetPayload(c, payload)
		return next(c)
	}
}

func (Reconciliation) All(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			query requestmodel.ReconciliationAll
		)

		if err := c.Bind(&query); err != nil {
			return response.R400(c, nil, "")
		}

		if err := query.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetQuery(c, query)
		return next(c)
	}
}

// Detail ...
func (Reconciliation) Detail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var id = c.Param("id")

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.ReconciliationInvalidID)
		}

		objID, _ := mongodb.NewIDFromString(id)

		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// ChangeStatus ...
func (Reconciliation) ChangeStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			payload requestmodel.ReconciliationPayloadChangeStatus
		)

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		// Check auth GG
		cs := echocontext.GetStaff(c).(externalauth.User)
		dataVerifyCode := authentication.StaffVerifyCodeBody{
			Source:  permission.AffiliateSource,
			StaffID: cs.ID,
			Code:    payload.CodeAuthGG,
		}
		if err := middleware.VerifyOtpCode(dataVerifyCode); err != nil {
			return response.R403(c, echo.Map{}, err.Error())
		}

		echocontext.SetPayload(c, payload)
		return next(c)
	}
}

// GetTransactionsByCondition ...
func (Reconciliation) GetTransactionsByCondition(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			query requestmodel.ReconciliationTransactionAll
		)

		if err := c.Bind(&query); err != nil {
			return response.R400(c, nil, "")
		}

		if err := query.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetQuery(c, query)
		return next(c)
	}
}

// GetStatistic ...
func (Reconciliation) GetStatistic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			query requestmodel.ReconciliationPayloadStatistic
		)

		if err := c.Bind(&query); err != nil {
			return response.R400(c, nil, "")
		}

		if err := query.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		echocontext.SetQuery(c, query)
		return next(c)
	}
}
