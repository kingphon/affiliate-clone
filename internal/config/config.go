package config

import (
	"os"

	"git.selly.red/Selly-Server/affiliate/external/constant"

	"git.selly.red/Selly-Modules/mongodb"
)

// ENV ...
type ENV struct {
	Env                            string
	ZookeeperPrefixExternal        string
	ZookeeperPrefixAffiliateCommon string
	ZookeeperPrefixAffiliateApp    string
	ZookeeperPrefixAffiliateAdmin  string
	MongoDB                        struct {
		URI    string
		DBName string

		ReplicaSet          string
		CAPem               string
		CertPem             string
		CertKeyFilePassword string
		ReadPrefMode        string
	}
	Admin struct {
		Server string
		Port   string
	}
	App struct {
		Server string
		Port   string
	}
	Nats       NatsConfig
	SecretKey  string
	MongoAudit MongoConfig `env:",prefix=MONGO_AUDIT_"`

	FileHost string

	HostShareURL string

	Deeplink struct {
		AccessTradeDeepLink string
	}

	// Crawl ...
	Crawl struct {
		URL  string
		Auth string
	}

	// ChecksumKey ...
	ChecksumKey string

	// Auth GG
	EnableAuthenticationService string
}

// NatsConfig ...
type NatsConfig struct {
	URL      string
	Username string
	Password string
	APIKey   string
}

// GetConnectOptions ...
func (dbCfg MongoConfig) GetConnectOptions() mongodb.Config {
	return mongodb.Config{
		Host:       dbCfg.Host,
		DBName:     dbCfg.DBName,
		Standalone: &mongodb.ConnectStandaloneOpts{},
		TLS:        &mongodb.ConnectTLSOpts{},
	}
}

// MongoConfig ...
type MongoConfig struct {
	Host   string
	DBName string
}

var env ENV

// GetENV ...
func GetENV() *ENV {
	return &env
}

// IsEnvDevelop ...
func IsEnvDevelop() bool {
	return env.Env == constant.EnvDevelop
}

// IsEnvStaging ...
func IsEnvStaging() bool {
	return env.Env == constant.EnvStaging
}

// IsEnvProduction ...
func IsEnvProduction() bool {
	return env.Env == constant.EnvProduction
}

// Init ...
func Init() {
	env = ENV{
		Env:                            os.Getenv("ENV"),
		ZookeeperPrefixExternal:        os.Getenv("ZOOKEEPER_PREFIX_EXTERNAL"),
		ZookeeperPrefixAffiliateCommon: os.Getenv("ZOOKEEPER_PREFIX_AFFILIATE_COMMON"),
		ZookeeperPrefixAffiliateApp:    os.Getenv("ZOOKEEPER_PREFIX_AFFILIATE_APP"),
		ZookeeperPrefixAffiliateAdmin:  os.Getenv("ZOOKEEPER_PREFIX_AFFILIATE_ADMIN"),
	}
}

// IsEnableAuthenticationService ...
func IsEnableAuthenticationService() bool {
	return env.EnableAuthenticationService == "true"
}
