package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"

	clientjstream "git.selly.red/Selly-Server/affiliate/internal/nats/jstream"

	jsmodel "git.selly.red/Selly-Modules/natsio/js/model"

	"git.selly.red/Selly-Modules/mongodb"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"

	"git.selly.red/Selly-Server/affiliate/external/constant"

	"git.selly.red/Selly-Modules/logger"

	"github.com/friendsofgo/errors"

	"github.com/panjf2000/ants/v2"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Server/affiliate/external/utils/prandom"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/database"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"

	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"

	httprequest "git.selly.red/Selly-Server/affiliate/external/utils/http"
	"git.selly.red/Selly-Server/affiliate/internal/config"
)

// TransactionCrawlInterface ...
type TransactionCrawlInterface interface {
	// Crawl ...
	Crawl(fromAt, toAt string)

	// CrawlTransactionData Crawl transaction data form system job...
	CrawlTransactionData()

	CheckAndUpdateTransactionsBySeller(ctx context.Context, doc responsemodel.AggregateTransactionTempGroupSeller) (errors []error)

	// CrawlTransactionDataFromAdmin ...
	CrawlTransactionDataFromAdmin()
}

// transactionCrawlImplement ...
type transactionCrawlImplement struct {
	TransactionService             transactionImplement
	TransactionHistoryService      transactionHistoryImplement
	SellerCampaignStatisticService sellerCampaignStatisticImplement
	SellerStatisticService         sellerStatisticImplement
	CampaignService                campaignImplement
	ClickService                   clickImplement
}

// TransactionCrawl ...
func TransactionCrawl() TransactionCrawlInterface {
	return transactionCrawlImplement{}
}

// CrawlTransactionDataFromAdmin ...
func (s transactionCrawlImplement) CrawlTransactionDataFromAdmin() {
	var (
		nowAt  = time.Now()
		toAt   = nowAt.Format(ptime.DateLayoutFull)
		fromAt = nowAt.Add(time.Duration(-240) * time.Minute).Format(ptime.DateLayoutFull)
	)

	s.Crawl(fromAt, toAt)
}

// CrawlTransactionData ...
func (s transactionCrawlImplement) CrawlTransactionData() {
	var (
		nowAt  = time.Now()
		toAt   = nowAt.Format(ptime.DateLayoutFull)
		fromAt = nowAt.Add(time.Duration(-120) * time.Minute).Format(ptime.DateLayoutFull)
	)

	s.Crawl(fromAt, toAt)
}

// Crawl ...
func (s transactionCrawlImplement) Crawl(fromAt, toAt string) {
	var (
		url    = fmt.Sprintf("%s/selly/transactions", config.GetENV().Crawl.URL)
		token  = config.GetENV().Crawl.Auth
		limit  = 500
		params = httprequest.RequestParam{
			"start": fromAt,
			"end":   toAt,
			"limit": strconv.Itoa(limit),
		}
		page = 0
	)

	fmt.Println("Start crawl transaction ...!")

	// 1. Generate name temp
	nameTemp := fmt.Sprintf("%s_%d", prandom.RandomStringWithLength(5), time.Now().Unix())
	dao := database.TransactionTempByNameCol(nameTemp)

	// 2. Insert transaction temps
	for {
		params["page"] = strconv.Itoa(page)
		_, resBody, err := httprequest.RequestGetWithParams(url, token, params)
		if err != nil {
			fmt.Println("RequestGetWithParams - error: ", err.Error())
			break
		}

		var r responsemodel.ResponseTransactionCrawlData
		if err = json.Unmarshal(resBody, &r); err != nil {
			fmt.Println("Err Unmarshal: ", err.Error())
			break
		}

		if err = s.insertTransactionTemps(dao, r.Data.Data); err != nil {
			fmt.Println("insertTransactionTemps - error: ", err.Error())
			break
		}

		if len(r.Data.Data) != limit {
			break
		}
		page++
	}

	// 3. Check and update transaction from transaction-temps
	if err := s.checkAndUpdateTransactionFromTemp(dao); err != nil {
		logger.Error("Error-checkAndUpdateTransactionFromTemp:", logger.LogData{
			Data: bson.M{
				"error": err.Error(),
			},
		})
	}

	// 4. Update campaign statistic
	ctx := context.Background()
	campaignIds := s.distinctCampaignID(ctx, dao)
	go Campaign(externalauth.User{}).UpdateStatisticListCampaign(ctx, campaignIds)

	// 5. Drop collection
	dao.Drop(ctx)

	fmt.Println("End crawl transaction ...")
}

