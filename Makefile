PACKAGES = \
	./internal/... \
	./packet/... \
	./client/... \
	./server/...

default: test lint

test:
	go test $(PACKAGES)

lint:
	go vet ./...
	golint ./...

deps:
	go get -v -u -d -t ./...

.PHONY: test lint deps
