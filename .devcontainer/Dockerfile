## ================================================================================================
# Utility versions
## ================================================================================================
ARG KUBECTL_VERSION=1.31.3@sha256:67096b043f84a462b349ebcdcb4258a4f18d1de93ea43cf0df4c908bb877a1ad
ARG K9S_VERSION=v0.32.7@sha256:a967991e635dc8327c25c9d20ea71089a9cad66793013c5675c81a1f6c79292b
ARG TASKFILE_VERSION=v3.38.0@sha256:308c4f5be86bffae3f956cbd7225c4fec69b0e7a89012120b818a10df45d7c59

ARG GOLINT_VERSION=v1.62.2-alpine@sha256:0f3af3929517ed4afa1f1bcba4eae827296017720e08ecd5c68b9f0640bc310d
ARG GORELEASER_VERSION=v2.4.8@sha256:e765952d9ed98e638a4d7755f96b29e5c746c808b54d21f04468e003800af0e3


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