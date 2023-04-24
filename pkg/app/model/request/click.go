package requestmodel

import (
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PayloadInsertClick ...
type PayloadInsertClick struct {
	SellerID           primitive.ObjectID
	CampaignID         primitive.ObjectID
	Device             mgaffiliate.Device
	Platform           mgaffiliate.Platform
	CampaignCommission mgaffiliate.CampaignCommission
	From               string
	ShareURL           string
}
