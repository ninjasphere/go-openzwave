GO_OPENZWAVE=github.com/ninjasphere/go-openzwave

all: build

# a rule that does the build in place here with whatever happens to be adjacent
here: deps fmt
	go install

clean: clean-src
	go clean -i
	rm -rf bin/* || true
	rm -rf .gopath || true

clean-src:
	find . -name '*~' -exec rm {} \;

fmt:
	gofmt -s -w *.go

# does a clean build of all dependencies from git hub
build:
	scripts/build.sh

#
# A raw go build will not build the required dependency, so
# we use make to achieve this.
#
deps:
	( go get -d $(GO_OPENZWAVE) ) || true
	cd $(GOPATH)/src/$(GO_OPENZWAVE) && make deps

test: install
	go test -v ./...

vet: install
	go vet ./...

run: here
	./go/bin/zwcli -monitor

.PHONY: all	dist clean test
