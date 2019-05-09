package stemcell

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/cf-platform-eng/isv-ci-toolkit/tileinspect"
	"github.com/cf-platform-eng/isv-ci-toolkit/tileinspect/metadata"
	. "github.com/pkg/errors"
)

//go:generate counterfeiter MetadataCmd
type MetadataCmd interface {
	WriteMetadata(out io.Writer) error
}

type Config struct {
	tileinspect.TileConfig
	MetadataCmd MetadataCmd
}

type TileMetadata struct {
	StemcellCriteria map[string]interface{} `json:"stemcell_criteria"`
}

func (cmd *Config) WriteStemcell(out io.Writer) error {
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
		return Wrap(err, "failed to decode stemcell criteria JSON")
	}

	version, ok := tileMetadata.StemcellCriteria["version"].(string)
	if !ok {
		return errors.New("could not convert stemcell criteria version to string")
	}
	fixed := strings.Contains(version, ".")

	tileMetadata.StemcellCriteria["floating"] = !fixed
	err = json.NewEncoder(out).Encode(tileMetadata.StemcellCriteria)
	if err != nil { // !branch-not-tested No good way to force this
		return Wrap(err, "failed to encode stemcell JSON")
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
	return cmd.WriteStemcell(os.Stdout)
}
