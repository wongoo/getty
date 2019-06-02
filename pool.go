package getty

import (
	"sync"
)

type pool struct {
	mu        sync.Mutex
	workers   int
	maxTasks  int
	taskQueue chan *session
	closed    bool
	done      chan struct{}
}

func newPool(size int, maxQueue int) *pool {
	return &pool{
		workers:   size,
		maxTasks:  maxQueue,
		taskQueue: make(chan *session, maxQueue),
		done:      make(chan struct{}),
	}
}

func (p *pool) Close() {
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

func (p *pool) addTask(ss *session) {
	if p.closed {
		return
	}
	p.taskQueue <- ss
}

func (p *pool) start() {
	for i := 0; i < p.workers; i++ {
		go p.startWorker()
	}
}

func (p *pool) startWorker() {
	for {
		select {
		case <-p.done:
			return
		case ss := <-p.taskQueue:
			if ss != nil {
				handleSession(ss)
			}
		}
	}
}

func handleSession(ss *session) {
	ss.handlePackage()
}
