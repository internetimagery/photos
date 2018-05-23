package utility

import (
	"testing"
)

func TestRoot(t *testing.T) {
	paths := []string{"/one/two/three/four", "/one/two/three/five", "/one/two/six/five"}
	if Root(paths) != "/one/two" {
		t.Fail()
	}
	paths = []string{"/one/two/three/four", "/six/two/three/five", "/one/two/six/five"}
	if Root(paths) != "/" {
		t.Fail()
	}
	paths = []string{"/one/two/three/four"}
	if Root(paths) != "/one/two/three/four" {
		t.Fail()
	}
	paths = []string{}
	if Root(paths) != "" {
		t.Fail()
	}

}
