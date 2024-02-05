include .test.env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

up: ### Run docker-compose
	docker compose up --build -d --abort-on-container-exit && docker compose logs -f
.PHONY: up

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

integration-test: ### run functional tests
	go clean -testcache && go test -v ./testing/...
.PHONY: integration-test

bin-deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/pressly/goose/v3/cmd/goose@v3.18.0
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest