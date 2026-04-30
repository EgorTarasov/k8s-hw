.PHONY: codegen codegen-clean proto proto-lint proto-clean oapi oapi-clean \
        fmt fmt-check vet staticcheck lint check tools

GO_PKGS := $(shell go list ./... | grep -v '/internal/generated/')
GO_FILES := $(shell find . -type f -name '*.go' -not -path './internal/generated/*' -not -path './.git/*')

codegen: proto oapi

codegen-clean: proto-clean oapi-clean

proto:
	buf generate

proto-lint:
	buf lint

proto-clean:
	rm -rf internal/generated/auth

oapi:
	mkdir -p internal/generated/api
	oapi-codegen -config api/cfg/server.yaml api/schema.yaml

oapi-clean:
	rm -rf internal/generated/api


tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest

fmt:
	gofmt -s -w $(GO_FILES)

fmt-check:
	@diff=$$(gofmt -s -l $(GO_FILES)); \
	if [ -n "$$diff" ]; then \
	  echo "Files need gofmt:"; echo "$$diff"; exit 1; \
	fi

vet:
	go vet $(GO_PKGS)

staticcheck:
	staticcheck $(GO_PKGS)

lint: fmt-check vet staticcheck

check: lint
	go build ./...
