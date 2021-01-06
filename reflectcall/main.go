package main

/*
这是一段关于rpc反射调用结构体所有方法的示例
*/

import (
	"fmt"
	"github.com/micro/go-micro/v2"
	"micro-test/reflectcall/handler"
	"micro-test/reflectcall/proto"
	"reflect"
)

type Person struct {
	Name string
	Age  int
	Sex  string
}

type tool struct {
	cap string
	key string
}

func (t *tool) print() {
	fmt.Println(t.cap, t.key)
}

func (p Person) Say(msg string) {
	fmt.Println("hello,", msg)
}

func (p Person) PrintInfo(t *tool) {
	//t.cap = "green"
	//t.key = "long"
	fmt.Printf("姓名:%s, 年龄:%s, 性别:%s, 参数tool内容:%s %s\n", p.Name, p.Age, p.Sex, t.key, t.cap)
}

/*
servers是方法名与方法的映射
rcvr保存结构体自身的reflect.Value值
typ保存结构体自身的reflect.Type值
*/
type service struct {
	servers map[string]reflect.Method
	rcvr    reflect.Value
	typ     reflect.Type
}

func MakeService(rep interface{}) *service {
	ser := service{}
	ser.typ = reflect.TypeOf(rep)
	ser.rcvr = reflect.ValueOf(rep)
	// name返回其包中的类型名称，举个例子，这里会返回Person，tool
	name := reflect.Indirect(ser.rcvr).Type().Name()
	fmt.Println(name)
	ser.servers = map[string]reflect.Method{}
	fmt.Println(ser.typ.NumMethod(), ser.typ.Name())
	for i := 0; i < ser.typ.NumMethod(); i++ {
		method := ser.typ.Method(i)
		//mtype := method.Type	// reflect.method
		mname := method.Name // string
		fmt.Println("mname : ", mname)
		ser.servers[mname] = method
	}
	return &ser
}

func main1() {
	p1 := Person{"Rbuy", 20, "男"}
	t1 := tool{
		cap: "cap1",
		key: "key1",
	}
	// 得到这个对象的全部方法，string对应reflect.method
	methods := MakeService(p1)
	// 利用得到的methods来调用其值
	methname := "PrintInfo"
	if method, ok := methods.servers[methname]; ok {
		// 得到第一个此method第1参数的Type，第零个当然就是结构体本身了
		/*replyType := method.Type.In(1)
		replyType = replyType.Elem()	// Elem会返回对
		// New returns a Value representing a pointer to a new zero value for the specified type.
		replyv := reflect.New(replyType)*/
		function := method.Func
		function.Call([]reflect.Value{methods.rcvr, reflect.ValueOf(&t1)})
		// 此时我们已经拿到了返回值
		print(methods.rcvr.String())
	}
}

func main() {
	service := micro.NewService(
		micro.Name("com.foo.service.hello"),
		micro.Version("latest"),
		//micro.Server(server.NewServer(server.TLSConfig(config))),
	)

	// Initialise service
	service.Init()

	// Register Handler
	proto.RegisterHelloHandler(service.Server(), new(handler.Hello))
}
