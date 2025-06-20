// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

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
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
)

// Setup prepares the development environment and downloads dependencies.
func Setup() {
	log.Info("setting up development environment...")
	start := time.Now()

	mg.Deps(installDevDeps, ensureAntlrJar, downloadDeps)

	elapsed := time.Since(start)
	log.Infof("setup completed in %.2fs", elapsed.Seconds())
}

var tools = []string{
	"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
	"golang.org/x/tools/cmd/goimports@latest",
	"github.com/goreleaser/goreleaser/v2@v2.10.2",
	"golang.org/x/tools/cmd/stringer@latest",
	"github.com/securego/gosec/v2/cmd/gosec@latest",
}

// installDevDeps installs development tools needed for building and testing.
func installDevDeps() error {
	log.Info("installing development dependencies...")
	start := time.Now()

	for _, tool := range tools {
		if err := runCmd("go", "install", tool); err != nil {
			return fmt.Errorf("failed to install %s: %v", tool, err)
		}
	}

	elapsed := time.Since(start)
	log.Infof("development dependencies installed after %.2fs", elapsed.Seconds())
	return nil
}

var (
	javaVersionRegex = regexp.MustCompile(`version "(.*)"`)

	ErrJavaNotFound           = errors.New("please install Java")
	ErrJavaVersion            = errors.New("unable to determine Java version")
	ErrJavaVersionUnsupported = errors.New("Java version is not supported")
)

// checkJavaVersion executes 'java -version', parses the version,
// and ensures it's Java 8 or higher. Returns error if checks fail.
func checkJavaVersion() error {
	log.Info("checking Java version...")
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

// ensureAntlrJar makes sure the ANTLR jar file exists, downloads it if missing.
func ensureAntlrJar() error {
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
			Info("downloading antlr.jar...")
		resp, err := http.Get(ANTLR_JAR_URL)
		if err != nil {
			log.WithError(err).Error("download antlr.jar failed")
			return err
		}
		defer resp.Body.Close()

		log.WithField("path", antlrPath).Info("writing to antlr.jar...")

		out, err := os.Create(antlrPath)
		if err != nil {
			log.WithError(err).Error("create file failed")
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			log.WithError(err).Error("save file failed")
			return err
		}
	}

	err := CopyFileIfNotExists(antlrPath, filepath.Join("internal", "unsafe", "antlr.jar"))
	if err != nil {
		return err
	}
	log.Info("antlr.jar is ready")
	return nil
}

// downloadDeps downloads Go module dependencies.
func downloadDeps() error {
	log.Info("downloading Go module dependencies")
	return runCmd("go", "mod", "download")
}
