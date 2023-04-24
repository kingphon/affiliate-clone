package responsemodel

import (
	"git.selly.red/Selly-Server/affiliate/external/utils/file"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseCampaignBrief ...
type ResponseCampaignBrief struct {
	ID                   primitive.ObjectID     `json:"_id"`
	Name                 string                 `json:"name"`
	Logo                 *file.FilePhoto        `json:"logo"`
	Commission           float64                `json:"commission"`
	Statistic            CampaignStatisticBrief `json:"statistic"`
	AllowShowShareAction bool                   `json:"allowShowShareAction"`
	From                 *ptime.TimeResponse    `json:"from,omitempty"`
	To                   *ptime.TimeResponse    `json:"to,omitempty"`
}

// CampaignStatisticBrief ...
type CampaignStatisticBrief struct {
	RewardTotal int64 `json:"rewardTotal"`
}

// ResponseCampaignDetail ...
type ResponseCampaignDetail struct {
	ID                   primitive.ObjectID     `json:"_id"`
	Name                 string                 `json:"name"`
	Logo                 *file.FilePhoto        `json:"logo"`
	Covers               []*file.FilePhoto      `json:"covers"`
	Desc                 string                 `json:"desc"`
	ShareDesc            string                 `json:"shareDesc"`
	Commission           float64                `json:"commission"`
	Statistic            CampaignStatisticBrief `json:"statistic"`
	From                 *ptime.TimeResponse    `json:"from,omitempty"`
	To                   *ptime.TimeResponse    `json:"to,omitempty"`
	Platforms            []PlatformBrief        `json:"platforms"`
	IsCanJoinCampaign    bool                   `json:"isCanJoinCampaign"`
	AllowShowShareAction bool                   `json:"allowShowShareAction"`
}

// PlatformBrief ...
type PlatformBrief struct {
	ID       primitive.ObjectID `json:"_id"`
	Code     string             `json:"code"`
	Platform string             `json:"platform"`
}

// ResponseGenerateShareURL ...
type ResponseGenerateShareURL struct {
	CampaignID primitive.ObjectID         `json:"campaignId"`
	Platforms  []PlatformGenerateShareURL `json:"platforms"`
}

// PlatformGenerateShareURL ...
type PlatformGenerateShareURL struct {
	ID        primitive.ObjectID `json:"_id"`
	Platform  string             `json:"platform"`
	ShareURL  string             `json:"shareURL"`
	Title     string             `json:"title"`
	ShareCode string             `json:"shareCode"`
	OpenVia   string             `json:"openVia"`
}

// ResponseAffiliateLink ...
type ResponseAffiliateLink struct {
	URL     string `json:"url"`
	ClickID string `json:"clickId"`
}

// ResponseCampaignShortInfo ...
type ResponseCampaignShortInfo struct {
	ID         primitive.ObjectID `json:"_id"`
	Name       string             `json:"name"`
	Commission float64            `json:"commission"`
	Logo       *file.FilePhoto    `json:"logo"`
}

// ResponseCampaignFilter ...
type ResponseCampaignFilter struct {
	Data []CampaignFilter `json:"data"`
}

// CampaignFilter ...
type CampaignFilter struct {
	Key   string     `json:"key"`
	Title string     `json:"title"`
	List  []KeyValue `json:"list"`
}

// KeyValue ...
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"text"`
}
