package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 读取 version && git tag ${version} && git push origin ${version}
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a package",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}
