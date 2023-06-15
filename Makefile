APP = catalog-cd
OUTPUT_DIR ?= bin
BIN = $(OUTPUT_DIR)/$(APP)

CMD ?= ./cmd/$(APP)/...
PKG ?= ./pkg/...

GOFLAGS ?= -v
GOFLAGS_TEST ?= -v -cover

ARGS ?=

.EXPORT_ALL_VARIABLES:

all: help

.PHONY: $(BIN)
$(BIN):
	go build -o $(BIN) $(CMD) $(ARGS)

.PHONY: build
build: $(BIN)

.PHONY: run
run:
	go run $(CMD) $(ARGS)

install:
	go install $(CMD)

test: test-unit

.PHONY: test-unit
test-unit:
	go test $(GOFLAGS_TEST) $(CMD) $(PKG) $(ARGS)

.PHONY: test-e2e/openshift
.PHONY: test-e2e-openshift
test-e2e/openshift: test-e2e-openshift
test-e2e-openshift: ## Run e2e tests on OpenShift.
	./automation/e2e-tests.sh openshift

.PHONY: test-e2e/kubernetes test-e2e-kubernetes
test-e2e/kubernetes: test-e2e-kubernetes
test-e2e-kubernetes: ## Run e2e tests on Kubernetes.
	./automation/e2e-tests.sh kubernetes

.PHONY: help
help:
	@grep -hE '^[ a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'
