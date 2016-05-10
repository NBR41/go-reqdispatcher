package dispatcher

import (
	"errors"
	"runtime"
	"testing"
)

func TestDispatcherOnNew(t *testing.T) {
	var goroutinecount = runtime.NumGoroutine()
	var proc = func(v []interface{}) error {
		return nil
	}
	var batchStep, nbWorkers = 1, 2
	var d = NewDispatcher(batchStep, nbWorkers, proc)

	if goroutinecount+nbWorkers != runtime.NumGoroutine() {
		t.Fatal("no new worker start")
	}

	d.Stop()

	checkThatGoRoutinesAreClosed(t)
}

func TestDispatcherProcessWithError(t *testing.T) {
	var experr = errors.New("test")
	var proc = func(v []interface{}) error {
		p, ok := v[0].(*int)
		if *p == 5 {
			return experr
		}

		if ok {
			*p++
		}
		return nil
	}
	var batchStep, nbWorkers = 1, 2
	var d = NewDispatcher(batchStep, nbWorkers, proc)

	var v1, v2, v3, v4 = 1, 5, 8, 10
	var err = d.Process([]interface{}{&v1, &v2, &v3, &v4})

	if err != experr {
		t.Fatalf("unexpected error: %s", err)
	}

	var closed bool
	if _, closed = <-d.in; closed {
		t.Fatal("in channel not closed")
	}

	if _, closed = <-d.out; closed {
		t.Fatal("in channel not closed")
	}

	checkThatGoRoutinesAreClosed(t)
}

func TestDispatcherProcess(t *testing.T) {
	var proc = func(v []interface{}) error {
		p, ok := v[0].(*int)
		if ok {
			*p++
		}
		return nil
	}
	var batchStep, nbWorkers = 2, 2
	var d = NewDispatcher(batchStep, nbWorkers, proc)

	var v1, v2, v3, v4, v5 = 1, 5, 8, 10, 13
	var err = d.Process([]interface{}{&v1, &v2, &v3, &v4, &v5})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if v1 != 2 || v2 != 5 || v3 != 9 || v4 != 10 || v5 != 14 {
		t.Fatalf("unexpected value: %d %d %d %d %d", v1, v2, v3, v4, v5)
	}

	var closed bool
	if _, closed = <-d.in; closed {
		t.Fatal("in channel not closed")
	}

	if _, closed = <-d.out; closed {
		t.Fatal("in channel not closed")
	}

	checkThatGoRoutinesAreClosed(t)
}
