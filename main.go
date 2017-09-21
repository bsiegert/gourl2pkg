// gourl2pkg is a tool to add or update Go packages in pkgsrc.
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var InfoLog = log.New(ioutil.Discard, "", log.LstdFlags)

var (
	index     = flag.Bool("index", false, "Print the reverse index of Go packages instead of adding any ports.")
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
	flag.Parse()
	if flag.NArg() == 0 && !*index {
		log.Fatal("Need at least one argument")
	}

	revIndex, err := FullScan(*pkgsrcdir)
	if err != nil {
		log.Fatal(err)
	}
	if *index {
		revIndex.WriteTo(os.Stdout)
		return
	}

	ToPackage := flag.Args()
	for len(ToPackage) > 0 {
		InfoLog.Printf("Remaining to package: %s", ToPackage)
		p := ToPackage[0]
		ToPackage = ToPackage[1:]
		if err := HandleURL(revIndex, p); err != nil {
			log.Fatal(err)
		}
	}
}

func HandleURL(r ReverseIndex, srcpath string) error {
	// Not a simple map lookup because of prefix matches
	for src, pkg := range r {
		if strings.HasPrefix(srcpath, src) {
			log.Printf("%s is already part of a pkgsrc package (%s)", srcpath, pkg.Path)
			return nil
		}
	}
	log.Printf("This is where we would package %s", srcpath)
	return nil
}
