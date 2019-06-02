package main

import (
	"github.com/dubbogo/getty"
	"github.com/dubbogo/getty/demo/hello"
	"github.com/dubbogo/getty/demo/hello/tcp"
)

func main() {
	client := getty.NewTCPClient(
		getty.WithServerAddress("127.0.0.1:8090"),
		getty.WithConnectionNumber(2),
	)

	client.RunEventLoop(tcp.NewHelloClientSession)

	go hello.HelloClientRequest()

	hello.WaitCloseSignals(client)
}
