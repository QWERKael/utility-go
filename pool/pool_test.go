package pool

import (
	"fmt"
	"testing"
)

func TestCheckPath(t *testing.T) {
	p := NewPool(3)
	for i := 0; i < 10; i++ {
		p.Dispatch(func(x ...interface{}) error {
			fmt.Printf("#[%d]: %v\n", i, x)
			return nil
		}, i, i*10)
	}
	p.Wait()
}
