package context

import (
	"fmt"
	"testing"

	"github.com/internetimagery/photos/sandbox"
)

func TestContext(t *testing.T) {
	sb := sandbox.NewSandBox(t)
	fmt.Println(sb)
}
