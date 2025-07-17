// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/util"
	"github.com/opencommand/tinge"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type initParameters struct {
	All   bool
	Files []string
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a package",
	Run: func(cmd *cobra.Command, args []string) {
		if util.Exists(config.HuloPkgFileName) {
			log.Fatal("package already exists")
		}

		log.Info(tinge.Styled().
			Bold("generating").
			Space().
			Text("hulo.pkg.yaml").
			String())

		pkg := config.NewHuloPkg()
		wd, err := os.Getwd()
		if err != nil {
			wd = "my-project"
		}
		pkg.Name = filepath.Base(wd)
		out, err := yaml.Marshal(&pkg)
		if err != nil {
			log.WithError(err).Fatal("fail to marshal package")
		}
		os.WriteFile(config.HuloPkgFileName, out, 0644)

		log.Info(tinge.Styled().Bold("setting up .gitignore").String())

		if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
			os.WriteFile(".gitignore", []byte("dist/"), 0644)
		} else {
			content, err := os.ReadFile(".gitignore")
			if err != nil {
				log.WithError(err).Fatal("fail to read .gitignore")
			}
			if !strings.Contains(string(content), "dist/") {
				os.WriteFile(".gitignore", append([]byte("dist/\n"), content...), 0644)
			}
		}

		log.Info(tinge.Styled().
			Bold("done!").
			Space().
			Text("please edit hulo.pkg.yaml and .gitignore accordingly.").
			String())

		log.Info("thanks for using hulo!")
	},
}
