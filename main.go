package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/wata727/packer-post-processor-amazon-ami-management/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}

	server.RegisterPostProcessor(new(amazonamimanagement.PostProcessor))
	server.Serve()
}
