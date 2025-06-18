// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"fmt"
	"time"

	"github.com/caarlos0/log"
)

// builds distributable archives into current directory.
func Pack() {
	log.Info("running pack...")
	start := time.Now()

	outputZip := fmt.Sprintf("%s-%s.zip", name, version)
	fileCount, err := ZipDirWithGitIgnore(".", outputZip)
	if err != nil {
		elapsed := time.Since(start)
		log.WithError(err).Fatalf("pack failed after %.2fs", elapsed.Seconds())
	}

	elapsed := time.Since(start)
	log.Infof("packed %d files in %.2fs", fileCount, elapsed.Seconds())
}
