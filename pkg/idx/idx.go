package idx

import "github.com/yitter/idgenerator-go/idgen"

// zh: 全局seq id 生成器
func NewIdGenera(workId uint16) {
	opt := idgen.NewIdGeneratorOptions(workId)
	idgen.SetIdGenerator(opt)
}

func NextId() int64 {
	return idgen.NextId()
}
