package service

import (
	"context"
	"fmt"
	"net/url"

	"git.selly.red/Selly-Server/affiliate/internal/config"

	"git.selly.red/Selly-Server/affiliate/external/constant"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// generateAffLink ...
func (s platformImplement) generateAffLink(ctx context.Context, platform mgaffiliate.Platform, clickID primitive.ObjectID) (url string) {
	switch platform.Partner.Source {
	case constant.PlatformPartnerSourceAccessTrade:
		url = s.generateFormatSourceAccessTrade(ctx, platform, clickID)
	case constant.PlatformPartnerSourceOther:
		url = platform.URL
	}
	return
}

// generateFormatSourceAccessTrade ...
func (platformImplement) generateFormatSourceAccessTrade(ctx context.Context, platform mgaffiliate.Platform, clickID primitive.ObjectID) string {
	var host = config.GetENV().Deeplink.AccessTradeDeepLink
	return fmt.Sprintf("%s/%s?url=%s&utm_campaign=%s&utm_content=%s",
		host,
		platform.Partner.CampaignID,
		url.QueryEscape(platform.URL),
		platform.CampaignID.Hex(),
		clickID.Hex(),
	)
}
