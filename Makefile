GO := go

GOLINT_VERSION:=1.19.1

lint: vet
test-all: test

test:
	$(GO) test -v ./...

format:
	gofmt -w ./..

#TODO:vnekhai do not download each time
vet:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $$(go env GOPATH)/bin v$(GOLINT_VERSION)
	$(BUILD_FLAGS) $$($(GO) env GOPATH)/bin/golangci-lint run

vendor-update:
	GO111MODULE=on go get -u ./...
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor