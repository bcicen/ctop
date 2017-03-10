FROM quay.io/vektorcloud/glibc:latest

ENV CTOP_VERSION 0.4
ENV CTOP_URL https://github.com/bcicen/ctop/releases/download/v${CTOP_VERSION}/ctop-${CTOP_VERSION}-linux-amd64

RUN wget -q $CTOP_URL -O /ctop && \
    chmod +x /ctop

ENTRYPOINT ["/ctop"]
