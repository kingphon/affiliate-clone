#!bin/bash

export DOMAIN_AFFILIATE_ADMIN=localhost:4000
export DOMAIN_AFFILIATE_APP=localhost:4001
export ENV=develop
export ZOOKEEPER_URI=127.0.0.1:2181
export ZOOKEEPER_PREFIX_EXTERNAL=/selly
export ZOOKEEPER_PREFIX_AFFILIATE_COMMON=/selly_affiliate/common
export ZOOKEEPER_PREFIX_AFFILIATE_ADMIN=/selly_affiliate/admin
export ZOOKEEPER_PREFIX_AFFILIATE_APP=/selly_affiliate/app


# make update-submodules branch=develop
update-submodules:
	git submodule update --init --recursive && \
	git submodule foreach git checkout $(branch) && \
	git submodule foreach git pull origin $(branch)

run-admin:
	go run cmd/admin/main.go


run-app:
	go run cmd/app/main.go

swagger-admin:
	swag init -d ./ -g cmd/admin/main.go \
    --exclude ./pkg/app \
    -o ./docs/admin --pd

swagger-app:
	swag init -d ./ -g cmd/app/main.go \
	--exclude ./pkg/admin \
	-o ./docs/app
# delete submodules folder in git cache
# git rm --cached submodules