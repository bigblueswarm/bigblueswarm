.DEFAULT_GOAL := help
SHELL := /bin/bash

VERSION = ""
SECRET = $(shell docker exec bbb1 sh -c "bbb-conf --secret" | grep -Po "Secret: (.*)" | cut -d: -f2 | xargs)

#help: @ list available tasks on this project
help:
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#'  | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#init: @ install project and init dependencies
init:
	@echo "[INIT] Install project and init dependencies"
	@echo "[INIT][1/3] install and setup pre-commit"
	curl https://pre-commit.com/install-local.py | python -
	source ~/.profile
	pip install pre-commit
	pre-commit --version
	pre-commit install
	@echo "[INIT][2/3] commitlint, conventional commit, husky and newman installation"
	npm install --save-dev @commitlint/{config-conventional,cli} husky, newman
	npx husky install
	npx husky add .husky/commit-msg "npx --no -- commitlint --edit \"$1\""
	@echo "[INIT][3/3] administration script dependencies installation"
	sudo apt-get install tidy jq

#scripts: @ download scripts
scripts:
	@echo "[SCRIPTS] install bigblueswarm scripts"
	git clone https://github.com/bigblueswarm/bbs-scripts scripts

#start: @ launch bigblueswarm using the default config file
start:
	go run cmd/bigblueswarm/main.go -config config.yml

#test.unit: @ run unit tests and coverage
test.unit:
	@echo "[TEST.UNIT] run unit tests and coverage"
	go test -covermode=atomic -coverprofile=coverage.out \
		github.com/bigblueswarm/bigblueswarm/v2/pkg/admin \
		github.com/bigblueswarm/bigblueswarm/v2/pkg/api \
		github.com/bigblueswarm/bigblueswarm/v2/pkg/app \
		github.com/bigblueswarm/bigblueswarm/v2/pkg/config \
		github.com/bigblueswarm/bigblueswarm/v2/pkg/utils \
		github.com/bigblueswarm/bigblueswarm/v2/pkg/balancer \
		github.com/bigblueswarm/bigblueswarm/v2/pkg/restclient

#test.integration: @ run integration tests
test.integration: build test.integration.cluster.start test.integration.bigblueswarm.run test.integration.launch test.integration.bigblueswarm.stop cluster.stop test.integration.cluster.remove

test.integration.cluster.start:
	@if [ "$(shell docker images | grep sledunois/bbb-dev | wc -l)" -eq "0" ]; then\
		@echo "BigBlueButton image not found, building it...";\
		make build.image;\
	fi
	@make cluster.start
	@sleep 5m
	@make cluster.init

test.integration.bigblueswarm.run:
	@echo "[RUN] start bigblueswarm"
	@nohup ./bin/bigblueswarm-${VERSION} --config config.yml &
	@sleep 15s

test.integration.bigblueswarm.stop:
	ps -ef | grep bigblueswarm-${VERSION} | grep -v grep | awk '{print $$2}' | xargs kill

test.integration.cluster.remove:
	docker rm -f bbb1 bbb2 redis influxdb

test.integration.launch:
	npm install newman
	node_modules/.bin/newman run ./test/bigblueswarm.postman_collection.json -e ./test/Integration\ test.postman_environment.json --env-var instance_secret="${SECRET}" --bail --verbose --ignore-redirects


#build.image: @ build custom bigbluebutton docker image
build.image:
	@make -f ./scripts/Makefile build.image

#build: @ build bigblueswarm binary
build:
	@echo "[BUILD] build bigblueswarm ${VERSION} binary"
	rm -rf bin
	go build -ldflags="-X 'main.version=${VERSION}' -X 'main.buildTime=$(shell date)' -X 'main.commitHash=$(shell git rev-parse HEAD)'" -o ./bin/bigblueswarm-${VERSION} ./cmd/bigblueswarm/main.go

#build.docker @ build bigblueswarm docker image
build.docker:
	@echo "[BUILD DOCKER] build bigblueswarm docker image"
	docker build . -t sledunois/bigblueswarm:${VERSION}

#cluster.init: @ initialize development cluster (initialize influxdb and telegraf)
cluster.init: cluster.influxdb cluster.telegraf

#cluster.start: @ start development cluster
cluster.start:
	@make -f ./scripts/Makefile cluster.start

#cluster.stop: @ stop development cluster
cluster.stop:
	@make -f ./scripts/Makefile cluster.stop

#cluster.influxdb: @ initialize influxdb database
cluster.influxdb:
	@make -f ./scripts/Makefile cluster.influxdb

#cluster.telegraf: @ initialize bigbluebutton telegraf configuration
cluster.telegraf:
	@make -f ./scripts/Makefile cluster.telegraf

#cluster.grafana: @ launch cluster with grafana
cluster.grafana:
	@make -f ./scripts/Makefile cluster.grafana

#cluster.consul: @ start development cluster using consul coniguration provider
cluster.consul:
	@make -f ./scripts/Makefile cluster.consul
