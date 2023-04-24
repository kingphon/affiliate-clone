package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// index ...
func index() {
	// transaction ...
	indexTransaction()

	// transaction history ...
	indexTransactionHistory()

	// seller statistic ...
	indexSellerStatistic()

	// seller campaign statistic ...
	indexSellerCampaignStatistic()

	// click ...
	indexClick()

	// share url ...
	indexShareURL()

	// platform ...
	indexPlatform()
}

// indexTransaction ...
func indexTransaction() {
	indexes := []mongo.IndexModel{
		newIndex("sellerId"),
		newIndex("code"),
		newIndex("reconciliationId"),
	}
	process(TransactionCol(), indexes)
}

// indexTransactionHistory ...
func indexTransactionHistory() {
	indexes := []mongo.IndexModel{
		newIndex("transactionId"),
	}
	process(TransactionHistoryCol(), indexes)
}

// indexSellerStatistic ...
func indexSellerStatistic() {
	indexes := []mongo.IndexModel{
		newIndex("sellerId"),
	}
	process(SellerStatisticCol(), indexes)
}

// indexSellerCampaignStatistic ...
func indexSellerCampaignStatistic() {
	indexes := []mongo.IndexModel{
		newIndex("sellerId"),
	}
	process(SellerCampaignStatisticCol(), indexes)
}

// indexClick ...
func indexClick() {
	indexes := []mongo.IndexModel{
		newIndex("sellerId"),
		newIndex("campaignId"),
		newIndex("platformId"),
	}
	process(ClickCol(), indexes)
}

// indexShareURL ...
func indexShareURL() {
	indexes := []mongo.IndexModel{
		newIndex("sellerId"),
		newUniqIndex("code"),
	}
	process(ShareURLCol(), indexes)
}

// indexPlatform ...
func indexPlatform() {
	indexes := []mongo.IndexModel{
		newUniqIndex("code"),
	}
	process(PlatformCol(), indexes)
}

//
// Method private
//

// process ...
func process(col *mongo.Collection, indexes []mongo.IndexModel) {
	opts := options.CreateIndexes().SetMaxTime(time.Minute * 30)
	_, err := col.Indexes().CreateMany(context.Background(), indexes, opts)
	if err != nil {
		fmt.Printf("Index collection %s err: %v", col.Name(), err)
	}
}

// newIndex ...
func newIndex(key ...string) mongo.IndexModel {
	var doc bsonx.Doc
	for _, s := range key {
		e := bsonx.Elem{
			Key:   s,
			Value: bsonx.Int32(1),
		}
		if strings.HasPrefix(s, "-") {
			e = bsonx.Elem{
				Key:   strings.Replace(s, "-", "", 1),
				Value: bsonx.Int32(-1),
			}
		}
		doc = append(doc, e)
	}

	return mongo.IndexModel{Keys: doc}
}

// newUniqIndex ...
func newUniqIndex(key ...string) mongo.IndexModel {
	var doc bsonx.Doc
	for _, s := range key {
		e := bsonx.Elem{
			Key:   s,
			Value: bsonx.Int32(1),
		}
		if strings.HasPrefix(s, "-") {
			e = bsonx.Elem{
				Key:   strings.Replace(s, "-", "", 1),
				Value: bsonx.Int32(-1),
			}
		}
		doc = append(doc, e)
	}
	opt := options.Index().SetUnique(true)
	return mongo.IndexModel{Keys: doc, Options: opt}
}
