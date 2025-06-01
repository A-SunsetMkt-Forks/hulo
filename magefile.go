//go:build mage

package main

// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// Run `mage setup` to setup the project.
//
// Run `mage build` to build the project.
//
// Run `mage clean` to clean the project.

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// donwloadDevDeps manages development dependencies
func donwloadDevDeps() error {
	log.Info("Installing development dependencies...")

	// Install other development tools
	tools := []string{
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		"golang.org/x/tools/cmd/goimports@latest",
	}

	for _, tool := range tools {
		if err := sh.Run("go", "install", tool); err != nil {
			return fmt.Errorf("failed to install %s: %v", tool, err)
		}
	}

	log.Info("Development dependencies installed successfully!")
	return nil
}

// Setup checks the environment and downloads the dependencies.
func Setup() error {
	log.Info("Start checking environment...")

	// Then setup project dependencies
	mg.Deps(donwloadDevDeps, resolveAntlrJar, downloadDeps)

	log.Info("Environment check completed!")
	return nil
}

var (
	javaVersionRegex = regexp.MustCompile(`version "(.*)"`)

	ErrJavaNotFound           = errors.New("please install Java")
	ErrJavaVersion            = errors.New("unable to determine Java version")
	ErrJavaVersionUnsupported = errors.New("Java version is not supported")
)

func checkJavaVersion() error {
	log.Info("Checking Java version...")
	cmd := exec.Command("java", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.WithError(err).Error("please install Java")
		return ErrJavaNotFound
	}

	// Parse Java version from output
	versionStr := string(output)
	if !strings.Contains(versionStr, "version") {
		log.WithError(ErrJavaVersion).Error("unable to determine Java version")
		return ErrJavaVersion
	}

	// Extract version number
	versionMatch := javaVersionRegex.FindStringSubmatch(versionStr)
	if len(versionMatch) < 2 {
		log.WithError(ErrJavaVersion).Error("unable to parse Java version")
		return ErrJavaVersion
	}

	ver, err := semver.NewVersion(versionMatch[1])
	if err != nil {
		log.WithError(err).Error("unable to parse Java version")
		return ErrJavaVersion
	}

	major := ver.Major()
	minor := ver.Minor()

	// ANTLR 4.13.2 requires Java 8 or higher
	if major < 8 || (major == 8 && minor < 0) {
		log.WithError(ErrJavaVersionUnsupported).
			Errorf("Java version must be 8 or higher, current version: %d.%d", major, minor)
		return ErrJavaVersionUnsupported
	}

	log.Infof("Java version: %s", ver.Original())
	return nil
}

const (
	// ANTLR 4.13.2 version
	ANTLR_JAR_URL = "https://www.antlr.org/download/antlr-4.13.2-complete.jar"
)

func resolveAntlrJar() error {
	mg.Deps(checkJavaVersion)
	log.Info("resolve antlr.jar")
	antlrPath := filepath.Join("syntax", "hulo", "parser", "antlr.jar")

	if _, err := os.Stat(antlrPath); os.IsNotExist(err) {
		log.Info("antlr.jar is not exist, downloading...")

		if err := os.MkdirAll(filepath.Dir(antlrPath), 0755); err != nil {
			log.WithError(err).Error("create directory failed")
			return err
		}

		log.WithField("url", ANTLR_JAR_URL).
			Info("Downloading antlr.jar...")
		resp, err := http.Get(ANTLR_JAR_URL)
		if err != nil {
			log.WithError(err).Error("download antlr.jar failed")
			return err
		}
		defer resp.Body.Close()

		log.WithField("file", antlrPath).Info("Creating file...")
		out, err := os.Create(antlrPath)
		if err != nil {
			log.WithError(err).Error("create file failed")
			return err
		}
		defer out.Close()

		log.WithField("file", antlrPath).Info("Saving file...")
		if _, err := io.Copy(out, resp.Body); err != nil {
			log.WithError(err).Error("save file failed")
			return err
		}

		log.Info("antlr.jar download completed")
	} else {
		log.WithField("file", antlrPath).Info("antlr.jar is exist")
	}
	return nil
}

func downloadDeps() error {
	log.Info("Downloading dependencies...")
	return sh.Run("go", "mod", "download")
}

// Clean cleans the project.
func Clean() error {
	log.Info("Start cleaning...")
	goreleaser := exec.Command("goreleaser", "build", "--clean")
	goreleaser.Stdout = os.Stdout
	goreleaser.Stderr = os.Stderr
	err := goreleaser.Run()
	if err != nil {
		log.WithError(err).Fatal("goreleaser build failed")
	}
	dirs := []string{"dist", "tmp"}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("clean %s failed: %v", dir, err)
		}
	}
	return nil
}

// Build builds the project.
func Build() error {
	log.Info("Start building...")

	mg.Deps(generateParser)

	goreleaser := exec.Command("goreleaser", "release", "--snapshot", "--clean")
	goreleaser.Stdout = os.Stdout
	goreleaser.Stderr = os.Stderr
	err := goreleaser.Run()
	if err != nil {
		log.WithError(err).Fatal("goreleaser release failed")
	}
	return nil
}

func generateParser() error {
	log.Info("Generating parser...")
	return sh.Run("java", "-jar", "syntax/hulo/parser/antlr.jar", "-Dlanguage=Go", "-visitor", "-no-listener", "-package", "parser", "syntax/hulo/parser/*.g4")
}

// Lint runs the linter
func Lint() error {
	log.Info("Running linter...")
	cmd := exec.Command("golangci-lint", "run", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.WithError(err).Fatal("lint failed")
	}
	return nil
}

// Format formats the code
func Format() error {
	log.Info("Formatting code...")
	return sh.Run("goimports", "-w", ".")
}

func Test() error {
	log.Info("Running tests...")
	cmd := exec.Command("go", "test", "./...", "-coverprofile=coverage.txt", "-covermode=atomic")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.WithError(err).Fatal("test failed")
	}
	return nil
}

// Dev starts the development server with hot reload
// func Dev() error {
// 	log.Info("Starting development server...")
// 	return sh.Run("air")
// }
