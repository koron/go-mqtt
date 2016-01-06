default: test lint

test:
	go test -v ./...

lint:
	go vet ./...
	golint ./...

.PHONY: test lint
