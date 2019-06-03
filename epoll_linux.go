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
	done        *chan struct{}
}

func createEpoll(doneChan *chan struct{}) (epoller, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		fd:          fd,
		lock:        &sync.RWMutex{},
		connections: make(map[int]*session),
		workPool:    newPool(100000, 200000, doneChan),
		done:        doneChan,
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

func (e *epoll) wait() error {
	events := make([]unix.EpollEvent, epollBatchSize)
	n, err := unix.EpollWait(e.fd, events, epollBatchSize)
	if err != nil {
		return err
	}
	e.lock.RLock()
	for i := 0; i < n; i++ {
		ss := e.connections[int(events[i].Fd)]
		if ss != nil {
			e.workPool.addTask(ss)
		}
	}
	e.lock.RUnlock()
	return nil
}

func (e *epoll) start() {
	go func() {
		for {
			select {
			case <-*e.done:
				return
			default:
				if err := e.wait(); err != nil {
					log.Warnf("failed to epoll wait %v", err)
					continue
				}
			}
		}
	}()
}
