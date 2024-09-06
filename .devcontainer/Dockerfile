## ================================================================================================
# Utility versions
## ================================================================================================
ARG KUBECTL_VERSION=1.31.0@sha256:7779e585c588adde5c8feafd584c1ffc76ffbe8475236aa14c13b8afe18d9af6
ARG K9S_VERSION=v0.32.4@sha256:32e0cf06b70f1b7e7576b64b378170ddda194b491ef4d04b7303f1b8ab81a771
ARG TASKFILE_VERSION=v3.38.0@sha256:308c4f5be86bffae3f956cbd7225c4fec69b0e7a89012120b818a10df45d7c59

ARG GOLINT_VERSION=v1.60.3-alpine@sha256:8326b535f6aff1995059a0025449e8fe150b3c8e48d7c310855302c638cb2ccf
ARG GORELEASER_VERSION=v2.2.0@sha256:847eba313710996f4eeeef1ac21a4fd5ea7527e3a500633ae78f40d9baefddac


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
FROM mcr.microsoft.com/devcontainers/go:1.23-bookworm@sha256:a3403234ec01e5383f63414aabfd3e2703a4c1e44b260354495c6f1e8cdc7594 AS workspace

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