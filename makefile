PRE_COMMIT_HOOK := ./.git/hooks/pre-commit
GOLANGCI_LINT_VER := v2.12.2

all: lint test hooks

hooks:
	echo "make lint test" > $(PRE_COMMIT_HOOK)
	chmod +x $(PRE_COMMIT_HOOK)
.PHONY: hooks

lint: bin/golangci-lint
	./bin/golangci-lint run
.PHONY: lint

test.coverage:
	go test -v -covermode=count -coverprofile=coverage.out ./...
.PHONY: test.coverage

test:
	go test ./...
.PHONY: test

bin/golangci-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ./bin $(GOLANGCI_LINT_VER)
