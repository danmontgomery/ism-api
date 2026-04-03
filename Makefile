.PHONY: build test lint run clean fmt vet docs docker docker-run

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

docker:
	docker build -t ism-api .

docker-run: docker
	docker run -p 8080:8080 ism-api

clean:
	rm -f $(BINARY)
