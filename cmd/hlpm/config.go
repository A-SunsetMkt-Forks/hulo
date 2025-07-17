// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config a package",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HULO_PATH:", os.Getenv("HULO_PATH"))
		fmt.Println("HULO_MODULES:", os.Getenv("HULO_MODULES"))
	},
}
