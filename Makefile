.PHONY: test
test:
	go test ./internal/... ./mqtopic ./packet ./client ./server ./itest
	go build -v ./examples/... ./cmd/...

.PHONY: test-full
test-full:
	go test -v -race ./...

.PHONY: tags
tags:
	gotags -f tags -R .

.PHONY: cover
cover:
	mkdir -p tmp
	go test -coverprofile tmp/_cover.out ./...
	go tool cover -html tmp/_cover.out -o tmp/cover.html

.PHONY: checkall
checkall: vet lint staticcheck

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: clean
clean:
	go clean
	rm -f tags
	rm -f tmp/_cover.out tmp/cover.html

# based on: github.com/koron-go/_skeleton/Makefile
