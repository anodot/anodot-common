GOLINT_VERSION:=1.19.1

vendor-update:
	GO111MODULE=on go get -u ./...
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

format:
	gofmt -w ./..

#TODO:vnekhai do not download each time
vet:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $$(go env GOPATH)/bin v$(GOLINT_VERSION)
	$(BUILD_FLAGS) $$(go env GOPATH)/bin/golangci-lint run