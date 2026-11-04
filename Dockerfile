# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOFLAGS=-p=2 go build -trimpath -o /out/ai-doc-gen ./cmd/ai-doc-gen

FROM alpine:3.22 AS typst
ARG TARGETARCH
RUN apk add --no-cache curl xz \
    && case "$TARGETARCH" in \
        arm64) TYPST_ARCH=aarch64; TYPST_SHA=5aa8d74a3d906e60ea12a66ac2f37f8eef1b14cbad7182a745e393a10c23dcee ;; \
        amd64) TYPST_ARCH=x86_64; TYPST_SHA=a6d077d0a95eed5a2eba715b2dae06be954f624ccbf85758a03f389ded33118c ;; \
        *) exit 1 ;; \
       esac \
    && curl -fsSL "https://github.com/typst/typst/releases/download/v0.15.1/typst-${TYPST_ARCH}-unknown-linux-musl.tar.xz" -o /tmp/typst.tar.xz \
    && echo "$TYPST_SHA  /tmp/typst.tar.xz" | sha256sum -c - \
    && tar -xJf /tmp/typst.tar.xz -C /tmp \
    && cp /tmp/typst-${TYPST_ARCH}-unknown-linux-musl/typst /usr/local/bin/typst

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 app \
    && mkdir -p /data/uploads /data/sources \
    && chown -R app:app /data

COPY --from=build /out/ai-doc-gen /usr/local/bin/ai-doc-gen
COPY --from=typst /usr/local/bin/typst /usr/local/bin/typst

USER app
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/ai-doc-gen"]
