FROM quay.io/vektorcloud/glibc:latest

ARG CTOP_VERSION=0.4.1
ENV CTOP_URL https://github.com/bcicen/ctop/releases/download/v${CTOP_VERSION}/ctop-${CTOP_VERSION}-linux-amd64

RUN echo $CTOP_URL && \
    wget -q $CTOP_URL -O /ctop && \
    chmod +x /ctop

ENTRYPOINT ["/ctop"]
