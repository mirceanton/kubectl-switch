## ================================================================================================
# Utility versions
## ================================================================================================
ARG KUBECTL_VERSION=1.31.0
ARG K9S_VERSION=v0.32.4
ARG TASKFILE_VERSION=v3.38.0

ARG GOLINT_VERSION=v1.60.3-alpine
ARG GORELEASER_VERSION=v2.2.0


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
FROM mcr.microsoft.com/devcontainers/go:1.23-bookworm AS workspace

COPY --from=k9s /bin/k9s /usr/local/bin/k9s

COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/kubectl
RUN kubectl completion bash | sudo tee /etc/bash_completion.d/kubectl.bash > /dev/null

COPY --from=golangci-lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint
RUN golangci-lint completion bash | sudo tee /etc/bash_completion.d/golangci-lint.bash > /dev/null

COPY --from=goreleaser /usr/bin/goreleaser /usr/local/bin/goreleaser
RUN goreleaser completion bash | sudo tee /etc/bash_completion.d/goreleaser.bash > /dev/null

COPY --from=taskfile /task /usr/local/bin/task

USER vscode
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]