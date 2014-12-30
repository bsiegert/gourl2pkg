package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"unicode"
	"unicode/utf8"
)

// This file contains code for scanning pkgsrc Makefiles and building a reverse
// index of packages and the Go importpaths they represent.
// TODO(bsiegert) eventually make this a separate package.

// Pkg is metadata about a package in pkgsrc.
type Pkg struct {
	// PkgPath is the path to the package (e.g. "lang/go").
	Path string
	// Name is the package name, without the version (e.g. "go").
	Name string
}

func (p *Pkg) String() string {
	return fmt.Sprintf("%s-*:../../%s", p.Name, p.Path)
}

// ReverseIndex maps Go importpaths to packages in pkgsrc.
type ReverseIndex map[string]*Pkg

// FullScan rebuilds the entire index by scanning all files in
func FullScan(basedir string) (ReverseIndex, error) {
	makefiles, err := filepath.Glob(filepath.Join(basedir, "*", "*", "Makefile"))
	if err != nil {
		return nil, err
	}
	r := ReverseIndex{}
	for _, m := range makefiles {
		importpath, ent, err := scanSingle(basedir, m)
		if err != nil {
			return nil, err
		}
		if importpath != "" {
			r[importpath] = ent
		}
	}
	return r, nil
}

// scanSingle scans a single Makefile and returns the import path (or "" if not
// a Go package) and a Pkg record.
func scanSingle(basedir, filename string) (string, *Pkg, error) {
	// TODO(bsiegert) this should re-use a buffer, using bytes.Buffer.ReadFrom.
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil, err
	}
	if !bytes.Contains(contents, []byte("GO_SRCPATH")) {
		// Not a Go package.
		return "", nil, nil
	}
	importpath := extractVar(contents, []byte("GO_SRCPATH"))
	if importpath == "" {
		importpath, err = extractVarMake(filename, "GO_SRCPATH")
		if err != nil {
			return "", nil, fmt.Errorf("getting GO_SRCPATH: %v", err)
		}
	}
	p := Pkg{}
	p.Path, err = filepath.Rel(basedir, filepath.Dir(filename))
	if err != nil {
		return "", nil, err
	}
	p.Name = extractVar(contents, []byte("PKGNAME"))
	if p.Name == "" {
		p.Name, err = extractVarMake(filename, "PKGNAME")
		if err != nil {
			return "", nil, fmt.Errorf("getting PKGNAME: %v", err)
		}
	}
	return importpath, &p, nil
}

// extractVar tries to extract the contents of the variable named varname,
// given the contents of a Makefile in contents. It returns the contents, or
// the empty string if the extraction failed.
func extractVar(contents []byte, varname []byte) string {
	n := bytes.LastIndex(contents, varname)
	if n < 0 {
		return ""
	}
	// There should be whitespace before varname, and a '=' after.
	if r, _ := utf8.DecodeLastRune(contents[:n]); n > 0 && !unicode.IsSpace(r) {
		return ""
	}
	contents = contents[n+len(varname):]
	if contents[0] != '=' {
		return ""
	}
	contents = bytes.TrimSpace(bytes.SplitN(contents[1:], []byte("\n"), 2)[0])
	// If it contains a $ sign, then it is not a simple string.
	if bytes.IndexByte(contents, '$') != -1 {
		return ""
	}
	return string(contents) // Make a copy.
}

// extractVarMake runs bmake on the Makefile to extract the variable name.
func extractVarMake(filename string, varname string) (string, error) {
	cmd := exec.Command("bmake", "show-var", "VARNAME="+varname)
	cmd.Dir = filepath.Dir(filename)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(output)), nil
}
