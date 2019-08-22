package makeconfig

import (
	"encoding/json"
	"github.com/cf-platform-eng/tileinspect"
	"github.com/cf-platform-eng/tileinspect/metadata"
	. "github.com/pkg/errors"
	"io"
	"os"
)

//go:generate counterfeiter MetadataCmd
type MetadataCmd interface {
	WriteMetadata(out io.Writer) error
}

type Config struct {
	tileinspect.TileConfig
	MetadataCmd MetadataCmd
}

type ConfigProperty struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
type ConfigTemplate struct {
	ProductProperties map[string]ConfigProperty `json:"product-properties"`
}

type BlueprintProperty struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Default      interface{} `json:"default"`
	Configurable bool        `json:"configurable"`
	Optional     bool        `json:"optional"`
	Label        string      `json:"label"`
}

type TileMetadata struct {
	PropertyBlueprints []BlueprintProperty `json:"property_blueprints"`
}

func (cmd *Config) WriteEmptyConfig(out io.Writer) error {
	var tileMetadata TileMetadata

	pr, pw := io.Pipe()
	errorChan := make(chan error)
	go func(errorChan chan error) {
		err := json.NewDecoder(pr).Decode(&tileMetadata)
		errorChan <- err
	}(errorChan)

	err := cmd.MetadataCmd.WriteMetadata(pw)
	if err != nil {
		return Wrap(err, "failed to read tile metadata")
	}

	err = <-errorChan
	if err != nil {
		return Wrap(err, "failed to decode property blueprints")
	}

	properties := ConfigTemplate{}
	properties.ProductProperties = map[string]ConfigProperty{}
	for _, property := range tileMetadata.PropertyBlueprints {
		if property.Configurable == false {
			continue
		}

		properties.ProductProperties[".properties." + property.Name] = ConfigProperty{
			Type:  property.Type,
			Value: property.Default,
		}
	}

	err = json.NewEncoder(out).Encode(properties)
	if err != nil { // !branch-not-tested No good way to force this
		return Wrap(err, "failed to encode property blueprints")
	}

	return nil
}

func (cmd *Config) Execute(args []string) error {
	cmd.MetadataCmd = &metadata.Config{
		TileConfig: tileinspect.TileConfig{
			Tile: cmd.Tile,
		},
		Format: "json",
	}
	return cmd.WriteEmptyConfig(os.Stdout)
}
