FROM gcr.io/distroless/static-debian12:nonroot
USER 8675:8675
COPY kubectl-switch /
ENTRYPOINT ["/kubectl-switch"]