// checkAndUpdateTransactionFromTemp ...
func (s transactionCrawlImplement) checkAndUpdateTransactionFromTemp(dao *mongo.Collection) (err error) {
	var ctx = context.Background()

	// 1. Aggregate group sellerId
	dataGroupSellers, err := s.aggregateTransactionTempGroupBySeller(ctx, dao)
	if err != nil {
		return
	}

	if len(dataGroupSellers) == 0 {
		err = errors.New("no doc aggregateTransactionTempGroupBySeller")
		return
	}

	// 2. worker
	// 2. Chạy cho từng seller // worker new pool cho 10 seller
	var (
		wg     sync.WaitGroup
		done   = make(chan bool)
		ch     = make(chan []error)
		result = make([]error, 0)
	)

	// prepare worker
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		defer wg.Done()
		doc, ok := i.(responsemodel.AggregateTransactionTempGroupSeller)
		if !ok {
			return
		}
		ch <- s.CheckAndUpdateTransactionsBySeller(ctx, doc)
	})
	defer p.Release()

	go func() {
		for response := range ch {
			if len(response) > 0 {
				result = append(result, response...)
			}
		}
		done <- true
	}()

	for _, datum := range dataGroupSellers {
		wg.Add(1)
		p.Invoke(datum)
	}

	wg.Wait()
	close(ch)
	<-done

	if len(result) > 0 {
		logger.Error("Crawl transactions errors:", logger.LogData{
			Data: bson.M{
				"errors": result,
			},
		})
	}

	return
}

// CheckAndUpdateTransactionsBySeller ...
func (s transactionCrawlImplement) CheckAndUpdateTransactionsBySeller(ctx context.Context, doc responsemodel.AggregateTransactionTempGroupSeller) (errors []error) {
	var (
		notifications           = make([]jsmodel.PushNotification, 0)
		isSellerUpdateStatistic bool
	)

	for _, temp := range doc.Temps {
		ns, isUpdateStatistic := s.processCheckAndUpsertTransactionDocTemp(ctx, temp)

		if len(ns) > 0 {
			notifications = append(notifications, ns...)
		}

		if isUpdateStatistic {
			isSellerUpdateStatistic = true
		}
	}

	// Update user statistic
	if isSellerUpdateStatistic {
		s.SellerStatisticService.UpsertStatisticBySellerID(ctx, doc.ID)
	}

	// Send notification
	if len(notifications) > 0 {
		go clientjstream.ClientJestStreamPull{}.PushNotifications(notifications)
	}
	return
}

// processCheckAndUpsertTransactionDocTemp ...
func (s transactionCrawlImplement) processCheckAndUpsertTransactionDocTemp(ctx context.Context, temp responsemodel.ResponseTransactionCrawl) (notifications []jsmodel.PushNotification, isUpdateStatistic bool) {
	// 1. check hash
	isUpdate, transaction := s.checkHashUpdateAndStatusTransaction(ctx, temp)
	if !isUpdate {
		return
	}

	if transaction.ID.IsZero() {
		// flow create new transaction
		notifications, _ = s.flowCreateNewTransaction(ctx, temp)
	} else {
		// flow update transaction
		notifications = s.flowUpdateTransaction(ctx, transaction, temp)
	}

	isUpdateStatistic = true
	return
}

