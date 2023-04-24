package service

import (
	"context"
	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/app/model/request"

	"git.selly.red/Selly-Server/affiliate/external/constant"

	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/affiliate/pkg/app/dao"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ClickInterface ...
type ClickInterface interface {
	// GenerateDocAndInsert ...
	GenerateDocAndInsert(ctx context.Context, payload requestmodel.PayloadInsertClick) (clickID primitive.ObjectID, err error)

	// UpdateByID ...
	UpdateByID(ctx context.Context, id primitive.ObjectID, payload interface{}) (err error)
}

// Click return click service
func Click() ClickInterface {
	return clickImplement{}
}

// clickImplement ...
type clickImplement struct{}

// UpdateByID ...
func (c clickImplement) UpdateByID(ctx context.Context, id primitive.ObjectID, payload interface{}) (err error) {
	err = dao.Click().UpdateOne(ctx, bson.D{{"_id", id}}, payload)
	return
}

// GenerateDocAndInsert ...
func (c clickImplement) GenerateDocAndInsert(ctx context.Context, payload requestmodel.PayloadInsertClick) (clickID primitive.ObjectID, err error) {
	newID := primitive.NewObjectID()
	raw := mgaffiliate.Click{
		ID:            newID,
		SearchString:  newID.Hex(),
		CampaignID:    payload.CampaignID,
		PlatformID:    payload.Platform.ID,
		SellerID:      payload.SellerID,
		Device:        payload.Device,
		CampaignURL:   payload.Platform.URL,
		ShareURL:      payload.ShareURL,
		PartnerSource: payload.Platform.Partner.Source,
		CreatedAt:     ptime.Now(),
		UpdatedAt:     ptime.Now(),
		Status:        constant.CampaignTransactionClickStatusPending,
		Commission:    payload.CampaignCommission,
		From:          payload.From,
	}

	// 2. Insert
	err = dao.Click().InsertOne(ctx, raw)
	clickID = raw.ID
	return
}
