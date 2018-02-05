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
	"log"
	"os"
	"os/exec"
)

var _ = fmt.Printf

// GoGet calls "go get" to download all the srcpaths. dir is the base directory
// of the gopath. It returns a list of downloaded repos and an error, if any.
func GoGet(srcpaths []string, dir string) ([]string, error) {
	args := []string{"get", "-v"}
	for _, s := range srcpaths {
		args = append(args, s+"/...")
	}
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), "GOPATH="+dir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// Read through the output to find download lines.
	var repos []string
	for _, line := range bytes.Split(output, []byte{'\n'}) {
		if n := bytes.Index(line, []byte(" (download)")); n != -1 {
			repos = append(repos, string(line[:n]))
			log.Printf("%s", line)
		}
	}
	return repos, nil
}

// GoGetResolve runs go get n times until there are no new repos.
func GoGetResolve(srcpaths []string, dir string) ([]string, error) {
	var ToPackage []string
	var err error
	for i := 1; ; i++ {
		InfoLog.Printf("Run %d", i)
		srcpaths, err = GoGet(srcpaths, dir)
		if err != nil {
			return nil, err
		}
		if len(srcpaths) == 0 {
			break
		}
		log.Println(srcpaths)
		ToPackage = append(ToPackage, srcpaths...)
	}
	return ToPackage, nil
}
