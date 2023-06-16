package main

import (
	"github.com/xyhubl/yim/internal/comet/grpc"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/xyhubl/yim/internal/comet"
	"github.com/xyhubl/yim/internal/comet/conf"
	"github.com/xyhubl/yim/pkg/vipers"
)

func main() {
	// go run main.go -config ./dev/comet.yaml
	config := &conf.Config{}
	vipers.InitViperConf(config, vipers.WithOpenWatching(true), vipers.WithConfigType("yaml"))
	server := comet.NewServer(config)
	// zh: 初始化websocket
	if err := comet.InitWebsocket(server, config.Websocket.Bind, runtime.NumCPU()); err != nil {
		panic(err)
	}
	// zh: rpc
	rpcSrv := grpc.New(config.RpcServer, server)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSrv.GracefulStop()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
