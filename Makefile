TEST_PACKAGE ?= ./...

.PHONY: build
build:
	go build -gcflags '-e' ./examples/... ./cmd/...

.PHONY: test
test:
	go test $(TEST_PACKAGE)

.PHONY: test-full
test-full:
	go test -v -race ./...

.PHONY: bench
bench:
	go test -bench $(TEST_PACKAGE)

.PHONY: tags
tags:
	gotags -f tags -R .

.PHONY: cover
cover:
	mkdir -p tmp
	go test -coverprofile tmp/_cover.out $(TEST_PACKAGE)
	go tool cover -html tmp/_cover.out -o tmp/cover.html

.PHONY: checkall
checkall: vet staticcheck

.PHONY: vet
vet:
	go vet $(TEST_PACKAGE)

.PHONY: staticcheck
staticcheck:
	staticcheck $(TEST_PACKAGE)

.PHONY: clean
clean:
	go clean
	rm -f tags
	rm -f tmp/_cover.out tmp/cover.html

# based on: github.com/koron-go/_skeleton/Makefile
