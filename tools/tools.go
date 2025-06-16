// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

// CopyFileIfNotExists copies a file from src to dst if dst doesn't exist.
// Returns error if copy fails or if src doesn't exist.
func CopyFileIfNotExists(src, dst string) error {
	// Check if destination file already exists
	if _, err := os.Stat(dst); err == nil {
		log.WithField("path", dst).Info("file already exists, skipping copy")
		return nil
	}

	// Check if source file exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", src)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dstFile.Close()

	// Copy file contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	log.WithField("src", src).WithField("dst", dst).Info("file copied successfully")
	return nil
}
