// gourl2pkg is a tool to add or update Go packages in pkgsrc.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var InfoLog *log.Logger

var (
	verbose   = flag.Bool("v", True, "Print verbose messages about what is happening.")
	pkgsrcdir = flag.String("pkgsrc", "", "Path to the top-level pkgsrc directory, will be taken from the PKGSRCDIR environment variable if not given.")
)

func main() {
	if flag.NArg() == 0 {
		Err("Need at least one argument")
	}
	for pkgpath := range flag.Arg() {
		if err := HandleURL(pkgpath); err != nil {
			Err(err)
			os.Exit(1)
		}
	}
}

func HandleURL(pkgpath string) error {
	return nil
}
