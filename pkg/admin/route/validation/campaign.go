package routevalidation

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"

	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
)

// Campaign ...
type Campaign struct {
}

// Create ...
func (Campaign) Create(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.CampaignCreate

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

// Update ...
func (Campaign) Update(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.CampaignUpdate

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

// Detail ...
func (Campaign) Detail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var id = c.Param("id")

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.CampaignInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// All ...
func (Campaign) All(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		var query requestmodel.CampaignAll

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

// ChangeStatus ...
func (Campaign) ChangeStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.CampaignChangeStatus

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
