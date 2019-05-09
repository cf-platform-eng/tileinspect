package metadata_test

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/cf-platform-eng/isv-ci-toolkit/tileinspect"
	"github.com/cf-platform-eng/isv-ci-toolkit/tileinspect/metadata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

type BadWriter struct{}

func (w *BadWriter) Write(p []byte) (int, error) {
	return 0, errors.New("I am a bad writer")
}

func CreateTestTileWithMetadata(metadata string) (*os.File, error) {
	file, err := ioutil.TempFile(".", "test-tile-*.pivotal")
	if err != nil {
		return nil, err
	}

	writer := zip.NewWriter(file)
	if metadata != "" {
		fileWriter, err := writer.Create("metadata/test-tile-metadata.yml")
		if err != nil {
			return nil, err
		}

		_, err = fileWriter.Write([]byte(metadata))
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	return file, err
}

var _ = Describe("WriteMetadata", func() {
	var (
		buffer *Buffer
		tile   *os.File
	)

	BeforeEach(func() {
		buffer = NewBuffer()
	})

	AfterEach(func() {
		err := buffer.Close()
		Expect(err).ToNot(HaveOccurred())

		if tile != nil {
			err = os.Remove(tile.Name())
			Expect(err).ToNot(HaveOccurred())
			tile = nil
		}
	})

	Context("Valid tile", func() {
		var config metadata.Config
		BeforeEach(func() {
			var err error
			tile, err = CreateTestTileWithMetadata(heredoc.Doc(`
			---
			metadata: content`))
			Expect(err).ToNot(HaveOccurred())

			config = metadata.Config{
				TileConfig: tileinspect.TileConfig{
					Tile: tile.Name(),
				},
			}
		})

		It("extracts the metadata file from the tile", func() {
			err := config.WriteMetadata(buffer)
			Expect(err).ToNot(HaveOccurred())
			Eventually(buffer).Should(Say(""))
		})

		Context("JSON format", func() {
			BeforeEach(func() {
				config.Format = "json"
			})
			It("extracts the metadata file from the tile and outputs it in JSON", func() {
				err := config.WriteMetadata(buffer)
				Expect(err).ToNot(HaveOccurred())
				Eventually(buffer).Should(Say(`{"metadata":"content"}`))
			})
		})

		Context("Bad output", func() {
			It("returns an error", func() {
				out := &BadWriter{}
				err := config.WriteMetadata(out)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("could not read from metadata/test-tile-metadata.yml (found inside %s): I am a bad writer", tile.Name())))
			})
		})
	})

	Context("Missing tile", func() {
		It("returns an error", func() {
			config := metadata.Config{}
			err := config.WriteMetadata(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("could not unzip : open : no such file or directory"))
		})
	})

	Context("Invalid tile path", func() {
		It("returns an error", func() {
			config := metadata.Config{
				TileConfig: tileinspect.TileConfig{
					Tile: "this/path/does/not/exist",
				},
			}
			err := config.WriteMetadata(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("could not unzip this/path/does/not/exist: open this/path/does/not/exist: no such file or directory"))
		})
	})

	Context("Invalid tile file", func() {
		BeforeEach(func() {
			var err error
			tile, err = ioutil.TempFile(".", "not-a-zip-file-*.pivotal")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error", func() {
			config := metadata.Config{
				TileConfig: tileinspect.TileConfig{
					Tile: tile.Name(),
				},
			}
			err := config.WriteMetadata(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf("could not unzip %s: zip: not a valid zip file", tile.Name())))
		})
	})

	Context("No metadata file inside tile", func() {
		BeforeEach(func() {
			var err error
			tile, err = CreateTestTileWithMetadata("")
			Expect(err).ToNot(HaveOccurred())
		})
		It("returns an error", func() {
			config := metadata.Config{
				TileConfig: tileinspect.TileConfig{
					Tile: tile.Name(),
				},
			}
			err := config.WriteMetadata(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("metadata file not found"))
		})
	})

	Context("Invalid metadata file inside tile", func() {
		var config metadata.Config
		BeforeEach(func() {
			var err error
			tile, err = CreateTestTileWithMetadata(": - this is not valid yaml")
			Expect(err).ToNot(HaveOccurred())

			config = metadata.Config{
				TileConfig: tileinspect.TileConfig{
					Tile: tile.Name(),
				},
			}
		})
		Context("Default format (yaml)", func() {
			It("returns the invalid yaml", func() {
				err := config.WriteMetadata(buffer)
				Expect(err).ToNot(HaveOccurred())
				Eventually(buffer).Should(Say(": - this is not valid yaml"))
			})
		})

		Context("JSON format", func() {
			It("returns an error", func() {
				config.Format = "json"
				err := config.WriteMetadata(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("could not read from metadata/test-tile-metadata.yml (found inside %s): yaml: did not find expected key", tile.Name())))
			})
		})
	})
})
