// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/tools/release"
	"github.com/magefile/mage/mg"
)

var (
	name = "hulo"
	// TODO build as tag to inject cmd/hulo.go
	version     = "1.0.0"
	description = ""
	homepage    = "https://hulo-lang.github.io/docs/"
	repository  = "https://github.com/hulo-lang/hulo"
	author      = "The Hulo Authors"
	license     = "MIT"
	keywords    = []string{}
)

type Release mg.Namespace

// builds and publishes releases to all targets (npm, pypi, homebrew, etc).
func (r Release) All() {
	mg.Deps(r.Npm, r.Pypi, r.Gofish)
}

// builds and publishes the PyPI package.
func (r Release) Pypi() {
	r.prepare()
}

// builds and publishes the npm package.
func (r Release) Npm() {
	r.prepare()
}

// generates and publishes GoFish manifest files.
func (r Release) Gofish() {
	r.prepare()
	food := release.Food{
		Name:        name,
		Description: description,
		Homepage:    homepage,
		Version:     version,
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

func (Release) prepare() {
	os.MkdirAll("dist/release", 0755)
}

func buildPackages() []release.Package {
	baseURL := "https://github.com/hulo-lang/hulo/releases/download/v" + version + "/"
	distDir := "dist"
	entries, err := os.ReadDir(distDir)
	if err != nil {
		log.Fatalf("failed to read dist dir: %v", err)
	}

	re := regexp.MustCompile(`^hulo_(\w+)_([\w\d]+)_v[\d.]+$`)
	var pkgs []release.Package

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		matches := re.FindStringSubmatch(entry.Name())
		if len(matches) != 3 {
			continue
		}
		goos := matches[1]
		arch := matches[2]

		var fileExt string
		if goos == "windows" {
			fileExt = ".zip"
		} else {
			fileExt = ".tar.gz"
		}

		archiveName := entry.Name() + fileExt
		url := baseURL + archiveName

		// ⚠️ 你需要计算实际文件的 sha256，这里是占位
		sha := "REPLACE_WITH_REAL_SHA256"

		binName := "hulo"
		if goos == "windows" {
			binName += ".exe"
		}

		pkg := release.Package{
			OS:     goos,
			Arch:   arch,
			URL:    url,
			SHA256: sha,
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
