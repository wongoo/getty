package main

import (
	"flag"
	"github.com/dubbogo/getty"
	"github.com/dubbogo/getty/demo/hello"
	"github.com/dubbogo/getty/demo/hello/tcp"
)

var (
	epollMode = flag.Bool("epoll", false, "epoll mode")
)

func main() {
	flag.Parse()

	client := getty.NewTCPClient(
		getty.WithServerAddress("127.0.0.1:8090"),
		getty.WithConnectionNumber(2),
	)

	client.RunEventLoop(tcp.NewHelloClientSession)

	go hello.HelloClientRequest()

	hello.WaitCloseSignals(client)
}
