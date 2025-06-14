// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"time"

	"github.com/caarlos0/log"
)

// Format automatically formats the source code using goimports.
func Format() {
	log.Info("formatting code...")
	start := time.Now()

	err := runCmd("goimports", "-w", ".")
	if err != nil {
		log.WithError(err).Fatal("format failed")
	}

	elapsed := time.Since(start)
	log.Infof("formatter completed in %.2fs", elapsed.Seconds())
}
