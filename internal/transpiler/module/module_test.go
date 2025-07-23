package transpiler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	"github.com/hulo-lang/hulo/internal/vfs/osfs"
)

func TestResolveAllDependencies(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	HULO_MODULES := filepath.Join(wd, `testdata/hulo_modules`)
	HULO_PATH := filepath.Join(wd, `testdata/hulo_path`)
	WORK_DIR := "./testdata"
	fs := memvfs.New()
	testdata := map[string]string{
		HULO_MODULES + "/hulo-lang/hello-world/1.0.0/index.hl": `
		import "hulo-lang/hello-world"`,
		HULO_MODULES + "/hulo-lang/hello-world/1.0.0/utils_calculate.hl": `
		pub fn add(a: num, b: num) -> num {
			return $a + $b
		}`,
		HULO_MODULES + "/hulo-lang/hello-world/1.0.0/hulo.pkg.yaml": `
		version: 1.0.0
		main: ./index.hl`,
		HULO_PATH + "/core/math/index.hl": `
pub fn add(a: num, b: num) => $a + $b

pub fn mul(a: num, b: num) => $a * $b`,
		WORK_DIR + "/src/main.hl": `
		import "hulo-lang/hello-world" as h
		import "hulo-lang/hello-world/utils_calculate" as uc
		import * from "math"
		import "math"

		`,
		WORK_DIR + "/hulo.pkg.yaml": `
		dependencies:
			hulo-lang/hello-world: 1.0.0
		main: ./src/main.hl
		`,
	}

	for path, content := range testdata {
		fs.WriteFile(path, []byte(content), 0644)
	}

	resolver := &DependecyResolver{
		fs:              osfs.New(),
		visited:         container.NewMapSet[string](),
		stack:           container.NewMapSet[string](),
		order:           []string{},
		modules:         make(map[string]*Module),
		pkgs:            make(map[string]*config.HuloPkg),
		options:         &config.Huloc{},
		huloModulesPath: HULO_MODULES,
		huloPath:        HULO_PATH,
	}

	err = ResolveAllDependencies(resolver, "./testdata/src/main.hl")
	if err != nil {
		t.Fatalf("ResolveAllDependencies: %v", err)
	}
}

// TODO Mock 工具自动生成代码进行单元测试
