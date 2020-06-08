#!/usr/bin/make -f

export GOPROXY=https://goproxy.io,direct
export GO111MODULE=on

.PHONY: build-linux go.mod install docker format test

build-linux:
	LEDGER_ENABLED=false CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o go-walletconnect-bridge-linux .

go.mod:
	@go mod tidy && go mod verify && go mod download

install:
	@go install -v .

docker: build-linux
	@docker rmi -f zhcppy/go-wallet-bridge:latest
	@docker build --no-cache -f Dockerfile -t zhcppy/go-wallet-bridge:latest .

format:
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s

test:
	@go test -short -cover
