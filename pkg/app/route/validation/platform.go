package routevalidation

import (
	"fmt"

	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/checksum"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/internal/config"
	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/request"
	"github.com/friendsofgo/errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// GenerateAffiliateLink ...
func (v Platform) GenerateAffiliateLink(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload requestmodel.GenerateAffiliateLinkBody

		if err := c.Bind(&payload); err != nil {
			return response.R400(c, nil, "")
		}

		if err := payload.Validate(); err != nil {
			return response.RouteValidation(c, err)
		}

		// Checksum check valid
		if isValid := v.isValidCheckSumGenerateAffiliateLink(payload); !isValid {
			return response.R400(c, nil, errors.New(errorcode.ChecksumInvalid).Error())
		}

		echocontext.SetPayload(c, payload)
		return next(c)
	}
}

// isValidCheckSumGenerateAffiliateLink ...
func (Platform) isValidCheckSumGenerateAffiliateLink(payload requestmodel.GenerateAffiliateLinkBody) (isValid bool) {
	var stringData = fmt.Sprintf("%s%s", payload.Code, config.GetENV().ChecksumKey)

	checksumString := checksum.Generate(stringData)
	if checksumString == payload.Checksum {
		isValid = true
	}
	return
}
