package service

import (
	"context"

	"git.selly.red/Selly-Modules/natsio/client"
	natsmodel "git.selly.red/Selly-Modules/natsio/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SellerInterface ...
type SellerInterface interface {
	GetSellerByIDs(ctx context.Context, sellerIDs []primitive.ObjectID) (sellers []natsmodel.ResponseSellerInfo, err error)
}

// sellerImplement ...
type sellerImplement struct{}

// Seller ...
func Seller() SellerInterface {
	return &sellerImplement{}
}

// GetSellerByIDs ...
func (s sellerImplement) GetSellerByIDs(ctx context.Context, sellerIDs []primitive.ObjectID) (sellers []natsmodel.ResponseSellerInfo, err error) {
	payload := natsmodel.GetListSellerByIDsRequest{SellerIDs: sellerIDs}
	result, err := client.GetSeller().GetListSellerInfoByIDs(payload)
	if err != nil {
		return
	}

	sellers = result.Sellers
	return
}
