// +build linux

/******************************************************
# MAINTAINER : wongoo
# LICENCE    : Apache License 2.0
# EMAIL      : gelnyang@163.com
# MOD        : 2019-06-11
******************************************************/

package getty

import (
	"sync"
	"syscall"
)

import (
	"golang.org/x/sys/unix"
)

type epoll struct {
	EpollOptions

	fd          int
	connections map[int]*session
	lock        *sync.RWMutex
	done        chan struct{}

	tPool *TaskPool
}

func NewEpoller(opts ...EpollOption) (Epoller, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	var options EpollOptions
	for _, opt := range opts {
		opt(&options)
	}

	options.validate()

	ep := &epoll{
		EpollOptions: options,
		fd:           fd,
		lock:         &sync.RWMutex{},
		connections:  make(map[int]*session),
		done:         make(chan struct{}),
		tPool:        CreateTaskPool(options.TaskPoolOptions),
	}

	return ep, nil
}

func (e *epoll) Add(s *session) error {
	conn := s.Conn()
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}

	e.lock.Lock()
	e.connections[fd] = s
	e.lock.Unlock()

	if len(e.connections)%100 == 0 {
		log.Infof("total number of connections: %v", len(e.connections))
	}

	return nil
}

func (e *epoll) Remove(s *session) error {
	conn := s.Conn()
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}

	e.lock.Lock()
	delete(e.connections, fd)
	e.lock.Unlock()

	if len(e.connections)%100 == 0 {
		log.Infof("total number of connections: %v", len(e.connections))
	}
	return nil
}

func (e *epoll) Wait() error {
	events := make([]unix.EpollEvent, e.maxEvents)
	n, err := unix.EpollWait(e.fd, events, e.maxEvents)
	if err != nil {
		return err
	}
	e.lock.RLock()
	for i := 0; i < n; i++ {
		ss := e.connections[int(events[i].Fd)]
		if ss != nil {
			e.tPool.AddTask(task{session: ss})
		}
	}
	e.lock.RUnlock()
	return nil
}

func (e *epoll) Start() {
	go func() {
		for {
			select {
			case <-e.done:
				return
			default:
				if err := e.Wait(); err != nil {
					log.Warnf("failed to epoll wait: %+v", err)
					continue
				}
			}
		}
	}()
}

func (e *epoll) Close() {
	close(e.done)
}
