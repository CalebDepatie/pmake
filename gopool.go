package main

import "sync"

type GoPool struct {
	maxJobs int
	curJobs int
	jobMut  sync.Mutex
	jobCond *sync.Cond
}

func NewGoPool(max int) *GoPool {
	p := &GoPool{
		maxJobs: max,
		curJobs: 0,
		jobMut:  sync.Mutex{},
	}
	p.jobCond = sync.NewCond(&p.jobMut)
	return p
}

// Blocks until there is room in the pool
func (p *GoPool) StartJob() {
	p.jobMut.Lock()
	for p.curJobs >= p.maxJobs {
		p.jobCond.Wait()
	}
	p.curJobs++
	p.jobMut.Unlock()
}

// Indicates to the go pool that the job as finished
func (p *GoPool) JobDone() {
	p.jobMut.Lock()
	p.curJobs--
	p.jobCond.Signal()
	p.jobMut.Unlock()
}
