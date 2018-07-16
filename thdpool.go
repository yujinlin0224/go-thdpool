package thdpool

import (
	"sync"
)

// ThdPool is a thread pool for Worker.
type ThdPool struct {
	wg      *sync.WaitGroup
	locker  sync.Locker
	workers chan Worker
	thdCnt  int
}

// New makes a thread pool with count of threads and a locker.
// Note that locker cannot be nil.
func New(thdCnt int, locker sync.Locker) *ThdPool {
	if locker == nil {
		panic("thdpool: locker should not be nil")
	}
	return &ThdPool{
		wg:      new(sync.WaitGroup),
		locker:  locker,
		workers: make(chan Worker, thdCnt),
		thdCnt:  thdCnt,
	}
}

// NewWithMutex makes a thread pool include a new locker with sync.Mutex type with count of threads.
func NewWithMutex(thdCnt int) *ThdPool {
	return New(thdCnt, new(sync.Mutex))
}

// NewWithRWMutex makes a thread pool include a new locker with sync.RWMutex type with count of threads.
func NewWithRWMutex(thdCnt int) *ThdPool {
	return New(thdCnt, new(sync.RWMutex))
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

// Run runs the thread pool.
func (tp *ThdPool) Run() (errs []error) {
	for i := 0; i < tp.thdCnt; i++ {
		tp.wg.Add(1)
		// Use goroutine to run all threads.
		go func(thdID int) {
			// Get next Worker when last one was done until all done.
			defer tp.wg.Done()
			for worker := range tp.workers {
				if err := worker.Work(thdID, tp.locker); err != nil {
					errs = append(errs, err)
				}
			}
		}(i)
	}
	return
}
