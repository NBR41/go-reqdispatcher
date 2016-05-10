package dispatcher

import (
	"sync"
)

// Not threadsafe

// Dispatcher struct for dispatcher
type Dispatcher struct {
	in        chan []interface{}
	out       chan error
	batchStep int
	workers   []*worker
	wg        sync.WaitGroup
}

// NewDispatcher returns new instance of dispatcher
func NewDispatcher(batchStep, nbWorkers int, proc Processor) *Dispatcher {
	var d = &Dispatcher{
		batchStep: batchStep,

		// channel to send bid chunk to workers
		in: make(chan []interface{}),

		// channel to receive error from worker
		out: make(chan error, nbWorkers),
	}

	// by default 5 workers
	d.workers = make([]*worker, nbWorkers, nbWorkers)
	var i int
	for i = 0; i < nbWorkers; i++ {
		d.workers[i] = newWorker(i, d.in, d.out, &d.wg, proc)
		d.workers[i].start()
	}
	return d
}

// Process process all the request of reqToProcess
// stop if one error is encountered
func (d *Dispatcher) Process(reqToProcess []interface{}) (err error) {
	var i, imax = 0, 0
Loop:
	for {
		select {
		//we receive one error from worker, need to stop
		case err = <-d.out:
			break Loop

		default:
			// iteration after all bids have been processed
			// need to close chanels and return
			if i > 0 && i >= len(reqToProcess) {
				break Loop
			}

			// send chunk of requests to process through channel
			imax = i + d.batchStep
			if imax > len(reqToProcess) {
				imax = len(reqToProcess)
			}
			d.wg.Add(1)
			d.in <- reqToProcess[i:imax]
			i = imax

			// we have sent all the datas, so we wait until all the processes are done
			// in the next iteration if there is an error it will be catch
			// else the for loop will be break
			if i >= len(reqToProcess) {
				d.wg.Wait()
			}
		}
	}

	d.wg.Wait()
	close(d.in)
	close(d.out)
	return
}

func (d *Dispatcher) Stop() {
	close(d.in)
	close(d.out)
}
