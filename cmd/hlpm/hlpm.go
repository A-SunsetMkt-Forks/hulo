// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/spf13/cobra"
)

var (
	huloPkg *config.HuloPkg
	hlpmCmd = &cobra.Command{
		Use:   "hlpm",
		Short: "Hulo Package Manager",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, World!")
		},
	}
)

func init() {
	hlpmCmd.AddCommand(versionCmd, buildCmd, runCmd,
		installCmd, uninstallCmd, publishCmd,
		listCmd, cacheCmd, configCmd, initCmd)
}

func main() {
	err := hlpmCmd.Execute()
	if err != nil {
		panic(err)
	}
}
