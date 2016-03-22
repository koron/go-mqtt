default: test

test:
	go test ./...

test-full:
	go test -v -race ./...

lint:
	go vet ./...
	@echo ""
	golint ./...

report:
	@echo "misspell"
	@find . -name *.go | xargs misspell
	@echo ""
	-gocyclo -over 9 -avg .
	@echo ""
	go vet ./...
	@echo ""
	golint ./...

deps:
	go get -v -u -d -t ./...

.PHONY: test test-full lint report deps
