package main

import (
	"github.com/xyhubl/yim/internal/comet/conf"
	"github.com/xyhubl/yim/pkg/vipers"
)

func main() {
	config := &conf.Config{}
	vipers.InitViperConf(config, vipers.WithOpenWatching(true), vipers.WithConfigType("yaml"))

}
