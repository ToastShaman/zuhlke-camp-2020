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
