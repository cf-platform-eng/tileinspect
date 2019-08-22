package checkconfig

import "github.com/cf-platform-eng/tileinspect"

type Config struct {
	tileinspect.TileConfig
	ConfigFile string `long:"config" short:"c" description:"path to config file" required:"true"`
}

func (cmd *Config) Execute(args []string) error {
	return nil
}
