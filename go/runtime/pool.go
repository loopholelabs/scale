package runtime

import "sync"

type Pool struct {
	pool sync.Pool
	New  func() (*Instantiated, error)
}

func NewPool(new func() (*Instantiated, error)) *Pool {
	return &Pool{
		New: new,
	}
}

func (p *Pool) Put(i *Instantiated) {
	if i != nil {
		p.pool.Put(i)
	}
}

func (p *Pool) Get() (*Instantiated, error) {
	rv, ok := p.pool.Get().(*Instantiated)
	if ok && rv != nil {
		return rv, nil
	}
	return p.New()
}
