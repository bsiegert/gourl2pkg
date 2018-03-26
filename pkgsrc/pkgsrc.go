/*-
* Copyright (c) 2018
*  Benny Siegert <bsiegert@gmail.com>
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

// Package pkgsrc contains functionality to handle new (Go) packages
// in pkgsrc, the NetBSD package collection.
package pkgsrc

import (
	"bytes"
	"os"
	"path/filepath"
)

// TODO(bsiegert) use PkgMeta in scan.go.

// PkgMeta is the metadata for a pkgsrc package.
type PkgMeta struct {
	// Path is the path to the package (e.g. "lang/go").
	Path string
	// Distname is the package name and version, as used in upstream.
	// For Github repo snapshots, this can be chosen freely, for example
	// with the date.
	Distname string
	// List of categories.
	Categories []string
	// Where to download the sources.
	MasterSites string
	// The top-level source path of the package.
	GoSrcpath string
	// List of all dependencies (test and non-test).
	AllDependencies []string
	// List of test-only dependencies.
	TestDependencies []string
}

func (p PkgMeta) MakefileContents() ([]byte, error) {
	var b bytes.Buffer
	err := makefileTmpl.Execute(&b, p)
	return b.Bytes(), err
}

func (p PkgMeta) CreatePackage(pkgsrcdir string) error {
	dir := filepath.Join(pkgsrcdir, p.Path)
	err := os.Mkdir(dir, 0777)
	if err != nil && !os.IsExist(err) {
		return err
	}

	touch(dir, "DESCR", nil)
	touch(dir, "PLIST", []byte("$NetBSD$"))

	m := filepath.Join(dir, "Makefile")
	c, err := p.MakefileContents()
	if err != nil { return err }
        os.Rename(m, m+".old") // Ignore errors
	return touch(dir, "Makefile", c)
}

// touch is like ioutil.WriteFile, except that it does not overwrite an
// existing file.
func touch(dir, filename string, contents []byte) error {
	fname := filepath.Join(dir, filename)
	if _, err := os.Stat(fname); err == nil {
		// Skip this file.
		return nil
	}
	f, err := os.OpenFile(fname, os.O_WRONLY /* not O_TRUNC */, 0666)
	if err != nil {
		return err
	}
	_, err = f.Write(contents)
	if err != nil {
		return err
	}
	return f.Close()
}
