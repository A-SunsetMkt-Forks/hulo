// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"os"
	"strings"
	"time"

	"github.com/caarlos0/log"
	build "github.com/hulo-lang/hulo/internal/build/vbs"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/vfs/osfs"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
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
			startTime := time.Now()
			if params.verbose {
				log.SetLevel(log.InfoLevel)
			}
			localFs := osfs.New()

			var huloc *config.Huloc

			if localFs.Exists(config.NAME) {
				log.Info("reading huloc.yaml")
				src, err := localFs.ReadFile(config.NAME)
				if err != nil {
					log.WithError(err).Fatal("fail to read file")
				}
				err = yaml.Unmarshal(src, huloc)
				if err != nil {
					log.WithError(err).Fatal("fail to unmarshal yaml")
				}

				err = huloc.Validate()
				if err != nil {
					log.WithError(err).Fatal("fail to validate huloc")
				}
			} else {
				huloc = &config.Huloc{Main: args[0], CompilerOptions: config.CompilerOptions{VBScript: &config.VBScriptOptions{CommentSyntax: "'"}}}
			}

			hulopath := os.Getenv("HULOPATH")
			if hulopath == "" {
				huloc.HULOPATH = "."
			} else {
				huloc.HULOPATH = hulopath
			}

			file := args[0]
			if file != "." {
				huloc.Main = file
			}

			if len(huloc.Language) == 0 {
				huloc.Language = config.L_VBSCRIPT
			}

			log.WithField("main", huloc.Main).
				WithField("langauge", huloc.Language).
				WithField("include", huloc.Include).
				WithField("exclude", huloc.Exclude).
				Info("starting compile")

			var popts []parser.ParserOptions

			if !params.verbose {
				popts = append(popts, parser.OptionDisableTracer())
			}

			results, err := build.Transpile(huloc.CompilerOptions.VBScript, huloc.Main, localFs, ".", huloc.HULOPATH, popts...)
			if err != nil {
				log.WithError(err).Info("fail to compile")
			}

			for file, code := range results {
				file = strings.Replace(file, ".hl", ".vbs", 1)
				err := localFs.WriteFile(file, []byte(code), 0644)
				if err != nil {
					log.WithError(err).Info("fail to write file")
				}
			}

			// err = compiler.Compile(huloc)
			// if err != nil {
			// 	log.WithError(err).Info("fail to compile")
			// }

			log.Infof("compile time: %s", time.Since(startTime))

			log.Info("thanks for using Hulo!")
		},
	}
)

func init() {
	huloCmd.PersistentFlags().StringVarP(&params.lang, "lang", "l", string(config.L_BASH), "specify a language to compile")
	huloCmd.PersistentFlags().BoolVarP(&params.verbose, "verbose", "v", false, "enables detailed log")
}

func main() {
	log.SetLevel(log.InfoLevel)
	if err := huloCmd.Execute(); err != nil {
		panic(err)
	}
}
