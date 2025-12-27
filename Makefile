BINARY_NAME=eogo

build:
	go build -o $(BINARY_NAME) cmd/eogo/main.go

install: build
	mv $(BINARY_NAME) /usr/local/bin/

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test ./...

server:
	go run cmd/server/main.go

air:
	air

dev: air
