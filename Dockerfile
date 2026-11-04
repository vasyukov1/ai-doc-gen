# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOFLAGS=-p=2 go build -trimpath -o /out/ai-doc-gen ./cmd/ai-doc-gen

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 app \
    && mkdir -p /data/uploads /data/sources \
    && chown -R app:app /data

COPY --from=build /out/ai-doc-gen /usr/local/bin/ai-doc-gen

USER app
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/ai-doc-gen"]
