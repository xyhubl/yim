package logic

import (
	"github.com/xyhubl/yim/internal/logic/conf"
	"github.com/xyhubl/yim/internal/logic/dao"
)

type Logic struct {
	c       *conf.Config
	DaoBase *dao.Base
}

func New(c *conf.Config) (l *Logic) {
	l = &Logic{
		c:       c,
		DaoBase: dao.New(c),
	}
	return l
}
