package interpreter

import (
	"testing"

	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser"
	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	script := `loop {
			echo "Hello, World!"
		}

		do {
			echo "Hello, World!"
		} loop ($a > 10)

		loop $i := 0; $i < 10; $i++ {
			echo "Hello, World!"
		}

		comptime {
			echo "Hello, World!"
		}`
	var node hast.Node
	var err error
	node, err = parser.ParseSourceScript(script)
	assert.NoError(t, err)
	hast.Print(node)
	interp := &Interpreter{}
	node = interp.Eval(node)
	hast.Print(node)
}
