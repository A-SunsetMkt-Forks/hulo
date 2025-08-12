package batch

import (
	"github.com/hulo-lang/hulo/internal/transpiler"
	bast "github.com/hulo-lang/hulo/syntax/batch/ast"
	btok "github.com/hulo-lang/hulo/syntax/batch/token"
	hast "github.com/hulo-lang/hulo/syntax/hulo/ast"
)

type RuleID uint

//go:generate stringer -type=RuleID -linecomment -output=hcr_string.go
const (
	RuleIllegal       RuleID = iota // illegal
	RuleCommentSyntax               // comment_syntax
	RuleBoolFormat                  // bool_format
)

type BoolAsNumberHandler struct{}

func (c *BoolAsNumberHandler) Strategy() string {
	return "number"
}

func (c *BoolAsNumberHandler) Apply(transpiler transpiler.Transpiler[bast.Node], node hast.Node) (bast.Node, error) {
	if _, ok := node.(*hast.TrueLiteral); ok {
		return &bast.Lit{Val: "1"}, nil
	}
	return &bast.Lit{Val: "0"}, nil
}

type BoolAsStringHandler struct{}

func (c *BoolAsStringHandler) Strategy() string {
	return "string"
}

func (c *BoolAsStringHandler) Apply(transpiler transpiler.Transpiler[bast.Node], node hast.Node) (bast.Node, error) {
	if _, ok := node.(*hast.TrueLiteral); ok {
		return &bast.Lit{Val: `"true"`}, nil
	}
	return &bast.Lit{Val: `"false"`}, nil
}

type BoolAsCmdHandler struct{}

func (c *BoolAsCmdHandler) Strategy() string {
	return "cmd"
}

func (c *BoolAsCmdHandler) Apply(transpiler transpiler.Transpiler[bast.Node], node hast.Node) (bast.Node, error) {
	return nil, nil
}

type MultiStringConvertor struct{}

type REMCommentHandler struct{}

func (c *REMCommentHandler) Strategy() string {
	return "rem"
}

func (c *REMCommentHandler) Apply(transpiler transpiler.Transpiler[bast.Node], node hast.Node) (bast.Node, error) {
	cg := node.(*hast.CommentGroup)
	docs := make([]*bast.Comment, len(cg.List))
	for i, d := range cg.List {
		docs[i] = &bast.Comment{Tok: btok.REM, Text: d.Text}
	}
	return &bast.CommentGroup{Comments: docs}, nil
}

type DoubleColonCommentHandler struct{}

func (c *DoubleColonCommentHandler) Strategy() string {
	return "::"
}

func (c *DoubleColonCommentHandler) Apply(transpiler transpiler.Transpiler[bast.Node], node hast.Node) (bast.Node, error) {
	cg := node.(*hast.CommentGroup)
	docs := make([]*bast.Comment, len(cg.List))
	for i, d := range cg.List {
		docs[i] = &bast.Comment{Tok: btok.DOUBLE_COLON, Text: d.Text}
	}
	return &bast.CommentGroup{Comments: docs}, nil
}
