package checkconfig_test

import (
    "github.com/MakeNowJust/heredoc"
    "github.com/cf-platform-eng/tileinspect/checkconfig"
    "github.com/cf-platform-eng/tileinspect/checkconfig/checkconfigfakes"
    "github.com/ghodss/yaml"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    . "github.com/onsi/gomega/gbytes"
    "io/ioutil"
    "os"
)

func makeConfigFile(contents string) (*os.File, error) {
    configFile, err := ioutil.TempFile("", "config-file")
    if err != nil {
        return nil, err
    }

    _, err = configFile.Write([]byte(contents))
    return configFile, err
}

var _ = Describe("CheckConfig", func() {

    var (
        buffer      *Buffer
        cmd         *checkconfig.Config
        configFile  *os.File
        metadataCmd *checkconfigfakes.FakeMetadataCmd
    )

    BeforeEach(func() {
        buffer = NewBuffer()
        metadataCmd = &checkconfigfakes.FakeMetadataCmd{}
        cmd = &checkconfig.Config{
            MetadataCmd: metadataCmd,
        }

        metadataCmd.LoadMetadataStub = func(target interface{}) error {
            err := yaml.Unmarshal([]byte(heredoc.Doc(`
            ---
            property_blueprints:
              - name: space
                type: string
                configurable: true
            `)), &target)

            Expect(err).ToNot(HaveOccurred())
            return nil
        }
    })

    AfterEach(func() {
        err := buffer.Close()
        Expect(err).ToNot(HaveOccurred())

        if configFile != nil {
            err = os.Remove(configFile.Name())
            Expect(err).ToNot(HaveOccurred())
        }
    })

    Context("config file does not exist", func() {
        BeforeEach(func() {
            cmd.ConfigFilePath = "/this/path/does/not/exist.json"
        })

        It("returns an error", func() {
            err := cmd.CheckConfig(buffer)
            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("config file does not exist"))
        })
    })

    Context("config file is not valid json or yaml", func() {
        BeforeEach(func() {
            var err error
            configFile, err = makeConfigFile("this is not valid anything")
            Expect(err).ToNot(HaveOccurred())

            cmd.ConfigFilePath = configFile.Name()
        })

        It("returns an error", func() {
            err := cmd.CheckConfig(buffer)
            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("config file is not valid JSON or YAML"))
        })
    })

    Context("config file is empty json", func() {
        BeforeEach(func() {
            var err error
            configFile, err = makeConfigFile("{}")
            Expect(err).ToNot(HaveOccurred())

            cmd.ConfigFilePath = configFile.Name()
        })

        It("says that there were not issues found", func() {
            err := cmd.CheckConfig(buffer)
            Expect(err).ToNot(HaveOccurred())
            Eventually(buffer).Should(Say("The config file appears to be valid"))
        })
    })

    Context("json config file", func() {
        BeforeEach(func() {
            var err error
            configFile, err = makeConfigFile(heredoc.Doc(`
			{
			    "product-properties": {
					".properties.space": {
						"type": "string",
						"value": "test-tile-space"
					}
				}
			}
			`))
            Expect(err).ToNot(HaveOccurred())

            cmd.ConfigFilePath = configFile.Name()
        })

        It("says that there were not issues found", func() {
            err := cmd.CheckConfig(buffer)
            Expect(err).ToNot(HaveOccurred())
            Eventually(buffer).Should(Say("The config file appears to be valid"))
        })
    })

    Context("json config file with a property that't not in the tile", func() {
        BeforeEach(func() {
            var err error
            configFile, err = makeConfigFile(heredoc.Doc(`
			{
			    "product-properties": {
					".properties.unknown": {
						"type": "string",
						"value": "this property does not exist in the tile"
					}
				}
			}
			`))
            Expect(err).ToNot(HaveOccurred())

            cmd.ConfigFilePath = configFile.Name()
        })

        It("says that there were not issues found", func() {
            err := cmd.CheckConfig(buffer)
            Expect(err).To(HaveOccurred())
            //Expect(err.Error()).To()
        })
    })

    Context("config file is empty yaml", func() {
        BeforeEach(func() {
            var err error
            configFile, err = makeConfigFile("---")
            Expect(err).ToNot(HaveOccurred())

            cmd.ConfigFilePath = configFile.Name()
        })

        It("says that there were not issues found", func() {
            err := cmd.CheckConfig(buffer)
            Expect(err).ToNot(HaveOccurred())
            Eventually(buffer).Should(Say("The config file appears to be valid"))
        })
    })

    Context("yaml config file", func() {
        BeforeEach(func() {
            var err error
            configFile, err = makeConfigFile(heredoc.Doc(`---
				product-properties:
				  ".properties.space":
				    value: test-tile-space
				    type: string
			`))
            Expect(err).ToNot(HaveOccurred())

            cmd.ConfigFilePath = configFile.Name()
        })

        It("says that there were no issues found", func() {
            err := cmd.CheckConfig(buffer)
            Expect(err).ToNot(HaveOccurred())
            Eventually(buffer).Should(Say("The config file appears to be valid"))
        })
    })
})

