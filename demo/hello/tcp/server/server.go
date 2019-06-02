package main

import (
	"github.com/dubbogo/getty"
	"github.com/dubbogo/getty/demo/hello"
	"github.com/dubbogo/getty/demo/hello/tcp"
)

func main() {
	server := getty.NewTCPServer(
		getty.WithLocalAddress(":8090"),
	)

	go server.RunEventLoop(tcp.NewHelloServerSession)

	hello.WaitCloseSignals(server)
}
