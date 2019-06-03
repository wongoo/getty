package getty

type pool struct {
	workers   int
	maxTasks  int
	taskQueue chan *session
	done      *chan struct{}
}

func newPool(size int, maxQueue int, doneChan *chan struct{}) *pool {
	return &pool{
		workers:   size,
		maxTasks:  maxQueue,
		taskQueue: make(chan *session, maxQueue),
		done:      doneChan,
	}
}

func (p *pool) Close() {
	close(p.taskQueue)
}

func (p *pool) addTask(ss *session) {
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
		case <-*p.done:
			p.Close()
			return
		case ss := <-p.taskQueue:
			if ss != nil {
				ss.handlePackage()
			}
		}
	}
}
