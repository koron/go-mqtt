default: test

test:
	go test -v ./...

lint:
	go vet ./...
	@echo ""
	golint ./...

report:
	echo "misspell"
	@find . -name *.go | xargs misspell
	@echo ""
	-gocyclo.exe -over 9 -avg .
	@echo ""
	go vet ./...
	@echo ""
	golint ./...

.PHONY: test lint
