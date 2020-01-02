GO := go
GO_ARGS:=GO111MODULE=on GOFLAGS=-mod=vendor

GOLINT_VERSION:=1.19.1

lint: vet
test-all: lint test

test:
	GOFLAGS=$(GO_ARGS) $(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic -timeout 10s ./pkg/...
format:
	gofmt -w ./pkg

#TODO:vnekhai do not download each time
vet:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $$(go env GOPATH)/bin v$(GOLINT_VERSION)
	$(GO_ARGS) $$($(GO) env GOPATH)/bin/golangci-lint run

vendor-update:
	GO111MODULE=on go get -u ./...
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor