package dispatcher

import (
	"sync"
)

type worker struct {
	ID   int
	in   chan []interface{}
	out  chan error
	wg   *sync.WaitGroup
	proc Processor
}

// Processor type for processor
type Processor func([]interface{}) error

func newWorker(ID int, in chan []interface{}, out chan error, wg *sync.WaitGroup, proc Processor) *worker {
	return &worker{ID, in, out, wg, proc}
}

func (w *worker) start() {
	go func() {
		for {
			requests, ok := <-w.in

			// the channel is close need to terminate the routine
			if !ok {
				return
			}

			// call adwords
			err := w.proc(requests)
			if err != nil {
				w.out <- err
				w.wg.Done()
				return
			}

			w.wg.Done()
		}
	}()
}
