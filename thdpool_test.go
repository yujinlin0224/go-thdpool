package thdpool

import (
	"bytes"
	"flag"
	`log`
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

// sleepWork is a struct implement Worker with work ID and sleep duration of milliseconds.
type sleepWork struct {
	workID, msDuration int
}

// Work does something need to do in thread pool.
func (sw *sleepWork) Work(thdID int, mutex sync.Locker) (err error) {
	var gID uint64
	if gID, err = getGoroutineID(); err != nil {
		return
	}
	log.Printf("thread %d (goroutine id %d) start work %v\n", thdID, gID, *sw)
	time.Sleep(time.Duration(sw.msDuration) * time.Millisecond)
	log.Printf("thread %d (goroutine id %d) done work %v\n", thdID, gID, *sw)
	return
}

// TestThdPool is a function for go test to test thdpool package.
func TestThdPool(t *testing.T) {
	var thdCnt = flag.Int("t", 100, "Count of threads.")
	var workCnt = flag.Int("w", 1000, "Count of works.")
	flag.Parse()
	
	// Make a new thread pool and run it.
	var thdPool = NewWithMutex(*thdCnt)
	defer thdPool.Close()
	if errs := thdPool.Run(); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
	
	// Add all works to the thread pool.
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < *workCnt; i++ {
		var work = &sleepWork{i, int(rand.Int63n(10000))}
		log.Printf("add worker (type sleepWork) %v to thread pool\n", *work)
		thdPool.AddWorker(work)
	}
	
}

// getGoroutineID gets goroutine id from the stack of runtime.
func getGoroutineID() (uint64, error) {
	var buf, prefix = make([]byte, 64), []byte("goroutine ")
	buf = bytes.TrimPrefix(buf[:runtime.Stack(buf, false)], prefix)
	return strconv.ParseUint(string(buf[:bytes.IndexByte(buf, ' ')]), 10, 64)
}
