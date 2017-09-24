package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/tools/go/loader"
)

func main() {
	flag.Parse()
	var conf loader.Config
	rest, err := conf.FromArgs(flag.Args(), true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rest)
}
