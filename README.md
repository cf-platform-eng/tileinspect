# Tileinspect

Tileinspect is a helpful utility to gather and validate information about a tile. Most of this data is available inside the tile's metadata file, but this utility simplifies the retrieval of that data.

## Commands

* `tileinspect metadata` - Prints the entire metadata file
* `tileinspect stemcell` - Prints the stemcell criteria information for this tile
* `tileinspect check-config` - Compares the tile's property blueprints and a config file (in JSON or YAML format) and checks if this config could be used to deploy this tile. Specifically, this check will:
  * Check the config file for syntax errors
  * Check the config file for a `product-properties` section
  * Check the config file for properties that are not defined in the tile
  * Check the config file for properties that are not available, because they are in a selector property option that was not selected
  * Check the tile for any required properties without defaults that are not supplied in the config file

## Developing

Utilize the Makefile for testing and building.

`make test` will execute the unit tests

`make build` will build the `tileinspect` binary