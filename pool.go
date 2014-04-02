package gotil

import (
	"container/list"
	"sync"
	)

type Pool struct {
	capacity int
	cons func() interface{}
	ready *list.List
	lock *sync.Mutex
}

func New(capacity int, cons func() interface{}) *Pool {
	ret = new(Pool)
	ret.capacity = capacity
	ret.cons = cons
	ret.ready = list.New()
	ret.ready.Init()
	ret.lock = new(sync.Mutex)
	return ret
}

func (p *Pool) Get() interface{} {
	p.lock.Lock()
	e := ready.Front()
	if e != nil {
		ready.Remove(e)
		p.lock.Unlock()
		return e.Value
	}
	p.lock.Unlock()

	// else create a new one and return that
	ret := p.cons()
	return n
}

func (p *Pool) Put(thing interface{}) {
	p.lock.Lock()
	if p.ready.Len() < p.capacity {
		p.ready.PushFront(thing)
	}
	p.lock.Unlock()
}