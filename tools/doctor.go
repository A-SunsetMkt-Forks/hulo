// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
	"github.com/opencommand/tinge"
)

var (
	ErrGoNotFound           = errors.New("please install Go")
	ErrGoVersion            = errors.New("unable to determine Go version")
	ErrGoVersionUnsupported = errors.New("Go version is not supported")
	ErrGitNotFound          = errors.New("please install Git")
	ErrGitVersion           = errors.New("unable to determine Git version")
)

// Doctor checks the development environment and reports any issues.
func Doctor() {
	log.Info(tinge.Styled().Bold("checking development environment...").String())
	start := time.Now()

	mg.Deps(
		checkGoVersion,
		checkGitVersion,
		checkJavaVersion,
		checkRequiredTools,
		checkAntlrJar,
	)

	elapsed := time.Since(start)
	log.Infof("environment check completed in %.2fs", elapsed.Seconds())
}

// checkGoVersion verifies that Go is installed and meets version requirements.
func checkGoVersion() error {
	log.Info("checking Go version...")
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.WithError(err).Error("please install Go")
		return ErrGoNotFound
	}

	versionStr := string(output)
	if !strings.Contains(versionStr, "go version") {
		log.WithError(ErrGoVersion).Error("unable to determine Go version")
		return ErrGoVersion
	}

	// Extract version number (e.g., "go version go1.21.0 darwin/amd64")
	parts := strings.Fields(versionStr)
	if len(parts) < 3 {
		log.WithError(ErrGoVersion).Error("unable to parse Go version")
		return ErrGoVersion
	}

	ver, err := semver.NewVersion(strings.TrimPrefix(parts[2], "go"))
	if err != nil {
		log.WithError(err).Error("unable to parse Go version")
		return ErrGoVersion
	}

	// Require Go 1.21 or higher
	if ver.Major() < 1 || (ver.Major() == 1 && ver.Minor() < 21) {
		log.WithError(ErrGoVersionUnsupported).
			Errorf("Go version must be 1.21 or higher, current version: %s", ver.Original())
		return ErrGoVersionUnsupported
	}

	log.Infof("Go version: %s", ver.Original())
	return nil
}

// checkGitVersion verifies that Git is installed.
func checkGitVersion() error {
	log.Info("checking Git version...")
	cmd := exec.Command("git", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.WithError(err).Error("please install Git")
		return ErrGitNotFound
	}

	versionStr := string(output)
	if !strings.Contains(versionStr, "git version") {
		log.WithError(ErrGitVersion).Error("unable to determine Git version")
		return ErrGitVersion
	}

	log.Infof("Git version: %s", strings.TrimSpace(versionStr))
	return nil
}

// checkRequiredTools verifies that all required development tools are installed.
func checkRequiredTools() error {
	log.Info("checking required development tools...")
	start := time.Now()

	for _, tool := range tools {
		parts := strings.Split(tool, "@")
		if len(parts) == 0 {
			continue
		}
		path := strings.TrimPrefix(parts[0], "github.com/")
		cmd := exec.Command("which", path)
		if err := cmd.Run(); err != nil {
			log.WithError(err).Errorf("missing required tool: %s", path)
			return fmt.Errorf("missing required tool: %s", path)
		}
		log.Infof("found tool: %s", path)
	}

	elapsed := time.Since(start)
	log.Infof("all required tools found in %.2fs", elapsed.Seconds())
	return nil
}

// checkAntlrJar verifies that the ANTLR jar file exists.
func checkAntlrJar() error {
	log.Info("checking ANTLR jar...")
	antlrPath := filepath.Join("syntax", "hulo", "parser", "antlr.jar")

	if _, err := os.Stat(antlrPath); os.IsNotExist(err) {
		log.WithError(err).Error("ANTLR jar not found")
		return fmt.Errorf("ANTLR jar not found at %s", antlrPath)
	}

	log.Info("ANTLR jar found")
	return nil
}
