FROM quay.io/vektorcloud/go:1.10

RUN apk add --no-cache make

COPY Gopkg.* /go/src/github.com/bcicen/ctop/
WORKDIR /go/src/github.com/bcicen/ctop/
RUN dep ensure -vendor-only

COPY . /go/src/github.com/bcicen/ctop
RUN make build && \
    mkdir -p /go/bin && \
    mv -v ctop /go/bin/

FROM scratch
ENV TERM=linux
COPY --from=0 /go/bin/ctop /ctop
ENTRYPOINT ["/ctop"]
