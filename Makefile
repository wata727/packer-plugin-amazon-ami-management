default: build

test:
	go test ./...

build: test
	go build -v

install: build
	mkdir -p ~/.packer.d/plugins
	install ./packer-post-processor-amazon-ami-management ~/.packer.d/plugins/

.PHONY: default test build install
