package thd_pool

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

// sleepWork is a struct implement Worker with
// work ID and sleep duration of milliseconds.
type sleepWork struct {
	workID, msDuration int
}

// sleepWork.Work does something need to do in thread pool.
func (sw *sleepWork) Work(thdID int, mutex *sync.Mutex) (err error) {
	var gID uint64
	if gID, err = getGoroutineID(); err != nil {
		return
	}
	fmt.Printf("thread %d (goroutine id %d) start work %v\n", thdID, gID, *sw)
	time.Sleep(time.Duration(sw.msDuration) * time.Millisecond)
	fmt.Printf("thread %d (goroutine id %d) done work %v\n", thdID, gID, *sw)
	return
}

// Test is a function for go test.
func Test(t *testing.T) {
	var thdCnt = flag.Int("t", 100, "Count of threads.")
	var workCnt = flag.Int("w", 1000, "Count of works.")
	flag.Parse()
	
	// Make thdPool (type ThdPool) and run it.
	var thdPool = New(*thdCnt)
	defer thdPool.Close()
	if errs := thdPool.Run(); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
	
	// Add all works to thdPool.
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < *workCnt; i++ {
		var work = &sleepWork{i, int(rand.Int63n(10000))}
		fmt.Printf("add worker (type sleepWork) %v to thread pool\n", *work)
		thdPool.AddWorker(work)
	}
	
}

// getGoroutineID gets goroutine id from the stack of runtime.
func getGoroutineID() (uint64, error) {
	var buf, prefix = make([]byte, 64), []byte("goroutine ")
	buf = bytes.TrimPrefix(buf[:runtime.Stack(buf, false)], prefix)
	return strconv.ParseUint(string(buf[:bytes.IndexByte(buf, ' ')]), 10, 64)
}
