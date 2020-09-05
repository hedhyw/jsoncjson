PRE_COMMIT_HOOK := ./.git/hooks/pre-commit

all: lint test hooks

hooks:
	echo "make lint test" > $(PRE_COMMIT_HOOK)
	chmod +x $(PRE_COMMIT_HOOK)
.PHONY: hooks

lint: bin/golangci-lint
	./bin/golangci-lint run
.PHONY: lint

test.coverage:
	go test -v -covermode=count -coverprofile=coverage.out
.PHONY: test.coverage

test:
	go test
.PHONY: test

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.30.0
