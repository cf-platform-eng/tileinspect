package main

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/cf-platform-eng/tileinspect/checkconfig"
	"github.com/cf-platform-eng/tileinspect/makeconfig"

	"github.com/cf-platform-eng/tileinspect/stemcell"
	"github.com/jessevdk/go-flags"

	"github.com/cf-platform-eng/tileinspect"
	"github.com/cf-platform-eng/tileinspect/metadata"
	"github.com/cf-platform-eng/tileinspect/version"
)

var checkConfigOpts checkconfig.Config
var makeConfigOpts makeconfig.Config
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
		"make-config",
		"Make a template config file",
		heredoc.Doc(`Make a template config file based on the property blueprints of this tile.
		The config file will contain a value for each selected, configurable property.
		The value will be, in order:
		* a value defined using the "--value" parameter
		* a default value defined in the tile
		* for dropdown_select and selector properties, the first option
		* a sample value that is meant to be replaced by the user
		
		Using the -v, --value parameter is useful for setting known values or for selecting a preferred option in a selector.
		
		Example: tileinspect make-config -t my-tile.pivotal -v .properties.network_selector:"Use TCP"`),
		&makeConfigOpts,
	)
	if err != nil {
		fmt.Println("Could not add make-config command")
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
