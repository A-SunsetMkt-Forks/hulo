// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/tools/release"
	"github.com/magefile/mage/mg"
)

var (
	name        = "hulo"
	version     = release.GitVersion{}
	description = "Hulo is a batch-oriented programming language."
	homepage    = "https://hulo-lang.github.io/docs/"
	repository  = "https://github.com/hulo-lang/hulo"
	author      = "The Hulo Authors"
	license     = "MIT"
	keywords    = []string{"hulo", "batch", "programming", "language", "vbs", "vbscript"}
)

type Release mg.Namespace

// builds and publishes releases to all targets (npm, pypi, homebrew, etc).
func (r Release) All() {
	mg.Deps(r.Gofish)
}

// generates and publishes GoFish manifest files.
func (r Release) Gofish() {
	mg.Deps(r.setup)
	food := release.Food{
		Name:        name,
		Description: description,
		Homepage:    homepage,
		Version:     version.SemVer,
		License:     license,
		Packages:    buildPackages(),
	}

	templatePath := filepath.Join("tools", "release", "templates", "gofish", "hulo.lua.tmpl")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}

	outputDir := filepath.Join("dist", "release", "gofish")
	os.MkdirAll(outputDir, 0755)
	outputPath := filepath.Join(outputDir, food.Name+".lua")
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, food)
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}
}

func (Release) setup() {
	mg.Deps(resolveVersion)
	err := os.Setenv("Version", version.SemVer)
	if err != nil {
		log.WithError(err).Fatal("failed to set Version environment variable")
	}

	err = os.Setenv("GOVERSION", runtime.Version())
	if err != nil {
		log.WithError(err).Fatal("failed to set GOVERSION environment variable")
	}

	err = runCmd("goreleaser", "release", "--snapshot", "--clean")
	if err != nil {
		log.WithError(err).Fatal("goreleaser build failed")
	}

	os.MkdirAll("dist/release", 0755)
}

type Checksum struct {
	SHA256  string `json:"sha256"`
	Archive string `json:"archive"`
}

func parseChecksums() ([]Checksum, error) {
	data, err := os.ReadFile("dist/checksums.txt")
	if err != nil {
		return nil, err
	}

	var checksums []Checksum
	lines := strings.SplitSeq(string(data), "\n")
	for line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "  ")
		if len(parts) != 2 {
			continue
		}
		checksums = append(checksums, Checksum{
			SHA256:  parts[0],
			Archive: parts[1],
		})
	}
	return checksums, nil
}

func buildPackages() []release.Package {
	baseURL := "https://github.com/hulo-lang/hulo/releases/download/v" + version.MajorMinorPatch + "/"

	checksums, err := parseChecksums()
	if err != nil {
		log.Fatalf("failed to parse checksums: %v", err)
	}

	var pkgs []release.Package

	for _, checksum := range checksums {
		matches := strings.SplitN(checksum.Archive, "_", 3)
		if len(matches) != 3 {
			log.WithField("archive", checksum.Archive).Warn("invalid archive name")
			continue
		}
		goos := strings.ToLower(matches[1])
		arch := matches[2]

		arch = strings.TrimSuffix(arch, ".tar.gz")
		arch = strings.TrimSuffix(arch, ".zip")

		url := baseURL + checksum.Archive

		binName := "hulo"
		if goos == "windows" {
			binName += ".exe"
		}

		pkg := release.Package{
			OS:     goos,
			Arch:   arch,
			URL:    url,
			SHA256: checksum.SHA256,
			Resources: []release.Resource{
				{
					Path:        binName,
					InstallPath: path.Join("bin", binName),
					Executable:  true,
				},
			},
		}

		pkgs = append(pkgs, pkg)
	}

	return pkgs
}
