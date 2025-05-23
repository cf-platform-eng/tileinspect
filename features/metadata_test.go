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
	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("tileinspect metadata", func() {
	steps := NewSteps()

	Scenario("yaml format", func() {
		steps.Given("I have a tile")
		steps.When("I run tileinspect metadata with -f yaml")
		steps.Then("I see the metadata in yaml format")
	})

	Scenario("json format", func() {
		steps.Given("I have a tile")
		steps.When("I run tileinspect metadata with -f json")
		steps.Then("I see the metadata in json format")
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

		define.When(`^I run tileinspect metadata with -f (.+)$`, func(format string) {
			cmd = exec.Command("go", "run", "../cmd/tileinspect/main.go", "metadata", "-f", format, "-t", tile.Name())
			var err error
			output, err = cmd.Output()
			Expect(err).ToNot(HaveOccurred())
		})

		define.Then("I see the metadata in (.+) format$", func(format string) {
			var err error
			value := make(map[string]interface{})
			if format == "json" {
				err = json.Unmarshal(output, &value)
			} else if format == "yaml" {
				err = yaml.Unmarshal(output, &value)
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(value).To(HaveKeyWithValue("name", "feature-test-tile"))
		})
	})
})
