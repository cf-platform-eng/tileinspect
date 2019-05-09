package main

import (
	"fmt"
	"os"

	"github.com/cf-platform-eng/isv-ci-toolkit/tileinspect/stemcell"
	flags "github.com/jessevdk/go-flags"

	"github.com/cf-platform-eng/isv-ci-toolkit/tileinspect"
	"github.com/cf-platform-eng/isv-ci-toolkit/tileinspect/metadata"
)

var metadataOpts metadata.Config
var stemcellOpts stemcell.Config
var config tileinspect.Config
var parser = flags.NewParser(&config, flags.Default)

func main() {
	_, err := parser.AddCommand(
		"metadata",
		"Dump metadata",
		"Dump tile metadata to stdout",
		&metadataOpts,
	)
	if err != nil {
		fmt.Println("Could not add metadata command")
		os.Exit(1)
	}

	_, err = parser.AddCommand(
		"stemcell",
		"Dump stemcell requirement",
		"Dump stemcell requirement to stdout",
		&stemcellOpts,
	)
	if err != nil {
		fmt.Println("Could not add stemcell command")
		os.Exit(1)
	}

	_, err = parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}
