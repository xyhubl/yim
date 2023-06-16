package conf

import "github.com/xyhubl/yim/pkg/vipers"

type Config struct {
	Base      vipers.Base
	Kafka     *Kafka     `mapstructure:"kafka"`
	Comet     *Comet     `mapstructure:"comet"`
	RpcClient *RpcClient `mapstructure:"rpc_client"`
}

type Kafka struct {
	Topic   string   `mapstructure:"topic"`
	Group   string   `mapstructure:"group"`
	Brokers []string `mapstructure:"brokers"`
}

type Comet struct {
	RoutineChan int `mapstructure:"routine_chan"`
	RoutineSize int `mapstructure:"routine_size"`
}

type RpcClient struct {
	Addr    string `mapstructure:"addr"`
	Dial    int    `mapstructure:"dial"`
	Timeout int    `mapstructure:"timeout"`
}
