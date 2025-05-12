package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/onelogin/terraform-provider-onelogin/cmd/tests/standalone"
	"github.com/onelogin/terraform-provider-onelogin/cmd/tests/v1"
	"github.com/onelogin/terraform-provider-onelogin/cmd/tests/v4"
)

func main() {
	// Define flags
	testName := flag.String("test", "", "Specify which test to run (v1_users, create_attribute, etc.)")
	flag.Parse()

	// Show help if no test is specified
	if *testName == "" {
		fmt.Println("Please specify a test to run with the -test flag")
		fmt.Println("Available tests:")
		fmt.Println("  - v1_users")
		fmt.Println("  - create_attribute")
		fmt.Println("  - v4_custom_attributes")
		os.Exit(1)
	}

	// Run the specified test
	switch *testName {
	case "v1_users":
		v1.TestV1Users()
	case "create_attribute":
		standalone.TestCreateAttribute()
	case "v4_custom_attributes":
		v4.TestV4CustomAttributes()
	default:
		fmt.Printf("Unknown test: %s\n", *testName)
		os.Exit(1)
	}
}