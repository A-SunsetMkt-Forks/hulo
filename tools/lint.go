// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"time"

	"github.com/caarlos0/log"
)

// Lint runs golangci-lint on the codebase.
func Lint() {
	log.Info("running linter...")
	start := time.Now()

	err := runCmd("golangci-lint", "run", "./...")
	if err != nil {
		log.WithError(err).Fatal("lint failed")
	}

	elapsed := time.Since(start)
	log.Infof("linter completed in %.2fs", elapsed.Seconds())
}
