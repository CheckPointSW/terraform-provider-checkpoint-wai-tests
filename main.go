package main

import (
	"context"
	"flag"
	"log"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/chkp-haimk/checkpoint-workforce-ai", // TODO: Replace with CheckPointSW before GA
		// Provider type name: checkpoint_workforce_ai
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
