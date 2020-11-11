.PHONY: pact
pact: test
	mkdir -p ../server/pacts
	cp target/pacts/* ../server/pacts

.PHONY: test
test:
	./mvnw clean verify