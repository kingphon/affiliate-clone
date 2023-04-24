package responsemodel

import (
	"git.selly.red/Selly-Server/affiliate/external/utils/file"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseShareURLInfo ...
type ResponseShareURLInfo struct {
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Cover      *file.FilePhoto    `json:"cover"`
	PlatformID primitive.ObjectID `json:"platformId"`
	CampaignID primitive.ObjectID `json:"campaignId"`
}
