package unsafe

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hulo-lang/hulo/internal/unsafe/generated"
)

func Parse(template string) {
	stream := antlr.NewInputStream(template)
	lex := generated.NewunsafeLexer(stream)
	tokens := antlr.NewCommonTokenStream(lex, antlr.TokenDefaultChannel)
	parser := generated.NewunsafeParser(tokens)

	fmt.Println(parser.Template().ToStringTree(nil, parser))
}
