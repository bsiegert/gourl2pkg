package main

import (
	"os"
	"os/exec"
)

// GoGet calls "go get" to download all the srcpaths. dir is the base directory
// of the gopath.
func GoGet(srcpaths []string, dir string) error {
	args := []string{"get", "-v"}
	for _, s := range srcpaths {
		args = append(args, s+"/...")
	}
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), "GOPATH="+dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
