package stemcell

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/cf-platform-eng/tileinspect"
	"github.com/cf-platform-eng/tileinspect/metadata"
	. "github.com/pkg/errors"
)

type Config struct {
	tileinspect.TileConfig
	MetadataCmd tileinspect.MetadataCmd
}

func (cmd *Config) WriteStemcell(out io.Writer) error {
	tileMetadata := &tileinspect.TileProperties{}
	err := cmd.MetadataCmd.LoadMetadata(&tileMetadata)
	if err != nil {
		return Wrap(err, "failed to load tile metadata")
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
	}
	return cmd.WriteStemcell(os.Stdout)
}
