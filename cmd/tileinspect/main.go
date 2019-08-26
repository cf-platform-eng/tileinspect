package main

import (
	"fmt"
	"os"

	"github.com/cf-platform-eng/tileinspect/checkconfig"

	"github.com/cf-platform-eng/tileinspect/stemcell"
	"github.com/jessevdk/go-flags"

	"github.com/cf-platform-eng/tileinspect"
	"github.com/cf-platform-eng/tileinspect/metadata"
	"github.com/cf-platform-eng/tileinspect/version"
)

var checkConfigOpts checkconfig.Config
var metadataOpts metadata.Config
var stemcellOpts stemcell.Config
var config tileinspect.Config
var parser = flags.NewParser(&config, flags.Default)

func main() {
	_, err := parser.AddCommand(
		"check-config",
		"Check config file",
		"Check that a config file for any issues with the given tile",
		&checkConfigOpts,
	)
	if err != nil {
		fmt.Println("Could not add check-config command")
		os.Exit(1)
	}

	_, err = parser.AddCommand(
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

	_, err = parser.AddCommand(
		"version",
		"print version",
		"print tileinspect version",
		&version.VersionOpt{})
	if err != nil {
		fmt.Println("Could not add version command")
		os.Exit(1)
	}

	_, err = parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}
