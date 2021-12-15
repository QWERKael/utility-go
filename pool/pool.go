package pool

import "sync"

type Work struct {
	F    func(...interface{}) error
	Args []interface{}
}

type Pool struct {
	Work       chan Work
	Cap        int
	Limiter    chan struct{}
	stopSignal chan struct{}
	wg         sync.WaitGroup
}

func NewPool(Cap int) Pool {
	return Pool{
		Work:       make(chan Work),
		Cap:        Cap,
		Limiter:    make(chan struct{}, Cap),
		stopSignal: make(chan struct{}, Cap),
	}
}

func (p *Pool) Dispatch(f func(...interface{}) error, args ...interface{}) {
	work := Work{
		F:    f,
		Args: args,
	}
	select {
	case p.Work <- work:
	case p.Limiter <- struct{}{}:
		go p.worker(f, args)
	}
}

func (p *Pool) worker(f func(...interface{}) error, args []interface{}) {
	defer func() { <-p.Limiter }()
	var work Work
	for {
		p.wg.Add(1)
		f(args...)
		p.wg.Done()
		select {
		case work = <-p.Work:
			f = work.F
			args = work.Args
		case <-p.stopSignal:
			break
		}
	}
}

func (p *Pool) StopAll() {
	for i := 0; i < p.Cap; i++ {
		p.stopSignal <- struct{}{}
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
