package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

// GetByCampaign ...
func (s platformImplement) GetByCampaign(ctx context.Context, id primitive.ObjectID) (result []responsemodel.ResponsePlatformInfo) {
	var (
		d    = dao.Platform()
		cond = bson.M{"campaignId": id}
	)

	docs := d.FindByCondition(ctx, cond)

	result = make([]responsemodel.ResponsePlatformInfo, 0)
	for _, doc := range docs {
		result = append(result, s.info(ctx, doc))
	}

	return
}

// info ...
func (s platformImplement) info(ctx context.Context, doc mgaffiliate.Platform) responsemodel.ResponsePlatformInfo {
	return responsemodel.ResponsePlatformInfo{
		ID:       doc.ID.Hex(),
		Code:     doc.Code,
		Status:   doc.Status,
		URL:      doc.URL,
		Platform: doc.PlatformType,
		Partner: responsemodel.ResponsePlatformPartner{
			Source:     doc.Partner.Source,
			CampaignID: doc.Partner.CampaignID,
		},
	}

}
