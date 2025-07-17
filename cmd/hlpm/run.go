// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/util"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [script|file] [-- hulo-flags...]",
	Short: "run a package or script",
	Long:  "Run a script or file, optionally passing flags to the hulo compiler.",
	Example: `  hlpm run main.hl                    # Run file with default options
  hlpm run main.hl -- --verbose       # Run file with verbose flag
  hlpm run main.hl -- --target=vbs    # Run file with target flag
  hlpm run dev -- --watch             # Run script with watch flag`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		huloPkg, err = util.LoadConfigure[*config.HuloPkg](config.HuloPkgFileName)
		if err != nil {
			log.WithError(err).Fatal("fail to load package")
		}

		hlpmArgs, huloArgs := separateArgs(args)

		if len(hlpmArgs) == 0 {
			log.Fatal("no arguments provided")
		}

		switch {
		case isScript(args[0]):
			runScript(args[0], huloArgs)
		case isFile(args[0]):
			runFile(args[0], huloArgs)
		default:
			log.Fatal("invalid argument")
		}
	},
}

func separateArgs(args []string) ([]string, []string) {
	for i, arg := range args {
		if arg == "--" {
			// 找到分隔符，分离参数
			hlpmArgs := args[:i]
			huloArgs := args[i+1:]
			return hlpmArgs, huloArgs
		}
	}

	// 没有分隔符，所有参数都是 hlpm 的
	return args, nil
}

func isScript(script string) bool {
	if _, ok := huloPkg.Scripts[script]; ok {
		return true
	}
	return false
}

func isFile(path string) bool {
	if !strings.HasSuffix(path, ".hl") {
		return false
	}

	_, err := os.Stat(path)
	return err == nil
}

func runScript(script string, huloArgs []string) {
	scr := huloPkg.Scripts[script]
	cmd := exec.Command(scr, huloArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.WithError(err).Fatal("fail to run script")
	}
}

func runFile(file string, huloArgs []string) {
	executable, err := findHuloExecutable()
	if err != nil {
		log.WithError(err).Fatal("fail to find hulo executable")
	}

	log.WithField("file", file).
		WithField("executable", executable).
		WithField("args", huloArgs).
		Infof("run script")

	cmd := exec.Command(executable, append([]string{file}, huloArgs...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.WithError(err).Fatal("fail to run file")
	}
}

func findHuloExecutable() (string, error) {
	// 1. 查找当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	exe := "hulo"
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}

	huloPath := filepath.Join(currentDir, exe)
	if _, err := os.Stat(huloPath); err == nil {
		return huloPath, nil
	}

	// 2. 查找 hlpm 同目录
	hlpmDir := filepath.Dir(os.Args[0])
	huloPath = filepath.Join(hlpmDir, exe)
	if _, err := os.Stat(huloPath); err == nil {
		return huloPath, nil
	}

	// 3. 查找系统 PATH
	return exec.LookPath("hulo")
}
