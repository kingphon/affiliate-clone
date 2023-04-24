package routevalidation

import (
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/request"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction ...
type Transaction struct{}

// All ...
func (Transaction) All(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var query requestmodel.TransactionAll

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
func (Transaction) Detail(next echo.HandlerFunc) echo.HandlerFunc {
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
