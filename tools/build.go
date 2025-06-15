// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"time"

	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
)

// Build compiles and packages the project using GoReleaser.
func Build() {
	log.Info("starting build...")
	start := time.Now()

	mg.Deps(generateParser)

	elapsed := time.Since(start)
	log.Infof("build completed in %.2fs", elapsed.Seconds())

	// TODO build:all build:single-platform
	err := runCmd("goreleaser", "build", "--snapshot", "--clean")
	if err != nil {
		log.WithError(err).Fatal("goreleaser build failed")
	}
}
