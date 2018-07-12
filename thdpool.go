package thdpool

import (
	"sync"
)

// ThdPool is a thread pool for Worker.
type ThdPool struct {
	wg      *sync.WaitGroup
	mutex   *sync.Mutex
	workers chan Worker
	thdCnt  int
}

// New makes a thread pool with count of threads and a mutex.
// if mutex is nil, it will be newed it.
func New(thdCnt int, mutex *sync.Mutex) *ThdPool {
	if mutex == nil {
		mutex = new(sync.Mutex)
	}
	return &ThdPool{
		wg:      new(sync.WaitGroup),
		mutex:   mutex,
		workers: make(chan Worker, thdCnt),
		thdCnt:  thdCnt,
	}
}

// ThdCnt gets count of threads.
func (tp *ThdPool) ThdCnt() int {
	return tp.thdCnt
}

// AddWorker adds the Worker to the workers channel to wait for calling it.
func (tp *ThdPool) AddWorker(worker Worker) {
	tp.workers <- worker
}

// Close closes the workers channel and wait them done.
func (tp *ThdPool) Close() {
	defer tp.wg.Wait()
	close(tp.workers)
}

// Work runs the thread pool.
func (tp *ThdPool) Run() (errs []error) {
	for i := 0; i < tp.thdCnt; i++ {
		tp.wg.Add(1)
		// Use goroutine to run all threads.
		go func(thdID int) {
			// Get next Worker when last one was done until all done.
			defer tp.wg.Done()
			for worker := range tp.workers {
				if err := worker.Work(thdID, tp.mutex); err != nil {
					errs = append(errs, err)
				}
			}
		}(i)
	}
	return
}
