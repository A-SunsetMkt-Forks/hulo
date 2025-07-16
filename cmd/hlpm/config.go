// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config a package",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}
