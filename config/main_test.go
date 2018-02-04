// Testing Config file.
package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type Temp struct {
	Name string
}

func (self Temp) Remove() {
	os.RemoveAll(self.Name)
}

func (self Temp) File(name string) string {
	return filepath.Join(self.Name, name)
}

func (self Temp) Missing(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func NewTemp() (*Temp, error) {
	tmp := new(Temp)
	dir, err := ioutil.TempDir("", "Go_Test")
	tmp.Name = dir
	return tmp, err
}

func TestNew(t *testing.T) {
	config1 := NewConfig()
	config2 := NewConfig()
	if config1.ID == config2.ID {
		t.Fail()
	}
}

func TestSave(t *testing.T) {
	dir, err := NewTemp()
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Remove()

	config := NewConfig()
	tmp := dir.File("config.json")
	err = config.Save(tmp)
	if err != nil {
		t.Error(err)
	}
	if dir.Missing(tmp) {
		t.Fail()
	}
}

func TestLoad(t *testing.T) {
	dir, err := NewTemp()
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Remove()

	config := NewConfig()
	tmp := dir.File("config.json")
	err = config.Save(tmp)
	if err != nil {
		t.Error(err)
	}
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
}

func TestUpdate(t *testing.T) {
	dir, err := NewTemp()
	if err != nil {
		t.Fatal(err)
	}
	defer dir.Remove()

	config := NewConfig()
	tmp := dir.File("config.json")
	err = config.Save(tmp)
	if err != nil {
		t.Error(err)
	}
	name := "TEST123"
	config.Name = name
	err = config.Save(tmp)
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
}
