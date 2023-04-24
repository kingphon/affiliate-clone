package service

import (
	"context"
	"fmt"
	"time"

	"git.selly.red/Selly-Modules/logger"
	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"

	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"

	requestmodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/request"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/constant"

	"git.selly.red/Selly-Modules/mongodb"

	"git.selly.red/Selly-Server/affiliate/external/utils/prandom"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
)

// MockDataInterface ...
type MockDataInterface interface {
	CreateMockData(ctx context.Context, campaignID, sellerID primitive.ObjectID)
	SandboxCrawl(d requestmodel.SandboxCrawlerTransaction)
}

// mockDataImplement ...
type mockDataImplement struct {
	CrawlService transactionCrawlImplement
}

// MockData ...
func MockData() MockDataInterface {
	return mockDataImplement{}
}

// SandboxCrawl ...
func (s mockDataImplement) SandboxCrawl(d requestmodel.SandboxCrawlerTransaction) {

	fmt.Println("Start crawl transaction ...!")
	t := responsemodel.ResponseTransactionCrawl{
		ClickID:         mongodb.ConvertStringToObjectID(d.ClickId),
		Campaign:        mongodb.ConvertStringToObjectID(d.Campaign),
		Source:          d.Source,
		Status:          d.Status,
		Commission:      d.Commission,
		TransactionID:   d.TransactionID,
		TransactionTime: ptime.TimeParseISODate(d.TransactionTime),
		RejectedReason:  d.RejectedReason,
		UpdatedHash:     d.UpdatedHash,
		UpdatedAt:       ptime.TimeParseISODate(d.UpdatedAt),
	}
	r := responsemodel.ResponseTransactionCrawlData{
		Code: 200,
		Data: struct {
			Data []responsemodel.ResponseTransactionCrawl `json:"data"`
		}{[]responsemodel.ResponseTransactionCrawl{t}},
		Message: "success",
	}

	// 1. Generate name temp
	nameTemp := fmt.Sprintf("%s_%d", prandom.RandomStringWithLength(5), time.Now().Unix())
	dao := database.TransactionTempByNameCol(nameTemp)

	// 2. Insert transaction temps
	if err := s.CrawlService.insertTransactionTemps(dao, r.Data.Data); err != nil {
		fmt.Println("insertTransactionTemps - error: ", err.Error())
		return
	}

	// 3. Check and update transaction from transaction-temps
	if err := s.CrawlService.checkAndUpdateTransactionFromTemp(dao); err != nil {
		logger.Error("Error-checkAndUpdateTransactionFromTemp:", logger.LogData{
			Data: bson.M{
				"error": err.Error(),
			},
		})
	}

	// 4. Update campaign statistic
	ctx := context.Background()
	campaignIds := s.CrawlService.distinctCampaignID(ctx, dao)
	fmt.Println("campaignIds : ", campaignIds)
	go Campaign(externalauth.User{}).UpdateStatisticListCampaign(ctx, campaignIds)

	// 5. Drop collection
	dao.Drop(ctx)
	fmt.Println("End crawl transaction ...")
}

// CreateMockData ...
func (s mockDataImplement) CreateMockData(ctx context.Context, campaignID, sellerID primitive.ObjectID) {
	platform := dao.Platform().FindOneByCondition(ctx, bson.D{
		{"campaignId", campaignID},
		{"status", "active"},
	})

	if platform.ID.IsZero() {
		return
	}

	// 1. Click
	newID := primitive.NewObjectID()
	click := mgaffiliate.Click{
		ID:         newID,
		CampaignID: campaignID,
		PlatformID: platform.ID,
		SellerID:   sellerID,
		Device: mgaffiliate.Device{
			Manufacturer: "Manufacturer",
			UserAgent:    "UserAgent",
			OSVersion:    "OSVersion",
			OSName:       "OSName",
			Model:        "Model",
		},
		CampaignURL:   platform.URL,
		PartnerSource: "access_trade",
		AffiliateURL:  platform.URL,
		CreatedAt:     ptime.Now(),
		UpdatedAt:     ptime.Now(),
		SearchString:  newID.Hex(),
	}
	dao.Click().InsertOne(ctx, click)

	// 2. Transaction
	code := prandom.RandomStringWithLength(10)
	transaction := mgaffiliate.Transaction{
		ID:                   primitive.NewObjectID(),
		SellerID:             sellerID,
		CampaignID:           campaignID,
		PlatformID:           platform.ID,
		SearchString:         mongodb.NonAccentVietnamese(code),
		Code:                 code,
		TransactionTime:      ptime.Now(),
		Source:               "access_trade",
		Commission:           75000,
		SellerCommissionRate: 80,
		SellerCommission:     60000,
		SellyCommission:      15000,
		EstimateCashbackAt:   ptime.Now().AddDate(0, 1, 0),
		Status:               constant.TransactionStatus.Pending.Key,
		Click: mgaffiliate.TransactionClick{
			ClickID:      click.ID,
			AffiliateURL: click.AffiliateURL,
		},
		Device: mgaffiliate.Device{
			Manufacturer: "Manufacturer",
			UserAgent:    "UserAgent",
			OSVersion:    "OSVersion",
			OSName:       "OSName",
			Model:        "Model",
		},
		UpdateHash: prandom.RandomStringWithLength(15),
		CreatedAt:  ptime.Now(),
		UpdatedAt:  ptime.Now(),
	}
	dao.Transaction().InsertOne(ctx, transaction)

	// 3. Transaction History
	transactionHistory := mgaffiliate.TransactionHistory{
		ID:            primitive.NewObjectID(),
		TransactionID: transaction.ID,
		SellerID:      transaction.SellerID,
		Status:        constant.TransactionStatus.Pending.Key,
		Desc:          "Ghi nhận",
		CreatedAt:     ptime.Now(),
	}
	dao.TransactionHistory().InsertOne(ctx, transactionHistory)

	// 4. Seller campaign statistic
	sellerCampaignStats := mgaffiliate.SellerCampaignStatistic{
		ID:         primitive.NewObjectID(),
		SellerID:   sellerID,
		CampaignID: campaignID,
		Statistic: mgaffiliate.Statistic{
			TransactionTotal:    1,
			TransactionApproved: 0,
			TransactionPending:  1,
			TransactionCashback: 0,
			TransactionRejected: 0,
			CommissionApproved:  0,
			CommissionPending:   60000,
			CommissionCashback:  0,
			CommissionRejected:  0,
		},
		CreatedAt: ptime.Now(),
		UpdatedAt: ptime.Now(),
	}
	dao.SellerCampaignStatistic().InsertOne(ctx, sellerCampaignStats)

	// 5. Seller statistic
	sellerStats := mgaffiliate.SellerStatistic{
		ID:       primitive.NewObjectID(),
		SellerID: sellerID,
		Statistic: mgaffiliate.Statistic{
			TransactionTotal:    1,
			TransactionApproved: 0,
			TransactionPending:  1,
			TransactionCashback: 0,
			TransactionRejected: 0,
			CommissionApproved:  0,
			CommissionPending:   60000,
			CommissionCashback:  0,
			CommissionRejected:  0,
		},
		CreatedAt: ptime.Now(),
		UpdatedAt: ptime.Now(),
	}
	dao.SellerStatistic().InsertOne(ctx, sellerStats)
}
