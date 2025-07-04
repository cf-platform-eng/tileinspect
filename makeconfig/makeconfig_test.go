package makeconfig_test

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/cf-platform-eng/tileinspect/makeconfig"
	"github.com/cf-platform-eng/tileinspect/tileinspectfakes"
	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MakeConfig", func() {
	var (
		cmd         *makeconfig.Config
		metadataCmd *tileinspectfakes.FakeMetadataCmd
	)

	BeforeEach(func() {
		metadataCmd = &tileinspectfakes.FakeMetadataCmd{}
		cmd = &makeconfig.Config{
			MetadataCmd: metadataCmd,
		}
	})

	Describe("product name", func() {
		BeforeEach(func() {
			metadataCmd.LoadMetadataStub = func(target interface{}) error {
				err := yaml.Unmarshal([]byte(heredoc.Doc(`
			---
			name: product
			property_blueprints:
            `)), &target)
				Expect(err).ToNot(HaveOccurred())
				return nil
			}
		})

		It("returns the product name", func() {
			config, err := cmd.MakeConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
			Expect(config.ProductName).To(Equal("product"))
		})
	})

	Describe("string properties", func() {
		BeforeEach(func() {
			metadataCmd.LoadMetadataStub = func(target interface{}) error {
				err := yaml.Unmarshal([]byte(heredoc.Doc(`
			---
			property_blueprints:
			  - name: basic-property
			    type: string
			    configurable: true
			  - name: optional-property
			    type: string
			    configurable: true
			    optional: true
			  - name: non-configurable-property
			    type: string
			    configurable: false
			  - name: property-with-default
			    type: string
			    configurable: true
			    default: awesome
            `)), &target)
				Expect(err).ToNot(HaveOccurred())
				return nil
			}
		})

		It("returns a config with string values", func() {
			config, err := cmd.MakeConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
			Expect(config.ProductProperties).ToNot(BeNil())
			Expect(config.ProductProperties).To(HaveKey(".properties.basic-property"))
			Expect(config.ProductProperties[".properties.basic-property"].Value).To(Equal("SAMPLE_STRING_VALUE"))
			Expect(config.ProductProperties).To(HaveKey(".properties.optional-property"))
			Expect(config.ProductProperties[".properties.optional-property"].Value).To(Equal("SAMPLE_STRING_VALUE"))
			Expect(config.ProductProperties).To(HaveKey(".properties.property-with-default"))
			Expect(config.ProductProperties[".properties.property-with-default"].Value).To(Equal("awesome"))
		})
	})

	Describe("boolean properties", func() {
		BeforeEach(func() {
			metadataCmd.LoadMetadataStub = func(target interface{}) error {
				err := yaml.Unmarshal([]byte(heredoc.Doc(`
			---
			property_blueprints:
			  - name: basic-property
			    type: boolean
			    configurable: true
			  - name: optional-property
			    type: boolean
			    configurable: true
			    optional: true
			  - name: non-configurable-property
			    type: boolean
			    configurable: false
			  - name: property-with-default
			    type: boolean
			    configurable: true
			    default: true
            `)), &target)
				Expect(err).ToNot(HaveOccurred())
				return nil
			}
		})

		It("returns a config with boolean values", func() {
			config, err := cmd.MakeConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
			Expect(config.ProductProperties).ToNot(BeNil())
			Expect(config.ProductProperties).To(HaveKey(".properties.basic-property"))
			Expect(config.ProductProperties[".properties.basic-property"].Value).To(Equal(false))
			Expect(config.ProductProperties).To(HaveKey(".properties.optional-property"))
			Expect(config.ProductProperties[".properties.optional-property"].Value).To(Equal(false))
			Expect(config.ProductProperties).To(HaveKey(".properties.property-with-default"))
			Expect(config.ProductProperties[".properties.property-with-default"].Value).To(Equal(true))
		})
	})

	Describe("integer properties", func() {
		BeforeEach(func() {
			metadataCmd.LoadMetadataStub = func(target interface{}) error {
				err := yaml.Unmarshal([]byte(heredoc.Doc(`
			---
			property_blueprints:
			  - name: basic-property
			    type: integer
			    configurable: true
			  - name: optional-property
			    type: integer
			    configurable: true
			    optional: true
			  - name: non-configurable-property
			    type: integer
			    configurable: false
			  - name: property-with-default
			    type: integer
			    configurable: true
			    default: 123
            `)), &target)
				Expect(err).ToNot(HaveOccurred())
				return nil
			}
		})

		It("returns a config with integer values", func() {
			config, err := cmd.MakeConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
			Expect(config.ProductProperties).ToNot(BeNil())
			Expect(config.ProductProperties).To(HaveKey(".properties.basic-property"))
			Expect(config.ProductProperties[".properties.basic-property"].Value).To(Equal(0))
			Expect(config.ProductProperties).To(HaveKey(".properties.optional-property"))
			Expect(config.ProductProperties[".properties.optional-property"].Value).To(Equal(0))
			Expect(config.ProductProperties).To(HaveKey(".properties.property-with-default"))
			Expect(config.ProductProperties[".properties.property-with-default"].Value).To(BeEquivalentTo(123))
		})
	})

	Describe("secret properties", func() {
		BeforeEach(func() {
			metadataCmd.LoadMetadataStub = func(target interface{}) error {
				err := yaml.Unmarshal([]byte(heredoc.Doc(`
			---
			property_blueprints:
			  - name: basic-property
			    type: secret
			    configurable: true
			  - name: optional-property
			    type: secret
			    configurable: true
			    optional: true
			  - name: non-configurable-property
			    type: secret
			    configurable: false
			  - name: property-with-default
			    type: secret
			    configurable: true
			    default: my-password
            `)), &target)
				Expect(err).ToNot(HaveOccurred())
				return nil
			}
		})

		It("returns a config with secret values", func() {
			config, err := cmd.MakeConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
			Expect(config.ProductProperties).ToNot(BeNil())
			Expect(config.ProductProperties).To(HaveKey(".properties.basic-property"))
			Expect(config.ProductProperties[".properties.basic-property"].Value).To(HaveKeyWithValue("secret", "SAMPLE_SECRET_VALUE"))
			Expect(config.ProductProperties).To(HaveKey(".properties.optional-property"))
			Expect(config.ProductProperties[".properties.optional-property"].Value).To(HaveKeyWithValue("secret", "SAMPLE_SECRET_VALUE"))
			Expect(config.ProductProperties).To(HaveKey(".properties.property-with-default"))
			Expect(config.ProductProperties[".properties.property-with-default"].Value).To(HaveKeyWithValue("secret", "my-password"))
		})
	})

	Describe("dropdown_select properties", func() {
		BeforeEach(func() {
			metadataCmd.LoadMetadataStub = func(target interface{}) error {
				err := yaml.Unmarshal([]byte(heredoc.Doc(`
			---
			property_blueprints:
			  - name: fruit
			    type: dropdown_select
			    default: lime
			    configurable: true
			    options:
			      - name: kiwi
			        label: Kiwi
			      - name: lime
			        label: Lime
			      - name: tomato
			        label: Tomato
			  - name: vegetable
			    type: dropdown_select
			    configurable: true
			    options:
			      - name: onion
			        label: Onion
			      - name: carrot
			        label: Carrot
			      - name: potato
			        label: Potato
            `)), &target)
				Expect(err).ToNot(HaveOccurred())
				return nil
			}
		})

		It("returns a config with selector values", func() {
			config, err := cmd.MakeConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
			Expect(config.ProductProperties).ToNot(BeNil())
			Expect(config.ProductProperties).To(HaveKey(".properties.fruit"))
			Expect(config.ProductProperties[".properties.fruit"].Value).To(Equal("lime"))
			Expect(config.ProductProperties).To(HaveKey(".properties.vegetable"))
			Expect(config.ProductProperties[".properties.vegetable"].Value).To(Equal("onion"))
		})
	})

	Describe("selector properties", func() {
		BeforeEach(func() {
			metadataCmd.LoadMetadataStub = func(target interface{}) error {
				err := yaml.Unmarshal([]byte(heredoc.Doc(`
			---
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
			  - name: browser
			    configurable: true
			    type: selector
			    default: Google Chrome
			    option_templates:
			      - name: explorer
			        select_value: Internet Explorer
			        property_blueprints:
			          - name: required-string
			            configurable: true
			            type: string
			      - name: chrome
			        select_value: Google Chrome
			        property_blueprints:
			          - name: required-string
			            configurable: true
			            type: string
            `)), &target)
				Expect(err).ToNot(HaveOccurred())
				return nil
			}
		})

		It("returns a config with selector values", func() {
			config, err := cmd.MakeConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
			Expect(config.ProductProperties).ToNot(BeNil())
			Expect(config.ProductProperties).To(HaveKey(".properties.continent"))
			Expect(config.ProductProperties[".properties.continent"].Value).To(Equal("North America"))
			Expect(config.ProductProperties).To(HaveKey(".properties.continent.north-america.required-string"))
			Expect(config.ProductProperties[".properties.continent.north-america.required-string"].Value).To(Equal("SAMPLE_STRING_VALUE"))
			Expect(config.ProductProperties).To(HaveKey(".properties.browser"))
			Expect(config.ProductProperties[".properties.browser"].Value).To(Equal("Google Chrome"))
			Expect(config.ProductProperties).To(HaveKey(".properties.browser.chrome.required-string"))
			Expect(config.ProductProperties[".properties.browser.chrome.required-string"].Value).To(Equal("SAMPLE_STRING_VALUE"))
		})
	})

	Describe("value overrides", func() {
		BeforeEach(func() {
			metadataCmd.LoadMetadataStub = func(target interface{}) error {
				err := yaml.Unmarshal([]byte(heredoc.Doc(`
			---
			property_blueprints:
			  - name: basic-property
			    type: string
			    configurable: true
			  - name: password
			    type: secret
			    configurable: true
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
			  - name: browser
			    configurable: true
			    type: selector
			    default: Google Chrome
			    option_templates:
			      - name: explorer
			        select_value: Internet Explorer
			        property_blueprints:
			          - name: required-string
			            configurable: true
			            type: string
			      - name: chrome
			        select_value: Google Chrome
			        property_blueprints:
			          - name: required-string
			            configurable: true
			            type: string
            `)), &target)
				Expect(err).ToNot(HaveOccurred())
				return nil
			}

			cmd.Values = make(map[string]string)
			cmd.Values[".properties.basic-property"] = "yes"
			cmd.Values[".properties.browser"] = "Internet Explorer"
			cmd.Values[".properties.continent"] = "Australia"
			cmd.Values[".properties.password"] = "super-secret"
		})

		It("returns a config with selector values", func() {
			config, err := cmd.MakeConfig()
			Expect(err).ToNot(HaveOccurred())
			Expect(config).ToNot(BeNil())
			Expect(config.ProductProperties).ToNot(BeNil())
			Expect(config.ProductProperties).To(HaveKey(".properties.basic-property"))
			Expect(config.ProductProperties[".properties.basic-property"].Value).To(Equal("yes"))
			Expect(config.ProductProperties).To(HaveKey(".properties.password"))
			Expect(config.ProductProperties[".properties.password"].Value).To(HaveKeyWithValue("secret", "super-secret"))
			Expect(config.ProductProperties).To(HaveKey(".properties.continent"))
			Expect(config.ProductProperties[".properties.continent"].Value).To(Equal("Australia"))
			Expect(config.ProductProperties).To(HaveKey(".properties.continent.australia.required-string"))
			Expect(config.ProductProperties[".properties.continent.australia.required-string"].Value).To(Equal("SAMPLE_STRING_VALUE"))
			Expect(config.ProductProperties).To(HaveKey(".properties.browser"))
			Expect(config.ProductProperties[".properties.browser"].Value).To(Equal("Internet Explorer"))
			Expect(config.ProductProperties).To(HaveKey(".properties.browser.explorer.required-string"))
			Expect(config.ProductProperties[".properties.browser.explorer.required-string"].Value).To(Equal("SAMPLE_STRING_VALUE"))
		})
	})
})
