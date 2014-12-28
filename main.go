// gourl2pkg is a tool to add or update Go packages in pkgsrc.
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var InfoLog = log.New(ioutil.Discard, "", log.LstdFlags)

var (
	verbose   = flag.Bool("v", true, "Print verbose messages about what is happening.")
	pkgsrcdir = flag.String("pkgsrc", "", "Path to the top-level pkgsrc directory, will be taken from the PKGSRCDIR environment variable if not given.")
)

func main() {
	if *verbose {
		InfoLog = log.New(os.Stderr, "", log.LstdFlags)
	}
	if flag.NArg() == 0 {
		log.Fatal("Need at least one argument")
	}
	for _, pkgpath := range flag.Args() {
		if err := HandleURL(pkgpath); err != nil {
			log.Fatal(err)
		}
	}
}

func HandleURL(pkgpath string) error {
	return nil
}
