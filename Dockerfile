FROM quay.io/vektorcloud/go:1.8

RUN apk add --no-cache make

COPY glide.* /go/src/github.com/bcicen/ctop/
WORKDIR /go/src/github.com/bcicen/ctop/
RUN glide install

COPY . /go/src/github.com/bcicen/ctop
RUN make build && \
    mkdir -p /go/bin && \
    mv -v ctop /go/bin/

FROM scratch
COPY --from=0 /go/bin/ctop /ctop
ENTRYPOINT ["/ctop"]
