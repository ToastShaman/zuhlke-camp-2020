DIR := ${CURDIR}
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(COMMIT)-$(DATE)
FLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: test
test:
	go test -v

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

.PHONY: newman
newman:
	docker run --rm -it --network="host" -v $(DIR)/documentation/postman:/etc/newman -t postman/newman \
		run todo_api.postman_collection.json --environment="docker.postman_environment.json"

.PHONY: perf
perf:
	docker run --rm -i loadimpact/k6 run -e MY_HOSTNAME=host.docker.internal:3000 - <perf/script.js
