ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.20

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

# Copy only what the Go build needs so unrelated changes (web/, nginx.conf,
# README, ...) do not bust the build cache.
COPY cmd ./cmd
COPY internal ./internal
COPY api ./api

ARG SERVICE
RUN test -n "${SERVICE}" || (echo "SERVICE build arg is required" && exit 1)
ENV GOFLAGS="-p=1" GOMEMLIMIT=512MiB GOGC=50
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/${SERVICE}

FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache ca-certificates && adduser -D -H app
USER app

COPY --from=builder /out/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]
