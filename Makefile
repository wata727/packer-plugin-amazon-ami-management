default: build

NAME=amazon-ami-management
BINARY=packer-plugin-${NAME}
MOCK_VERSION?=$(shell go list -m github.com/golang/mock | cut -d " " -f2)
SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)
PLUGIN_FQN=$(shell grep -E '^module' <go.mod | sed -E 's/module \s*//')

deps:
	go install github.com/golang/mock/mockgen@${MOCK_VERSION}
	go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${SDK_VERSION}

generate: deps
	go generate ./...

test: deps
	go test ./...

build: test
	go build -v

install: build
	packer plugins install --path ${BINARY} "$(shell echo "${PLUGIN_FQN}" | sed 's/packer-plugin-//')"

plugin-check: deps build
	packer-sdc plugin-check ${BINARY}

.PHONY: default deps generate test build install plugin-check
