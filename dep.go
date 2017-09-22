package main

import (
	"bytes"
	"fmt"
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
		//fmt.Printf("%s\n", line)
		if n := bytes.Index(line, []byte(" (download)")); n != -1 {
			repos = append(repos, string(line[:n]))
		}
	}
	return repos, nil
}
