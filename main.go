package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/onelogin/terraform-provider-onelogin/onelogin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: onelogin.Provider,
	})
}
