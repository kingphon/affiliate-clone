package zk

import (
	"fmt"
	"os"

	zk "git.selly.red/Selly-Modules/zookeeper"

	"git.selly.red/Selly-Server/affiliate/internal/config"
)

// Connect ...
func Connect() {
	var (
		uri = os.Getenv("ZOOKEEPER_URI")
	)

	// Connect
	if err := zk.Connect(uri); err != nil {
		panic(err)
	}

	envVars := config.GetENV()
	server(envVars)
	commonValues(envVars)
}

// get value in zookeeper
func server(envVars *config.ENV) {
	// Admin
	adminPrefix := getAdminPrefix("")
	envVars.Admin.Server = zk.GetStringValue(fmt.Sprintf("%s/server", adminPrefix))
	envVars.Admin.Port = zk.GetStringValue(fmt.Sprintf("%s/port", adminPrefix))

	// App
	appPrefix := getAppPrefix("")
	envVars.App.Server = zk.GetStringValue(fmt.Sprintf("%s/server", appPrefix))
	envVars.App.Port = zk.GetStringValue(fmt.Sprintf("%s/port", appPrefix))

	// MongoDB
	mongodbPrefix := getExternalPrefix("/mongodb/affiliate")
	envVars.MongoDB.URI = zk.GetStringValue(fmt.Sprintf("%s/uri", mongodbPrefix))
	envVars.MongoDB.DBName = zk.GetStringValue(fmt.Sprintf("%s/db_name", mongodbPrefix))

	envVars.MongoDB.ReplicaSet = zk.GetStringValue(fmt.Sprintf("%s/replica_set", mongodbPrefix))
	envVars.MongoDB.CAPem = zk.GetStringValue(fmt.Sprintf("%s/ca_pem", mongodbPrefix))
	envVars.MongoDB.CertPem = zk.GetStringValue(fmt.Sprintf("%s/cert_pem", mongodbPrefix))
	envVars.MongoDB.CertKeyFilePassword = zk.GetStringValue(fmt.Sprintf("%s/cert_key_file_password", mongodbPrefix))
	envVars.MongoDB.ReadPrefMode = zk.GetStringValue(fmt.Sprintf("%s/read_pref_mode", mongodbPrefix))

	// NATS
	natsPrefix := getExternalPrefix("/nats/affiliate")
	envVars.Nats.URL = zk.GetStringValue(fmt.Sprintf("%s/uri", natsPrefix))
	envVars.Nats.Username = zk.GetStringValue(fmt.Sprintf("%s/user", natsPrefix))
	envVars.Nats.Password = zk.GetStringValue(fmt.Sprintf("%s/password", natsPrefix))
	envVars.Nats.APIKey = zk.GetStringValue(fmt.Sprintf("%s/api_key", natsPrefix))

	// MongoDB_AUDIT
	mongoAuditPrefix := getExternalPrefix("/mongodb/affiliate_audit")
	envVars.MongoAudit.Host = zk.GetStringValue(fmt.Sprintf("%s/host", mongoAuditPrefix))
	envVars.MongoAudit.DBName = zk.GetStringValue(fmt.Sprintf("%s/db_name", mongoAuditPrefix))

	// FileHost
	commonPrefix := getCommonPrefix("")
	envVars.FileHost = zk.GetStringValue(fmt.Sprintf("%s/file_host", commonPrefix))
	envVars.HostShareURL = zk.GetStringValue(fmt.Sprintf("%s/aff_campaign_host_share_url", commonPrefix))

	// Authentication
	authPrefix := getAppPrefix("/authentication")
	envVars.SecretKey = zk.GetStringValue(fmt.Sprintf("%s/auth_secretkey", authPrefix))

	// Deeplink
	deeplinkPrefix := getCommonPrefix("/deeplink")
	envVars.Deeplink.AccessTradeDeepLink = zk.GetStringValue(fmt.Sprintf("%s/access_trade", deeplinkPrefix))

	// Crawl
	crawlPrefix := getAdminPrefix("/crawl")
	envVars.Crawl.URL = zk.GetStringValue(fmt.Sprintf("%s/url", crawlPrefix))
	envVars.Crawl.Auth = zk.GetStringValue(fmt.Sprintf("%s/auth", crawlPrefix))

	// Checksum
	checksumPrefix := getAppPrefix("")
	envVars.ChecksumKey = zk.GetStringValue(fmt.Sprintf("%s/checksum", checksumPrefix))

	// AUTH GG
	envVars.EnableAuthenticationService = zk.GetStringValue(fmt.Sprintf("%s/authentication_google/enable", adminPrefix))
}

func getExternalPrefix(group string) string {
	return fmt.Sprintf("%s%s", config.GetENV().ZookeeperPrefixExternal, group)
}

func commonValues(envVars *config.ENV) {
}

func getCommonPrefix(group string) string {
	return fmt.Sprintf("%s%s", config.GetENV().ZookeeperPrefixAffiliateCommon, group)
}

func getAppPrefix(group string) string {
	return fmt.Sprintf("%s%s", config.GetENV().ZookeeperPrefixAffiliateApp, group)
}

func getAdminPrefix(group string) string {
	return fmt.Sprintf("%s%s", config.GetENV().ZookeeperPrefixAffiliateAdmin, group)
}
