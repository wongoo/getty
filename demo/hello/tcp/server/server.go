package main

import (
	"github.com/dubbogo/getty"
	"github.com/dubbogo/getty/examples/hello"
	"github.com/dubbogo/getty/examples/hello/tcp"
)

func main() {
	server := getty.NewTCPServer(
		getty.WithLocalAddress(":8090"),
	)

	go server.RunEventLoop(tcp.NewHelloServerSession)

	hello.WaitCloseSignals(server)
}
