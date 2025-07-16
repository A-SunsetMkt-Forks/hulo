// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"fmt"

	"github.com/caarlos0/log"
	"github.com/google/go-github/v73/github"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/util"
	"github.com/spf13/cobra"
)

type installParams struct {
	Packages []string
	Proxies  []string
}

var installCmd = &cobra.Command{
	Use:     "install [packages...]",
	Short:   "install a package",
	Aliases: []string{"i", "add", "get", "download"},
	Run: func(cmd *cobra.Command, args []string) {
		pkg, err := util.LoadConfigure[config.HuloPkg](config.HuloPkgFileName)
		if err != nil {
			log.WithError(err).Fatal("fail to load package")
		}
		fmt.Println(pkg)
		client := github.NewClient(nil)

		release, _, err := client.Repositories.GetLatestRelease(context.Background(), "hulo-lang", "hulo")
		if err != nil {
			log.WithError(err).Fatal("fail to get latest release")
		}
		fmt.Println(release)
	},
}

func init() {}
