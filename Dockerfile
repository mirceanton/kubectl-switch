FROM alpine:3.20.3@sha256:a8f120106f5549715aa966fd7cefaf3b7045f6414fed428684de62fec8c2ca4b
USER 8675:8675
COPY kube-switcher /
ENTRYPOINT ["/kube-switcher"]