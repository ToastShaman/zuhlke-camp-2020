pact: test
	mkdir -p ../server/pacts
	cp target/pacts/* ../server/pacts
.PHONY: pact

test:
	./mvnw clean verify
.PHONY: test

publish:
	./mvnw pact:publish
.PHONY: publish
