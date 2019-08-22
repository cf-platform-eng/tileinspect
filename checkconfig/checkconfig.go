package checkconfig

import (
    "github.com/cf-platform-eng/tileinspect"
    "github.com/ghodss/yaml"
    "github.com/pkg/errors"
    "io"
    "io/ioutil"
    "os"
)

type Config struct {
    tileinspect.TileConfig
    ConfigFilePath string `long:"config" short:"c" description:"path to config file" required:"true"`
}

type ConfigFile struct {
    ProductProperties map[string]interface{} `json:"product-properties"`
}

func (cmd *Config) CheckConfig(out io.Writer) error {
    configFileContents, err := ioutil.ReadFile(cmd.ConfigFilePath)

    if err != nil {
        return errors.Wrap(err, "config file does not exist")
    }

    configFile := &ConfigFile{}
    err = yaml.Unmarshal(configFileContents, configFile)
    if err != nil {
        return errors.Wrap(err, "config file is not valid JSON or YAML")
    }

    _, err = out.Write([]byte("The config file appears to be valid"))

    return nil
}

func (cmd *Config) Execute(args []string) error {
    return cmd.CheckConfig(os.Stdout)
}
