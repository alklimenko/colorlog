package rotator

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRotator(t *testing.T) {
	r := NewBuilder().WithStrategy(StrategySize).WithCount(3).WithSize(125).WithCheckPeriod(time.Millisecond * 10).Build()
	wg := &sync.WaitGroup{}
	wg.Add(10)
	for g := 0; g < 10; g++ {
		go func() {
			for i := 0; i < 10000; i++ {
				r.Write([]byte(fmt.Sprintf("%d\n", i)))
				time.Sleep(time.Millisecond * 10)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
