# Makefile to interact with tektoncd-catalog
all: help

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

