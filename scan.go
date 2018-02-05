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

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"unicode"
	"unicode/utf8"
)

// This file contains code for scanning pkgsrc Makefiles and building a reverse
// index of packages and the Go importpaths they represent.
// TODO(bsiegert) eventually make this a separate package.

var makeProg = func() string {
	if _, err := exec.LookPath("bmake"); err == nil {
		return "bmake"
	}
	return "make"
}()

// Pkg is metadata about a package in pkgsrc.
type Pkg struct {
	// PkgPath is the path to the package (e.g. "lang/go").
	Path string
}

func (p *Pkg) String() string {
	return p.Path
}

// ReverseIndex maps Go importpaths to packages in pkgsrc.
type ReverseIndex map[string]*Pkg

// The ReverseIndex Singleton.
var revIndex ReverseIndex

// WriteTo prints the reverse index to w.
func (r ReverseIndex) WriteTo(w io.Writer) error {
	var list []string
	for path, pkg := range r {
		list = append(list, fmt.Sprintf("%s\t%s\n", path, pkg))
	}
	sort.Strings(list)

	tw := tabwriter.NewWriter(w, 2, 1, 1, ' ', 0)
	for _, line := range list {
		tw.Write([]byte(line))
	}
	return tw.Flush()
}

// PrefixMatch returns a prefix match of m in r and whether there was a match.
func (r ReverseIndex) PrefixMatch(m string) (string, bool) {
	// Not a simple map lookup because of prefix matches.
	// TODO(bsiegert) return longest
	for src, pkg := range r {
		if strings.HasPrefix(m, src) {
			return pkg.Path, true
		}
	}
	return "", false
}

// FullScan rebuilds the entire index by scanning all files in the pkgsrc dir.
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
	cmd := exec.Command(makeProg, "show-var", "VARNAME="+varname)
	cmd.Dir = filepath.Dir(filename)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(output)), nil
}
