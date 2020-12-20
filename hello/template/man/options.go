package man

import "time"

/**
这种option注入方式是micro常用的方式，以插拔方式为对象设置参数
*/

type Options struct {
	Sex         string
	Age         int
	Wrapper     []Wrapper
	WorkOptions WorkOptions
}

type WorkOptions struct {
	ManHour     int
	DutyTime    time.Duration
	WorkWrapper []WorkWrapper
}

func NewOptions(opts ...Option) Options {
	opt := Options{
		Sex: DefaultSex,
		Age: DefaultAge,
		WorkOptions: WorkOptions{
			ManHour:  DefaultManHour,
			DutyTime: DefaultDuty,
		},
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Sex(sex string) Option {
	return func(options *Options) {
		options.Sex = sex
	}
}

func Age(age int) Option {
	return func(options *Options) {
		options.Age = age
	}
}

func Wrap(wrapper Wrapper) Option {
	return func(options *Options) {
		options.Wrapper = append(options.Wrapper, wrapper)
	}
}

func WrapWork(wrapper WorkWrapper) Option {
	return func(options *Options) {
		options.WorkOptions.WorkWrapper = append(options.WorkOptions.WorkWrapper, wrapper)
	}
}

func WithManHour(manHour int) WorkOption {
	return func(options *WorkOptions) {
		options.ManHour = manHour
	}
}

func WithDutyTime(dutyTime time.Duration) WorkOption {
	return func(options *WorkOptions) {
		options.DutyTime = dutyTime
	}
}