// flowCreateNewTransaction ...
func (s transactionCrawlImplement) flowCreateNewTransaction(ctx context.Context, transactionTemp responsemodel.ResponseTransactionCrawl) (notifications []jsmodel.PushNotification, err error) {
	// 1. Insert transaction
	transactionRaw := s.convertDataCrawl(ctx, transactionTemp)
	if err = s.TransactionService.InsertOne(ctx, transactionRaw); err != nil {
		return
	}

	// 2. Insert transaction histories
	s.TransactionHistoryService.GenerateAndInsertHistoriesByNewTransaction(ctx, transactionRaw)

	// 3. Update seller campaign statistic
	s.SellerCampaignStatisticService.UpsertStatisticBySellerAndCampaignID(ctx, transactionRaw.SellerID, transactionTemp.Campaign)

	// 4. Update status click
	s.ClickService.UpdateStatusByID(ctx, transactionRaw.Click.ClickID, constant.CampaignTransactionClickStatusCompleted)

	// 5. Check cashback
	if transactionRaw.IsValidStatus(constant.TransactionStatus.Cashback.Key) {
		s.flowInsertCashflowWhenTransactionCashback(ctx, transactionRaw)
	}

	// 6. Generate notifications
	notifications = s.getDataPushNotificationByTransaction(ctx, true, transactionRaw)
	return
}

// flowInsertCashflowWhenTransactionCashback ...
func (s transactionCrawlImplement) flowInsertCashflowWhenTransactionCashback(ctx context.Context, transaction mgaffiliate.Transaction) {
	var (
		clientJSStream = clientjstream.ClientJestStreamPull{}
		payload        = jsmodel.PayloadCashflowsBySeller{
			SellerID: transaction.SellerID.Hex(),
			List:     make([]jsmodel.CashflowSeller, 0),
		}
	)

	campaign, _ := s.CampaignService.FindByID(ctx, transaction.CampaignID)

	// Assign cashflows
	cashflow := jsmodel.CashflowSeller{
		Value:    transaction.SellerCommission,
		Action:   constant.CashflowAction.AffiliateTransactionCashback,
		Category: constant.CashflowCategory.Affiliate,
		TargetID: transaction.ID.Hex(),
		Options: &jsmodel.CashFlowOptions{
			AffiliateTransactionCode: transaction.Code,
			AffiliateCampaignID:      transaction.CampaignID.Hex(),
			AffiliateCampaignName:    campaign.Name,
		},
	}
	payload.List = append(payload.List, cashflow)

	// Pull
	clientJSStream.InsertCashflowBySeller(payload)
}

// flowUpdateTransaction ...
func (s transactionCrawlImplement) flowUpdateTransaction(ctx context.Context, transaction mgaffiliate.Transaction, transactionTemp responsemodel.ResponseTransactionCrawl) (notifications []jsmodel.PushNotification) {
	if transaction.IsValidStatus(constant.TransactionStatus.Cashback.Key) {
		return
	}

	if transaction.Code != transactionTemp.TransactionID {
		return
	}

	if transaction.UpdateHash == transactionTemp.UpdatedHash {
		return
	}

	var (
		isChangeStatus bool
		updateData     = bson.M{
			"updateHash": transactionTemp.UpdatedHash,
			"updatedAt":  ptime.Now(),
		}
	)

	if transaction.Status != transactionTemp.Status {
		isChangeStatus = true
		updateData["status"] = transactionTemp.Status
		transaction.Status = transactionTemp.Status
	}

	if transaction.Commission != transactionTemp.Commission {
		newCommission := transactionTemp.Commission
		commissionInfo := calculateTransactionCommission(transaction.SellerCommissionRate, newCommission)
		updateData["commission"] = newCommission
		updateData["sellerCommissionRate"] = commissionInfo.SellerPercent
		updateData["sellerCommission"] = commissionInfo.Seller
		updateData["sellyCommission"] = commissionInfo.Selly
	}

	if transaction.Status == constant.TransactionStatus.Cashback.Key {
		updateData["reconciliationId"] = transactionTemp.ReconciliationID
	}

	// Update status transaction
	s.TransactionService.UpdateByID(ctx, transaction.ID, updateData)

	// Update seller campaign statistic
	s.SellerCampaignStatisticService.UpsertStatisticBySellerAndCampaignID(ctx, transaction.SellerID, transactionTemp.Campaign)

	// Check cashback
	if transaction.IsValidStatus(constant.TransactionStatus.Cashback.Key) {
		s.flowInsertCashflowWhenTransactionCashback(ctx, transaction)
	}

	// History
	s.checkAndInsertHistory(ctx, isChangeStatus, transaction)

	// Generate notification
	if isChangeStatus {
		notifications = s.getDataPushNotificationByTransaction(ctx, false, transaction)
	}
	return
}

