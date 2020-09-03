package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type (
	Work func()

	Pool struct {
		c      chan Work
		wg     sync.WaitGroup
		ctx    context.Context
		cancel context.CancelFunc
	}
)

// New creates a new worker pool with `workers` parallel worker goroutines, and
// a queue size of 'queueSize'.
func New(workers int, queueSize int) (*Pool, error) {

	if workers < 1 {
		return nil, fmt.Errorf("number of workers (%d) must be greater than 0", workers)
	}
	if queueSize < 1 {
		return nil, fmt.Errorf("size of the work queue (%d) must be greater than 0", queueSize)
	}

	ctx, cancel := context.WithCancel(context.Background())
	p := Pool{
		ctx:    ctx,
		cancel: cancel,
		c:      make(chan Work, queueSize),
	}

	p.spawn(workers)

	return &p, nil
}

func (p *Pool) spawn(workers int) {
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// Add places 'work' on the queue to be processed by one of the workers in the `Pool`.
// If the number of items in the queue exceeds `queueSize` then the Add function blocks,
// providing back pressure.
func (p *Pool) Add(work Work) {
	p.c <- work
}

func (p *Pool) Close() {
	// NB: we explicitly dont close the channel, so that any concurrent calls to Add() wont assert
	p.cancel()
	p.wg.Wait()

	// TODO: since other goroutines can still be submitting work it can be the case that
	// one or more threads are blocking in `Add()`, because the `queueSize` is exceeded.
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case w := <-p.c:
			w()
		case <-p.ctx.Done():
			return
		}
	}

}
