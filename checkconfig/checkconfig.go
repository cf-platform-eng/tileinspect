package checkconfig

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cf-platform-eng/tileinspect"
	"github.com/cf-platform-eng/tileinspect/metadata"
	"github.com/ghodss/yaml"
	. "github.com/pkg/errors"
)

//go:generate counterfeiter MetadataCmd
type MetadataCmd interface {
	LoadMetadata(target interface{}) error
}

type Config struct {
	tileinspect.TileConfig
	MetadataCmd    MetadataCmd
	ConfigFilePath string `long:"config" short:"c" description:"path to config file" required:"true"`
}

type ConfigFile struct {
	ProductProperties map[string]*ConfigFileProperty `json:"product-properties"`
}

type ConfigFileProperty struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type TileProperty struct {
	Name            string           `json:"name"`
	Type            string           `json:"type"`
	Configurable    bool             `json:"configurable"`
	Default         interface{}      `json:"default"`
	Optional        bool             `json:"optional"`
	ChildProperties []TileProperties `json:"option_templates"`
}
type TileProperties struct {
	Name               string         `json:"name"`
	PropertyBlueprints []TileProperty `json:"property_blueprints"`
	SelectValue        string         `json:"select_value"`
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func checkTileProperties(checkForRequiredProperties bool, propertyPrefix string, configValues map[string]*ConfigFileProperty, tileProperties []TileProperty) ([]string, []error) {
	var errs []error
	var validKeys []string

	for _, property := range tileProperties {
		propertyKey := propertyPrefix + "." + property.Name
		validKeys = append(validKeys, propertyKey)

		if configValues[propertyKey] != nil {
			if !property.Configurable {
				errs = append(errs, fmt.Errorf("the config file contains a property (%s) that is not configurable", propertyKey))
			}
		}

		if checkForRequiredProperties {
			if property.Configurable && !property.Optional && property.Default == nil {
				if configValues[propertyKey] == nil {
					errs = append(errs, fmt.Errorf("the config file is missing a required property (%s)", propertyKey))
				}
			}
		}

		if property.Type == "selector" {
			for _, option := range property.ChildProperties {
				isSelected := configValues[propertyKey] != nil && configValues[propertyKey].Value == option.SelectValue

				childPrefix := propertyKey + "." + option.Name
				childKeys, childErrs := checkTileProperties(isSelected, childPrefix, configValues, option.PropertyBlueprints)
				validKeys = append(validKeys, childKeys...)
				errs = append(errs, childErrs...)

				if !isSelected {
					for _, childKey := range childKeys {
						if configValues[childKey] != nil {
							errs = append(errs, fmt.Errorf("the config file contains a property (%s) that is not selected", childKey))
						}
					}
				}
			}
		}
	}

	return validKeys, errs
}

func (cmd *Config) CompareProperties(configFile *ConfigFile, tileProperties *TileProperties) []error {
	prefix := ".properties"
	validKeys, errs := checkTileProperties(true, prefix, configFile.ProductProperties, tileProperties.PropertyBlueprints)

	for key := range configFile.ProductProperties {
		if strings.Index(key, prefix) != 0 {
			errs = append(errs, fmt.Errorf("the config file contains a property (%s) that does not start with %s", key, prefix))
		} else if !stringInSlice(key, validKeys) {
			errs = append(errs, fmt.Errorf("the config file contains a property (%s) that is not defined in the tile", key))
		}
	}

	return errs
}

func (cmd *Config) CheckConfig(out io.Writer) error {
	configFileContents, err := ioutil.ReadFile(cmd.ConfigFilePath)
	if err != nil {
		return Wrapf(err, "failed to read the config file: %s", cmd.ConfigFilePath)
	}

	configFile := &ConfigFile{}
	err = yaml.Unmarshal(configFileContents, configFile)
	if err != nil {
		return Wrap(err, "the config file does not contain valid JSON or YAML")
	}

	if configFile.ProductProperties == nil {
		return errors.New(`the config file is missing a "product-properties" section`)
	}

	tileProperties := &TileProperties{}
	err = cmd.MetadataCmd.LoadMetadata(tileProperties)
	if err != nil {
		return Wrap(err, "failed to load metadata from the tile")
	}

	errs := cmd.CompareProperties(configFile, tileProperties)
	if len(errs) > 0 {
		errorStrings := make([]string, len(errs))
		for i := range errs {
			errorStrings[i] = errs[i].Error()
		}
		return errors.New(strings.Join(errorStrings, "\n"))
	}

	_, _ = out.Write([]byte("The config file appears to be valid\n"))
	return nil
}

func (cmd *Config) Execute(args []string) error {
	cmd.MetadataCmd = &metadata.Config{
		TileConfig: tileinspect.TileConfig{
			Tile: cmd.Tile,
		},
	}
	return cmd.CheckConfig(os.Stdout)
}
