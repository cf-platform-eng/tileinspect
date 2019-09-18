package makeconfig

import (
	"encoding/json"
	"os"

	"github.com/cf-platform-eng/tileinspect"
	"github.com/cf-platform-eng/tileinspect/metadata"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

//go:generate counterfeiter MetadataCmd
type MetadataCmd interface {
	LoadMetadata(target interface{}) error
}

type Config struct {
	tileinspect.TileConfig
	Format      string `long:"format" short:"f" description:"output file type" choice:"yaml" choice:"json" default:"yaml"`
	MetadataCmd MetadataCmd
}

type TileProperty struct {
	Name            string      `json:"name"`
	Type            string      `json:"type"`
	Configurable    bool        `json:"configurable"`
	Default         interface{} `json:"default"`
	Optional        bool        `json:"optional"`
	Options         []Option
	ChildProperties []TileProperties `json:"option_templates"`
}
type TileProperties struct {
	Name               string         `json:"name"`
	PropertyBlueprints []TileProperty `json:"property_blueprints"`
	SelectValue        string         `json:"select_value"`
}
type Option struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

var SampleValues = map[string]interface{}{
	"boolean":            false,
	"disk_type_dropdown": "{disk_type}",
	"integer":            int(0),
	"network_address":    "SAMPLE_NETWORK_ADDRESS",
	"port":               int(0),
	"secret": map[string]string{
		"secret": "SAMPLE_SECRET_VALUE",
	},
	"string":           "SAMPLE_STRING_VALUE",
	"vm_type_dropdown": "{vm_type}",
}

func configForProperties(propertyPrefix string, tileProperties []TileProperty) (*tileinspect.ConfigFile, error) {
	config := &tileinspect.ConfigFile{
		ProductProperties: make(map[string]*tileinspect.ConfigFileProperty),
	}

	for _, property := range tileProperties {
		propertyKey := propertyPrefix + "." + property.Name
		if !property.Configurable {
			continue
		}

		defaultValue := SampleValues[property.Type]
		if property.Type == "dropdown_select" {
			defaultValue = property.Options[0].Name
		}
		if property.Type == "selector" {
			defaultValue = property.ChildProperties[0].SelectValue
		}

		configProperty := &tileinspect.ConfigFileProperty{
			Value: defaultValue,
			Type:  property.Type,
		}
		config.ProductProperties[propertyKey] = configProperty

		if property.Type == "secret" {
			if property.Default != nil {
				configProperty.Value = map[string]string{
					"secret": property.Default.(string),
				}
			}
		} else {
			if property.Default != nil {
				configProperty.Value = property.Default
			}
		}

		if property.Type == "selector" {
			for _, option := range property.ChildProperties {
				if configProperty.Value == option.SelectValue {
					childConfigs, _ := configForProperties(propertyKey+"."+option.Name, option.PropertyBlueprints)
					for k, v := range childConfigs.ProductProperties {
						config.ProductProperties[k] = v
					}
				}
			}
		}
	}

	return config, nil
}

func (cmd *Config) MakeConfig() (*tileinspect.ConfigFile, error) {
	tileProperties := &TileProperties{}
	err := cmd.MetadataCmd.LoadMetadata(tileProperties)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load metadata from the tile")
	}

	return configForProperties(".properties", tileProperties.PropertyBlueprints)
}

func (cmd *Config) Execute(args []string) error {
	cmd.MetadataCmd = &metadata.Config{
		TileConfig: tileinspect.TileConfig{
			Tile: cmd.Tile,
		},
	}

	config, err := cmd.MakeConfig()
	if config != nil {
		var bytes []byte
		if cmd.Format == "yaml" {
			bytes, err = yaml.Marshal(config)
		} else if cmd.Format == "json" {
			bytes, err = json.Marshal(config)
		}
		if err != nil {
			return errors.Wrap(err, "failed to output config file")
		}
		_, err := os.Stdout.Write(bytes)
		if err != nil {
			return errors.Wrap(err, "failed to print config file")
		}
	}

	return err
}
