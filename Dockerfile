# =================================================================================================
# BUILDER STAGE
# =================================================================================================
FROM golang:1.23.0-alpine@sha256:ac8b5667d9e4b3800c905ccd11b27a0f7dcb2b40b6ad0aca269eab225ed5584e AS builder

ARG PKG=github.com/mirceanton/kube-switcher
ARG VERSION=dev

WORKDIR /build
COPY . .

RUN go build -ldflags "-s -w -X github.com/mirceanton/kube-switcher/cmd.version=${VERSION}" -o kube-switcher


# =================================================================================================
# PRODUCTION STAGE
# =================================================================================================
FROM scratch
USER 8675:8675
COPY --from=builder --chmod=555 /build/kube-switcher /kube-switcher
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/kube-switcher"]
