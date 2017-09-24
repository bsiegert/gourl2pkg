package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var stdLib = stdLibPackages()

func main() {
	flag.Parse()
	for _, a := range flag.Args() {
		if err := ShowImportsRecursive(a); err != nil {
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

// ShowImportsRecursive prints dependencies for srcpath and every one of
// its subdirectories.
func ShowImportsRecursive(srcpath string) error {
	ctx := build.Default
	basedir := filepath.Join(filepath.SplitList(ctx.GOPATH)[0], "src")
	return filepath.Walk(filepath.Join(basedir, srcpath), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") || base == "testdata" {
			return filepath.SkipDir
		}
		// TODO(bsiegert) what should be the behavior for "vendor"
		// and "internal" trees be?

		srcpath, err := filepath.Rel(basedir, path)
		if err != nil {
			return err
		}
		return ShowImports(srcpath)
	})
}

func ShowImports(srcpath string) error {
	ctx := build.Default
	// TODO set GOPATH
	pkg, err := ctx.Import(srcpath, "", 0)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Imports of %s:\n", srcpath)
	printImports(pkg.Imports)
	fmt.Printf("Test imports of %s:\n", srcpath)
	printImports(pkg.TestImports)
	return nil
}

func printImports(imports []string) {
	for _, imp := range imports {
		if _, ok := stdLib[imp]; !ok {
			fmt.Printf(" - %s\n", imp)
		}
	}
}
