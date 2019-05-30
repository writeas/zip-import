package main

import (
	"fmt"
	"os"

	"github.com/writeas/zip-import"
)

func main() {
	err := zipimport.Parse("posts.zip")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse archive: %s\n", err)
		os.Exit(1)
	}
}
