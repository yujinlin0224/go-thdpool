package mul_thd

import (
	"testing"
	"flag"
	"math/rand"
	"time"
	"fmt"
	"sync"
	"bytes"
	"runtime"
	"strconv"
)

// SleepWork is a struct implement Runner with
// work ID and sleep duration of milliseconds.
type sleepWork struct {
	workID, msDuration int
}

// SleepWork.Run run what SleepWork need to do.
func (sw *sleepWork) Run(thdID int, mutex *sync.Mutex) (err error) {
	var gID, _ = getGoroutineID()
	fmt.Printf("thread %d (goroutine id %d) start work %v\n", thdID, gID, *sw)
	time.Sleep(time.Duration(sw.msDuration) * time.Millisecond)
	fmt.Printf("thread %d (goroutine id %d) done work %v\n", thdID, gID, *sw)
	return nil
}

// Test is a function for go test.
func Test(t *testing.T) {
	var thdCnt = flag.Int("t", 10, "Count of thread.")
	var workCnt = flag.Int("w", 100, "Count of work.")
	flag.Parse()
	
	// Make mulThd (type MulThd) and run it.
	var mulThd = New(*thdCnt)
	defer mulThd.Close()
	if errs := mulThd.Run(); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
	
	// Add all works to mulThd.
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < *workCnt; i++ {
		var work = &sleepWork{i, int(rand.Int63n(10000))}
		fmt.Printf("add runner (type sleepWork) %v to mulThd\n", *work)
		mulThd.AddRunner(work)
	}
	
}

// getGoroutineID get goroutine id from the stack of runtime.
func getGoroutineID() (uint64, error) {
	var buf = make([]byte, 64)
	var prefix = []byte("goroutine ")
	buf = bytes.TrimPrefix(buf[:runtime.Stack(buf, false)], prefix)
	return strconv.ParseUint(string(buf[:bytes.IndexByte(buf, ' ')]), 10, 64)
}
