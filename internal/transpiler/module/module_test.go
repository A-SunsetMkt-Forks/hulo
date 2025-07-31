package module

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hulo-lang/hulo/internal/config"
	"github.com/hulo-lang/hulo/internal/container"
	"github.com/hulo-lang/hulo/internal/vfs/memvfs"
	"github.com/hulo-lang/hulo/internal/vfs/osfs"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.NoError(t, err)
}

// TODO Mock 工具自动生成代码进行单元测试

func TestMangle(t *testing.T) {
	code := `
	var a = 1
	loop {
		echo $a
		let a = 2
	}

	loop {
		let a = 3
	}

	loop {
		echo $a
	}

	loop $a := 0; $a < 10; $a++ {
		echo $a
	}

	fn test() {}`

	module := createMangleModule(t, code)
	err := MangleModule(module)
	require.NoError(t, err)

	ast.Print(module.AST)
	module.PrintSymbolTable()
}

func TestMangleLoop(t *testing.T) {
	code := `
	loop $a := 0; $a < 10; $a++ {
		echo $a
	}`

	module := createMangleModule(t, code)
	err := MangleModule(module)
	require.NoError(t, err)

	ast.Print(module.AST)
	module.PrintSymbolTable()
}

func createMangleModule(t *testing.T, code string) *Module {
	file, err := parser.ParseSourceScript(code)
	require.NoError(t, err)

	module := &Module{
		ModuleID: 0,
		Name:     "test",
		Path:     "/test/test.hl",
		AST:      file,
		Exports:  make(map[string]*ExportInfo),
		Imports:  []*ImportInfo{},
		Symbols:  nil,
		state:    ModuleStateUnresolved,
	}
	require.NoError(t, module.BuildSymbolTable())

	return module
}
