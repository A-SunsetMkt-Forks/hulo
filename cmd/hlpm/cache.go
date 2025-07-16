package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 全局依赖缓存

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "cache a package",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

var cacheListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all cached packages",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

var cacheCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean all cached packages",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}
