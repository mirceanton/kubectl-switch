FROM alpine:3.21.1@sha256:b97e2a89d0b9e4011bb88c02ddf01c544b8c781acf1f4d559e7c8f12f1047ac3
USER 8675:8675
COPY kubectl-switch /
ENTRYPOINT ["/kubectl-switch"]