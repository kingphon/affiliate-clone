package handler

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/service"
)

// Platform ...
type Platform struct{}

// Update godoc
// @tags Platform
// @summary Update
// @id app-platform-update
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Platform id"
// @param payload body requestmodel.PlatformCreate true "Payload"
// @success 200 {object} responsemodel.ResponseUpdate
// @router /campaign-platforms/{id} [put]
func (Platform) Update(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Platform(cs)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.PlatformCreate)
	)

	result, err := s.Update(ctx, id, payload)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}
	return response.R200(c, responsemodel.ResponseUpdate{ID: result}, "")
}

// ChangeStatus godoc
// @tags Platform
// @summary Change Status
// @id app-platform-change status
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Platform id"
// @param payload body requestmodel.PlatformChangeStatus true "Payload"
// @success 200 {object} responsemodel.ResponseChangeStatus
// @router /campaign-platforms/{id}/status [patch]
func (Platform) ChangeStatus(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Platform(cs)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.PlatformChangeStatus)
	)

	result, err := s.ChangeStatus(ctx, id, payload)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}
	return response.R200(c, result, "")
}
