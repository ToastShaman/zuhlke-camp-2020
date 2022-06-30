PWD := $(shell pwd)

pact: test
	mkdir -p ../server/pacts
	cp build/pacts/* ../server/pacts
.PHONY: pact

test:
	./gradlew clean check
.PHONY: test

publish:
	./gradlew pactPublish
.PHONY: publish

stub: pact
	docker run -t -p 8080:8080 -v "$(PWD)/build/pacts/:/app/pacts" pactfoundation/pact-stub-server -p 8080 -d pacts
.PHONY: stub