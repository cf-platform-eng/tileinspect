package makeconfig

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/cf-platform-eng/tileinspect"
	"github.com/cf-platform-eng/tileinspect/checkconfig"
	"github.com/cf-platform-eng/tileinspect/metadata"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

type Config struct {
	tileinspect.TileConfig
	Format      string            `long:"format" short:"f" description:"output file type" choice:"yaml" choice:"json" default:"yaml"`
	Values      map[string]string `long:"value" short:"v" description:"set a value for a given property with the format: .properties.key:value"`
	MetadataCmd tileinspect.MetadataCmd
}

var SampleValues = map[string]interface{}{
	"boolean":            false,
	"disk_type_dropdown": "{disk_type}",
	"integer":            int(0),
	"network_address":    "SAMPLE_NETWORK_ADDRESS",
	"port":               int(0),
	"secret": map[string]interface{}{
		"secret": "SAMPLE_SECRET_VALUE",
	},
	"string":           "SAMPLE_STRING_VALUE",
	"vm_type_dropdown": "{vm_type}",
}

func (cmd *Config) getValueForProperty(property tileinspect.TileProperty, valueOverride string) interface{} {
	if valueOverride != "" {
		property.Default = valueOverride
	}

	if property.Default != nil {
		if property.Type == "secret" {
			return map[string]interface{}{
				"secret": property.Default.(string),
			}
		} else {
			return property.Default
		}
	}

	if property.Type == "dropdown_select" {
		return property.Options[0].Name
	} else if property.Type == "selector" {
		return property.ChildProperties[0].SelectValue
	}

	return SampleValues[property.Type]
}

func (cmd *Config) setValuesForProperties(config *tileinspect.ConfigFile, propertyPrefix string, tileProperties []tileinspect.TileProperty) {
	for _, property := range tileProperties {
		propertyKey := propertyPrefix + "." + property.Name
		if !property.Configurable {
			continue
		}

		if config.ProductProperties[propertyKey] == nil {
			config.ProductProperties[propertyKey] = &tileinspect.ConfigFileProperty{
				Value: cmd.getValueForProperty(property, cmd.Values[propertyKey]),
				Type:  property.Type,
			}
		}

		if property.Type == "selector" {
			for _, option := range property.ChildProperties {
				if config.ProductProperties[propertyKey].Value == option.SelectValue {
					cmd.setValuesForProperties(config, propertyKey+"."+option.Name, option.PropertyBlueprints)
				}
			}
		}
	}
}

func (cmd *Config) MakeConfig() (*tileinspect.ConfigFile, error) {
	tileProperties := &tileinspect.TileProperties{}
	err := cmd.MetadataCmd.LoadMetadata(tileProperties)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load metadata from the tile")
	}

	config := &tileinspect.ConfigFile{
		ProductProperties: make(map[string]*tileinspect.ConfigFileProperty),
	}

	cmd.setValuesForProperties(config, ".properties", tileProperties.PropertyBlueprints)

	check := &checkconfig.Config{}
	errs := check.CompareProperties(config, tileProperties)
	if len(errs) > 0 {
		errorStrings := make([]string, len(errs))
		for i := range errs {
			errorStrings[i] = errs[i].Error()
		}
		return nil, errors.Errorf("failed to construct a valid config file:\n%s\n", strings.Join(errorStrings, "\n"))
	}

	return config, nil
}

func (cmd *Config) Execute(args []string) error {
	cmd.MetadataCmd = &metadata.Config{
		TileConfig: tileinspect.TileConfig{
			Tile: cmd.Tile,
		},
	}

	config, err := cmd.MakeConfig()
	if err != nil {
		return err
	}

	var bytes []byte
	if cmd.Format == "yaml" {
		bytes, err = yaml.Marshal(config)
	} else if cmd.Format == "json" {
		bytes, err = json.Marshal(config)
	}
	if err != nil {
		return errors.Wrap(err, "failed to convert config file")
	}

	_, err = os.Stdout.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "failed to print config file")
	}

	return err
}
