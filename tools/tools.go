// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/caarlos0/log"
	ignore "github.com/sabhiram/go-gitignore"
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

// ZipDirWithGitIgnore creates a zip file from the source directory,
// excluding files matching .gitignore patterns.
func ZipDirWithGitIgnore(sourceDir, outputZip string) (int, error) {
	log.WithField("dir", sourceDir).
		WithField("output", outputZip).
		Info("packing directory")

	// Load .gitignore if exists
	ig := ignore.CompileIgnoreLines()
	if data, err := os.ReadFile(filepath.Join(sourceDir, ".gitignore")); err == nil {
		lines := strings.Split(string(data), "\n")
		ig = ignore.CompileIgnoreLines(lines...)
		log.Info(".gitignore loaded")
	} else {
		log.Info("no .gitignore found or failed to read, using default ignore rules")
	}

	// Create the zip file
	outFile, err := os.Create(outputZip)
	if err != nil {
		log.WithError(err).Errorf("failed to create zip file: %s", outputZip)
		return 0, err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	var fileCount int
	err = filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.WithError(err).Errorf("walk error at %s", path)
			return err
		}

		// Skip the output zip itself
		if path == outputZip {
			return nil
		}

		// Skip ignored paths
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			log.WithError(err).Errorf("failed to get relative path for %s", path)
			return err
		}
		if relPath == "." || relPath == "" {
			return nil
		}
		relPath = filepath.ToSlash(relPath) // normalize to forward slashes for gitignore

		if ig.MatchesPath(relPath) {
			if d.IsDir() {
				log.Debugf("skip dir (gitignore): %s", relPath)
				return filepath.SkipDir
			}
			log.Debugf("skip file (gitignore): %s", relPath)
			return nil
		}

		// Directories are not added to zip
		if d.IsDir() {
			return nil
		}

		// Add file to zip
		err = addFileToZip(zipWriter, path, relPath)
		if err != nil {
			log.WithError(err).Errorf("failed to add file: %s", relPath)
			return err
		}
		fileCount++
		log.Debugf("added: %s", relPath)
		return nil
	})

	return fileCount, err
}

func addFileToZip(zipWriter *zip.Writer, fullPath, relPath string) error {
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	w, err := zipWriter.Create(relPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, file)
	return err
}

func resolveVersion() {
	data, err := exec.Command("gitversion").Output()
	if err != nil {
		log.WithError(err).Fatal("failed to get git version")
	}
	err = json.Unmarshal(data, &version)
	if err != nil {
		log.WithError(err).Fatal("failed to unmarshal git version")
	}
}
