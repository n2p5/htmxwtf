
SERVER_NAME=todo.server

.PHONY: generate build test run

generate:
	templ generate

build: generate
	go build -o $(SERVER_NAME)

test:
	go test -v ./...

run: build
	./$(SERVER_NAME)

clean:
	rm -f $(SERVER_NAME)

dev:
	@find . -type f \( \
	-name '*.go' ! -name '*_templ.go' -o -name '*.templ' \) \
  | entr -r sh -c 'make run'