#!/usr/bin/make -f

export GOPROXY=https://goproxy.io

.PHONY: build-linux
build-linux: go.mod
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 go build -mod=readonly -o go-walletconnect-bridge-linux .

.PHONY: go.mod
go.mod:
	@go mod tidy && go mod verify && go mod download

.PHONY: install
install: go.mod
	@go install -v -mod=readonly .

.PHONY: format
format:
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s

.PHONY: test
test:
	@VERSION=$(VERSION) go test -short -cover -mod=readonly -tags='ledger test_ledger_mock'
