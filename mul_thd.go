package mul_thd

import (
	"sync"
	"fmt"
	"bytes"
	"runtime"
	"strconv"
)

// Runner is a interface can call Run to run Runner with multi-threading.
type Runner interface {
	Run(thdID int, mutex *sync.Mutex) error
}

type MulThd struct {
	wg    *sync.WaitGroup
	mutex *sync.Mutex
	rChan chan Runner
	
	ThdCnt int
}

func New(thdCnt int) (mulThd *MulThd) {
	mulThd = new(MulThd)
	mulThd.ThdCnt = thdCnt
	mulThd.wg = new(sync.WaitGroup)
	mulThd.wg.Add(thdCnt)
	mulThd.mutex = new(sync.Mutex)
	mulThd.rChan = make(chan Runner, thdCnt)
	return
}

func (mt *MulThd) AddRunner(runner Runner) {
	mt.rChan <- runner
}

func (mt *MulThd) Close() {
	defer mt.wg.Wait()
	close(mt.rChan)
}

func (mt *MulThd) Run() (errs []error) {
	for i := 0; i < mt.ThdCnt; i++ {
		// Use goroutine to run all threads.
		go func(thdID int) {
			// Use endless loop to get next runner when last one was done.
			for {
				if runner, ok := <-mt.rChan; ok {
					if err := runner.Run(thdID, mt.mutex); err != nil {
						errs = append(errs, err)
					}
				} else {
					// When no runner waiting for run, done this thread.
					fmt.Printf("thread %d done.\n", thdID)
					mt.wg.Done()
					return
				}
			}
		}(i)
	}
	return
}

// GetGoroutineID get goroutine id from the stack of runtime.
func GetGoroutineID() (uint64, error) {
	var buf = make([]byte, 64)
	var prefix = []byte("goroutine ")
	buf = bytes.TrimPrefix(buf[:runtime.Stack(buf, false)], prefix)
	return strconv.ParseUint(string(buf[:bytes.IndexByte(buf, ' ')]), 10, 64)
}
