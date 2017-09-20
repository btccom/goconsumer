
# determine number of cores so we can create equivelant amount of DBs for tests
CORES=$(shell cat /proc/cpuinfo | grep processor | wc -l)

# gather options for tests
TESTARGS=$(TESTOPTIONS)

# gather options for coverage
COVERAGEARGS=$(COVERAGEOPTIONS)

test: test-cleanup test-root test-redis test-mock
test-race: test-race-root test-race-redis test-race-mock

test-cleanup:
		rm -rf coverage/ 2>> /dev/null || exit 0 &&       \
		mkdir coverage

test-root:
		go test -coverprofile=coverage/goconsumer.out -v  \
		github.com/btccom/goconsumer                      \
		$(TESTARGS)

test-redis:
		go test -coverprofile=coverage/redis.out -v       \
		github.com/btccom/goconsumer/redis                \
		$(TESTARGS)

test-mock:
		go test -coverprofile=coverage/mock.out -v        \
		github.com/btccom/goconsumer/mock                 \
		$(TESTARGS)

test-race-root:
		go test -test -v                                  \
		github.com/btccom/goconsumer                      \
		$(TESTARGS)

test-race-redis:
		go test -race -v                                  \
		github.com/btccom/goconsumer/redis                \
		$(TESTARGS)

test-race-mock:
		go test -race -v                                  \
		github.com/btccom/goconsumer/mock                 \
		$(TESTARGS)


# concat all coverage reports together
coverage-concat:
	echo "mode: set" > coverage/full && \
    grep -h -v "^mode:" coverage/*.out >> coverage/full

# full coverage report
coverage: coverage-concat
	go tool cover -func=coverage/full $(COVERAGEARGS)

# full coverage report
coverage-html: coverage-concat
	go tool cover -html=coverage/full $(COVERAGEARGS)

minimum-coverage: coverage-concat
	./tools/minimum-coverage.sh 90

benchmark-mock:
	cd mock && go test -bench=.

benchmark-redis:
	cd redis && go test -bench=.
