package service

import (
	"context"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"

	"git.selly.red/Selly-Server/affiliate/external/utils/file"

	"git.selly.red/Selly-Server/affiliate/pkg/app/errorcode"
	"github.com/friendsofgo/errors"

	"git.selly.red/Selly-Server/affiliate/pkg/app/dao"
	"go.mongodb.org/mongo-driver/bson"

	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/response"
)

// ShareURLInterface ...
type ShareURLInterface interface {
	// GetInfoShareURL ...
	GetInfoShareURL(ctx context.Context, code string) (result *responsemodel.ResponseShareURLInfo, err error)

	// GetSellerShareURLByCode ...
	GetSellerShareURLByCode(ctx context.Context, code string) (result *mgaffiliate.SellerShareURL, err error)
}

// ShareURL return platform service
func ShareURL() ShareURLInterface {
	return shareURLImplement{}
}

// shareURLImplement ...
type shareURLImplement struct{}

// GetSellerShareURLByCode ...
func (s shareURLImplement) GetSellerShareURLByCode(ctx context.Context, code string) (result *mgaffiliate.SellerShareURL, err error) {
	var (
		d    = dao.SellerShareURL()
		cond = bson.D{{"code", code}}
	)

	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.CodeInvalid)
		return
	}

	result = &doc
	return
}

// GetInfoShareURL ...
func (s shareURLImplement) GetInfoShareURL(ctx context.Context, code string) (result *responsemodel.ResponseShareURLInfo, err error) {
	var (
		d    = dao.SellerShareURL()
		cond = bson.D{{"code", code}}
	)

	doc := d.FindOneByCondition(ctx, cond)
	if doc.ID.IsZero() {
		err = errors.New(errorcode.CodeInvalid)
		return
	}

	campaign := campaignImplement{}.FindCampaignActiveByID(ctx, doc.CampaignID)
	if campaign.ID.IsZero() {
		err = errors.New(errorcode.CodeInvalid)
		return
	}

	var cover *file.FilePhoto
	if len(campaign.Covers) > 0 {
		cover = fileResponseImplement{}.ConvertResponseFilePhoto(campaign.Covers[0]).GetResponseData()
	}

	result = &responsemodel.ResponseShareURLInfo{
		Title:      campaign.Name,
		Content:    campaign.ShareDesc,
		Cover:      cover,
		PlatformID: doc.PlatformID,
		CampaignID: campaign.ID,
	}

	return
}
