package testutil

import (
	"fmt"
	"os"
	"testing"
)

func TestTempDir(t *testing.T) {
	tmpDir := NewTempDir(t, "TestTempDir")
	if _, err := os.Stat(tmpDir.Dir); err != nil {
		fmt.Println(err)
		t.Fail()
	}
	tmpDir.Close()
	if _, err := os.Stat(tmpDir.Dir); !os.IsNotExist(err) {
		fmt.Println("Tempdir not removed.")
		t.Fail()
	}

}

func TestUserInput(t *testing.T) {
	testMessage := "Hello"
	defer UserInput(t, testMessage+"\n")()

	resultMessage := ""
	_, err := fmt.Scanln(&resultMessage)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if resultMessage != testMessage {
		fmt.Println("Expected", testMessage)
		fmt.Println("Got", resultMessage)
		t.Fail()
	}

}