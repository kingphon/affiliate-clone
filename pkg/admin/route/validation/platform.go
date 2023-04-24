package routevalidation

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/errorcode"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
)

// Platform ...
type Platform struct{}

// Detail ...
func (Platform) Detail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var id = c.Param("id")

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.PlatformNotFound)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		echocontext.SetParam(c, "id", objID)
		return next(c)
	}
}

// Update ...
func (Platform) Update(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.PlatformCreate

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

// ChangeStatus ...
func (Platform) ChangeStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.PlatformChangeStatus

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

// Create ...
func (Platform) Create(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			payload requestmodel.PlatformCreate
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
