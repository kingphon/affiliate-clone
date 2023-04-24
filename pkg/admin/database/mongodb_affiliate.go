package database

import (
	"fmt"

	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/affiliate/internal/config"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
)

var affiliateDB *mongo.Database

// ConnectMongoDBAffiliate ...
func ConnectMongoDBAffiliate() {
	var (
		cfg = config.GetENV().MongoDB
		err error
		tls *mongodb.ConnectTLSOpts
	)

	if cfg.ReplicaSet != "" {
		tls = &mongodb.ConnectTLSOpts{
			ReplSet:             cfg.ReplicaSet,
			CaFile:              cfg.CAPem,
			CertKeyFile:         cfg.CertPem,
			CertKeyFilePassword: cfg.CertKeyFilePassword,
			ReadPreferenceMode:  cfg.ReadPrefMode,
		}
	}

	// Connect
	affiliateDB, err = mongodb.Connect(mongodb.Config{
		Host:       cfg.URI,
		DBName:     cfg.DBName,
		Monitor:    apmmongo.CommandMonitor(),
		TLS:        tls,
		Standalone: &mongodb.ConnectStandaloneOpts{},
	})
	if err != nil {
		panic(err)
	}

	// Index ...
	index()
}

// GetMongoDBAffiliate ...
func GetMongoDBAffiliate() *mongo.Database {
	return affiliateDB
}

// CampaignCol ...
func CampaignCol() *mongo.Collection {
	return affiliateDB.Collection(colAffiliateCampaign)
}

// PlatformCol ...
func PlatformCol() *mongo.Collection {
	return affiliateDB.Collection(colAffiliatePlatform)
}

// TransactionCol ...
func TransactionCol() *mongo.Collection {
	return affiliateDB.Collection(colAffiliateTransaction)
}

// TransactionTempByNameCol ...
func TransactionTempByNameCol(name string) *mongo.Collection {
	return affiliateDB.Collection(fmt.Sprintf("transaction-temps-%s", name))
}

// ClickCol ...
func ClickCol() *mongo.Collection {
	return affiliateDB.Collection(colAffiliateClick)
}

// TransactionHistoryCol ...
func TransactionHistoryCol() *mongo.Collection {
	return affiliateDB.Collection(colAffiliateTransactionHistory)
}

// SellerCampaignStatisticCol ...
func SellerCampaignStatisticCol() *mongo.Collection {
	return affiliateDB.Collection(colAffiliateSellerCampaignStatistic)
}

// SellerStatisticCol ...
func SellerStatisticCol() *mongo.Collection {
	return affiliateDB.Collection(colAffiliateSellerStatistic)
}

// ShareURLCol ...
func ShareURLCol() *mongo.Collection {
	return affiliateDB.Collection(colAffiliateSellerShareURL)
}

// ReconciliationCol ...
func ReconciliationCol() *mongo.Collection {
	return affiliateDB.Collection(colAffiliateReconciliation)
}

// ReconcilicationTempByNameCol ...
func ReconcilicationTempByNameCol(name string) *mongo.Collection {
	return affiliateDB.Collection(fmt.Sprintf("reconcilication-temps-%s", name))
}
