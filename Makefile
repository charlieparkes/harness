INSTALL_DIR?=$(HOME)/.local/bin

define announce
	@printf '\033[36m⌁ %s\033[0m\n' "$(1)"
endef

.PHONY: default
default: all

.PHONY: all
all: generate mocks lint test build

SOURCES := $(shell find cmd internal -type f -name '*.go') go.mod $(wildcard go.sum)
bin/harness: $(SOURCES)
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o ./bin/harness ./cmd/harness

.PHONY: generate
generate:
	$(call announce,$@)
	go generate ./...

.PHONY: build
build: generate bin/harness
	$(call announce,$@)

.PHONY: install
install: bin/harness
	$(call announce,$@)
	mkdir -p $(INSTALL_DIR)
	cp ./bin/harness $(INSTALL_DIR)/harness

.PHONY: lint
lint:
	$(call announce,$@)
	golangci-lint run --fix ./...

.PHONY: test
test:
	$(call announce,$@)
	go test ./...

.PHONY: mocks
mocks:
	$(call announce,$@)
	go tool github.com/vektra/mockery/v3
