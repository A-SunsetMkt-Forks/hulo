// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"os"
	"runtime"

	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
)

// Project metadata for release configuration
var (
	name        = "hulo"
	version     = GitVersion{}
	description = "Hulo is a batch-oriented programming language."
	homepage    = "https://hulo-lang.github.io/docs/"
	repository  = "https://github.com/hulo-lang/hulo"
	author      = "The Hulo Authors"
	license     = "MIT"
	keywords    = []string{"hulo", "batch", "programming", "language", "vbs", "vbscript"}
)

// Release creates a new release using GoReleaser with snapshot mode.
// This function sets up environment variables and runs the release process.
func Release() {
	mg.Deps(resolveVersion)

	// Set version environment variable for GoReleaser
	err := os.Setenv("Version", version.SemVer)
	if err != nil {
		log.WithError(err).Fatal("failed to set Version environment variable")
	}

	// Set Go version environment variable
	err = os.Setenv("GOVERSION", runtime.Version())
	if err != nil {
		log.WithError(err).Fatal("failed to set GOVERSION environment variable")
	}

	// Run GoReleaser in snapshot mode with clean flag
	err = runCmd("goreleaser", "release", "--snapshot", "--clean")
	if err != nil {
		log.WithError(err).Fatal("goreleaser build failed")
	}
}

// GitVersion represents the version information extracted from Git tags
// using GitVersion tool. This struct contains all the version metadata
// needed for building and releasing the project.
type GitVersion struct {
	AssemblySemFileVer        string `json:"AssemblySemFileVer"`
	AssemblySemVer            string `json:"AssemblySemVer"`
	BranchName                string `json:"BranchName"`
	BuildMetaData             any    `json:"BuildMetaData"`
	CommitDate                string `json:"CommitDate"`
	CommitsSinceVersionSource int    `json:"CommitsSinceVersionSource"`
	EscapedBranchName         string `json:"EscapedBranchName"`
	FullBuildMetaData         string `json:"FullBuildMetaData"`
	FullSemVer                string `json:"FullSemVer"`
	InformationalVersion      string `json:"InformationalVersion"`
	Major                     int    `json:"Major"`
	MajorMinorPatch           string `json:"MajorMinorPatch"`
	Minor                     int    `json:"Minor"`
	Patch                     int    `json:"Patch"`
	PreReleaseLabel           string `json:"PreReleaseLabel"`
	PreReleaseLabelWithDash   string `json:"PreReleaseLabelWithDash"`
	PreReleaseNumber          int    `json:"PreReleaseNumber"`
	PreReleaseTag             string `json:"PreReleaseTag"`
	PreReleaseTagWithDash     string `json:"PreReleaseTagWithDash"`
	SemVer                    string `json:"SemVer"`
	Sha                       string `json:"Sha"`
	ShortSha                  string `json:"ShortSha"`
	UncommittedChanges        int    `json:"UncommittedChanges"`
	VersionSourceSha          string `json:"VersionSourceSha"`
	WeightedPreReleaseNumber  int    `json:"WeightedPreReleaseNumber"`
}
