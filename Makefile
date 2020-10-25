VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/hayashiki/tarsier.revision=$(CURRENT_REVISION)"
GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE ?= on

.PHONY: test
test:
	go test -v ./...

.PHONY: show-version
show-version: $(GOBIN)/gobump
	@gobump show -r .

$(GOBIN)/gobump:
	@cd && go get github.com/x-motemen/gobump/cmd/gobump

.PHONY: tag
tag:
	git tag -a "v$(VERSION)" -m "Release $(VERSION)"
	git push --tags

.PHONY: build
build:
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/tarsier

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: vet
vet:
	go vet ./...
