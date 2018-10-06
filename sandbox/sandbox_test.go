package sandbox

import (
	"fmt"
	"testing"
)

func TestSandbox(t *testing.T) {
	sandbox := NewSandBox(t)
	fmt.Println(sandbox)
	defer sandbox.Close()
}
