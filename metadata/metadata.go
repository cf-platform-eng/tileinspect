package metadata

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"regexp"

	"github.com/cf-platform-eng/tileinspect"
	"github.com/ghodss/yaml"
	. "github.com/pkg/errors"
)


type Config struct {
	tileinspect.TileConfig
	// duplicate choice required by go-flags
	// nolint:staticcheck
	Format string `long:"format" short:"f" description:"output file type" choice:"yaml" choice:"json" default:"yaml"`
}

func findInZip(term string, zip *zip.ReadCloser) *zip.File {
	re := regexp.MustCompile(term)
	for _, f := range zip.File {
		if re.Match([]byte(f.Name)) {
			return f
		}
	}
	return nil
}

func (cmd *Config) dumpFile(zipFile *zip.File, out io.Writer) error {
	inFile, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer inFile.Close()

	if cmd.Format == "json" {
		buf, err := io.ReadAll(inFile)
		if err != nil {
			return err
		}

		j, err := yaml.YAMLToJSON(buf)
		if err != nil {
			return err
		}

		_, err = out.Write(j)
		return err
	} else {
		_, err = io.Copy(out, inFile)
		return err
	}
}

func (cmd *Config) findMetadataFile() (*zip.File, error) {
	tile, err := zip.OpenReader(cmd.Tile)
	if err != nil {
		return nil, Wrapf(err, "could not unzip %s", cmd.Tile)
	}

	metadataFile := findInZip(`metadata/.*\.yml`, tile)
	if metadataFile == nil {
		return nil, errors.New("metadata file not found")
	}

	return metadataFile, nil
}

func (cmd *Config) LoadMetadata(target interface{}) error {
	metadataFile, err := cmd.findMetadataFile()
	if err != nil {
		return err
	}

	file, err := metadataFile.Open()
	if err != nil {
		return Wrap(err, "could not open the metadata file")
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return Wrap(err, "could not read the metadata file")
	}

	err = yaml.Unmarshal(buf, &target)
	if err != nil {
		return Wrap(err, "could not load the metadata file")
	}

	return nil
}

func (cmd *Config) WriteMetadata(out io.Writer) error {
	metadataFile, err := cmd.findMetadataFile()
	if err != nil {
		return err
	}

	err = cmd.dumpFile(metadataFile, out)
	if err != nil {
		return Wrapf(err, "could not read from %s (found inside %s)", metadataFile.Name, cmd.Tile)
	}

	return nil
}

func (cmd *Config) Execute(args []string) error {
	return cmd.WriteMetadata(os.Stdout)
}
