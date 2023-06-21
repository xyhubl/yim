package conf

import "github.com/xyhubl/yim/pkg/vipers"

type Config struct {
	Base       vipers.Base `mapstructure:"base"`
	Redis      *Redis      `mapstructure:"redis"`
	Kafka      *Kafka      `mapstructure:"kafka"`
	Mongo      *Mongo      `mapstructure:"mongo"`
	RPCServer  *RPCServer  `mapstructure:"rpc_server"`
	HttpServer *HTTPServer `mapstructure:"http_server"`
}

type HTTPServer struct {
	Network      string `yaml:"network"`
	Addr         string `yaml:"addr"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type RPCServer struct {
	Network           string `yaml:"network"`
	Addr              string `yaml:"addr"`
	Timeout           int    `yaml:"timeout"`
	IdleTimeout       int    `yaml:"idle_timeout"`
	ForceCloseWait    int    `yaml:"force_close_wait"`
	KeepAliveInterval int    `yaml:"keep_alive_interval"`
	KeepAliveTimeout  int    `yaml:"keep_alive_timeout"`
}

type Redis struct {
	Network      string `yaml:"network"`
	Addr         string `yaml:"addr"`
	UserName     string `yaml:"user_name"`
	Password     string `yaml:"password"`
	Active       int    `yaml:"active"`
	Idle         int    `yaml:"idle"`
	DialTimeout  int    `yaml:"dial_timeout"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	IdleTimeout  int    `yaml:"idle_timeout"`
	Expire       int    `yaml:"expire"`
}

type Kafka struct {
	Topic   string
	Brokers []string
}

type Mongo struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
