package vipers

type Base struct {
	Env
}

type Env struct {
	Host string `mapstructure:"host" yaml:"host" toml:"host" json:"host"`
}
