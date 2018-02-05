// Testing Config file.
package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func tempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "Go_Test")
	if err != nil {
		t.Error(err)
	}
	return dir
}

func isMissing(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func TestNew(t *testing.T) {
	config1 := NewConfig()
	config2 := NewConfig()
	if config1.ID == config2.ID {
		t.Fail()
	}
}

func TestSave(t *testing.T) {
	dir := tempDir(t)
	defer os.RemoveAll(dir)

	config := NewConfig()
	tmp := filepath.Join(dir, "config.json")
	err := config.Save(tmp)
	if err != nil {
		t.Error(err)
	}
	if isMissing(tmp) {
		t.Fail()
	}
}

func TestLoad(t *testing.T) {
	dir := tempDir(t)
	defer os.RemoveAll(dir)

	config := NewConfig()
	tmp := filepath.Join(dir, "config.json")
	config.Save(tmp)

	config2, err := LoadConfig(tmp)
	if err != nil {
		t.Error(err)
	}
	if config.ID != config2.ID {
		t.Fail()
	}
	if config2.Root != tmp {
		t.Fail()
	}
	if config2.ID == "" {
		t.Fail()
	}

	_, err = LoadConfig(tmp + "nothere")
	if err == nil {
		t.Fail()
	}

	tmp2 := filepath.Join(dir, "bad.json")
	ioutil.WriteFile(tmp2, []byte("{ this is bad json"), 644)
	_, err = LoadConfig(tmp2)
	if err == nil {
		t.Fail()
	}

	tmp3 := filepath.Join(dir, "incomplete.json")
	ioutil.WriteFile(tmp3, []byte("{\"name\" : \"This is incomplete.\"}"), 644)
	_, err = LoadConfig(tmp3)
	if err == nil {
		t.Fail()
	}

}

func TestUpdate(t *testing.T) {
	dir := tempDir(t)
	defer os.RemoveAll(dir)

	config := NewConfig()
	tmp := filepath.Join(dir, "config.json")
	config.Save(tmp)

	name := "TEST123"
	config.Name = name
	err := config.Save(tmp)
	if err != nil {
		t.Error(err)
	}
	config2, err := LoadConfig(tmp)
	if err != nil {
		t.Error(err)
	}
	if config.Name != config2.Name {
		t.Fail()
	}
	config.ID = "something else"
	err = config.Save(tmp)
	if err == nil {
		t.Fail()
	}
}
