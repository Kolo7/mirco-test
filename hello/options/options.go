package options

/**
这种option注入方式是micro常用的方式，以插拔方式为对象设置参数
*/

type Options struct {
	A string
	B string
	C string
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		A: "default A",
		B: "default B",
		C: "default C",
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func A(a string) Option {
	return func(options *Options) {
		options.A = a
	}
}

func B(a string) Option {
	return func(options *Options) {
		options.A = a
	}
}
