// +build linux

package getty

import (
	"sync"
	"syscall"
)

import (
	"golang.org/x/sys/unix"
)

type epoll struct {
	fd          int
	connections map[int]*session
	lock        *sync.RWMutex
	workPool    *pool
}

func createEpoll() (epoller, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		fd:          fd,
		lock:        &sync.RWMutex{},
		connections: make(map[int]*session),
		workPool:    newPool(100000, 200000),
	}, nil
}

func (e *epoll) add(s *session) error {
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

func (e *epoll) remove(s *session) error {
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

func (e *epoll) wait() ([]*session, error) {
	events := make([]unix.EpollEvent, epollBatchSize)
	n, err := unix.EpollWait(e.fd, events, epollBatchSize)
	if err != nil {
		return nil, err
	}
	e.lock.RLock()
	connections := make([]*session, n)
	for i := 0; i < n; i++ {
		connections[i] = e.connections[int(events[i].Fd)]
	}
	e.lock.RUnlock()
	return connections, nil
}

func (e *epoll) start() {
	go func() {
		for {
			sessions, err := e.wait()
			if err != nil {
				log.Warnf("failed to epoll wait %v", err)
				continue
			}
			for _, ss := range sessions {
				if ss == nil {
					break
				}

				e.workPool.addTask(ss)
			}
		}
	}()
}
