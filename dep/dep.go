/*-
 * Copyright (c) 2017, 2018
 *	Benny Siegert <bsiegert@gmail.com>
 *
 * Provided that these terms and disclaimer and all copyright notices
 * are retained or reproduced in an accompanying document, permission
 * is granted to deal in this work without restriction, including un-
 * limited rights to use, publicly perform, distribute, sell, modify,
 * merge, give away, or sublicence.
 *
 * This work is provided "AS IS" and WITHOUT WARRANTY of any kind, to
 * the utmost extent permitted by applicable law, neither express nor
 * implied; without malicious intent or gross negligence. In no event
 * may a licensor, author or contributor be held liable for indirect,
 * direct, other damage, loss, or other issues arising in any way out
 * of dealing in the work, even if advised of the possibility of such
 * damage or existence of a defect, except proven that it results out
 * of said person's immediate fault when using the work as intended.
 */

package dep

import (
	"bufio"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bsiegert/gourl2pkg/pkgsrc"
)

var stdLib = stdLibPackages()

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

type prefixMatcher interface {
	PrefixMatch(string) (string, bool)
}

// FindImportsRecursive finds dependencies for srcpath and every one of
// its subdirectories and stores them in meta.
func FindImportsRecursive(gopath string, revIndex prefixMatcher, meta *pkgsrc.PkgMeta) error {
	ctx := build.Default
	if gopath != "" {
		ctx.GOPATH = gopath
	}
	basedir := filepath.Join(filepath.SplitList(ctx.GOPATH)[0], "src")
	depends := make(map[string]struct{})
	testDepends := make(map[string]struct{})
	srcpath := meta.GoSrcpath

	err := filepath.Walk(filepath.Join(basedir, srcpath), func(path string, info os.FileInfo, err error) error {
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

		relpath, err := filepath.Rel(basedir, path)
		if err != nil {
			return err
		}
		pkg, _ := ctx.Import(relpath, "", 0)
		for _, d := range pkg.Imports {
			if skipImport(d, revIndex, basedir, srcpath) {
				continue
			}
			depends[d] = struct{}{}
		}
		for _, d := range pkg.TestImports {
			if skipImport(d, revIndex, basedir, srcpath) {
				continue
			}
			if _, ok := depends[d]; ok {
				continue
			}
			testDepends[d] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return err
	}

	meta.Dependencies = pkgsForImports(depends, revIndex)
	meta.TestDependencies = pkgsForImports(testDepends, revIndex)
	return nil
}

func skipImport(dep string, revIndex prefixMatcher, basedir, srcpath string) bool {
	// Depends on another package from the same base.
	if strings.HasPrefix(dep, srcpath) {
		//log.Printf("Self dependency %s -> %s", srcpath, dep)
		return true
	}
	// Depends on a package in the standard library.
	if _, ok := stdLib[dep]; ok {
		return true
	}
	// Vendored dependency.
	for srcpath != "." {
		vendor := filepath.Join(basedir, srcpath, "vendor", dep)
		if _, err := os.Stat(vendor); err == nil {
			// log.Printf("Dependency on vendored package %s", vendor)
			return true
		}
		srcpath = filepath.Dir(srcpath)
	}
	// Unresolved dependency.
	if _, ok := revIndex.PrefixMatch(dep); !ok {
		fmt.Printf("%s (UNRESOLVED)\n", dep)
		return true
	}
	// cgo.
	return dep == "C"
}

func pkgsForImports(imports map[string]struct{}, revIndex prefixMatcher) []string {
	pkgs := make(map[string]struct{})
	for imp := range imports {
		pkg, ok := revIndex.PrefixMatch(imp)
		if !ok {
			fmt.Printf("Unresolved dependency: %s\n", imp)
			continue
		}
		pkgs[pkg] = struct{}{}
	}
	pkgList := []string{}
	for pkg := range pkgs {
		pkgList = append(pkgList, pkg)
	}
	sort.Strings(pkgList)
	return pkgList
}
