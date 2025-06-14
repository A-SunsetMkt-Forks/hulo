// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package tools

import (
	"time"

	"github.com/caarlos0/log"
)

// Test runs all unit tests with coverage.
func Test() {
	log.Info("running tests...")
	start := time.Now()

	err := runCmd("go", "test", "./...", "-coverprofile=coverage.txt", "-covermode=atomic")
	if err != nil {
		log.WithError(err).Fatal("test failed")
	}

	elapsed := time.Since(start)
	log.Infof("tests completed in %.2fs", elapsed.Seconds())
}