var _ = Describe("CompareProperties", func() {
    var (
        checkConfig    *checkconfig.Config
        configFile     *checkconfig.ConfigFile
        tileProperties *checkconfig.TileProperties
    )

    BeforeEach(func() {
        checkConfig = &checkconfig.Config{}
    })

    Context("Tile with no properties", func() {
        BeforeEach(func() {
            tileProperties = &checkconfig.TileProperties{
                PropertyBlueprints: []checkconfig.TileProperty{},
            }
        })

        Context("Empty config file", func() {
            It("should pass", func() {
                configFile = &checkconfig.ConfigFile{}
                errs := checkConfig.CompareProperties(configFile, tileProperties)
                Expect(errs).To(BeEmpty())
            })
        })
        Context("Non-empty config file", func() {
            It("should return an error and print the extra config parameters", func() {
                configFile = &checkconfig.ConfigFile{
                    ProductProperties: map[string]*checkconfig.ConfigFileProperty{
                        ".properties.my-property": {
                            Type:  "string",
                            Value: "hi",
                        },
                        ".properties.my-other-property": {
                            Type:  "string",
                            Value: "hi",
                        },
                    },
                }
                errs := checkConfig.CompareProperties(configFile, tileProperties)
                Expect(errs).To(HaveLen(2))
                Expect(errs[0].Error()).To(ContainSubstring("the config file contains a property (.properties.my-property) that is not defined in the tile"))
                Expect(errs[1].Error()).To(ContainSubstring("the config file contains a property (.properties.my-other-property) that is not defined in the tile"))
            })
        })
    })

    Context("Tile with only simple properties with defaults", func() {
        BeforeEach(func() {
            tileProperties = &checkconfig.TileProperties{
                PropertyBlueprints: []checkconfig.TileProperty{
                    {
                        Name:         "property-one",
                        Type:         "string",
                        Configurable: true,
                        Default:      "on",
                    },
                    {
                        Name:         "property-two",
                        Type:         "string",
                        Configurable: false,
                    },
                    {
                        Name:         "property-three",
                        Type:         "boolean",
                        Configurable: true,
                        Default:      true,
                    },
                    {
                        Name:         "property-four",
                        Type:         "boolean",
                        Configurable: false,
                    },
                },
            }
        })

        Context("Empty config file", func() {
            It("should pass", func() {
                configFile = &checkconfig.ConfigFile{}
                errs := checkConfig.CompareProperties(configFile, tileProperties)
                Expect(errs).To(BeEmpty())
            })
        })

        Context("Config file has properties that don't start with the right prefix", func() {
            It("should return an error and print the bad config parameter", func() {
                configFile = &checkconfig.ConfigFile{
                    ProductProperties: map[string]*checkconfig.ConfigFileProperty{
                        "property-one": {
                            Type:  "string",
                            Value: "hi",
                        },
                    },
                }
                errs := checkConfig.CompareProperties(configFile, tileProperties)
                Expect(errs).To(HaveLen(1))
                Expect(errs[0].Error()).To(ContainSubstring("the config file contains a property (property-one) that does not start with .properties"))
            })
        })

        Context("Valid config file", func() {
            It("should pass", func() {
                configFile = &checkconfig.ConfigFile{
                    ProductProperties: map[string]*checkconfig.ConfigFileProperty{
                        ".properties.property-one": {
                            Type:  "string",
                            Value: "hi",
                        },
                    },
                }
                errs := checkConfig.CompareProperties(configFile, tileProperties)
                Expect(errs).To(BeEmpty())
            })
        })
        Context("Config file overrides a non-configurable parameter", func() {
            It("should return an error and print the bad config parameter", func() {
                configFile = &checkconfig.ConfigFile{
                    ProductProperties: map[string]*checkconfig.ConfigFileProperty{
                        ".properties.property-two": {
                            Type:  "string",
                            Value: "hi",
                        },
                        ".properties.property-four": {
                            Type:  "boolean",
                            Value: true,
                        },
                    },
                }
                errs := checkConfig.CompareProperties(configFile, tileProperties)
                Expect(errs).To(HaveLen(2))
                Expect(errs[0].Error()).To(ContainSubstring("the config file contains a property (.properties.property-two) that is not configurable"))
                Expect(errs[1].Error()).To(ContainSubstring("the config file contains a property (.properties.property-four) that is not configurable"))
            })
        })
    })

    Context("Tile with a selector property", func() {
        BeforeEach(func() {
            tileProperties = &checkconfig.TileProperties{
                PropertyBlueprints: []checkconfig.TileProperty{
                    {
                        Name:         "simple-property",
                        Type:         "string",
                        Configurable: true,
                    },
                    {
                        Name:         "selector-property",
                        Type:         "selector",
                        Configurable: true,
                        ChildProperties: []checkconfig.TileProperties{
                            {
                                Name:        "option-one",
                                SelectValue: "Option One",
                                PropertyBlueprints: []checkconfig.TileProperty{
                                    {
                                        Name:         "option-one-property-one",
                                        Type:         "string",
                                        Configurable: true,
                                    },
                                    {
                                        Name:         "option-one-property-two",
                                        Type:         "string",
                                        Configurable: true,
                                    },
                                },
                            },
                            {
                                Name:        "option-two",
                                SelectValue: "Option Two",
                                PropertyBlueprints: []checkconfig.TileProperty{
                                    {
                                        Name:         "option-two-property-one",
                                        Type:         "boolean",
                                        Configurable: true,
                                    },
                                },
                            },
                        },
                    },
                },
            }
        })

        Context("Valid config file", func() {
            It("should pass", func() {
                configFile = &checkconfig.ConfigFile{
                    ProductProperties: map[string]*checkconfig.ConfigFileProperty{
                        ".properties.simple-property": {
                            Type:  "string",
                            Value: "hi",
                        },
                        ".properties.selector-property": {
                            Type:  "selector",
                            Value: "Option Two",
                        },
                        ".properties.selector-property.option-two.option-two-property-one": {
                            Type:  "boolean",
                            Value: false,
                        },
                    },
                }
                errs := checkConfig.CompareProperties(configFile, tileProperties)
                Expect(errs).To(BeEmpty())
            })
        })
    })
})
