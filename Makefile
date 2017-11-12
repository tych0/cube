BENCH_FLAGS ?= -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem
PKGS ?= $(shell glide novendor | grep -v examples)
PKG_FILES ?= $(shell find . --name *.go | grep -v vendor)
GO_VERSION := $(shell go version | cut -d " " -f 3)

.PHONY: all
all: lint test

.PHONY: dependencies
dependencies:
	@echo "Installing Glide and locked dependencies..."
	glide --version || go get -u -f github.com/Masterminds/glide
	glide install
	@echo "Installing gometalineter"
	gometalinter --version || go get -u -f github.com/alecthomas/gometalinter && gometalinter --install
	gometalinter --install

.PHONY: lint
lint:
	@rm -rf lint.log
	@echo "Checking formatting..."
	@$(foreach dir,$(PKGS_FILES),gofmt -d -s $(dir) 2>&1 | tee -a lint.log;)
	@echo "Running gometalinter"
	@$(foreach dir,$(PKGS),gometalinter --disable-all --enable=golint --enable=vet $(dir) 2>&1 | tee -a lint.log;)
	@[ ! -s lint.log ]

.PHONY: test
test:
	@.build/test.sh

.PHONY: ci
ci: SHELL := /bin/bash
ci: test
	bash <(curl -s https://codecov.io/bash)

.PHONY: bench
BENCH ?= .
bench:
	@$(foreach pkg,$(PKGS),go test -bench=$(BENCH) -run="^$$" $(BENCH_FLAGS) $(pkg);)