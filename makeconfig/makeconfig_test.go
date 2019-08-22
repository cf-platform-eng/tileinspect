package makeconfig_test

import (
	"github.com/cf-platform-eng/tileinspect/makeconfig/makeconfigfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("WriteEmptyConfig", func() {
	Context("Valid tile", func() {
		var (
			buffer      *Buffer
			metadataCmd *makeconfigfakes.FakeMetadataCmd
		)

		BeforeEach(func() {
			buffer = NewBuffer()
			metadataCmd = &makeconfigfakes.FakeMetadataCmd{}
		})

		AfterEach(func() {
			err := buffer.Close()
			Expect(err).ToNot(HaveOccurred())
		})

	})
})
