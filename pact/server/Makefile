COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(COMMIT)-$(DATE)
FLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build
build: clean
	env GOOS=darwin GOARCH=amd64 go build $(FLAGS) -o build/todo-server-${VERSION}-darwin-amd64
	env GOOS=linux GOARCH=amd64 go build $(FLAGS) -o build/todo-server-${VERSION}-linux-amd64

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: run
run:
	go run $(FLAGS) main.go
