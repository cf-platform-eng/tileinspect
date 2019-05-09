package tileinspect

type Config struct {
	Debug bool `long:"debug" description:"Outputs more info than usual"`
}

type TileConfig struct {
	Tile string `long:"tile" short:"t" description:"path to product file" required:"true"`
}
