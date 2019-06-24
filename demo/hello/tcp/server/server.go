/******************************************************
# MAINTAINER : wongoo
# LICENCE    : Apache License 2.0
# EMAIL      : gelnyang@163.com
# MOD        : 2019-06-11
******************************************************/

package main

import (
	"flag"
	"github.com/dubbogo/getty"
	"github.com/dubbogo/getty/demo/hello/tcp"
	"github.com/dubbogo/getty/demo/util"
)

var (
	epollMode            = flag.Bool("epoll", false, "epoll mode")
	epollMaxEvents       = flag.Int("epoll_max_events", 0, "epoll max events")
	epollTaskQueueLength = flag.Int("epoll_task_queue_length", 50, "epoll task queue length")
	epollTaskQueueNumber = flag.Int("epoll_task_queue_number", 2, "epoll task queue number")
	epollTaskPollSize    = flag.Int("epoll_task_pool_size", 1000, "epoll task poll size")

	taskPollMode        = flag.Bool("taskPool", false, "task pool mode")
	taskPollQueueLength = flag.Int("task_queue_length", 100, "task queue length")
	taskPollQueueNumber = flag.Int("task_queue_number", 4, "task queue number")
	taskPollSize        = flag.Int("task_pool_size", 2000, "task poll size")
)

var (
	ep       getty.Epoller
	taskPoll *getty.TaskPool
)

func main() {
	var (
		err error
	)

	flag.Parse()

	util.SetLimit()

	options := []getty.ServerOption{getty.WithLocalAddress(":8090")}
	if *epollMode {
		ep, err = getty.NewEpoller(
			getty.WithEpollMaxEvents(*epollMaxEvents),
			getty.WithEpollTaskQueueLength(*epollTaskQueueLength),
			getty.WithEpollTaskQueueNumber(*epollTaskQueueNumber),
			getty.WithEpollTaskPoolSize(*epollTaskPollSize),
		)
		if err != nil {
			panic(err)
		}
		ep.Start()
	}

	if *taskPollMode {
		taskPoll = getty.NewTaskPool(
			getty.WithTaskPoolTaskQueueLength(*taskPollQueueLength),
			getty.WithTaskPoolTaskQueueNumber(*taskPollQueueNumber),
			getty.WithTaskPoolTaskPoolSize(*taskPollSize),
		)
	}

	server := getty.NewTCPServer(options...)

	go server.RunEventLoop(NewHelloServerSession)

	util.WaitCloseSignals(server)
	if ep != nil {
		ep.Close()
	}
}

func NewHelloServerSession(session getty.Session) (err error) {
	err = tcp.InitialSession(session)
	if err != nil {
		return
	}
	session.SetEpoller(ep)
	session.SetTaskPool(taskPoll)
	return
}
