// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package strategies

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/caarlos0/log"

	bash "github.com/hulo-lang/hulo/internal/build/bash"
	vbs "github.com/hulo-lang/hulo/internal/build/vbs"
	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/vfs/osfs"
	"sigs.k8s.io/yaml"
)

var _ Strategy = (*CompileStrategy)(nil)

type CompileStrategy struct{}

func (c *CompileStrategy) CanHandle(params *Parameters) bool {
	return len(params.Args) > 0
}

func (c *CompileStrategy) Execute(params *Parameters, args []string) error {
	startTime := time.Now()

	verbose := os.Getenv("HULO_VERBOSE")
	if verbose == "true" {
		params.Verbose = true
	}

	if params.Verbose {
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
		err = yaml.Unmarshal(src, &huloc)
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

	hulopath := os.Getenv("HULO_PATH")
	if hulopath == "" {
		execPath, err := os.Executable()
		if err != nil {
			log.WithError(err).Fatal("fail to get executable path")
		}
		huloc.HuloPath = filepath.Dir(execPath)
	} else {
		huloc.HuloPath = hulopath
	}

	file := args[0]
	if file != "." {
		huloc.Main = file
	}

	// 优先级：命令行参数 > huloc.yaml
	if len(params.Langs) > 0 {
		huloc.Targets = params.Langs
	}

	// 如果 targets 中没有指定语言，则默认编译所有
	if len(huloc.Targets) == 0 {
		huloc.Targets = append(huloc.Targets, config.L_VBSCRIPT, config.L_BASH)
	}

	if len(huloc.OutDir) == 0 {
		huloc.OutDir = "."
	}

	if len(params.OutDir) > 0 {
		huloc.OutDir = params.OutDir
	}

	log.WithField("main", huloc.Main).
		WithField("include", huloc.Include).
		WithField("exclude", huloc.Exclude).
		Info("starting compile")

	wg := sync.WaitGroup{}

	for _, lang := range huloc.Targets {

		log.IncreasePadding()
		log.WithField("lang", lang).Info("compiling")
		log.DecreasePadding()

		wg.Add(1)
		go func() {
			defer wg.Done()
			var results map[string]string
			var err error
			switch lang {
			case config.L_VBSCRIPT:
				results, err = vbs.Transpile(huloc.CompilerOptions.VBScript, huloc.Main, localFs, ".", huloc.HuloPath)
				if err != nil {
					log.WithError(err).Info("fail to compile")
				}
			case config.L_BASH:
				results, err = bash.Transpile(huloc, localFs, ".", huloc.HuloPath, huloc.Main)
				if err != nil {
					log.WithError(err).Info("fail to compile")
				}
			}

			for file, code := range results {
				if strings.Contains(file, "std") {
					continue
				}
				os.MkdirAll(filepath.Join(huloc.OutDir, filepath.Dir(file)), 0755)
				err := localFs.WriteFile(filepath.Join(huloc.OutDir, file), []byte(code), 0644)
				if err != nil {
					log.WithError(err).Info("fail to write file")
				}
			}
		}()
	}

	wg.Wait()

	log.Infof("compile time: %s", time.Since(startTime))

	log.Info("thanks for using Hulo!")
	return nil
}

func (c *CompileStrategy) Priority() int {
	return PriorityCompile
}
