package main

import (
	"github.com/xyhubl/yim/internal/job"
	"github.com/xyhubl/yim/internal/job/conf"
	"github.com/xyhubl/yim/pkg/vipers"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// go run main.go -config ./dev/job.yaml
	config := &conf.Config{}
	vipers.InitViperConf(config, vipers.WithOpenWatching(true), vipers.WithConfigType("yaml"))
	j := job.New(config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			j.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
