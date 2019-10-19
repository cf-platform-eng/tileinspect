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

type Config struct {
	tileinspect.TileConfig
	MetadataCmd    tileinspect.MetadataCmd
	ConfigFilePath string `long:"config" short:"c" description:"path to config file" required:"true"`
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func checkCollectionProperties(checkForRequiredProperties bool, propertyPrefix string, configValues *[]interface{}, tileProperties *[]tileinspect.TileProperty) ([]string, []error) {
	var errs []error
	var validKeys []string

	if len(*configValues) == 0 {
		for _, prop := range *tileProperties {
			if prop.Configurable && !prop.Optional {
				errs = append(errs, fmt.Errorf("collection (%s) is missing required property %s", propertyPrefix, prop.Name))
			}
		}
	} else {
		for _, valueInterface := range *configValues {
			if value, ok := valueInterface.(map[string]interface{}); ok {
				for _, prop := range *tileProperties {
					if _, ok := value[prop.Name]; ok {
						if !prop.Configurable {
							errs = append(errs, fmt.Errorf("collection (%s) contains unconfigurable property %s", propertyPrefix, prop.Name))
						}
					} else {
						if prop.Configurable && !prop.Optional {
							errs = append(errs, fmt.Errorf("collection (%s) is missing required property %s", propertyPrefix, prop.Name))
						}
					}
				}
			} else {
				errs = append(errs, fmt.Errorf("collection (%s) contains invalid item %v", propertyPrefix, valueInterface))
			}
		}
	}

	return validKeys, errs
}

func checkTileProperties(checkForRequiredProperties bool, propertyPrefix string, configValues map[string]*tileinspect.ConfigFileProperty, tileProperties []tileinspect.TileProperty) ([]string, []error) {
	var errs []error
	var validKeys []string

	for _, property := range tileProperties {
		propertyKey := propertyPrefix + "." + property.Name
		validKeys = append(validKeys, propertyKey)
		hasValue := configValues[propertyKey] != nil

		if hasValue && !property.Configurable {
			errs = append(errs, fmt.Errorf("the config file contains a property (%s) that is not configurable", propertyKey))
		}

		if property.Type == "secret" && hasValue {
			var value map[string]interface{}
			var ok bool
			var secret string
			if value, ok = configValues[propertyKey].Value.(map[string]interface{}); !ok {
				errs = append(errs, fmt.Errorf("the config file value for property (%s) is not in the right format. Should be {\"secret\": \"<SECRET VALUE>\"}", propertyKey))
			} else if secret, ok = value["secret"].(string); !ok {
				errs = append(errs, fmt.Errorf("the config file value for property (%s) is not in the right format. Should be {\"secret\": \"<SECRET VALUE>\"}", propertyKey))
			} else if secret == "" {
				hasValue = false
			}
		}

		if checkForRequiredProperties {
			if property.Configurable && !property.Optional {
				if property.Default == nil && property.Type != "dropdown_select" {
					if !hasValue {
						errs = append(errs, fmt.Errorf("the config file is missing a required property (%s)", propertyKey))
					}
				}
			}
		}

		if property.Type == "selector" {
			for _, option := range property.ChildProperties {
				isSelected := (hasValue && configValues[propertyKey].Value == option.SelectValue) || (!hasValue && property.Default == option.SelectValue)

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

		if property.Type == "dropdown_select" && hasValue {
			validValue := false
			for _, option := range property.Options {
				if configValues[propertyKey].Value == option.Name {
					validValue = true
				}
			}
			if !validValue {
				errs = append(errs, fmt.Errorf("the config file value for property (%s) is invalid: %v", propertyKey, configValues[propertyKey].Value))
			}
		}

		if property.Type == "collection" && hasValue {
			var values []interface{}
			var ok bool
			if values, ok = configValues[propertyKey].Value.([]interface{}); ok {
				childKeys, childErrs := checkCollectionProperties(checkForRequiredProperties, propertyKey, &values, &(property.PropertyBlueprints))
				validKeys = append(validKeys, childKeys...)
				errs = append(errs, childErrs...)
			} else {
				errs = append(errs, fmt.Errorf("the config file value for the collection blueprints (%s) is not in the right format. Should be [ { \"name\": \"value\", ... }, ... ]", propertyKey))
			}
		}
	}

	return validKeys, errs
}

func (cmd *Config) CompareProperties(configFile *tileinspect.ConfigFile, tileProperties *tileinspect.TileProperties) []error {
	prefix := ".properties"
	validKeys, errs := checkTileProperties(true, prefix, configFile.ProductProperties, tileProperties.PropertyBlueprints)

	for _, jobProperties := range tileProperties.JobTypes {
		jobKeys, jobErrs := checkTileProperties(true, fmt.Sprintf(".%s", jobProperties.Name), configFile.ProductProperties, jobProperties.PropertyBlueprints)
		validKeys = append(validKeys, jobKeys...)
		errs = append(errs, jobErrs...)
	}

	for key := range configFile.ProductProperties {
		// if strings.Index(key, prefix) != 0 {
		// 	errs = append(errs, fmt.Errorf("the config file contains a property (%s) that does not start with %s", key, prefix))
		// } else
		if !stringInSlice(key, validKeys) {
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

	configFile := &tileinspect.ConfigFile{}
	err = yaml.Unmarshal(configFileContents, configFile)
	if err != nil {
		return Wrap(err, "the config file does not contain valid JSON or YAML")
	}

	if configFile.ProductProperties == nil {
		return errors.New(`the config file is missing a "product-properties" section`)
	}

	tileProperties := &tileinspect.TileProperties{}
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
