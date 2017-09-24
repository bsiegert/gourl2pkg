package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
)

func main() {
	flag.Parse()
	for _, a := range flag.Args() {
		if err := ShowImports(a); err != nil {
			log.Fatal(err)
		}
	}
}

func ShowImports(srcpath string) error {
	ctx := build.Default
	// TODO set GOPATH
	pkg, err := ctx.Import(srcpath, "", 0)
	if err != nil {
		return err
	}
	fmt.Printf("Imports: %v\nTest Imports: %v\n", pkg.Imports, pkg.TestImports)
	return nil
}
