package tfd_test

import (
	"fmt"

	"github.com/akatsuki-py/tfd"
)

func ExampleCreateSelectDialog() {
	extensions := []string{"md", "txt"} // select *.md or *.txt file
	hiddenDisplay := false              // if hidden-file displays
	path, err := tfd.CreateSelectDialog(extensions, hiddenDisplay)
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
	// Output: /usr/local/go/README.md
}
