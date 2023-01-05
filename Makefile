default: build

MOCK_VERSION?=$(shell go list -m github.com/golang/mock | cut -d " " -f2)
SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

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
	mkdir -p ~/.packer.d/plugins
	mv ./packer-plugin-amazon-ami-management ~/.packer.d/plugins/

.PHONY: default deps test build install
