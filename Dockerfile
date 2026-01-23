FROM cgr.dev/chainguard/go:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev-build
ARG BUILD=none

RUN VERSION=${VERSION} BUILD=${BUILD} make build && \
    mkdir -p /go/bin && \
    mv -v ctop /go/bin/

FROM scratch
ENV TERM=linux
COPY --from=builder /go/bin/ctop /ctop
ENTRYPOINT ["/ctop"]
