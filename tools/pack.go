// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"fmt"
	"time"

	"github.com/caarlos0/log"
)

func Pack() {
	log.Info("running pack...")
	start := time.Now()

	err := ZipDirWithGitIgnore(".", fmt.Sprintf("%s-%s.zip", name, version))
	if err != nil {
		elapsed := time.Since(start)
		log.WithError(err).Fatalf("pack failed after %.2fs", elapsed.Seconds())
	}

	elapsed := time.Since(start)
	log.Infof("pack completed in %.2fs", elapsed.Seconds())
}
