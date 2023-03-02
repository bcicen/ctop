NAME=ctop
VERSION=$(shell cat VERSION)
BUILD=$(shell git rev-parse --short HEAD)
LD_FLAGS="-w -X main.version=$(VERSION) -X main.build=$(BUILD)"

clean:
	rm -rf _build/ release/

build:
	go mod download
	CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o ctop

build-all:
	mkdir -p _build
	GOOS=darwin  GOARCH=amd64   CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-darwin-amd64
	GOOS=darwin  GOARCH=arm64   CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-darwin-arm64
	GOOS=linux   GOARCH=amd64   CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-linux-amd64
	GOOS=linux   GOARCH=arm     CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-linux-arm
	GOOS=linux   GOARCH=arm64   CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-linux-arm64
	GOOS=linux   GOARCH=ppc64le CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-linux-ppc64le
	GOOS=windows GOARCH=amd64   CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-windows-amd64
	cd _build; sha256sum * > sha256sums.txt

run-dev:
	rm -f ctop.sock ctop
	go build -ldflags $(LD_FLAGS) -o ctop
	CTOP_DEBUG=1 ./ctop

image:
	docker build -t ctop -f Dockerfile .

release:
	mkdir release
	cp _build/* release
	cd release; sha256sum --quiet --check sha256sums.txt && \
	gh release create v$(VERSION) -d -t v$(VERSION) *

.PHONY: build
