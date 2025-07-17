// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/util"
	"github.com/spf13/cobra"
)

type uninstallParameters struct {
	RemoveFiles bool
}

var (
	uninstallParams uninstallParameters
	uninstallCmd    = &cobra.Command{
		Use:     "uninstall [packages...]",
		Short:   "uninstall a package",
		Aliases: []string{"rm", "remove", "del", "delete", "un"},
		Run: func(cmd *cobra.Command, args []string) {
			if !util.Exists(config.HuloPkgFileName) {
				log.Fatal("package file not found, please run `hlpm init` first")
			}

			if len(args) == 0 {
				log.Fatal("no packages specified")
			}

			huloModulesDir, err := getHuloModulesDir()
			if err != nil {
				log.WithError(err).Fatal("failed to get hulo modules directory")
			}

			pkg, err := util.LoadConfigure[config.HuloPkg](config.HuloPkgFileName)
			if err != nil {
				log.WithError(err).Fatal("fail to load package")
			}

			for _, pkgName := range args {
				if err := uninstallPackage(pkgName, uninstallParams.RemoveFiles, huloModulesDir, &pkg); err != nil {
					log.WithError(err).WithField("package", pkgName).Fatal("failed to uninstall package")
				}
			}

			util.SaveConfigure(config.HuloPkgFileName, pkg)
			log.Info("uninstall completed")
		},
	}
)

func init() {
	uninstallCmd.Flags().BoolVarP(&uninstallParams.RemoveFiles, "remove-files", "f", false, "remove package files from disk")
}

// uninstallPackage removes a package from dependencies and optionally removes files
func uninstallPackage(pkgName string, removeFiles bool, huloModulesDir string, pkg *config.HuloPkg) error {
	owner, repo, _, err := resolvePackage(pkgName)
	if err != nil {
		return fmt.Errorf("failed to resolve package name: %w", err)
	}

	// Remove from dependencies
	packageKey := fmt.Sprintf("%s/%s", owner, repo)
	if _, exists := pkg.Dependencies[packageKey]; !exists {
		return fmt.Errorf("%s is not installed", pkgName)
	}

	delete(pkg.Dependencies, packageKey)
	log.WithField("package", pkgName).Info("removed from dependencies")

	// Optionally remove files
	if removeFiles {
		packageDir := filepath.Join(huloModulesDir, owner, repo)
		if err := os.RemoveAll(packageDir); err != nil {
			return fmt.Errorf("failed to remove package files: %w", err)
		}
		log.WithField("directory", packageDir).Info("removed package files")
	} else {
		log.Info("package files preserved (use --remove-files to delete them)")
	}

	return nil
}
