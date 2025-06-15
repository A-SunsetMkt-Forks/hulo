// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"os"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/compiler"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type huloCmdParameters struct {
	lang    string
	verbose bool
	version bool
}

var (
	params  = huloCmdParameters{}
	huloCmd = &cobra.Command{
		Use:   "hulo filename",
		Short: "Hulo is a batch-oriented programming language.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if params.verbose {
				log.SetLevel(log.InfoLevel)
			}

			log.Info("reading huloc.yaml")
			src, err := os.ReadFile(config.FILE)
			if err != nil {
				log.WithError(err).Fatal("fail to read file")
			}

			huloc := &config.Huloc{}

			err = yaml.Unmarshal(src, huloc)
			if err != nil {
				log.WithError(err).Fatal("fail to unmarshal yaml")
			}

			file := args[0]
			if file != "." {
				huloc.Main = file
			}

			if len(huloc.Language) == 0 {
				huloc.Language = config.L_BASH
			}

			log.WithField("main", huloc.Main).
				WithField("langauge", huloc.Language).
				WithField("include", huloc.Include).
				WithField("exclude", huloc.Exclude).
				Info("starting compile")

			err = compiler.Compile(huloc)
			if err != nil {
				log.WithError(err).Info("fail to compile")
			}

			log.Info("thanks for using Hulo!")
		},
	}
)

func init() {
	huloCmd.PersistentFlags().StringVarP(&params.lang, "lang", "l", string(config.L_BASH), "specify a language to compile")
	huloCmd.PersistentFlags().BoolVarP(&params.verbose, "verbose", "v", false, "enables detailed log")
}

func main() {
	log.SetLevel(log.WarnLevel)
	if err := huloCmd.Execute(); err != nil {
		panic(err)
	}
}
