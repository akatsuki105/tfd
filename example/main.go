package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/akatsuki-py/tfd"
)

func main() {
	fmt.Println("Please select *.md or *.txt file...")
	time.Sleep(time.Second * 1)

	extensions := []string{"md", "txt"} // select *.md or *.txt file
	hiddenDisplay := false              // if hidden-file displays
	path, err := tfd.CreateSelectDialog(extensions, hiddenDisplay)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
