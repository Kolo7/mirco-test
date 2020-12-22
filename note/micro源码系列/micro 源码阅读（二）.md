# micro 源码阅读（二）

这部分重心落在server上，目的就是要弄明白micro抽象出来的server包含了一些什么，以及搞清楚对grpc的运用设计。

## Start

```go
// service.go
// 在service.Start()真正主要调用的方法
server.Start()

// options.go
// server早在调用NewService时，options构造使用了默认实现
func newOptions(opts ...Option) Options {
    opt := Options{
        ...
        Server: server.DefaultServer,
    }
}

// server/server.go
var (
    DefaultServer   Server = newRpcServer()
)

// server/rpc_server.go
func newRpcServer(opts ...Option) Server {
	options := newOptions(opts...)
	...
}
```

总体脉络延续上一节的结束，service.Start()真正的是调用了opts.Server.Start()，opts在`micro.NewService`方法中被构造，而构造opts时主动的使用了`server.DefaultServer`，逐步追踪可以看到Server的默认实现。

从上面也可以有一些了解，DefaultServer是单例。

遇到了熟悉的结构，newOptions中都是注入默认实现，不再跟进查看；和service不同的是在server中要多出来了Router的构造，目前也不了解，先放一边。

```
// server/rpc_server.go
func (s *rpcServer) Start() error {
    ...
    // start listening on the transport
    ts, err := config.Transport.Listen(config.Address)
    if err != nil {
        return err
    }
}
```
调用了Transport的Listen方法，而我很好奇这块的实现，这已经是对传输层的抽象了，进一步追踪过去。

## Transport

### 传输层的配置
```go
// transport/http_transport.go
func (h *httpTransport) Listen(addr string, opts ...ListenOption) (Listener, error) {
    ...
    var l net.Listener
    ...
    // TODO: support use of listen options
    if h.opts.Secure || h.opts.TLSConfig != nil {
        config := h.opts.TLSConfig
        fn := func(addr string) (net.Listener, error) {
            if config == nil{
            ...
            // generate a certificate
            cert, err := mls.Certificate(hosts...)
            ...
            config = &tls.Config{Certificates: []tls.Certificate{cert}}
            }
            return tls.Listen("tcp", addr, config)
        }

        l, err = mnet.Listen(addr, fn)
    } else {
    	fn := func(addr string) (net.Listener, error) {
    		return net.Listen("tcp", addr)
        }
        
    	l, err = mnet.Listen(addr, fn)
    }
}
```

源码在这块调用了Listen方法，它实际上是对服务启动的封装，返回值是一个封装对象`Listener`，这个对象持有一个字段`listener`，是golang的SDK的官方net.listener，持有它就是持有一个http连接。

DefaultTransport是httpTransport，追踪到httpTransport的Listen实现可以看到主要工作是设置一个启动http服务的lister。而这块又考虑到了是否采用TLS方式，如果启动了TLS方式却没有配置文件，甚至有自生成密钥和证书。其实这很坑，如果我配置没配好，却发现TLS启用了，在没有文档的情况下，可能排除问题会花一个很长的时间。

TODO:这里有关一个知识点，就是读入证书密钥文件，启动http监听，暂时做个记录，不进一步去分析。

[tls服务启动](https://colobu.com/2016/06/07/simple-golang-tls-examples/)

这部分代码有关设置tls配置，瞎摸索了一下，还需要官方默认实现验证一下。

### 阻塞发生

走到这已经初见底层，而阻塞方法的调用似乎就发生在这里面`mnet.Listen(addr, fn)`。

```go
// util/net/net.go
func Listen(addr string, fn func(string) (net.Listener, error)) (net.Listener, error) {
    if strings.Count(addr, ":") == 1 && strings.Count(addr, "-") == 0 {
		return fn(addr)
    }
    ...
    // range the ports
	for port := min; port <= max; port++ {
		// try bind to host:port
		ln, err := fn(HostPort(host, port))
		if err == nil {
			return ln, nil
		}

		// hit max port
		if port == max {
			return nil, err
		}
	}
}
```

这是一个工具函数，主要职能就是多样定制化的去监听端口，上面是最基本的`localhost:port`，查看代码可以发现还支持:`localhost:minport-maxport`。它将会从minport开始尝试，直到成功启动或者是尝试maxport也失败了。

## DefaultServer的多次设定
在我以为DefaultServer就是rpcServer的时候，我一路追踪过去，发现httpTransport中的断点并没有生效，我返回前面查看，发现DefaultServer确实有一瞬间是rpcServer，但在golang有一个不太友好地设计`init()`，而micro源码作者居然在一处初始化，一处又对DefaultServer覆盖赋值了。

[son of a bitch](https://baike.baidu.com/item/%E8%8E%AB%E7%94%9F%E6%B0%94/7381700?fr=aladdin)
```go
// defaults.go
func init() {
    // default client
    client.DefaultClient = gcli.NewClient()
    // default server
    server.DefaultServer = gsrv.NewServer()
    // default store
    store.DefaultStore = memoryStore.NewStore()
    // set default trace
    trace.DefaultTracer = memTrace.NewTracer()
}
```

## grpcServer的实现

前面的分析只是将Transport替换成了http方式，实际上也是合理的设计，不过这里再对grpc来次全新的分析。

### 额外的Option

grpcServer实现了Server接口，在原Options的基础上，多了一些参数，而这些多出来的参数并不是采用组合或者是继承的方式加上去，而是作为额外option，保存到server.Options.Context当中。

```go
// server/options.go
type Options struct{
    ...
    Context context.Context
}

// server/grpc/options.go
type codecsKey struct{}
...

func Codec(contentType string, c encoding.Codec) server.Option {
    return func(o *server.Options) {
    	codecs := make(map[string]encoding.Codec)
    	if o.Context == nil {
    		o.Context = context.Background()
    	}
    	if v, ok := o.Context.Value(codecsKey{}).(map[string]encoding.Codec); ok && v != nil {
    		codecs = v
    	}
    	codecs[contentType] = c
    	o.Context = context.WithValue(o.Context, codecsKey{}, codecs)
    }
}
```