// gourl2pkg is a tool to add or update Go packages in pkgsrc.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var InfoLog = log.New(ioutil.Discard, "", log.LstdFlags)

var (
	verbose   = flag.Bool("v", true, "Print verbose messages about what is happening.")
	pkgsrcdir = flag.String("pkgsrc", "", "Path to the top-level pkgsrc directory, will be taken from the PKGSRCDIR environment variable if not given.")
)

func init() {
	if *verbose {
		InfoLog = log.New(os.Stderr, "", log.LstdFlags)
	}
	if *pkgsrcdir == "" {
		*pkgsrcdir = os.Getenv("PKGSRC")
	}
	if *pkgsrcdir == "" {
		*pkgsrcdir = "/usr/pkgsrc"
	}
}

func main() {
	/*if flag.NArg() == 0 {
		log.Fatal("Need at least one argument")
	}
	for _, pkgpath := range flag.Args() {
		if err := HandleURL(pkgpath); err != nil {
			log.Fatal(err)
		}
	}*/
	r, err := FullScan(*pkgsrcdir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
}

func HandleURL(pkgpath string) error {
	return nil
}
