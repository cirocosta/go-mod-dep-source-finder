package main

import (
	"flag"
	"fmt"
)

var (
	goModFilepath = flag.String("file", "./go.mod", "location of the go.mod file")
)

func main() {
	flag.Parse()
	fmt.Println("vim-go")
}
