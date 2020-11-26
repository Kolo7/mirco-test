package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	hello "client/proto/hello"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("com.foo.service.hello.client"),
	)

	helloClient := hello.NewHelloService("com.foo.service.hello", service.Client())
	resp, err := helloClient.Call(context.TODO(), &hello.Request{
		Name: "kuangle",
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Msg)
}
