package main

import (
	"flag"
	"math/rand"
	"fmt"
	"time"
	"sync"
	"runtime"
	"bytes"
	"strconv"
	"os"
	"reflect"
	"errors"
)

// Runner is a interface can call Run to run Runner with multi-threading.
type Runner interface {
	Run(thdID int) error
}

// GetRunners convert any array or slice type to Runner slice with error.
func GetRunners(origValues interface{}) (runners []Runner, err error) {
	switch reflect.TypeOf(origValues).Kind() {
	case reflect.Array, reflect.Slice:
		var values = reflect.ValueOf(origValues)
		for i := 0; i < values.Len(); i++ {
			var val = values.Index(i)
			if val.Type().Implements(reflect.TypeOf((*Runner)(nil)).Elem()) {
				runners = append(runners, val.Interface().(Runner))
			}
		}
		if values.Len() == 0 && len(runners) == 0 {
			err = errors.New("origValues in GetRunners " +
				"do not have any values implement Runners")
		}
	default:
		err = errors.New("origValues in GetRunners is not a array or slice")
	}
	return
}

// RunRunnersMulThd call all Run in Runner slice with multi-threading.
func RunRunnersMulThd(threadLen int, runners []Runner) (errs []error) {
	var wg sync.WaitGroup
	wg.Add(threadLen)
	defer wg.Wait()
	
	var runnersChan = make(chan Runner, threadLen)
	defer close(runnersChan)
	
	for i := 0; i < threadLen; i++ {
		// Use goroutine to run all threads.
		go func(thdID int) {
			// Use endless loop to get next runner when last one was done.
			for {
				if runner, ok := <-runnersChan; ok {
					if err := runner.Run(thdID); err != nil {
						errs = append(errs, err)
					}
				} else {
					// When no runner waiting for run, done this thread.
					wg.Done()
					return
				}
			}
		}(i)
	}
	
	// Add runner to channel of runners.
	for _, runner := range runners {
		fmt.Printf("add runner %v\n", runner)
		runnersChan <- runner
	}
	
	return
}

// SleepWork is a struct implement Runner with
// work ID and sleep duration of milliseconds.
type SleepWork struct {
	WorkID, MsDuration int
}

// SleepWork.Run run what need SleepWork to do.
func (sw SleepWork) Run(thdID int) (err error) {
	var gID, _ = GetGoroutineID()
	fmt.Printf("thread %d (goroutine id %d) start work %v\n", thdID, gID, sw)
	time.Sleep(time.Duration(sw.MsDuration) * time.Millisecond)
	fmt.Printf("thread %d (goroutine id %d) done work %v\n", thdID, gID, sw)
	return nil
}

func main() {
	var threadLen = flag.Int("t", 0, "Length of thread.")
	var workLen = flag.Int("w", 0, "Length of work.")
	flag.Parse()
	
	// Make a slice of SleepWork with random sleep duration.
	rand.Seed(time.Now().UnixNano())
	var works = make([]*SleepWork, *workLen)
	for i := 0; i < *workLen; i++ {
		var work = &SleepWork{i, int(rand.Int63n(10000))}
		fmt.Printf("make runner (type work) %v\n", work)
		works[i] = work
	}
	
	// Run all SleepWork in works with multi-threading.
	if runners, err := GetRunners(works); err != nil {
		panic(err)
	} else {
		if errs := RunRunnersMulThd(*threadLen, runners); len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}
	
}

// GetGoroutineID get goroutine id from the stack of runtime.
func GetGoroutineID() (uint64, error) {
	var buf = make([]byte, 64)
	var prefix = []byte("goroutine ")
	buf = bytes.TrimPrefix(buf[:runtime.Stack(buf, false)], prefix)
	return strconv.ParseUint(string(buf[:bytes.IndexByte(buf, ' ')]), 10, 64)
}
