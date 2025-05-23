package stemcell_test

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/cf-platform-eng/tileinspect/stemcell"
	"github.com/cf-platform-eng/tileinspect/tileinspectfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/pkg/errors"
)

var _ = Describe("WriteStemcell", func() {
	Context("Valid tile", func() {
		var (
			buffer      *Buffer
			metadataCmd *tileinspectfakes.FakeMetadataCmd
		)

		BeforeEach(func() {
			buffer = NewBuffer()
			metadataCmd = &tileinspectfakes.FakeMetadataCmd{}
		})

		AfterEach(func() {
			err := buffer.Close()
			Expect(err).ToNot(HaveOccurred())
		})

		Context("Floating stemcell version", func() {
			BeforeEach(func() {
				metadataCmd.LoadMetadataStub = func(target interface{}) error {
					err := json.Unmarshal([]byte(heredoc.Doc(`{
					  "stemcell_criteria": {
						"os": "ubuntu-xenial",
						"requires_cpi": false,
  						"version": "170"
					  }
					}`)), &target)
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
				metadataCmd.LoadMetadataStub = func(target interface{}) error {
					err := json.Unmarshal([]byte(heredoc.Doc(`{
					  "stemcell_criteria": {
						"os": "ubuntu-xenial",
						"requires_cpi": false,
  						"version": "170.1234"
					  }
					}`)), &target)
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
				metadataCmd.LoadMetadataReturns(errors.New("write metadata error"))
			})

			It("returns an error", func() {
				config := stemcell.Config{
					MetadataCmd: metadataCmd,
				}
				err := config.WriteStemcell(buffer)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to load tile metadata: write metadata error"))
			})
		})
	})
})
