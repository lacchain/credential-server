# This makefile defines the following targets
#
#   - all (default) - builds all targets and runs all tests

PROJECT_NAME = credential-mother
BASE_VERSION = 0.0.9
PREV_VERSION = 0.0.9-rc1
IS_RELEASE = false

ARCH=$(shell go env GOARCH)
MARCH=$(shell go env GOOS)-$(shell go env GOARCH)
STABLE_TAG ?= $(ARCH)-$(BASE_VERSION)-stable

all: rename docker unit-tests

credential-provider-server: bin/credential-provider-server

bin/%: $(GO_SOURCE)
	@echo "Building ${@F} in bin directory ..."
	@mkdir -p bin && go build -o bin/${@F} -tags "pkcs11" -ldflags "$(GO_LDFLAGS)" $(PKGNAME)/$(path-map.${@F})
	@echo "Built bin/${@F}"

all-tests: checks credential-provider-server
	@scripts/run_unit_tests
	@scripts/run_integration_tests