.PHONY: codegen codegen-clean proto proto-lint proto-clean oapi oapi-clean \
        fmt fmt-check vet staticcheck lint check tools \
        docker-login images images-push image-% push-%

REGISTRY ?= ghcr.io
OWNER    ?= egortarasov
REPO     ?= k8s-hw
TAG      ?= dev
PLATFORM ?= linux/amd64
SERVICES := auth-service shop-backend order-worker web
IMAGE    = $(REGISTRY)/$(OWNER)/$(REPO)/$*:$(TAG)

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

# Login to ghcr.io. Expects GITHUB_TOKEN (PAT with write:packages) in env.
docker-login:
	@test -n "$$GITHUB_TOKEN" || (echo "GITHUB_TOKEN is required (PAT with write:packages)"; exit 1)
	@echo "$$GITHUB_TOKEN" | docker login $(REGISTRY) -u $(OWNER) --password-stdin

image-web:
	docker buildx build \
	  --platform $(PLATFORM) \
	  -f Dockerfile.web \
	  -t $(REGISTRY)/$(OWNER)/$(REPO)/web:$(TAG) \
	  -t $(REGISTRY)/$(OWNER)/$(REPO)/web:latest \
	  --load \
	  .

image-%:
	docker buildx build \
	  --platform $(PLATFORM) \
	  --build-arg SERVICE=$* \
	  -t $(IMAGE) \
	  -t $(REGISTRY)/$(OWNER)/$(REPO)/$*:latest \
	  --load \
	  .

push-web:
	docker buildx build \
	  --platform $(PLATFORM) \
	  -f Dockerfile.web \
	  -t $(REGISTRY)/$(OWNER)/$(REPO)/web:$(TAG) \
	  -t $(REGISTRY)/$(OWNER)/$(REPO)/web:latest \
	  --push \
	  .

push-%:
	docker buildx build \
	  --platform $(PLATFORM) \
	  --build-arg SERVICE=$* \
	  -t $(IMAGE) \
	  -t $(REGISTRY)/$(OWNER)/$(REPO)/$*:latest \
	  --push \
	  .

images: $(addprefix image-,$(SERVICES))

images-push: $(addprefix push-,$(SERVICES))
