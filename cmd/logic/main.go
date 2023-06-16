package main

import (
	"fmt"
	"github.com/xyhubl/yim/internal/logic"
	"github.com/xyhubl/yim/internal/logic/conf"
	"github.com/xyhubl/yim/internal/logic/grpc"
	"github.com/xyhubl/yim/internal/logic/http"
	"github.com/xyhubl/yim/pkg/vipers"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// go run main.go -config ./dev/logic.yaml
	config := &conf.Config{}
	vipers.InitViperConf(config, vipers.WithOpenWatching(true), vipers.WithConfigType("yaml"))
	fmt.Println(config.Base.Host, config.Base.DebugModule)
	logic.Login = logic.New(config)
	http.Server(config)
	grpc.NewRpcSrv(config.RPCServer, logic.Login)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			http.Stop()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
