package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os/exec"
)

var stdLib = stdLibPackages()

func main() {
	flag.Parse()
	for _, a := range flag.Args() {
		if err := ShowImports(a); err != nil {
			log.Fatal(err)
		}
	}
}

// stdlibPackages returns the set of packages in the standard Go library.
// The expansion of "std" is done inside the go tool, so shell out.
func stdLibPackages() map[string]struct{} {
	pkgs := make(map[string]struct{})
	cmd := exec.Command("go", "list", "std")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(stdout)
	for s.Scan() {
		pkgs[s.Text()] = struct{}{}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return pkgs
}

func ShowImports(srcpath string) error {
	ctx := build.Default
	// TODO set GOPATH
	pkg, err := ctx.Import(srcpath, "", 0)
	if err != nil {
		return err
	}
	fmt.Println("Imports:")
	printImports(pkg.Imports)
	fmt.Println("Test imports:")
	printImports(pkg.TestImports)
	return nil
}

func printImports(imports []string) {
	for _, imp := range imports {
		std := ""
		if _, ok := stdLib[imp]; ok {
			std = " (std lib)"
		}
		fmt.Printf(" - %s%s\n", imp, std)
	}
}