// checkAndInsertHistory ...
func (s transactionCrawlImplement) checkAndInsertHistory(ctx context.Context, isChangeStatus bool, transaction mgaffiliate.Transaction) {
	if isChangeStatus {
		s.TransactionHistoryService.GenerateAndInsertHistoriesByUpdateTransaction(ctx, transaction)
	}
}

// insertTransactionTemps ...
func (s transactionCrawlImplement) insertTransactionTemps(dao *mongo.Collection, data []responsemodel.ResponseTransactionCrawl) (err error) {
	if len(data) == 0 {
		return
	}

	var (
		wg     sync.WaitGroup
		done   = make(chan bool)
		ch     = make(chan []mongo.WriteModel)
		ctx    = context.Background()
		result = make([]mongo.WriteModel, 0)
	)

	// prepare worker
	p, _ := ants.NewPoolWithFunc(50, func(i interface{}) {
		defer wg.Done()
		doc, ok := i.(responsemodel.ResponseTransactionCrawl)
		if !ok {
			return
		}
		ch <- s.checkAndConvertTransactionTempByDocCrawl(ctx, doc)
	})
	defer p.Release()

	go func() {
		for response := range ch {
			if len(response) > 0 {
				result = append(result, response...)
			}
		}
		done <- true
	}()

	for _, datum := range data {
		wg.Add(1)
		p.Invoke(datum)
	}

	wg.Wait()
	close(ch)
	<-done

	// Insert bulk write
	if len(result) > 0 {
		_, err = dao.BulkWrite(ctx, result)
	}
	return
}

// checkHashUpdateAndStatusTransaction ...
func (transactionCrawlImplement) checkHashUpdateAndStatusTransaction(ctx context.Context, doc responsemodel.ResponseTransactionCrawl) (result bool, transaction mgaffiliate.Transaction) {
	// 1. check hash transaction
	var transactionSv = transactionImplement{}

	transaction = transactionSv.FindByCode(ctx, doc.TransactionID)

	if !transaction.ID.IsZero() && transaction.UpdateHash == doc.UpdatedHash {
		result = false
		return
	}

	// Transaction status cashback => not crawl
	if transaction.IsValidStatus(constant.TransactionStatus.Cashback.Key) {
		result = false
		return
	}

	result = true
	return
}

// checkAndConvertTransactionTempByDocCrawl ...
func (s transactionCrawlImplement) checkAndConvertTransactionTempByDocCrawl(ctx context.Context, doc responsemodel.ResponseTransactionCrawl) (result []mongo.WriteModel) {
	// 1. check hash
	// Transaction status cashback => not crawl
	if isUpdate, _ := s.checkHashUpdateAndStatusTransaction(ctx, doc); !isUpdate {
		return
	}

	// 2. Find click get seller id
	var clickSv = clickImplement{}
	clickRaw := clickSv.FindByID(ctx, doc.ClickID)
	if clickRaw.ID.IsZero() {
		return
	}

	// 3. Assign seller
	doc.ID = mongodb.NewObjectID()
	doc.SellerID = clickRaw.SellerID

	// 4. Response
	insertModel := mongo.NewInsertOneModel().SetDocument(doc)
	result = append(result, insertModel)
	return
}

//
// PRIVATE METHOD
//

// aggregateTransactionTempGroupBySeller...
func (transactionCrawlImplement) aggregateTransactionTempGroupBySeller(ctx context.Context, dao *mongo.Collection) (result []responsemodel.AggregateTransactionTempGroupSeller, err error) {
	// Aggregate group sellerId
	var (
		pipeline = bson.A{
			bson.D{
				{"$group",
					bson.D{
						{"_id", "$sellerId"},
						{"temps", bson.D{{"$addToSet", "$$ROOT"}}},
					},
				},
			},
		}
	)

	cursor, err := dao.Aggregate(ctx, pipeline)
	if err != nil {
		return
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, &result)
	return
}

