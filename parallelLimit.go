package parallelLimit

import (
	"fmt"
	"sync"
)

type errQty struct {
	errors int
	mu     sync.RWMutex
}

func (e *errQty) addError() {
	e.mu.Lock()
	e.errors++
	e.mu.Unlock()
}

func (e *errQty) getErrors() int {
	e.mu.RLock()
	res := e.errors
	e.mu.RUnlock()
	return res
}

func ParallelLimit(funcs []func() error, maxWorkers int, maxErrors int) error {
	e := errQty{}
	wg := sync.WaitGroup{}
	tasksQty := len(funcs)
	limit := make(chan struct{}, maxWorkers)
	for i, f := range funcs {
		limit <- struct{}{}
		wg.Add(1)
		go func(f func() error) {
			if err := f(); err != nil {
				e.addError()
			}
			wg.Done()
			<-limit
		}(f)
		fmt.Println("scheduled task:", i+1, "/", tasksQty, "| errors:", e.getErrors(), "/", maxErrors)
		if e.getErrors() > maxErrors {
			fmt.Println("stop scheduling after task number", i+1)
			break
		}
	}
	wg.Wait()
	if e.getErrors() > maxErrors {
		return fmt.Errorf("too much errors: %d > %d", e.getErrors(), maxErrors)
	}
	return nil
}
