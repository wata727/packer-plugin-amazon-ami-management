default: build

test:
	go test ./...

build: test
	go build -v

install: build
	mkdir -p ~/.packer.d/plugins
	install ./packer-post-processor-amazon-ami-management ~/.packer.d/plugins/

release: test
	go get -u github.com/mitchellh/gox
	mkdir -p dist releases
	gox --output 'dist/{{.OS}}_{{.Arch}}/{{.Dir}}' -ldflags="-w -s"
	zip -j releases/packer-post-processor-amazon-ami-management_darwin_386.zip    dist/darwin_386/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_darwin_amd64.zip  dist/darwin_amd64/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_freebsd_386.zip   dist/freebsd_386/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_freebsd_amd64.zip dist/freebsd_amd64/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_freebsd_arm.zip   dist/freebsd_arm/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_linux_386.zip     dist/linux_386/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_linux_amd64.zip   dist/linux_amd64/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_linux_arm.zip     dist/linux_arm/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_netbsd_386.zip    dist/netbsd_386/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_netbsd_amd64.zip  dist/netbsd_amd64/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_netbsd_arm.zip    dist/netbsd_arm/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_openbsd_386.zip   dist/openbsd_386/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_openbsd_amd64.zip dist/openbsd_amd64/packer-post-processor-amazon-ami-management
	zip -j releases/packer-post-processor-amazon-ami-management_windows_386.zip   dist/windows_386/packer-post-processor-amazon-ami-management.exe
	zip -j releases/packer-post-processor-amazon-ami-management_windows_amd64.zip dist/windows_amd64/packer-post-processor-amazon-ami-management.exe

clean:
	rm -rf dist/
	rm -rf releases/

mock:
	go get -u github.com/golang/mock/mockgen
	go generate ./...

.PHONY: default test build install release clean
