package routevalidation

import (
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/external/utils/parray"
	appconstant "git.selly.red/Selly-Server/affiliate/pkg/app/constant"
	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/request"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Campaign ...
type Campaign struct{}

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

// ClickDetail ...
func (Campaign) ClickDetail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var id = c.Param("clickID")

		if !primitive.IsValidObjectID(id) {
			return response.R404(c, nil, errorcode.CampaignInvalidID)
		}

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return response.R400(c, nil, "")
		}

		echocontext.SetParam(c, "clickId", objID)
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

// GroupDetail ...
func (Campaign) GroupDetail(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id   = c.Param("id")
			list = []string{
				appconstant.CampaignGroupType.CampaignList.Key,
			}
		)

		if !parray.ContainsStr(list, id) {
			return response.R404(c, nil, errorcode.CampaignGroupNotFound)
		}

		echocontext.SetParam(c, "id", id)
		return next(c)
	}
}

// GroupAll ...
func (Campaign) GroupAll(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var query requestmodel.CampaignGroupAll

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

// UpdateClickBody ...
func (c Campaign) UpdateClickBody(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.ClickUpdateBody

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
