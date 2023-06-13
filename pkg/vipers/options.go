package vipers

type options struct {
	// 开启监视
	openWatching bool
	// 配置文件类型
	configType string
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithOpenWatching(boolean bool) Option {
	return optionFunc(func(o *options) {
		o.openWatching = boolean
	})
}

func WithConfigType(typ string) Option {
	return optionFunc(func(o *options) {
		o.configType = typ
	})
}
