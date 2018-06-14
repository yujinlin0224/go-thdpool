package mul_thd

import (
	"sync"
)

// Runner is a interface used in MulThd.
type Runner interface {
	Run(thdID int, mutex *sync.Mutex) error
}

// MulThd can call each Run of all Runner interfaces with multi-threading.
type MulThd struct {
	wg    *sync.WaitGroup
	mutex *sync.Mutex
	rChan chan Runner
	
	ThdCnt int
}

// New make a MulThd with count of threads.
func New(thdCnt int) (mulThd *MulThd) {
	mulThd = new(MulThd)
	mulThd.ThdCnt = thdCnt
	mulThd.wg = new(sync.WaitGroup)
	mulThd.wg.Add(thdCnt)
	mulThd.mutex = new(sync.Mutex)
	mulThd.rChan = make(chan Runner, thdCnt)
	return
}

// AddRunner add the runner to the channel to wait for calling it.
func (mt *MulThd) AddRunner(runner Runner) {
	mt.rChan <- runner
}

// Close close the channel and wait for each Run of all Runner interfaces done.
func (mt *MulThd) Close() {
	defer mt.wg.Wait()
	close(mt.rChan)
}

// Run call each Run of all Runner interfaces with multi-threading.
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
					// When no runner waiting to run, done this thread.
					mt.wg.Done()
					return
				}
			}
		}(i)
	}
	return
}

