default: test

test:
	go test ./...

lint:
	go vet ./...
	@echo ""
	golint ./...

report:
	@echo "misspell"
	@find . -name *.go | xargs misspell
	@echo ""
	-gocyclo.exe -over 9 -avg .
	@echo ""
	go vet ./...
	@echo ""
	golint ./...

deps:
	go get -v -u -d -t ./...

.PHONY: test lint deps
