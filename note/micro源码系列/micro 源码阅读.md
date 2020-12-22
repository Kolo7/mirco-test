# micro 源码阅读

这是micro源码阅读后所作的笔记。

## 版本

>go get github.com/micro/micro/v2@2.9.1

## 启动自动生成的模板

自生成模板

>micro new --namespace=com.foo --gopath=false hello

项目结构

- hello
  - handler
    hello.go
  - proto
    - hello
      hello.proto
  - subscriber
    hello.go
  Dockerfile
  generate.go
  go.mod
  main.go
  Makefile
  plugin.go
  README.md
  .gitgnore

启动的方法在此略过不做介绍，直接进入正题，从main.go开始摸清micro的结构。


```go
// main.go

func main(){
    service := micro.NewService(micro.Name("com.foo.service.hello"),micro.Version("latest"),)
    service.Init()
    hello.RegisterHelloHandler(service.Server(), new(handler.Hello))
    micro.RegisterSubscriber("com.foo.service.hello", service.Server(), new(subscriber.Hello))
    if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

代码结构放得稍微有些紧凑，减少篇幅长度。

这一部分启动的是一个service，是server端的代码。

看起来这部分启动做的事情稀松平常，初始化一个`service`,设置了一些参数；随后调用了初始话方法；接下来是注册handler和subscriber；最后主goroutine阻塞在了Run方法。

## 头绪的开始
基本的启动步骤体现了service的重要意义，为了找到一点感觉，先查看下service的结构。

```go
//micro.go
type Service interface {
	Name() string
	Init(...Option)
	Options() Options
	Client() client.Client
	Server() server.Server
	Run() error
	String() string
}

// service.go
type service struct {
    opts Options
    once sync.Once
}

func newService()...
func Init()...
func Run()...
func Server()...
...
```

service结构体包含东西简单明了，Options和一个保证只执行一次的小工具（单例）。

作者的代码风格也算是始终如一，分析起来不算晦涩。

Options和struct的关系

Options的注入方式比较新颖，通过返回注入函数的方式，将注入的操作交给了Option而不是Options。反向注入的好处显而易见，Options事先不需要全部的Option，需要某个Option或者是需要自己定制某个Option时可以加入，可插拔性可扩展性很强。

```go
// Options.go

type Options struct {
    Client client.Client
    Server server.Server
    ...
}

func newOptions(opts ...Option) {
    opt := Options{
        Client: DefaultClient,
        Server: DefaultServer,
        ...
    }
    for _, o := range opts {
        o(&opt)
    }
    return opt
}

func Client(client client.Client) Option {
    return func(opt Options){
        opt.Client = client
    }
}

func Server(server server.Server) Option {
    return func(opt Options) {
        opt.Server = server
    }
}

// micro.go
type Option func(*Options)
```

Options是struct，而Option是语法糖，它是一个注入函数。

一般会怎么写，java转型的我只会setter/getter，这种pojo的方式在这样的框架源码中就很不适用。

## Wrapper链

micro中有很多的功能都是通过包装的方式来实现，结合上面Option的特殊输入方式，能够在需要的时候为service包装好不同的功能实现。

```go
// options.go

func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			o.Client = w[i-1](o.Client)
		}
	}
}

// client/Wrapper.go
type Wrapper func(Client) Client
```

Wrapper同样是一个函数换名的语法糖，不同点在于Wrapper是运用了代理模式的包装器。这里的client.warpper运用的比较简单。下面是一个复杂一点的例子。

```go
// options.go

func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		o.Client.Init(client.WrapCall(w...))
	}
}

// client/wrapper.go
type CallFunc func(ctx context.Context, node *registry.Node, req Request, rsp interface{}, opts CallOptions) error
type CallWrapper func(CallFunc) CallFunc

```

Init方法盲猜也能知道参数是...Option，而且这个Option只不过是client的一个返回Options注入函数的换名函数。



上述的Options、Option、Wrapper运用类似的go特性就是换名。

## 服务的启动

略过一些次要的东西，这篇主要是分析一下micro的启动的主体脉络，直接看`Run`方法。

```go
// service.go

func (s *service) Run() error {
    ...
    if err := s.Start(); err != nil {
		return err
    }
    ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, signalutil.Shutdown()...)
	}

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}

func (s *service) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}
```

Run方法干的事情也不少，省略号的地方包括了一些server.Handler的注册，Profile服务启动（debug的信息，例如追踪），还有默认日志模块的设置。  

随后通过调用Start方法，调用了Server.Start()，也就是真正的服务，在Start之前处理了两种插件方法。

然后线程阻塞等待kill信号或者是Context的Done信号；结束时调用Stop完成收尾工作，而Stop实际上也是调用了Server的Stop方法，由Server完成收尾工作。

整个启动流程至此完成，期间认识到Service的意义，整个可插拔模块化就是针对Service的定制化和默认实现。下一步重点分析Server模块。

## 补充

DefaultServer之类的默认实现并不是那么简单的只指定了一次，而具体的踩坑在下一篇分析当中我也才发现。
