package service

import (
	"context"
	"fmt"

	adminconstant "git.selly.red/Selly-Server/affiliate/pkg/admin/constant"

	"git.selly.red/Selly-Server/affiliate/pkg/admin/locale"

	"git.selly.red/Selly-Server/affiliate/external/constant"

	jsmodel "git.selly.red/Selly-Modules/natsio/js/model"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
)

// getDataPushNotificationByTransaction ...
func (s transactionCrawlImplement) getDataPushNotificationByTransaction(ctx context.Context, isNewTransaction bool, t mgaffiliate.Transaction) (result []jsmodel.PushNotification) {
	campaign, _ := s.CampaignService.FindByID(ctx, t.CampaignID)

	var (
		statusConst = constant.TransactionStatus
		status      string
	)

	if isNewTransaction && t.IsValidStatuses([]string{statusConst.Pending.Key, statusConst.Approved.Key}) {
		status = statusConst.New.Key
	} else {
		if t.IsValidStatus(statusConst.Approved.Key) {
			return
		}
		status = t.Status
	}

	n := s.convertToJestStreamPayloadNotification(t, campaign.Name, status)
	result = append(result, n)
	return
}

// convertToJestStreamPayloadNotification ...
func (s transactionCrawlImplement) convertToJestStreamPayloadNotification(t mgaffiliate.Transaction, campaignName, status string) jsmodel.PushNotification {
	return jsmodel.PushNotification{
		User:     t.SellerID.Hex(),
		Type:     adminconstant.NotificationTypeAffiliateTransactionChangeStatus,
		TargetID: t.ID.Hex(),
		Category: adminconstant.NotificationCategoryAffiliate,
		Options:  s.getTitleAndContentNotification(status, t.Code, campaignName),
	}
}

// getTitleAndContentNotification ...
func (transactionCrawlImplement) getTitleAndContentNotification(typeStr, code, campaignName string) (result jsmodel.NotificationOptions) {
	switch typeStr {
	case constant.TransactionStatus.New.Key:
		result.Title = locale.NotificationTitleTransactionNew
		result.Content = fmt.Sprintf(locale.NotificationContentTransactionNew, code, campaignName)
	case constant.TransactionStatus.Cashback.Key:
		result.Title = locale.NotificationTitleTransactionCashback
		result.Content = fmt.Sprintf(locale.NotificationContentTransactionCashback, code, campaignName)
	case constant.TransactionStatus.Rejected.Key:
		result.Title = locale.NotificationTitleTransactionRejected
		result.Content = fmt.Sprintf(locale.NotificationContentTransactionRejected, code, campaignName)
	}
	return
}
