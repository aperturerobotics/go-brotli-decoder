SHELL:=bash
GOIMPORTS=hack/bin/goimports
GOFUMPT=hack/bin/gofumpt
GOLANGCI_LINT=hack/bin/golangci-lint
GOLIST=go list -f "{{ .Dir }}" -m

export GO111MODULE=on
undefine GOOS
undefine GOARCH

all:

vendor:
	go mod vendor

$(GOIMPORTS):
	cd ./hack; \
	go build -v \
		-o ./bin/goimports \
		golang.org/x/tools/cmd/goimports

$(GOFUMPT):
	cd ./hack; \
	go build -v \
		-o ./bin/gofumpt \
		mvdan.cc/gofumpt

$(GOLANGCI_LINT):
	cd ./hack; \
	go build -v \
		-o ./bin/golangci-lint \
		github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --timeout=10m

.PHONY: fix
fix: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fix --timeout=10m

.PHONY: format
format: $(GOFUMPT) $(GOIMPORTS)
	$(GOIMPORTS) -w ./
	$(GOFUMPT) -w ./

.PHONY: test
test:
	go test -v ./...
