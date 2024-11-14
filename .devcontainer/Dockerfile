## ================================================================================================
# Utility versions
## ================================================================================================
ARG KUBECTL_VERSION=1.31.2@sha256:0eab9ec8f5e0f75271277467ebb7513b36dea0122bc68615d18c47fece4fd82c
ARG K9S_VERSION=v0.32.6@sha256:2d2410af2f8166dd2b67a63d2f172404a49d9f1aec9fc3719c85ccedc8bc660d
ARG TASKFILE_VERSION=v3.38.0@sha256:308c4f5be86bffae3f956cbd7225c4fec69b0e7a89012120b818a10df45d7c59

ARG GOLINT_VERSION=v1.62.0-alpine@sha256:a94c2da655a1f3e66bb7d2c9232324def349f2ba56a3d715d89e83d98f5bd933
ARG GORELEASER_VERSION=v2.4.5@sha256:123935a56bb422aee8c7eafd9c48868d03ffc84bc599eac28251a60ee393ac57


## ================================================================================================
# "Build" stage for utilities with docker images already present
## ================================================================================================
FROM bitnami/kubectl:${KUBECTL_VERSION} AS kubectl
FROM docker.io/derailed/k9s:${K9S_VERSION} AS k9s
FROM ghcr.io/mirceanton/taskfile:${TASKFILE_VERSION} AS taskfile

FROM golangci/golangci-lint:${GOLINT_VERSION} AS golangci-lint
FROM goreleaser/goreleaser:${GORELEASER_VERSION} AS goreleaser


## ================================================================================================
## Main image
## ================================================================================================
FROM mcr.microsoft.com/devcontainers/go:1.23-bookworm@sha256:2e00578e7c526e76ebd51597260928d8542d27c17ffe6390d79200638c51a3e6 AS workspace

# Install additional tools
COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/kubectl
COPY --from=k9s /bin/k9s /usr/local/bin/k9s
COPY --from=taskfile /task /usr/local/bin/task
COPY --from=golangci-lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint
COPY --from=goreleaser /usr/bin/goreleaser /usr/local/bin/goreleaser

# Generate bash completions
RUN kubectl completion bash | sudo tee /etc/bash_completion.d/kubectl.bash > /dev/null
RUN golangci-lint completion bash | sudo tee /etc/bash_completion.d/golangci-lint.bash > /dev/null
RUN goreleaser completion bash | sudo tee /etc/bash_completion.d/goreleaser.bash > /dev/null

USER vscode
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]