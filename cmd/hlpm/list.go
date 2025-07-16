package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 列出当前项目的依赖树

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all packages",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}
