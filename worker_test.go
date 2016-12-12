package goreqdispatcher

import (
	_ "errors"
	"runtime"
	_ "sync"
	"testing"
	"time"
)

/*
func TestWorkerWhenGoRoutineCreated(t *testing.T) {
	var proc = func(v []interface{}) error {
		return nil
	}
	var in = make(chan []interface{}, 1)
	var out = make(chan error, 1)
	var wg = &sync.WaitGroup{}

	w := newWorker(1, in, out, wg, proc)

	var goroutinecount = runtime.NumGoroutine()
	w.start()
	if goroutinecount+1 != runtime.NumGoroutine() {
		t.Fatalf("no new go routine on start, current go routine count %d", runtime.NumGoroutine())
	}
	close(in)
	checkThatGoRoutinesAreClosed(t)
}

func TestWorkerWhenReceiveRequestApplyFunc(t *testing.T) {
	var proc = func(v []interface{}) error {
		p, ok := v[0].(*int)

		if ok {
			*p++
		}
		return nil
	}

	var v = 0
	var in = make(chan []interface{}, 1)
	var out = make(chan error, 1)
	var wg = &sync.WaitGroup{}

	w := newWorker(1, in, out, wg, proc)
	w.start()
	wg.Add(1)
	in <- []interface{}{&v}
	wg.Wait()

	if v != 1 {
		t.Fatalf("unexpected increment %d", v)
	}
	close(in)
	checkThatGoRoutinesAreClosed(t)
}

func TestWorkerWhenReceiveRequestApplyFuncPutErrorOnErrorChannel(t *testing.T) {
	var experr = errors.New("test")
	var proc = func(v []interface{}) error {
		return experr
	}

	var v = 1
	var in = make(chan []interface{}, 1)
	var out = make(chan error, 1)
	var wg = &sync.WaitGroup{}

	w := newWorker(1, in, out, wg, proc)
	w.start()
	wg.Add(1)
	in <- []interface{}{v}
	wg.Wait()

	err := <-out
	if err != experr {
		t.Fatalf("unexpected error: %s", err)
	}
	close(in)
	checkThatGoRoutinesAreClosed(t)
}
*/
var initialGoRoutineCount = runtime.NumGoroutine() + 1

func checkThatGoRoutinesAreClosed(t *testing.T) {
	var i int
	for {

		if initialGoRoutineCount == runtime.NumGoroutine() {
			t.Logf("waiting stop %d == %d", initialGoRoutineCount, runtime.NumGoroutine())
			break
		}

		t.Logf("waiting %d goroutines", runtime.NumGoroutine())
		if i == 10 {
			t.Fatal("all go routines are not closed")
		}
		i++
		time.Sleep(5 * time.Millisecond)
	}
}
