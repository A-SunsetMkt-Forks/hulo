// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/cmd/hulo/strategies"
	"github.com/spf13/cobra"
)

var (
	registry = strategies.NewStrategyRegistry()
	params   = &strategies.Parameters{}
	huloCmd  = &cobra.Command{
		Use:   "hulo [filename]",
		Short: "Hulo is a batch-oriented programming language.",
		Example: `
# print the version
hulo -V

# compile main.hl to current directory
hulo main.hl

# compile main.hl to ./output
hulo -d ./output main.hl

# compile main.hl to ./output with VBScript language
hulo -l vbs -d ./output main.hl

# enable verbose mode
hulo main.hl -v
`,
		Run: func(cmd *cobra.Command, args []string) {
			params.Args = args

			strategy := registry.FindStrategy(params)
			if strategy == nil {
				cmd.Help()
				return
			}

			err := strategy.Execute(params, args)
			if err != nil {
				log.WithError(err).Fatal("fail to execute strategy")
			}
		},
	}
)

func init() {
	huloCmd.PersistentFlags().BoolVarP(&params.Version, "version", "V", false, "print the version")
	huloCmd.PersistentFlags().StringSliceVarP(&params.Langs, "lang", "l", []string{}, "specify languages to compile")
	huloCmd.PersistentFlags().BoolVarP(&params.Verbose, "verbose", "v", false, "enables detailed log")
	huloCmd.PersistentFlags().StringVarP(&params.OutDir, "outdir", "o", ".", "specify a directory to write the output")
}

func main() {
	registry.Register(&strategies.VersionStrategy{})
	registry.Register(&strategies.CompileStrategy{})
	log.SetLevel(log.InfoLevel)
	if err := huloCmd.Execute(); err != nil {
		panic(err)
	}
}
