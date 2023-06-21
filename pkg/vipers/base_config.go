package vipers

type Base struct {
	DebugModule string `mapstructure:"debug_module" yaml:"debug_module" toml:"debug_module" json:"debug_module"`
	Host        string `mapstructure:"host" yaml:"host" toml:"host" json:"host"`
	WorkId      uint16 `mapstructure:"work_id" yaml:"work_id" json:"work_id"`
}
