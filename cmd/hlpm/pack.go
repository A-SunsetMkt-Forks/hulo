package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "pack a package",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}
