// Testing Add command

package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"testing"

	"github.com/corona10/goimagehash"
	"github.com/internetimagery/photos/sandbox"
)

func TestTest(t *testing.T) {
	dir := sandbox.NewSandbox(t)
	defer dir.Close()
	// Check asset is there
	asset := dir.Get("img1.jpg")
	if _, err := os.Stat(asset); os.IsNotExist(err) {
		t.Fail()
	}

	Fhandle, err := os.Open(asset)
	if err != nil {
		t.Fatal(err)
	}
	defer Fhandle.Close()

	img1, err := jpeg.Decode(Fhandle)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := goimagehash.DifferenceHash(img1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("dHash", hash.GetHash(), hash.ToString())

	hash, err = goimagehash.PerceptionHash(img1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pHash", hash.GetHash(), hash.ToString())
}
