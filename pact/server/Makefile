DIR := ${CURDIR}
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(COMMIT)-$(DATE)
FLAGS := -ldflags "-X main.version=$(VERSION)"

test:
	@go test -v
.PHONY: test

lbuild: clean test
	@env GOOS=darwin GOARCH=amd64 go build $(FLAGS) -o out/todo-server-${VERSION}-darwin-amd64
.PHONY: build

build: clean test
	@env DOCKER_BUILDKIT=1 docker build --build-arg VERSION="$(VERSION)" -o out .
.PHONY: dbuild

clean:
	@rm -rf ./out
.PHONY: clean

check: spellcheck staticcheck securitycheck vet
.PHONY: check

spellcheck:
	@command -v misspell > /dev/null 2>&1 || go get github.com/client9/misspell/cmd/misspell
	@misspell -locale="UK" -i="color" -error -source="text" main.go
.PHONY: spellcheck

staticcheck:
	@command -v staticcheck > /dev/null 2>&1 || go get honnef.co/go/tools/cmd/staticcheck
	@staticcheck -checks="all" -tests $(GOFMT_FILES)
.PHONY: staticcheck

securitycheck:
	@command -v gosec > /dev/null 2>&1 || go get github.com/securego/gosec/v2/cmd/gosec
	@gosec ./...
.PHONY: securitycheck

vet:
	@go vet ./...
.PHONY: vet

# example: make run ARGS="-cors"
run:
	@go run $(FLAGS) main.go $(ARGS)
.PHONY: run

newman:
	@docker run --rm -it \
		--network="host" \
		-v $(DIR)/documentation/postman:/etc/newman \
		-t postman/newman \
		run todo_api.postman_collection.json --environment="docker.postman_environment.json"
.PHONY: newman

perf:
	@docker run --rm -i \
		loadimpact/k6 \
		run -e MY_HOSTNAME=host.docker.internal:3000 - <perf/script.js
.PHONY: perf

swagger:
	@docker run --rm \
		-p 9002:8080 \
		-v $(DIR)/documentation/api.yml:/api.yaml \
		-e SWAGGER_JSON=/api.yaml \
		swaggerapi/swagger-ui
.PHONY: swagger

