// +build feature

package features_test

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	. "github.com/bunniesandbeatings/goerkin"
	"github.com/cf-platform-eng/tileinspect/features"
	. "github.com/onsi/ginkgo"
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

		define.Then(`^it says that the secret is invalid$`, func() {
			Expect(output).To(ContainSubstring(`the config file value for property (.properties.my-secret) is not in the right format. Should be {"secret": "<SECRET VALUE>"}`))
			Expect(exitError).To(HaveOccurred())
		})
	})
})
