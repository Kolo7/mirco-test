package httpserver

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
)

func main() {
	// new 一个全新的http复用器
	mux := http.NewServeMux()
	// 为mux设置一个handleFunc
	var buf *bufio.ReadWriter
	var con net.Conn
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		b, _ := ioutil.ReadAll(request.Body)
		// 不让request.Body.Close()有效，这样也就无法断开传输层的连接
		request.Body = ioutil.NopCloser(bytes.NewReader(b))
		// 劫持，获取到运输层的net.Conn，以及读写流(其实也是向net.Con读写)
		hj := writer.(http.Hijacker)
		con, buf, _ = hj.Hijack()
	})

	server := http.Server{
		// 这里接收一个Handler接口，而http复用器实现了该接口
		Handler: mux,
	}
	var l net.Listener
	l, _ = net.Listen("tcp", ":http")

	go func() {
		// http连接内部做了循环等待传输层的消息到来
		_ = server.Serve(l)
	}()

}
