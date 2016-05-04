package dispatcher

import (
	"runtime"
	"sync"
	"testing"
)

type incCount struct {
	i int
}

func (t *incCount) inc() {
	t.i++
}

func TestWorkerWhenGoRoutineCreated(t *testing.T) {
	var proc = func(v []interface{}) error {
		return nil
	}
	var in = make(chan []interface{}, 1)
	var out = make(chan error, 1)
	var wg = &sync.WaitGroup{}

	w := newWorker(in, out, wg, proc)

	var goroutinecount = runtime.NumGoroutine()

	w.start()
	if goroutinecount+1 != runtime.NumGoroutine() {
		t.Fatal("no new go routine on start")
	}
}

func TestWorkerWhenReceiveRequestApplyFunc(t *testing.T) {

	var proc = func(v []interface{}) error {
		p, ok := v[0].(*incCount)

		if ok {
			p.inc()
		}
		return nil
	}

	var v = &incCount{}
	var in = make(chan []interface{}, 1)
	var out = make(chan error, 1)
	var wg = &sync.WaitGroup{}

	w := newWorker(in, out, wg, proc)
	w.start()
	wg.Add(1)
	in <- []interface{}{v}
	wg.Wait()

	if v.i != 1 {
		t.Fatalf("unexpected increment %d", v.i)
	}

}

/*
func TestWorkerWhenReceiveRequestApplyFuncPutErrorOnErrorChannel(*testing.T) {

}

func TestWorkerOnClosedInChannelStopRoutine(*testing.T) {

}
*/
