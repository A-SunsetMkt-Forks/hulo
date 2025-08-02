// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"os"
	"time"

	"github.com/caarlos0/log"
	"github.com/opencommand/tinge"
)

// Clean removes generated files and build artifacts.
func Clean() {
	log.Info(tinge.Styled().Bold("starting clean...").String())
	start := time.Now()

	dirs := []string{"dist"}
	for _, dir := range dirs {
		log.WithField("dir", dir).Info("remove files")
		if err := os.RemoveAll(dir); err != nil {
			log.WithError(err).WithField("dir", dir).Fatal("clean failed")
		}
	}

	elapsed := time.Since(start)
	log.Infof("clean completed in %.2fs", elapsed.Seconds())
}
