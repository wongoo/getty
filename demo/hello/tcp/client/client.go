package main

import (
	"flag"
	"github.com/dubbogo/getty"
	"github.com/dubbogo/getty/demo/hello"
	"github.com/dubbogo/getty/demo/hello/tcp"
)

var (
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 1, "number of tcp connections")
)

func main() {
	flag.Parse()
	hello.SetLimit()

	client := getty.NewTCPClient(
		getty.WithServerAddress(*ip+":8090"),
		getty.WithConnectionNumber(*connections),
	)

	client.RunEventLoop(tcp.NewHelloClientSession)

	go hello.HelloClientRequest()

	hello.WaitCloseSignals(client)
}
