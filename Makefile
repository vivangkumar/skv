.PHONY: test build

fmt:
	gofmt -s -w .

build:
	rm -rf ./build
	go build -o build/skv

run:
	LOG_LEVEL=debug ./build/skv

test:
	go test ./... -race -v