func (transactionCrawlImplement) distinctCampaignID(ctx context.Context, dao *mongo.Collection) []primitive.ObjectID {
	var (
		results = make([]primitive.ObjectID, 0)
	)
	data, _ := dao.Distinct(ctx, "campaign", bson.M{})
	for _, d := range data {
		id, ok := d.(primitive.ObjectID)
		if !ok {
			continue
		}
		results = append(results, id)
	}
	return results
}

// convertDataCrawlToRawByList ...
func (s transactionCrawlImplement) convertDataCrawlToRawByList(dataCrawl []responsemodel.ResponseTransactionCrawl) (result []mgaffiliate.Transaction) {
	result = make([]mgaffiliate.Transaction, 0)

	var total = len(dataCrawl)
	if total == 0 {
		return
	}

	var (
		wg  = sync.WaitGroup{}
		ctx = context.Background()
	)

	wg.Add(total)
	result = make([]mgaffiliate.Transaction, total)
	for i, data := range dataCrawl {
		go func(index int, crawl responsemodel.ResponseTransactionCrawl) {
			defer wg.Done()
			result[index] = s.convertDataCrawl(ctx, crawl)
		}(i, data)
	}
	wg.Wait()

	return
}

// convertDataCrawl ...
func (s transactionCrawlImplement) convertDataCrawl(ctx context.Context, data responsemodel.ResponseTransactionCrawl) mgaffiliate.Transaction {
	var (
		wg       = sync.WaitGroup{}
		clickRaw mgaffiliate.Click
		campaign mgaffiliate.Campaign
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		clickRaw = clickImplement{}.FindByID(ctx, data.ClickID)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Find campaign
		var campaignService = campaignImplement{
			CurrentStaff: externalauth.User{},
		}
		campaign, _ = campaignService.FindByID(ctx, data.Campaign)
	}()

	wg.Wait()

	// Commission
	commissionInfo := calculateTransactionCommission(clickRaw.Commission.SellerPercent, data.Commission)
	estimateSellerCommission := clickRaw.Commission.CommissionCalculateEstimateSellerCommission()

	// EstimateCashbackAt
	estimateCashbackAt := campaign.CalculateTransactionEstimateCashbackAt(data.TransactionTime)
	return mgaffiliate.Transaction{
		ID:                       data.ID,
		SellerID:                 data.SellerID,
		SearchString:             mongodb.NonAccentVietnamese(fmt.Sprintf("%s", data.TransactionID)),
		CampaignID:               data.Campaign,
		PlatformID:               clickRaw.PlatformID,
		Code:                     data.TransactionID,
		TransactionTime:          data.TransactionTime,
		Source:                   data.GetSourceTransaction(),
		Commission:               data.Commission,
		SellerCommissionRate:     commissionInfo.SellerPercent,
		SellerCommission:         commissionInfo.Seller,
		SellyCommission:          commissionInfo.Selly,
		EstimateCashbackAt:       estimateCashbackAt,
		EstimateSellerCommission: estimateSellerCommission,
		Status:                   data.Status,
		Click: mgaffiliate.TransactionClick{
			ClickID:      clickRaw.ID,
			AffiliateURL: data.ClickURL,
		},
		Device:     clickRaw.Device,
		UpdateHash: data.UpdatedHash,
		CreatedAt:  ptime.Now(), //data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
		Category:   data.Category,
		From:       clickRaw.From,
	}
}

func calculateTransactionCommission(sellerPercent float64, partnerCommission float64) mgaffiliate.CampaignCommission {
	sellerCommission := math.Round(math.Max(partnerCommission*sellerPercent/float64(100), 0))
	return mgaffiliate.CampaignCommission{
		Real:          partnerCommission,
		SellerPercent: sellerPercent,
		Seller:        sellerCommission,
		Selly:         partnerCommission - sellerCommission,
	}
}
