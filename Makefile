.PHONY: build test lint run clean fmt vet

BINARY := ism-api
PKG := ./...

build:
	go build -o $(BINARY) ./cmd/server

run: build
	./$(BINARY)

test:
	gotestsum --format pkgname-and-test-fails -- -race $(PKG)

lint: vet
	@echo "Lint passed (go vet)"

vet:
	go vet $(PKG)

fmt:
	gofmt -w .

clean:
	rm -f $(BINARY)
