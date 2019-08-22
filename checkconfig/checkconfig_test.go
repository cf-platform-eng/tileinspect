package checkconfig_test

import (
    "github.com/MakeNowJust/heredoc"
    "github.com/cf-platform-eng/tileinspect/checkconfig"
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
        buffer     *Buffer
        cmd        *checkconfig.Config
        configFile *os.File
    )

    BeforeEach(func() {
        buffer = NewBuffer()
        cmd = &checkconfig.Config{}
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
