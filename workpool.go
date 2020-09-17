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

// New creates a new worker pool with `workers` parallel goroutines to act as workers. Work is
// posted to a queue of maximum size 'queueSize'. The worker pool is valid until a call to `Close()`.
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

// Add places 'work' on the queue to be processed by one of the workers in the `Pool`.
// If the number of items in the queue exceeds `queueSize` then the Add function blocks,
// providing back-pressure. To limit blocking, set the `queueSize` to an appropriately high
// number.
func (p *Pool) Add(work Work) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to submit work to the pool : %s", r)
		}
	}()

	// never schedule `nil` work as it is meaningless and we use `nil` to signify that the chan is closed
	if work != nil {
		p.c <- work
	}
	return nil
}

// Complete prevents any further jobs being queued and waits to complete all queued work
// Following a call to Close(), any calls to Add() will fail with an error.
func (p *Pool) Complete() {
	close(p.c)
	p.wg.Wait()
}

// Cancel prevents any further jobs being queued and completes all work currently being executed. It cancels all other work which has not been started.
// Following a call to Close(), any calls to Add() will fail with an error.
func (p *Pool) Cancel() {
	p.cancel()
	p.wg.Wait()

	// since other goroutines could still be calling `Add()`, we want to inform them that the
	// pool is closed. The simplest way to do this is to close the channel. This also will release any
	// goroutines that are currently blocking because the chan is full.
	close(p.c)
}

// Close calls Cancel() and is here primarily to support io.Closer()
func (p *Pool) Close() error {
	p.Cancel()
	return nil
}

func (p *Pool) spawn(workers int) {
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		// prefer the cancellable context so we dont drain the queue
		select {
		case <-p.ctx.Done():
			return
		default:
		}

		select {
		case w := <-p.c:
			if w == nil {
				// nil signifies close of the channel, as it is the zero value for func type
				return
			}
			w()
		case <-p.ctx.Done():
			return
		}
	}
}
