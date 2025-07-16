// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"runtime"

	"github.com/hulo-lang/hulo/cmd/meta"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version",
	Run: func(cmd *cobra.Command, args []string) {
		// Print logo
		fmt.Print(meta.HLPMLOGO)
		fmt.Println()

		// Print version information
		fmt.Printf("hlpm %s\n", meta.Version)
		fmt.Printf("Hulo Package Manager (build %s)\n", meta.Date)
		fmt.Printf("Hulo %s %s Package Manager (build %s, %s)\n", runtime.GOOS, runtime.GOARCH, meta.Commit, meta.GoVersion)
	},
}
