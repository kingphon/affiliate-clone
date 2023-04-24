package handler

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"

	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/service"
)

// Campaign ...
type Campaign struct{}

// Create godoc
// @tags Campaign
// @summary Create
// @id app-campaign-create
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload body requestmodel.CampaignCreate true "Payload"
// @success 200 {object} responsemodel.ResponseCreate
// @router /campaigns [post]
func (Campaign) Create(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		payload = echocontext.GetPayload(c).(requestmodel.CampaignCreate)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Campaign(cs)
	)

	result, err := s.CreateWithClientData(ctx, payload)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}

	return response.R200(c, responsemodel.ResponseCreate{ID: result}, "")
}

// Update godoc
// @tags Campaign
// @summary Update
// @id app-campaign-update
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Campaign id"
// @param payload body requestmodel.CampaignUpdate true "Payload"
// @success 200 {object} responsemodel.ResponseUpdate
// @router /campaigns/{id} [put]
func (Campaign) Update(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Campaign(cs)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.CampaignUpdate)
	)

	result, err := s.Update(ctx, id, payload)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}
	return response.R200(c, responsemodel.ResponseUpdate{ID: result}, "")
}

// All godoc
// @tags Campaign
// @summary All
// @id admin-campaign-all
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.CampaignAll true "Query"
// @success 200 {object} responsemodel.ResponseCampaignAll
// @router /campaigns [get]
func (Campaign) All(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		qParams = echocontext.GetQuery(c).(requestmodel.CampaignAll)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Campaign(cs)
		q       = mgquery.AppQuery{
			Page:    qParams.Page,
			Limit:   qParams.Limit,
			SortStr: qParams.Sort,
			Affiliate: mgquery.Affiliate{
				Status:  qParams.Status,
				Keyword: qParams.Keyword,
				FromAt:  ptime.TimeParseISODate(qParams.FromAt),
				ToAt:    ptime.TimeParseISODate(qParams.ToAt),
			},
		}
	)

	result := s.All(ctx, q)
	return response.R200(c, result, "")
}

// ChangeStatus godoc
// @tags Campaign
// @summary Change Status
// @id app-campaign-change status
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Campaign id"
// @param payload body requestmodel.CampaignChangeStatus true "Payload"
// @success 200 {object} responsemodel.ResponseChangeStatus
// @router /campaigns/{id}/status [patch]
func (Campaign) ChangeStatus(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		cs      = echocontext.GetStaff(c).(externalauth.User)
		s       = service.Campaign(cs)
		id      = echocontext.GetParam(c, "id").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.CampaignChangeStatus)
	)

	result, err := s.ChangeStatus(ctx, id, payload)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}
	return response.R200(c, result, "")
}

// Detail godoc
// @tags Campaign
// @summary Detail
// @id admin-campaign-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Campaign id"
// @success 200 {object} responsemodel.ResponseCampaignDetail
// @router /campaigns/{id} [get]
func (Campaign) Detail(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.Campaign(externalauth.User{})
		id  = echocontext.GetParam(c, "id").(primitive.ObjectID)
	)

	result, err := s.Detail(ctx, id)
	if err != nil {
		response.R400(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}

// GetPlatformByCampaign godoc
// @tags Campaign
// @summary Detail
// @id admin-campaign-getPlatformByCampaign
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Campaign id"
// @success 200 {object} responsemodel.ResponsePlatformInfos
// @router /campaigns/{id}/platforms [get]
func (Campaign) GetPlatformByCampaign(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.Campaign(externalauth.User{})
		id  = echocontext.GetParam(c, "id").(primitive.ObjectID)
	)

	result := s.GetPlatformByCampaign(ctx, id)
	return response.R200(c, result, "")
}

// CreatePlatform godoc
// @tags Campaign
// @summary Create platform
// @id admin-campaign-create-platform
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Campaign id"
// @param payload body requestmodel.PlatformCreate true "Payload"
// @success 200 {object} responsemodel.ResponseCreate
// @router /campaigns/{id}/platforms [post]
func (Campaign) CreatePlatform(c echo.Context) error {
	var (
		ctx        = echocontext.GetContext(c)
		payload    = echocontext.GetPayload(c).(requestmodel.PlatformCreate)
		cs         = echocontext.GetStaff(c).(externalauth.User)
		s          = service.Campaign(cs)
		campaignID = echocontext.GetParam(c, "id").(primitive.ObjectID)
	)

	result, err := s.CreatePlatformWithClientData(ctx, campaignID, payload)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}

	return response.R200(c, responsemodel.ResponseCreate{ID: result}, "")
}
