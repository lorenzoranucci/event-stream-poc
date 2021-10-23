APP_VERSION ?= dev
KUBEVAL_ARCH ?= darwin-amd64

_GOLANG_IMAGE = golang:1.17.2
_PROJECT_NAME = event-stream-golang
_PROJECT_DIR = $(patsubst %/,%,$(dir $(realpath $(firstword $(MAKEFILE_LIST)))))
_PROJECT_PREFIX = github.com/ProntoPro

ifneq ("$(wildcard ${_PROJECT_DIR}/.env)", "")
	_ENV_FILE=.env
else
	_ENV_FILE=.env.dist
endif

# Helpers

# $1: image
# $2: command
define run-docker
	@docker run --rm --network local_default \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e CGO_ENABLED=0 \
		-e GOPRIVATE=${_PROJECT_PREFIX}/gollum \
		--env-file=${_PROJECT_DIR}/${_ENV_FILE} \
		-w /srv \
		-v ${_PROJECT_DIR}/:/srv \
		-v ${_PROJECT_DIR}/../.netrc:/root/.netrc \
		$1 $2
endef

# Env

.PHONY: env env-test env-ci clear-env

env:
	@grep -h -s -v '^#' ${_ENV_FILE} | sed 's/^/export /' | sed 's/=/="/' | sed 's/$$/"/'
	@echo '# run this command to configure your shell:'
	@echo '# eval "$$(make env)"'

env-test:
	@grep -h -s -v '^#' .env.test | sed 's/^/export /' | sed 's/=/="/' | sed 's/$$/"/'
	@echo 'export CGO_ENABLED=0'
	@echo 'export GOMAXPROCS=1'
	@echo '# run this command to configure your shell:'
	@echo '# eval "$$(make env-test)"'

env-ci:
	@grep -h -s -v '^#' .env.ci | sed 's/^/export /' | sed 's/=/="/' | sed 's/$$/"/'
	@echo 'export CGO_ENABLED=0'
	@echo 'export GOARCH=amd64'
	@echo 'export GOMAXPROCS=1'
	@echo 'export GOMODCACHE="${_PROJECT_DIR}/.ci-gomodcache"'
	@echo 'export GOOS=linux'
	@echo '# run this command to configure your shell:'
	@echo '# eval "$$(make env-ci)"'

clear-env:
	@grep -h -s -v '^#' .env.dist .env .env.ci .env.test | sed 's/\=.*//g' | sort | uniq | sed 's/^/unset /'

# Vendors

.PHONY: mod-vendor mod-tidy mod-check-upgrades mod-upgrade

mod-vendor:
	@go mod vendor

mod-tidy:
	@go mod tidy

mod-check-upgrades:
	@go list -mod=readonly -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all

mod-upgrade:
	@go get -u ./... && make mod-tidy

# Code style

.PHONY: fmt imports cs lint

fmt:
	@find . -name *.go -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} bash -c 'echo "formatting {}.." && gofmt -w -s {}'

imports:
	@find . -name *.go -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} bash -c 'echo "fixing imports for {}.." && goimports --local "${_PROJECT_PREFIX}" -v -w {}'

cs-fix: fmt imports

lint:
	@golangci-lint run \
		--deadline=5m0s \
		--out-format=line-number \
		-exclude-use-default=false \
		./...

deps:
	@go-cleanarch -domain domain -application application -infrastructure infrastructure

# Testing

.PHONY: test-unit test-unit-with-race-check test-integration test-integration-with-race-check

test-unit:
	@go test -v -covermode=count -coverprofile=build/coverage-unit.out -tags=unit -count 1 `go list ./... | grep -v /vendor/`

test-unit-with-race-check:
	@CGO_ENABLED=1 \
	go test -v -covermode=atomic -coverprofile=build/coverage-unit-race-check.out -race -tags=unit -count 1 `go list ./... | grep -v /vendor/`

test-integration:
	@go test -v -covermode=count -coverprofile=build/coverage-integration.out -tags=integration -p 1 `go list ./... | grep -v /vendor/`

test-integration-with-race-check:
	@CGO_ENABLED=1 \
	go test -v -covermode=atomic -coverprofile=build/coverage-integration-race-check.out -race -tags=integration -p 1 `go list ./... | grep -v /vendor/`

.PHONY: test-all test-all-with-race-check test-all-fail-fast test-all-fail-fast-with-race-check

test-all:
	@go test -v -covermode=count -coverprofile=build/coverage-all.out -p 1 `go list ./... | grep -v /vendor/`

test-all-with-race-check:
	@CGO_ENABLED=1 \
	go test -v -covermode=count -coverprofile=build/coverage-all-race-check.out -race -p 1 `go list ./... | grep -v /vendor/`

test-all-fail-fast:
	@go test -v -failfast -covermode=count -coverprofile=build/coverage-all-fail-fast.out -p 1 `go list ./... | grep -v /vendor/`

test-all-fail-fast-with-race-check:
	@CGO_ENABLED=1 \
	go test -v -failfast -covermode=atomic -coverprofile=build/coverage-all-fail-fast-race-check.out -race -p 1 `go list ./... | grep -v /vendor/`

# Build

.PHONY: build build-all clean docker-image

build:
	@go build -a -ldflags "-X main.version=${APP_VERSION}" -o bin/${_PROJECT_NAME} cmd/${_PROJECT_NAME}/main.go

build-all: clean mod-vendor build

clean:
	@rm -rf bin/*

docker-image: build-docker
	@docker build -f Dockerfile -t prontopro/${_PROJECT_NAME}:latest .

# Run

.PHONY: run-server run-consumer-update-tags run-consumer-update run-consumer-delete
.PHONY: run-consumer-sorting-scores run-consumer-sorting-default-scores run-consumer-sorting-geo-clusters

run-system:
	@docker-compose down && docker-compose up -d

post-review:
	@curl -H "Content-Type: application/json" \
     -X POST http://localhost:8080/reviews \
     -d '{"comment": "This is a comment", "rating": 10 }'

log-consumer:
	@docker logs -f pp_event_stream_golang_review_created_consumer

log-server:
	@docker logs -f pp_event_stream_golang_review_created_consumer

# Utils

%-docker:
	$(call run-docker,${_GOLANG_IMAGE},make $*)
