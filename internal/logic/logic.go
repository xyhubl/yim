package logic

import (
	"github.com/xyhubl/yim/internal/logic/conf"
	"github.com/xyhubl/yim/internal/logic/dao"
	"github.com/xyhubl/yim/pkg/idx"
	"github.com/xyhubl/yim/pkg/log"
)

var Login *Logic

type Logic struct {
	C *conf.Config
}

func New(c *conf.Config) (l *Logic) {
	idx.NewIdGenera(c.Base.WorkId)
	dao.New(c)
	log.InitLog("info")
	l = &Logic{
		C: c,
	}
	return l
}
