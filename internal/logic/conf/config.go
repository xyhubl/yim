package conf

type Config struct {
	Redis     *Redis
	RPCServer *RPCServer `yaml:"rpc_server"`
}

type RPCServer struct {
	Network           string
	Addr              string
	Timeout           int
	IdleTimeout       int
	MaxLifeTime       int
	ForceCloseWait    int
	KeepAliveInterval int
	KeepAliveTimeout  int
}

type Redis struct {
	Network      string
	Addr         string
	Auth         string
	Active       int
	Idle         int
	DialTimeout  int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	Expire       int
}
