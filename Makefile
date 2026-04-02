.PHONY: build test lint run clean fmt vet docs

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

docs: build
	@echo "Starting server... API docs at http://localhost:8080/docs"
	./$(BINARY)

clean:
	rm -f $(BINARY)
