package testutils

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

const (
	sleepInterval = 500 * time.Millisecond
)

// ConcurrentCases returns an array of integer
// which is used for number of cores when running
// concurrent benchmarks
func ConcurrentCases() []int {
	procPow := int(math.Log2(float64(runtime.NumCPU())))
	cases := make([]int, 0, procPow+1)
	cases = append(cases, 1)
	for proc := 0; proc < procPow; proc++ {
		cases = append(cases, 2<<proc)
	}
	return cases
}

type CheckFunc func() (bool, error)

func Wait(timeout time.Duration, f CheckFunc) error {
	start := time.Now()
	for time.Now().Add(-timeout).Before(start) {
		ok, err := f()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		time.Sleep(sleepInterval)
	}
	return fmt.Errorf("timeout")
}
