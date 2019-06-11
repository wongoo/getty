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
	hello.SetLimit()

	options := []getty.ServerOption{getty.WithLocalAddress(":8090")}
	if *epollMode {
		options = append(options, getty.WithEpollMode())
	}

	server := getty.NewTCPServer(options...)

	go server.RunEventLoop(tcp.NewHelloServerSession)

	hello.WaitCloseSignals(server)
}
