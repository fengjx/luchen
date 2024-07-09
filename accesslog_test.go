package luchen_test

import (
	"sync"
	"testing"

	"github.com/fengjx/luchen"
)

func TestAccessLog(t *testing.T) {
	l := luchen.NewAccessLog(10, 10, 10)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			l.Print(map[string]any{
				"idx": i,
				"foo": "bar",
			})
			wg.Done()
		}()
	}
	wg.Wait()
}
