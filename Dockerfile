# =================================================================================================
# BUILDER STAGE
# =================================================================================================
FROM golang:1.22.5-alpine AS builder

ARG PKG=github.com/mirceanton/kube-switcher
ARG VERSION=dev
ARG REVISION=dev

WORKDIR /build
COPY . .

RUN go build -ldflags "-s -w -X main.Version=${VERSION} -X main.Gitsha=${REVISION}" -o kube-switcher


# =================================================================================================
# PRODUCTION STAGE
# =================================================================================================
FROM scratch
USER 8675:8675
COPY --from=builder --chmod=555 /build/kube-switcher /kube-switcher
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/kube-switcher"]