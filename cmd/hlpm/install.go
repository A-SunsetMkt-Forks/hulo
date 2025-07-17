// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"archive/zip"

	"github.com/caarlos0/log"

	"github.com/google/go-github/v73/github"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/util"
	"github.com/spf13/cobra"
)

type installParameters struct {
	Proxies  []string
	Registry string
	Retries  int
}

var (
	installParams installParameters
	installCmd    = &cobra.Command{
		Use:     "install [packages...]",
		Short:   "install a package",
		Aliases: []string{"i", "add", "get", "download"},
		Run: func(cmd *cobra.Command, args []string) {
			if !util.Exists(config.HuloPkgFileName) {
				log.Fatal("package file not found, please run `hlpm init` first")
			}

			if len(args) == 0 {
				log.Fatal("no packages specified")
			}

			if util.Exists(config.HuloRCFileName) {
				rc, err := util.LoadConfigure[config.HuloRC](config.HuloRCFileName)
				if err != nil {
					log.WithError(err).Fatal("fail to load hulorc")
				}
				if len(rc.Registry.Default) > 0 {
					installParams.Registry = rc.Registry.Default
				}
				if rc.Network.Retries > 0 {
					installParams.Retries = rc.Network.Retries
				}
			}

			huloModulesDir, err := getHuloModulesDir()
			if err != nil {
				log.WithError(err).Fatal("failed to get hulo modules directory")
			}
			pkg, err := util.LoadConfigure[config.HuloPkg](config.HuloPkgFileName)
			if err != nil {
				log.WithError(err).Fatal("fail to load package")
			}
			switch installParams.Registry {
			case "github.com":
				client := github.NewClient(nil)

				for _, pkgName := range args {
					owner, repo, version, err := resolvePackage(pkgName)
					if err != nil {
						log.WithError(err).Fatal("failed to resolve package name")
					}

					log.WithField("owner", owner).
						WithField("repo", repo).
						WithField("version", version).
						Info("installing package")

					// If no version specified, get the latest release
					if version == "" {
						release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
						if err != nil {
							log.WithError(err).Warn("failed to get latest release, trying default branch")
							// If no release, try to get the default branch
							repository, _, err := client.Repositories.Get(context.Background(), owner, repo)
							if err != nil {
								log.WithError(err).Fatal("failed to get repository info")
							}
							version = repository.GetDefaultBranch()
						} else {
							version = release.GetTagName()
						}
					}

					// Create download URL for source code
					// GitHub provides source code downloads at: https://github.com/owner/repo/archive/refs/heads/main.zip
					// or for tags: https://github.com/owner/repo/archive/refs/tags/v1
					var downloadURL string
					if strings.HasPrefix(version, "v") || strings.Contains(version, ".") {
						// It's a tag/version
						downloadURL = fmt.Sprintf("https://github.com/%s/%s/archive/refs/tags/%s.zip", owner, repo, version)
					} else {
						// Its a branch
						downloadURL = fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.zip", owner, repo, version)
					}

					log.WithField("url", downloadURL).Info("downloading source code")

					tempFile, err := os.CreateTemp("", "hulo-install-*")
					if err != nil {
						log.WithError(err).Fatal("failed to create temporary file")
					}
					defer os.Remove(tempFile.Name())

					for i := range installParams.Retries {
						err = util.Download(downloadURL, tempFile)
						if err != nil {
							log.WithError(err).Warnf("failed to download source code, retrying %d/%d", i+1, installParams.Retries)
							continue
						}
						break
					}
					if err != nil {
						log.WithError(err).Fatal("failed to download source code")
					}

					if err := extractZipToModules(tempFile.Name(), owner, repo, version, huloModulesDir); err != nil {
						log.WithError(err).Fatal("failed to extract package")
					}

					log.Infof("Successfully installed %s/%s@%s to %s",
						owner, repo, version, filepath.Join(huloModulesDir, owner, repo, version))
					pkg.Dependencies[fmt.Sprintf("%s/%s", owner, repo)] = version
				}

			default:
				log.WithField("registry", installParams.Registry).Fatal("unsupported registry")
			}

			util.SaveConfigure(config.HuloPkgFileName, pkg)
		},
	}
)

func init() {
	installCmd.Flags().StringVarP(&installParams.Registry, "registry", "r", "github.com", "registry to use")
	installCmd.Flags().IntVarP(&installParams.Retries, "retries", "t", 3, "number of retries")
}

// ansurfen/docwiz@0.10
func resolvePackage(pkgName string) (owner, repo, version string, err error) {
	// Handle version specification like owner/repo@v1.0.0
	if strings.Contains(pkgName, "@") {
		parts := strings.Split(pkgName, "@")
		if len(parts) == 2 {
			pkgName = parts[0]
			version = parts[1]
		} else {
			return "", "", "", fmt.Errorf("invalid package name format: %s", pkgName)
		}
	}

	// Parse owner/repo
	if strings.Contains(pkgName, "/") {
		parts := strings.Split(pkgName, "/")
		if len(parts) == 2 {
			owner = parts[0]
			repo = parts[1]
		} else {
			return "", "", "", fmt.Errorf("invalid package name format: %s", pkgName)
		}
	} else {
		return "", "", "", fmt.Errorf("invalid package name format: %s", pkgName)
	}

	return owner, repo, version, nil
}

// extractZipToModules extracts a zip file to the modules directory
func extractZipToModules(zipPath, owner, repo, version, modulesDir string) error {
	// Create the target directory: modulesDir/owner/repo/version
	targetDir := filepath.Join(modulesDir, owner, repo, version)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Open the zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// Find the root directory name (usually repo-version)
	var rootDir string
	if len(reader.File) > 0 {
		firstFile := reader.File[0]
		parts := strings.Split(firstFile.Name, "/")
		if len(parts) > 0 {
			rootDir = parts[0]
		}
	}

	// Extract each file
	for _, file := range reader.File {
		// Skip directories
		if file.FileInfo().IsDir() {
			continue
		}

		// Remove the root directory prefix from the file path
		relativePath := file.Name
		if rootDir != "" && strings.HasPrefix(relativePath, rootDir+"/") {
			relativePath = strings.TrimPrefix(relativePath, rootDir+"/")
		}

		// Skip if the file is in the root directory itself
		if relativePath == "" || relativePath == file.Name {
			continue
		}

		// Create the file path
		filePath := filepath.Join(targetDir, relativePath)

		// Create the directory for the file if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for file %s: %w", relativePath, err)
		}

		// Create the file
		outFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", relativePath, err)
		}

		// Open the zip file entry
		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open zip entry %s: %w", relativePath, err)
		}

		// Copy the content
		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return fmt.Errorf("failed to copy file %s: %w", relativePath, err)
		}
	}

	log.WithField("target", targetDir).Info("extracted package successfully")
	return nil
}

// getHuloModulesDir returns the path to .hulo_modules directory in user's home
func getHuloModulesDir() (string, error) {
	modulesDir := os.Getenv("HULO_MODULES")
	if modulesDir != "" {
		return modulesDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	huloModulesDir := filepath.Join(homeDir, ".hulo_modules")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(huloModulesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create .hulo_modules directory: %w", err)
	}

	return huloModulesDir, nil
}
