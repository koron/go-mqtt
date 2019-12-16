default: test

test:
	go test ./internal/... ./packet ./client ./server ./itest
	go build -v ./examples/... ./cmd/...

test-full:
	go test -v -race ./...

lint:
	go vet ./...
	@echo ""
	golint ./...

cyclo:
	-gocyclo -top 10 -avg .

report:
	@echo "misspell"
	@find . -name *.go | xargs misspell
	@echo ""
	-gocyclo -over 14 -avg .
	@echo ""
	go vet ./...
	@echo ""
	golint ./...

deps:
	go get -v -u -d -t ./...

tags:
	gotags -f tags -R .

.PHONY: test test-full lint cyclo report deps tags
