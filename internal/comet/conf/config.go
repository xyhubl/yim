package conf

import "github.com/xyhubl/yim/pkg/vipers"

type Config struct {
	Base      vipers.Base `mapstructure:"base"`
	TCP       *TCP        `mapstructure:"tcp"`
	Protocol  *Protocol   `mapstructure:"protocol"`
	Bucket    *Bucket     `mapstructure:"bucket"`
	Websocket *Websocket  `mapstructure:"websocket"`
	RpcClient *RpcClient  `mapstructure:"rpc_client"`
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
	Timer            int `mapstructure:"timer"`
	TimerSize        int `mapstructure:"timer_size"`
	SvrProto         int `mapstructure:"svr_proto"`
	CliProto         int `mapstructure:"cli_proto"`
	HandshakeTimeout int `mapstructure:"handshake_timeout"`
}

type Bucket struct {
	Size          int    `mapstructure:"size"`
	Channel       int    `mapstructure:"channel"`
	Room          int    `mapstructure:"room"`
	RoutineAmount uint64 `mapstructure:"routine_amount"`
	RoutineSize   int    `mapstructure:"routine_size"`
}

type Websocket struct {
	Bind []string `mapstructure:"bind"`
}

type RpcClient struct {
	Addr    string `mapstructure:"addr"`
	Dial    int    `mapstructure:"dial"`
	Timeout int    `mapstructure:"timeout"`
}
