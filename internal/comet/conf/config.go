package conf

import "github.com/xyhubl/yim/pkg/vipers"

type Config struct {
	Base     vipers.Base `mapstructure:"base"`
	TCP      *TCP        `mapstructure:"tcp"`
	Protocol *Protocol   `mapstructure:"protocol"`
	Bucket   *Bucket     `mapstructure:"bucket"`
}

type TCP struct {
	Bind         []string `mapstructure:"bind"`
	SndBuf       int      `mapstructure:"snd_buf"`
	RcvBuf       int      `mapstructure:"rcv_buf"`
	KeepAlive    bool     `mapstructure:"keep_alive"`
	Reader       int      `mapstructure:"reader"`
	ReadBuf      int      `mapstructure:"read_buf"`
	ReadBufSize  int      `mapstructure:"read_buf_size"`
	Writer       int      `mapstructure:"writer"`
	WriteBuf     int      `mapstructure:"write_buf"`
	WriteBufSize int      `mapstructure:"write_buf_size"`
}

type Protocol struct {
	Timer            int
	TimerSize        int
	SvrProto         int
	CliProto         int
	HandshakeTimeout int
}

type Bucket struct {
	Size          int    `mapstructure:"size"`
	Channel       int    `mapstructure:"channel"`
	Room          int    `mapstructure:"room"`
	RoutineAmount uint64 `mapstructure:"routine_amount"`
	RoutineSize   int    `mapstructure:"routine_size"`
}
