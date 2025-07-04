//go:build feature
// +build feature

package features_test

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	. "github.com/bunniesandbeatings/goerkin/v2"
	"github.com/cf-platform-eng/tileinspect/features"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("tileinspect check-config", func() {
	steps := NewSteps()

	Describe("Secrets", func() {
		Scenario("Value given for a required property", func() {
			steps.Given("I have a tile with a required secret property")
			steps.And("I have a config file with a secret value of \"secrets!\"")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is valid")
		})

		Scenario("Value missing for a required property", func() {
			steps.Given("I have a tile with a required secret property")
			steps.And("I have an empty config file")
			steps.When("I run tileinspect check-config")
			steps.Then("it says that the secret is missing")
		})

		Scenario("Value empty for a required property", func() {
			steps.Given("I have a tile with a required secret property")
			steps.And("I have a config file with a secret value of \"\"")
			steps.When("I run tileinspect check-config")
			steps.Then("it says that the secret is missing")
		})

		Scenario("Value missing for an optional property", func() {
			steps.Given("I have a tile with an optional secret property")
			steps.And("I have an empty config file")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is valid")
		})

		Scenario("Value empty for an optional property", func() {
			steps.Given("I have a tile with an optional secret property")
			steps.And("I have a config file with a secret value of \"\"")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is valid")
		})

		Scenario("Invalid value", func() {
			steps.Given("I have a tile with a required secret property")
			steps.And("I have a config file with an invalid secret value")
			steps.When("I run tileinspect check-config")
			steps.Then("it says that the secret is invalid")
		})
	})

	Describe("Dropdown Select", func() {
		Scenario("Valid string values", func() {
			steps.Given("I have a tile with a dropdown_select property")
			steps.And("I have a config file with a valid dropdown_select value")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is valid")
		})

		Scenario("Invalid string values", func() {
			steps.Given("I have a tile with a dropdown_select property")
			steps.And("I have a config file with an invalid dropdown_select value")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the dropdown_select value is invalid")
		})

		Scenario("Numeric values", func() {
			steps.Given("I have a tile with a numeric dropdown_select property")
			steps.And("I have a config file with a numeric dropdown_select value")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is valid")
		})
	})

	Describe("Collection", func() {
		Scenario("Valid collection in config", func() {
			steps.Given("I have a tile with a collection property")
			steps.And("I have a config with matching collection")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is valid")
		})
		Scenario("Missing collection", func() {
			steps.Given("I have a tile with a collection property")
			steps.And("I have a config without collection")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file missing collection")
		})
		Scenario("Non-configurable collection items", func() {
			steps.Given("I have a tile with a collection property")
			steps.And("I have a config with a collection with non-configurable value")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file configures unconfigurable collection item")
		})
		Scenario("Default values", func() {
			steps.Given("I have a tile with a collection property and a default value")
			steps.And("I have a config collection without the value")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is valid")
		})
		Scenario("Empty value with optional property", func() {
			steps.Given("I have a tile with an optional collection property")
			steps.And("I have a config with an empty collection value")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is valid")
		})
		Scenario("Empty value with required property", func() {
			steps.Given("I have a tile with a collection property")
			steps.And("I have a config with an empty collection value")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file has an empty value for a required collection")
		})
	})

	Describe("per job blueprints", func() {
		Scenario("Valid per job config", func() {
			steps.Given("I have a tile with a per job properties")
			steps.And("I have a config with per job values")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is valid")
		})
		Scenario("Missing per job config value", func() {
			steps.Given("I have a tile with a per job properties")
			steps.And("I have a config missing per job values")
			steps.When("I run tileinspect check-config")
			steps.Then("it says the config file is missing per job value")
		})
	})

	steps.Define(func(define Definitions) {
		var (
			tile       *os.File
			configFile *os.File
			cmd        *exec.Cmd
			output     string
			exitError  error
		)

		AfterEach(func() {
			if tile != nil {
				err := os.Remove(tile.Name())
				Expect(err).ToNot(HaveOccurred())
			}
			if configFile != nil {
				err := os.Remove(configFile.Name())
				Expect(err).ToNot(HaveOccurred())
			}
		})

		define.Given(`^I have a tile with a required secret property$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
            ---
            name: feature-test-tile
            property_blueprints:
              - name: my-secret
                configurable: true
                type: secret
                optional: false
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a tile with an optional secret property$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
            ---
            name: feature-test-tile
            property_blueprints:
              - name: my-secret
                configurable: true
                type: secret
                optional: true
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have an empty config file$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(`{"product-properties": {}}`)
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config file with a secret value of "(.*)"$`, func(passwordValue string) {
			var err error
			configFile, err = features.MakeConfigFile(fmt.Sprintf(heredoc.Doc(`
            {
              "product-properties": {
                ".properties.my-secret": {
                  "value": {
                    "secret" : "%s"
                  }
                }
              }
            }
            `), passwordValue))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config file with an invalid secret value$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
            {
              "product-properties": {
                ".properties.my-secret": {
                  "value": "secret"
                }
              }
            }`))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a tile with a dropdown_select property$`, func() {

			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
            ---
            name: feature-test-tile
            property_blueprints:
              - name: my-dropdown
                configurable: true
                type: dropdown_select
                optional: false
                options:
                  - label: Low
                    name: low
                  - label: Medium
                    name: medium
                  - label: High
                    name: high
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a tile with a numeric dropdown_select property$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
            ---
            name: feature-test-tile
            property_blueprints:
              - name: my-numeric-dropdown
                configurable: true
                type: dropdown_select
                optional: false
                options:
                  - label: 1
                    name: 1
                  - label: 2
                    name: 2
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a tile with a collection property$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
            ---
            name: feature-test-tile
            property_blueprints:
              - name: my-collection
                configurable: true
                type: collection
                optional: false
                property_blueprints:
                  - name: property-1
                    optional: false
                    type: string
                    configurable: true
                  - name: property-2
                    optional: false
                    type: string
                    configurable: false
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a tile with an optional collection property$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
            ---
            name: feature-test-tile
            property_blueprints:
              - name: my-collection
                configurable: true
                type: collection
                optional: true
                property_blueprints:
                  - name: property-1
                    optional: false
                    type: string
                    configurable: true
                  - name: property-2
                    optional: false
                    type: string
                    configurable: false
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a tile with a collection property and a default value$`, func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
			---
			name: feature-test-tile
			property_blueprints:
			  - name: my-collection
			    configurable: true
			    type: collection
			    optional: false
			    property_blueprints:
			      - name: property-1
			        optional: false
			        type: string
			        configurable: true
			        default: ""
			      - name: property-2
			        optional: false
			        type: string
			        configurable: false
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given("^I have a tile with a per job properties$", func() {
			var err error
			tile, err = features.MakeTileWithMetadata(heredoc.Doc(`
			---
			name: feature-test-tile
			job_types:
			- name: job_type_1
			  property_blueprints:
			  - name: job_1_value1
			    type: string
			    optional: false
			    configurable: true  
			`))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config file with a valid dropdown_select value$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
            {
              "product-properties": {
                ".properties.my-dropdown": {
                  "value": "medium"
                }
              }
            }
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config file with an invalid dropdown_select value$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
            {
              "product-properties": {
                ".properties.my-dropdown": {
                  "value": "this is not a valid value"
                }
              }
            }
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config file with a numeric dropdown_select value$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
            {
              "product-properties": {
                ".properties.my-numeric-dropdown": {
                  "value": 2
                }
              }
            }
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config with matching collection$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
            {
              "product-properties": {
                ".properties.my-collection": {
                  "value": [
                    { 
                        "property-1": "value-1"
                    }
                  ]
                }
              }
            }
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config with a collection with non-configurable value$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
            {
              "product-properties": {
                ".properties.my-collection": {
                  "value": [
                    { 
                      "property-1", "value-1",
                      "property-2": "value-2"
                    }
                  ]
                }
              }
            }
            `))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config without collection$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
            {
              "product-properties": {
			  }
            }`))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config collection without the value$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
            {
              "product-properties": {
                ".properties.my-collection": {
                  "value": [
                    { 
					}
				  ]
				}
		      }
            }`))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config with an empty collection value$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
            {
              "product-properties": {
                ".properties.my-collection": {
                  "value": []
				}
		      }
            }`))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config with per job values$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
			{
			  "product-properties": {
				".job_type_1.job_1_value1": {
				  "value": "value"
				}
			  }
			}
			`))
			Expect(err).ToNot(HaveOccurred())
		})

		define.Given(`^I have a config missing per job values$`, func() {
			var err error
			configFile, err = features.MakeConfigFile(heredoc.Doc(`
			{
			  "product-properties": {
			  }
			}
			`))
			Expect(err).ToNot(HaveOccurred())

		})

		define.When(`^I run tileinspect check-config$`, func() {
			cmd = exec.Command("go", "run", "../cmd/tileinspect/main.go", "check-config", "-c", configFile.Name(), "-t", tile.Name())
			var outputBytes []byte
			outputBytes, exitError = cmd.CombinedOutput()
			output = string(outputBytes)
		})

		define.Then(`^it says the config file is valid$`, func() {
			Expect(output).To(ContainSubstring("The config file appears to be valid"))
			Expect(exitError).ToNot(HaveOccurred())
		})

		define.Then(`^it says that the secret is missing$`, func() {
			Expect(output).To(ContainSubstring("the config file is missing a required property (.properties.my-secret)"))
			Expect(exitError).To(HaveOccurred())
		})

		define.Then(`^it says the config file is missing per job value$`, func() {
			Expect(output).To(ContainSubstring("the config file is missing a required property (.job_type_1.job_1_value1)"))
			Expect(exitError).To(HaveOccurred())
		})

		define.Then(`^it says that the secret is invalid$`, func() {
			Expect(output).To(ContainSubstring(`the config file value for property (.properties.my-secret) is not in the right format. Should be {"secret": "<SECRET VALUE>"}`))
			Expect(exitError).To(HaveOccurred())
		})

		define.Then(`^it says the dropdown_select value is invalid$`, func() {
			Expect(output).To(ContainSubstring(`the config file value for property (.properties.my-dropdown) is invalid: this is not a valid value`))
			Expect(exitError).To(HaveOccurred())
		})

		define.Then(`^it says the config is invalid$`, func() {
			Expect(output).To(ContainSubstring(`the config file value for property (.properties.my-dropdown) is invalid: this is not a valid value`))
			Expect(exitError).To(HaveOccurred())
		})

		define.Then(`^it says the config file missing collection$`, func() {
			Expect(output).To(ContainSubstring(`the config file is missing a required property (.properties.my-collection)`))
			Expect(exitError).To(HaveOccurred())
		})

		define.Then(`^it says the config file configures unconfigurable collection item$`, func() {
			Expect(output).To(ContainSubstring(`collection (.properties.my-collection) contains unconfigurable property property-2`))
			Expect(exitError).To(HaveOccurred())
		})

		define.Then(`^it says the config file has an empty value for a required collection$`, func() {
			Expect(output).To(ContainSubstring(`collection (.properties.my-collection) is missing required property property-1`))
			Expect(exitError).To(HaveOccurred())
		})
	})
})
