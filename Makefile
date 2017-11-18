NAME=ctop
VERSION=$(shell cat VERSION)
BUILD=$(shell git rev-parse --short HEAD)
EXT_LD_FLAGS="-Wl,--allow-multiple-definition"
LD_FLAGS="-w -X main.version=$(VERSION) -X main.build=$(BUILD) -extldflags=$(EXT_LD_FLAGS)"

clean:
	rm -rf _build/ release/

build:
	dep ensure
	CGO_ENABLED=0 go build -tags release -ldflags $(LD_FLAGS) -o ctop

build-dev:
	go build -ldflags "-w -X main.version=$(VERSION)-dev -X main.build=$(BUILD) -extldflags=$(EXT_LD_FLAGS)"

build-all:
	mkdir -p _build
	GOOS=darwin GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-darwin-amd64
	GOOS=linux  GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-linux-amd64
	GOOS=linux  GOARCH=arm   go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-linux-arm
	GOOS=linux  GOARCH=arm64 go build -tags release -ldflags $(LD_FLAGS) -o _build/ctop-$(VERSION)-linux-arm64
	cd _build; sha256sum * > sha256sums.txt

image:
	docker build -t ctop -f Dockerfile .

release:
	mkdir release
	go get github.com/progrium/gh-release/...
	cp _build/* release
	gh-release create bcicen/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

.PHONY: build
