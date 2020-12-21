package man

import (
	"context"
	"hello/template/man/fakeman"
	"time"
)

/**
定义模块的主要接口
*/
type Man interface {
	Init(...Options)
	Eat(string)
	Motion()
	Think() string
	Play(Pet)
	Work(ctx context.Context, req string, msg string, opt WorkOptions) string
}

/**
可以定义一些辅助接口
*/

type Pet interface {
	Sound() string
	Fawn()
}

/**
给fun(*Options)换名，Option就是成为了Options的作用器
*/
type Option func(*Options)

type WorkOption func(*WorkOptions)

/**
提供一些默认参数值
*/
var (
	DefaultSex     = "male"
	DefaultAge     = 18
	DefaultManHour = 100
	DefaultDuty    = time.Hour * 8
	DefaultMan     = fakeman.NewFakeMan()

	NewMan func(...Option) Man = fakeman.NewFakeMan
)
