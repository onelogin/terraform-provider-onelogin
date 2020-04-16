package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/onelogin/onelogin-terraform-provider/onelogin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: onelogin.Provider,
	})
}
