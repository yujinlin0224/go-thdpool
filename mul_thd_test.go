package mul_thd

import (
	"testing"
	"flag"
	"math/rand"
	"time"
	"fmt"
	"os"
	"sync"
)

// SleepWork is a struct implement Runner with
// work ID and sleep duration of milliseconds.
type sleepWork struct {
	workID, msDuration int
}

// SleepWork.Run run what need SleepWork to do.
func (sw *sleepWork) Run(thdID int, mutex *sync.Mutex) (err error) {
	var gID, _ = GetGoroutineID()
	fmt.Printf("thread %d (goroutine id %d) start work %v\n", thdID, gID, sw)
	time.Sleep(time.Duration(sw.msDuration) * time.Millisecond)
	fmt.Printf("thread %d (goroutine id %d) done work %v\n", thdID, gID, sw)
	return nil
}

func Test(t *testing.T) {
	var thdCnt = flag.Int("t", 10, "Count of thread.")
	var workCnt = flag.Int("w", 100, "Count of work.")
	flag.Parse()
	
	// Make a slice of SleepWork with random sleep duration.
	rand.Seed(time.Now().UnixNano())
	var mulThd = New(*thdCnt)
	defer mulThd.Close()
	
	// Run all SleepWork in works with multi-threading.
	if errs := mulThd.Run(); len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	
	for i := 0; i < *workCnt; i++ {
		var work = &sleepWork{i, int(rand.Int63n(10000))}
		fmt.Printf("make runner (type work) %v\n", work)
		mulThd.AddRunner(work)
	}
	
}
