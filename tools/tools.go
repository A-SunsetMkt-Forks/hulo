// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"os"
	"os/exec"
	"strings"

	"github.com/caarlos0/log"
)

// runCmd executes a system command specified by name and arguments,
// using the current working directory.
// It attaches the standard input, output, and error streams to the process,
// and logs the command being executed.
func runCmd(name string, args ...string) error {
	return runCmdInDir("", name, args...)
}

// runCmdInDir executes a system command specified by name and arguments,
// in the given working directory 'dir'.
// If 'dir' is empty, it runs in the current directory.
// The command's standard input, output, and error streams are connected
// to those of the current process, enabling interactive input/output.
// The function logs the executed command and directory before running.
// Returns an error if the command fails to start or execute.
func runCmdInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.WithField("cmd", strings.Join(append([]string{name}, args...), " ")).
		WithField("dir", dir).
		Info("execute command")
	return cmd.Run()
}
