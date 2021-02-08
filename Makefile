BIN := "./bin/bannerrotation"

generate:
	go generate ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test -v -count=100 -race -gcflags=-l -timeout=5m ./internal/...

test-integration:
	ENV_FILE=.env go test -v -count=1 ./test/...

run:
	docker-compose up --build

stop:
	docker-compose stop

build:
	go build -v -o $(BIN) ./cmd

migrate:
	goose -dir=migrations postgres "$(DB_DSN)" up

.PHONY: build test