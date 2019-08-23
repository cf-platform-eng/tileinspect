package checkconfig

import (
    "errors"
    "fmt"
    "github.com/cf-platform-eng/tileinspect"
    "github.com/cf-platform-eng/tileinspect/metadata"
    "github.com/ghodss/yaml"
    . "github.com/pkg/errors"
    "io"
    "io/ioutil"
    "os"
    "strings"
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

func checkTileProperties(propertyPrefix string, configFile *ConfigFile, tileProperties []TileProperty) ([]string, []error) {
    var errs []error
    var keys []string

    for _, property := range tileProperties {
        key := propertyPrefix + "." + property.Name
        keys = append(keys, key)

        if configFile.ProductProperties[key] != nil {
            if !property.Configurable {
                errs = append(errs, fmt.Errorf("the config file contains a property (%s) that is not configurable", key))
            }
        }

        if property.Type == "selector" {
            for _, option := range property.ChildProperties {
                childPrefix := key + "." + option.Name
                childKeys, childErrs := checkTileProperties(childPrefix, configFile, option.PropertyBlueprints)
                keys = append(keys, childKeys...)
                errs = append(errs, childErrs...)
            }
        }
    }

    return keys, errs
}

func (cmd *Config) CompareProperties(configFile *ConfigFile, tileProperties *TileProperties) []error {
    prefix := ".properties"
    keys, errs := checkTileProperties(prefix, configFile, tileProperties.PropertyBlueprints)

    for key := range configFile.ProductProperties {
        if strings.Index(key, prefix) != 0 {
            errs = append(errs, fmt.Errorf("the config file contains a property (%s) that does not start with %s", key, prefix))
        } else if !stringInSlice(key, keys) {
            errs = append(errs, fmt.Errorf("the config file contains a property (%s) that is not defined in the tile", key))
        }
    }

    return errs
}

func (cmd *Config) CheckConfig(out io.Writer) error {
    configFileContents, err := ioutil.ReadFile(cmd.ConfigFilePath)

    if err != nil {
        return Wrap(err, "config file does not exist")
    }

    configFile := &ConfigFile{}
    err = yaml.Unmarshal(configFileContents, configFile)
    if err != nil {
        return Wrap(err, "config file is not valid JSON or YAML")
    }

    tileProperties := &TileProperties{}
    err = cmd.MetadataCmd.LoadMetadata(tileProperties)
    if err != nil {
        return Wrap(err, "failed to load tile metadata")
    }

    errs := cmd.CompareProperties(configFile, tileProperties)
    if len(errs) > 0 {
        return errors.New("Lots of errors")
    }

    _, _ = out.Write([]byte("The config file appears to be valid"))
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
