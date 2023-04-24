- External data

```shell
create /selly
```

```shell
create /selly_affiliate/admin
create /selly_affiliate/admin/server affiliate_admin
create /selly_affiliate/admin/port :4000
```

```shell
create /selly_affiliate/app
create /selly_affiliate/app/server affiliate_app
create /selly_affiliate/app/port :4001
```

```shell
create /selly/mongodb
create /selly/mongodb/affiliate
create /selly/mongodb/affiliate/uri "mongodb://localhost:27017"
create /selly/mongodb/affiliate/db_name "affiliate"
```

------------------------

```shell
create /selly/nats
create /selly/nats/affiliate
create /selly/nats/affiliate/uri "localhost:4222"
create /selly/nats/affiliate/api_key "selly"
create /selly/nats/affiliate/user ""
create /selly/nats/affiliate/password ""
create /selly/mongodb/affiliate/replica_set 
create /selly/mongodb/affiliate/ca_pem 
create /selly/mongodb/affiliate/cert_pem 
create /selly/mongodb/affiliate/cert_key_file_password 
create /selly/mongodb/affiliate/read_pref_mode 
```

------------------------

```shell
create /selly/mongodb
create /selly/mongodb/affiliate_audit
create /selly/mongodb/affiliate_audit/host "mongodb://localhost:27017"
create /selly/mongodb/affiliate_audit/db_name "audit-affiliate"

``` 

------------------------

- Location data

```shell
create /selly_affiliate
create /selly_affiliate/common
create /selly_affiliate/common/file_host https://media.selly.vn/unibag-develop
create /selly_affiliate/common/aff_campaign_host_share_url https://unibag.xyz

create /selly_affiliate/common/deeplink
create /selly_affiliate/common/deeplink/access_trade https://pub2.accesstrade.vn/tool/deep_link 

```

```shell
create /selly_affiliate
create /selly_affiliate/app
create /selly_affiliate/app/authentication
create /selly_affiliate/app/authentication/auth_secretkey "authsecretkey"
create /selly_affiliate/app/checksum checksumsecretkey
```

```shell
create /selly_affiliate/admin/crawl
create /selly_affiliate/admin/crawl/url 54.251.228.51:3050
create /selly_affiliate/admin/crawl/auth Y2FzaGJhZ2dvZ2xvYmFsMjAyMA==

```

```shell
create /selly_affiliate/admin/authentication_google
create /selly_affiliate/admin/authentication_google/enable "true"


```