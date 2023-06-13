package conf

type Config struct {
	RPCServer `yaml:"rpc_server"`
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
