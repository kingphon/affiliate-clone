package handler

import (
	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/affiliate/external/constant"
	"git.selly.red/Selly-Server/affiliate/external/response"
	"git.selly.red/Selly-Server/affiliate/external/utils/echocontext"
	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	"git.selly.red/Selly-Server/affiliate/external/utils/routemiddleware"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/request"
	"git.selly.red/Selly-Server/affiliate/pkg/app/service"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

// Campaign ...
type Campaign struct{}

// GroupItems godoc
// @tags Campaign
// @summary GroupItems
// @id app-campaign-group-list-items
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param id path string true "Campaign group type"
// @param payload query requestmodel.CampaignAll false "Query"
// @success 200 {object} responsemodel.ResponseList
// @router /campaigns/groups/{id}/items [get]
func (cp Campaign) GroupItems(c echo.Context) error {
	return cp.All(c)
}

// GroupAll godoc
// @tags Campaign
// @summary Group Campaign All
// @id app-campaign-group-list
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.CampaignGroupAll false "Query"
// @success 200 {object} responsemodel.ResponseListCampaignGroupAll
// @router /campaigns/groups [get]
func (Campaign) GroupAll(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		s       = service.Campaign()
		qParams = echocontext.GetQuery(c).(requestmodel.CampaignGroupAll)
		q       = mgquery.AppQuery{
			Screen: qParams.Screen,
		}
	)

	result := s.GroupAll(ctx, q)
	return response.R200(c, result, "")
}

// GroupDetail godoc
// @tags Campaign
// @summary Group Detail
// @id app-campaign-group-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param id path string true "Campaign group type"
// @success 200 {object} responsemodel.ResponseCampaignGroupBrief
// @router /campaigns/groups/{id} [get]
func (Campaign) GroupDetail(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.Campaign()
		id  = echocontext.GetParam(c, "id").(string)
	)

	result := s.GroupDetail(ctx, id)
	return response.R200(c, result, "")
}

// All godoc
// @tags Campaign
// @summary All
// @id app-campaign-list
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @param payload query requestmodel.CampaignAll false "Query"
// @success 200 {object} responsemodel.ResponseList
// @router /campaigns [get]
func (Campaign) All(c echo.Context) error {
	var (
		ctx       = echocontext.GetContext(c)
		s         = service.Campaign()
		qParams   = echocontext.GetQuery(c).(requestmodel.CampaignAll)
		pageToken = echocontext.PageTokenDecode(qParams.PageToken)
		q         = mgquery.AppQuery{
			Page:    pageToken.Page,
			Limit:   int64(constant.Limit20),
			Keyword: qParams.Keyword,
			SortInterface: bson.D{
				{"order", -1},
				{"createdAt", -1},
			},
			SortStr: qParams.Sort,
		}
	)

	result := s.All(ctx, q)
	return response.R200(c, result, "")
}

// UpdateClick godoc
// @tags Campaign
// @summary UpdateClick
// @id update-campaign-click
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  clickID path string true "Click id"
// @param payload body requestmodel.ClickUpdateBody true "Payload"
// @success 200 {object} nil
// @router /campaigns/click/{clickID} [put]
func (Campaign) UpdateClick(c echo.Context) error {
	var (
		ctx     = echocontext.GetContext(c)
		s       = service.Campaign()
		clickID = echocontext.GetParam(c, "clickId").(primitive.ObjectID)
		payload = echocontext.GetPayload(c).(requestmodel.ClickUpdateBody)
	)
	err := s.UpdateClick(ctx, clickID, payload)
	if err != nil {
		return response.R400(c, nil, err.Error())
	}
	return response.R200(c, nil, "")
}

// Detail godoc
// @tags Campaign
// @summary Detail
// @id app-campaign-detail
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Campaign id"
// @success 200 {object} responsemodel.ResponseDetail
// @router /campaigns/{id} [get]
func (Campaign) Detail(c echo.Context) error {
	var (
		ctx         = echocontext.GetContext(c)
		s           = service.Campaign()
		id          = echocontext.GetParam(c, "id").(primitive.ObjectID)
		sellerID, _ = mongodb.NewIDFromString(echocontext.GetCurrentUserID(c))
	)

	result, err := s.Detail(ctx, id, sellerID)
	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}

// GenerateShareURL godoc
// @tags Campaign
// @summary GenerateShareURL
// @id app-campaign-generate-share-url
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Campaign id"
// @success 200 {object} responsemodel.ResponseGenerateShareURL
// @router /campaigns/{id}/share-url [patch]
func (Campaign) GenerateShareURL(c echo.Context) error {
	var (
		ctx    = echocontext.GetContext(c)
		s      = service.Campaign()
		id     = echocontext.GetParam(c, "id").(primitive.ObjectID)
		userID = echocontext.GetCurrentUserID(c)
		osName = echocontext.GetHeaders(c).Get(routemiddleware.HeaderOSName)
	)
	result, err := s.GenerateShareURL(ctx, id, userID)
	if strings.ToLower(osName) == constant.OSNameAndroid {
		for i, p := range result.Platforms {
			if p.Platform == constant.PlatformAll {
				result.Platforms[i].Platform = constant.PlatformAndroid
			}
		}
	}

	if err != nil {
		return response.R404(c, echo.Map{}, err.Error())
	}
	return response.R200(c, result, "")
}

// Filter godoc
// @tags Campaign
// @summary Filter
// @id app-campaign-filter
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @success 200 {object} responsemodel.ResponseCampaignFilter
// @router /campaigns/filter [get]
func (Campaign) Filter(c echo.Context) error {
	var (
		ctx = echocontext.GetContext(c)
		s   = service.Campaign()
	)

	result := s.GetCampaignFilter(ctx)
	return response.R200(c, result, "")
}

// GetSellerCampaignStatistic godoc
// @tags Campaign
// @summary Get seller campaign statistic
// @id app-get-seller-campaign-statistic
// @security ApiKeyAuth
// @accept json
// @produce json
// @param device-id header string true "Device id"
// @Param  id path string true "Campaign id"
// @success 200 {object} responsemodel.ResponseGetSellerCampaignStatistic
// @router /campaigns/{id}/statistic [get]
func (Campaign) GetSellerCampaignStatistic(c echo.Context) error {
	var (
		ctx         = echocontext.GetContext(c)
		s           = service.Campaign()
		id          = echocontext.GetParam(c, "id").(primitive.ObjectID)
		sellerID, _ = mongodb.NewIDFromString(echocontext.GetCurrentUserID(c))
	)

	result := s.GetSellerCampaignStatistic(ctx, id, sellerID)
	return response.R200(c, result, "")
}
