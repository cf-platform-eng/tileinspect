//go:build feature
// +build feature

package features_test

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	. "github.com/bunniesandbeatings/goerkin/v2"
	"github.com/cf-platform-eng/tileinspect/features"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("tileinspect stemcell", func() {
	steps := NewSteps()

	Scenario("yaml format", func() {
		steps.Given("I have a tile")
		steps.When("I run tileinspect stemcell")
		steps.Then("I see the stemcell criteria")
	})

	steps.Define(func(define Definitions) {
		var (
			tile   *os.File
			cmd    *exec.Cmd
			output []byte
		)

		define.Given(`^I have a tile$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
			---
			name: feature-test-tile
			stemcell_criteria:
			  os: ubuntu-xenial
			  requires_cpi: false
			  version: '97.32'
			property_blueprints:
			  - name: simple-string
			    configurable: true
			    type: string
			`))
			Expect(err).ToNot(HaveOccurred())
		}, func() {
			err := os.Remove(tile.Name())
			Expect(err).ToNot(HaveOccurred())
		})

		define.When(`^I run tileinspect stemcell$`, func() {
			cmd = exec.Command("go", "run", "../cmd/tileinspect/main.go", "stemcell", "-t", tile.Name())
			var err error
			output, err = cmd.Output()
			Expect(err).ToNot(HaveOccurred())
		})

		define.Then("I see the stemcell criteria$", func() {
			type StemcellCriteria struct {
				Floating    bool   `json:"floating"`
				OS          string `json:"os"`
				RequiresCPI bool   `json:"requires_cpi"`
				Version     string `json:"version"`
			}
			result := &StemcellCriteria{}
			err := json.Unmarshal(output, result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).ToNot(BeNil())
			Expect(result.Floating).To(BeFalse())
			Expect(result.OS).To(Equal("ubuntu-xenial"))
			Expect(result.RequiresCPI).To(BeFalse())
			Expect(result.Version).To(Equal("97.32"))
		})
	})
})
