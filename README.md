# Tileinspect

Tileinspect is a helpful utility to gather and validate information about a tile. Most of this data is available inside the tile's metadata file, but this utility simplifies the retrieval of that data.

## Commands

### `tileinspect metadata`

Prints the entire metadata file.

### `tileinspect stemcell`

Prints the stemcell criteria information for this tile


### `tileinspect check-config`
Compares the tile's property blueprints and a config file (in JSON or YAML format) and checks if this config could be used to deploy this tile.

Specifically, this check will:
* Check the config file for syntax errors
* Check the config file for a `product-properties` section
* Check the config file for properties that are not defined in the tile
* Check the config file for properties that are not available, because they are in a selector property option that was not selected
* Check the tile for any required properties without defaults that are not supplied in the config file

### `tileinspect make-config`

Creates a valid config file for this tile. This will provide a quick starting point for making config files for repeated testing.

Tileinspect will pick a value for the properties in this order:
* A value provided with the `-v|--value` CLI option
* A default value provided specified by the tile
* For `dropdown_select` and `selector` properties, the first value
* A sample value (e.g. `SAMPLE_STRING_VALUE`) that is meant to be replaced

For tiles with selectors, non-selected options will not have any values for their properties in the config file. Use the `-v` flag to set a value for that selector and `tileinspect make-config` will populate the config with the properties for the selected option.

Example:
```
tileinspect make-config -t my-tile.pivotal -v .properties.network_selector:"Use TCP"
``` 

### `tileinspect version`

Prints the current version of Tileinspect.

## Developing

When making changes, please utilize the Makefile for testing and building:

`make test` will execute the unit tests

`make test-features` will execute the feature tests

`make build` will build the `tileinspect` binary
