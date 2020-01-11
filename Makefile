#!/usr/bin/make -f

GOPROXY=GO111MODULE=on GOPROXY=https://goproxy.io

.PHONY: build-linux
build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 go build -mod=readonly -o go-walletconnect-bridge-linux .

.PHONY: go.sum
go.sum: go.mod
	@$(GOPROXY) go mod tidy
	@$(GOPROXY) go mod verify
	@$(GOPROXY) go mod download

.PHONY: install
install: go.sum
	@go install -v -mod=readonly .

.PHONY: format
format:
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s

.PHONY: test
test:
	@VERSION=$(VERSION) go test -short -cover -mod=readonly -tags='ledger test_ledger_mock'
