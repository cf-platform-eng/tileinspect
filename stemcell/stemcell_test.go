package stemcell_test

import (
	"encoding/json"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/cf-platform-eng/isv-ci-toolkit/tileinspect/stemcell"
	"github.com/cf-platform-eng/isv-ci-toolkit/tileinspect/stemcell/stemcellfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/pkg/errors"
)

var _ = Describe("WriteStemcell", func() {
	Context("Valid tile", func() {
		var (
			buffer      *Buffer
			metadataCmd *stemcellfakes.FakeMetadataCmd
		)

		BeforeEach(func() {
			buffer = NewBuffer()
			metadataCmd = &stemcellfakes.FakeMetadataCmd{}
		})

		AfterEach(func() {
			err := buffer.Close()
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Floating stemcell version", func() {
			BeforeEach(func() {
				metadataCmd.WriteMetadataStub = func(out io.Writer) error {
					_, err := out.Write([]byte(heredoc.Doc(`{
					  "stemcell_criteria": {
						"os": "ubuntu-xenial",
						"requires_cpi": false,
  						"version": "170"
					  }
					}`)))
					Expect(err).ToNot(HaveOccurred())
					return nil
				}
			})

			It("returns the stemcell requirements", func() {
				config := stemcell.Config{
					MetadataCmd: metadataCmd,
				}
				err := config.WriteStemcell(buffer)
				Expect(err).ToNot(HaveOccurred())

				var stemcellCriteria map[string]interface{}
				err = json.Unmarshal(buffer.Contents(), &stemcellCriteria)
				Expect(err).ToNot(HaveOccurred())

				Expect(stemcellCriteria["os"]).To(Equal("ubuntu-xenial"))
				Expect(stemcellCriteria["requires_cpi"]).To(Equal(false))
				Expect(stemcellCriteria["version"]).To(Equal("170"))
				Expect(stemcellCriteria["floating"]).To(Equal(true))
			})
		})
		Context("Fixed stemcell version", func() {
			BeforeEach(func() {
				metadataCmd.WriteMetadataStub = func(out io.Writer) error {
					_, err := out.Write([]byte(heredoc.Doc(`{
					  "stemcell_criteria": {
						"os": "ubuntu-xenial",
						"requires_cpi": false,
  						"version": "170.1234"
					  }
					}`)))
					Expect(err).ToNot(HaveOccurred())
					return nil
				}
			})

			It("returns the stemcell requirements", func() {
				config := stemcell.Config{
					MetadataCmd: metadataCmd,
				}
				err := config.WriteStemcell(buffer)
				Expect(err).ToNot(HaveOccurred())

				var stemcellCriteria map[string]interface{}
				err = json.Unmarshal(buffer.Contents(), &stemcellCriteria)
				Expect(err).ToNot(HaveOccurred())

				Expect(stemcellCriteria["os"]).To(Equal("ubuntu-xenial"))
				Expect(stemcellCriteria["requires_cpi"]).To(Equal(false))
				Expect(stemcellCriteria["version"]).To(Equal("170.1234"))
				Expect(stemcellCriteria["floating"]).To(Equal(false))
			})
		})

		Context("Failed to get metadata", func() {
			BeforeEach(func() {
				metadataCmd.WriteMetadataReturns(errors.New("write metadata error"))
			})

			It("returns an error", func() {
				config := stemcell.Config{
					MetadataCmd: metadataCmd,
				}
				err := config.WriteStemcell(buffer)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to read tile metadata: write metadata error"))
			})
		})

		Context("Invalid stemcell criteria JSON", func() {
			BeforeEach(func() {
				metadataCmd.WriteMetadataStub = func(out io.Writer) error {
					_, err := out.Write([]byte(heredoc.Doc(`{
					  "stemcell_criteria": "this is not a map"
					}`)))
					Expect(err).ToNot(HaveOccurred())
					return nil
				}
			})

			It("returns an error", func() {
				config := stemcell.Config{
					MetadataCmd: metadataCmd,
				}
				err := config.WriteStemcell(buffer)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to decode stemcell criteria JSON: json: cannot unmarshal string into Go struct field TileMetadata.stemcell_criteria of type map[string]interface {}"))
			})
		})

		Context("Invalid stemcell criteria JSON", func() {
			BeforeEach(func() {
				metadataCmd.WriteMetadataStub = func(out io.Writer) error {
					_, err := out.Write([]byte(heredoc.Doc(`{
					  "stemcell_criteria": "this is not a map"
					}`)))
					Expect(err).ToNot(HaveOccurred())
					return nil
				}
			})

			It("returns an error", func() {
				config := stemcell.Config{
					MetadataCmd: metadataCmd,
				}
				err := config.WriteStemcell(buffer)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to decode stemcell criteria JSON: json: cannot unmarshal string into Go struct field TileMetadata.stemcell_criteria of type map[string]interface {}"))
			})
		})
	})
})
