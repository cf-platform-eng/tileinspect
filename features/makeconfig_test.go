// +build feature

package features_test

import (
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	. "github.com/bunniesandbeatings/goerkin"
	"github.com/cf-platform-eng/tileinspect"
	"github.com/cf-platform-eng/tileinspect/features"
	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("tileinspect make-config", func() {
	steps := NewSteps()

	Scenario("Simple tile properties", func() {
		steps.Given("I have a tile with simple properties")

		steps.When("I run tileinspect make-config")

		steps.Then("I see the template config file")
		steps.And("the config file contains sample values for the properties")
	})

	Scenario("Simple tile properties with defaults", func() {
		steps.Given("I have a tile with simple properties with defaults")

		steps.When("I run tileinspect make-config")

		steps.Then("I see the template config file")
		steps.And("the config file contains values for the properties, using the defaults")
	})

	Scenario("Tile with a selector", func() {
		steps.Given("I have a tile with a selector")

		steps.When("I run tileinspect make-config")

		steps.Then("I see the template config file")
		steps.And("the config file uses the first option for the selector")
	})

	Scenario("Tile with a selector and value override", func() {
		steps.Given("I have a tile with a selector")

		steps.When("I run tileinspect make-config with a custom value selected")

		steps.Then("I see the template config file")
		steps.And("the config file uses the option for the selector that I gave on the cli")
	})


	steps.Define(func(define Definitions) {
		var (
			tile       *os.File
			cmd        *exec.Cmd
			output     []byte
			configFile *tileinspect.ConfigFile
		)

		AfterEach(func() {
			if tile != nil {
				err := os.Remove(tile.Name())
				Expect(err).ToNot(HaveOccurred())
			}
		})

		define.Given(`^I have a tile with simple properties$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
			---
			name: feature-test-tile
			property_blueprints:
			  - name: simple-string
			    configurable: true
			    type: string
			  - name: simple-integer
			    configurable: true
			    type: integer
			`))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Then(`^the config file contains sample values for the properties$`, func() {
			Expect(configFile).ToNot(BeNil())
			Expect(configFile.ProductProperties).ToNot(BeNil())
			Expect(configFile.ProductProperties).To(HaveKey(".properties.simple-string"))
			Expect(configFile.ProductProperties[".properties.simple-string"].Value).To(Equal("SAMPLE_STRING_VALUE"))
			Expect(configFile.ProductProperties).To(HaveKey(".properties.simple-integer"))
			Expect(configFile.ProductProperties[".properties.simple-integer"].Value).To(BeEquivalentTo(0))
		})

		define.Given(`^I have a tile with simple properties with defaults$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
			---
			name: feature-test-tile
			property_blueprints:
			  - name: flavor
			    configurable: true
			    default: vanilla
			    type: string
			  - name: chocolate-sauce
			    configurable: true
			    default: true
			    type: boolean
			`))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Then(`^the config file contains values for the properties, using the defaults$`, func() {
			Expect(configFile).ToNot(BeNil())
			Expect(configFile.ProductProperties).ToNot(BeNil())
			Expect(configFile.ProductProperties).To(HaveKey(".properties.flavor"))
			Expect(configFile.ProductProperties[".properties.flavor"].Value).To(Equal("vanilla"))
			Expect(configFile.ProductProperties).To(HaveKey(".properties.chocolate-sauce"))
			Expect(configFile.ProductProperties[".properties.chocolate-sauce"].Value).To(BeTrue())
		})

		define.Given(`^I have a tile with a selector$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
			---
			name: feature-test-tile
			property_blueprints:
			  - name: continent
			    configurable: true
			    type: selector
			    option_templates:
			      - name: north-america
			        select_value: North America
			        property_blueprints:
			          - name: required-string
			            configurable: true
			            type: string
			      - name: australia
			        select_value: Australia
			        property_blueprints:
			          - name: required-string
			            configurable: true
			            type: string
			`))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Then(`^the config file uses the first option for the selector$`, func() {
			Expect(configFile).ToNot(BeNil())
			Expect(configFile.ProductProperties).ToNot(BeNil())
			Expect(configFile.ProductProperties).To(HaveKey(".properties.continent"))
			Expect(configFile.ProductProperties[".properties.continent"].Value).To(Equal("North America"))
			Expect(configFile.ProductProperties).To(HaveKey(".properties.continent.north-america.required-string"))
			Expect(configFile.ProductProperties[".properties.continent.north-america.required-string"].Value).To(Equal("SAMPLE_STRING_VALUE"))
		})

		define.When(`^I run tileinspect make-config$`, func() {
			cmd = exec.Command("go", "run", "../cmd/tileinspect/main.go", "make-config", "-f", "yaml", "-t", tile.Name())
			var err error
			output, err = cmd.Output()
			Expect(err).ToNot(HaveOccurred())
		})

		define.When(`^I run tileinspect make-config with a custom value selected$`, func() {
			cmd = exec.Command("go", "run", "../cmd/tileinspect/main.go", "make-config", "-f", "yaml", "-t", tile.Name(), "-v", ".properties.continent:Australia")
			var err error
			output, err = cmd.Output()
			Expect(err).ToNot(HaveOccurred())
		})

		define.Then(`^I see the template config file$`, func() {
			err := yaml.Unmarshal(output, &configFile)
			Expect(err).ToNot(HaveOccurred())
		})

		define.Then(`^the config file uses the option for the selector that I gave on the cli$`, func() {
			Expect(configFile).ToNot(BeNil())
			Expect(configFile.ProductProperties).ToNot(BeNil())
			Expect(configFile.ProductProperties).To(HaveKey(".properties.continent"))
			Expect(configFile.ProductProperties[".properties.continent"].Value).To(Equal("Australia"))
			Expect(configFile.ProductProperties).To(HaveKey(".properties.continent.australia.required-string"))
			Expect(configFile.ProductProperties[".properties.continent.australia.required-string"].Value).To(Equal("SAMPLE_STRING_VALUE"))
		})
	})
})
