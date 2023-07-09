package ebus

// ThreadPool use to execute handler asynchronous
type ThreadPool interface {
	Submit(task func()) error
}

var defaultThreadPool *_threadPool

type _threadPool struct{}

func (p *_threadPool) Submit(task func()) error {
	go task()
	return nil
}
