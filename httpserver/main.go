package httpserver

import (
	"net"
	"net/http"
)

func main() {
	// new 一个全新的http复用器
	mux := http.NewServeMux()
	// 为mux设置一个handleFunc
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

	})

	server := http.Server{
		// 这里接收一个Handler接口，而http复用器实现了该接口
		Handler: mux,
	}

	l, _ := net.Listen("tcp", ":http")

	_ = server.Serve(l)
	_ = server.ListenAndServe()
}
