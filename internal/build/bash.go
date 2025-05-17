package build

import (
	"fmt"

	bast "github.com/hulo-lang/hulo/syntax/bash/ast"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/internal/config"
)

func GenerateBash(opts *config.BashOptions, node hast.Node) error {
	bnode := translate2Bash(opts, node)

	fmt.Println("write to file system", bnode)
	return nil
}

func translate2Bash(opts *config.BashOptions, node hast.Node) bast.Node {
	return nil
}
