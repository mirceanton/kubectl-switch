FROM alpine:3.21.0@sha256:eb37f58646a901dc7727cf448cae36daaefaba79de33b5058dab79aa4c04aefb
USER 8675:8675
COPY kubectl-switch /
ENTRYPOINT ["/kubectl-switch"]