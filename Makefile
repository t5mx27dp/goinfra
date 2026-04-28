GO := GOPROXY=https://goproxy.cn,direct go

.DEFAULT_GOAL := all

.PHONY: all
all: test

.PHONY: test
test:
	$(GO) test -count=1 -race ./...

.PHONY: test-verbose
test-verbose:
	$(GO) test -count=1 -race -v ./...
