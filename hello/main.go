package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"hello/handler"
	"hello/subscriber"

	hello "hello/proto/hello"
)

func main() {
	// New Service
	/*cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}*/
	service := micro.NewService(
		micro.Name("com.foo.service.hello"),
		micro.Version("latest"),
		//micro.Server(server.NewServer(server.TLSConfig(config))),
	)

	// Initialise service
	service.Init()

	// Register Handler
	hello.RegisterHelloHandler(service.Server(), new(handler.Hello))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("com.foo.service.hello", service.Server(), new(subscriber.Hello))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
