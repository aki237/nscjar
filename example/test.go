package main

import (
	"fmt"
	"os"

	"github.com/aki237/nscjar"
)

func main() {
	if len(os.Args) != 2 {
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	jar := nscjar.Parser{}

	c, err := jar.Unmarshal(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, val := range c {
		jar.Marshal(os.Stdout, val)
	}
}
