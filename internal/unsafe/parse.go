package unsafe

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

func Parse(template string) {
	stream := antlr.NewInputStream(template)
	lex := NewunsafeLexer(stream)
	tokens := antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
	parser := NewUnsafeParser(tokens)

	fmt.Println(parser.Template().ToStringTree(nil, parser))
}
