// +build feature

package features

import (
	"archive/zip"
	"io/ioutil"
	"os"
)

func MakeTileWithMetadata(metadata string) (*os.File, error) {
	file, err := ioutil.TempFile("", "feature-test-tile-*.pivotal")
	if err != nil {
		return nil, err
	}

	writer := zip.NewWriter(file)

	metadataWriter, err := writer.Create("metadata/metadata.yml")
	if err != nil {
		return nil, err
	}

	_, err = metadataWriter.Write([]byte(metadata))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	return file, err
}

func MakeConfigFile(data string) (*os.File, error) {
	file, err := ioutil.TempFile("", "config")
	if err != nil {
		return nil, err
	}

	_, err = file.Write([]byte(data))
	if err != nil {
		return nil, err
	}
	file.Close()

	return file, err
}
