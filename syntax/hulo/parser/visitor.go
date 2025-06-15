package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser/generated"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

type Visitor struct {
	*generated.BasehuloParserVisitor
}

func accept[T any](tree antlr.ParseTree, visitor antlr.ParseTreeVisitor) (T, bool) {
	if tree == nil {
		return *new(T), false
	}
	t, ok := tree.Accept(visitor).(T)
	return t, ok
}

// VisitIdentifier implements the Visitor interface for Identifier
func (v *Visitor) VisitIdentifier(node antlr.TerminalNode) interface{} {
	if node == nil {
		return nil
	}

	ident := &ast.Ident{
		NamePos: token.Pos(node.GetSymbol().GetStart()),
		Name:    node.GetText(),
	}

	return ident
}

// VisitLiteral implements the Visitor interface for Literal
func (v *Visitor) VisitLiteral(ctx *generated.LiteralContext) interface{} {
	if ctx == nil {
		return nil
	}

	if ctx.NumberLiteral() != nil {
		return &ast.BasicLit{
			Kind:     token.NUM,
			Value:    ctx.NumberLiteral().GetText(),
			ValuePos: token.Pos(ctx.NumberLiteral().GetSymbol().GetStart()),
		}
	}
	if ctx.BoolLiteral() != nil {
		return &ast.BasicLit{
			Kind:     token.IDENT, // Using IDENT for boolean literals
			Value:    ctx.BoolLiteral().GetText(),
			ValuePos: token.Pos(ctx.BoolLiteral().GetSymbol().GetStart()),
		}
	}
	if ctx.StringLiteral() != nil {
		return &ast.BasicLit{
			Kind:     token.STR,
			Value:    ctx.StringLiteral().GetText(),
			ValuePos: token.Pos(ctx.StringLiteral().GetSymbol().GetStart()),
		}
	}
	if ctx.NULL() != nil {
		return &ast.BasicLit{
			Kind:     token.IDENT, // Using IDENT for null
			Value:    "null",
			ValuePos: token.Pos(ctx.NULL().GetSymbol().GetStart()),
		}
	}
	return nil
}

// VisitConditionalExpression implements the Visitor interface for ConditionalExpression
func (v *Visitor) VisitConditionalExpression(ctx *generated.ConditionalExpressionContext) interface{} {
	if ctx == nil {
		return nil
	}

	cond, _ := accept[ast.Expr](ctx.ConditionalBoolExpression(), v)
	then, _ := accept[ast.Expr](ctx.ConditionalExpression(0), v)
	els, _ := accept[ast.Expr](ctx.ConditionalExpression(1), v)

	return &ast.BinaryExpr{
		X:     cond,
		OpPos: token.Pos(ctx.QUEST().GetSymbol().GetStart()),
		Op:    token.QUEST,
		Y: &ast.BinaryExpr{
			X:     then,
			OpPos: token.Pos(ctx.COLON().GetSymbol().GetStart()),
			Op:    token.COLON,
			Y:     els,
		},
	}
}

// VisitLogicalExpression implements the Visitor interface for LogicalExpression
func (v *Visitor) VisitLogicalExpression(ctx *generated.LogicalExpressionContext) interface{} {
	if ctx == nil {
		return nil
	}

	x, _ := accept[ast.Expr](ctx.ShiftExpression(0), v)
	y, _ := accept[ast.Expr](ctx.ShiftExpression(1), v)

	return &ast.BinaryExpr{
		X:     x,
		OpPos: token.Pos(ctx.GetLogicalOp().GetStart()),
		Op:    token.Token(ctx.GetLogicalOp().GetTokenType()),
		Y:     y,
	}
}

// VisitShiftExpression implements the Visitor interface for ShiftExpression
func (v *Visitor) VisitShiftExpression(ctx *generated.ShiftExpressionContext) interface{} {
	if ctx == nil {
		return nil
	}

	x, _ := accept[ast.Expr](ctx.AddSubExpression(0), v)
	y, _ := accept[ast.Expr](ctx.AddSubExpression(1), v)

	return &ast.BinaryExpr{
		X:     x,
		OpPos: token.Pos(ctx.GetShiftOp().GetStart()),
		Op:    token.Token(ctx.GetShiftOp().GetTokenType()),
		Y:     y,
	}
}

// VisitAddSubExpression implements the Visitor interface for AddSubExpression
func (v *Visitor) VisitAddSubExpression(ctx *generated.AddSubExpressionContext) interface{} {
	if ctx == nil {
		return nil
	}

	x, _ := accept[ast.Expr](ctx.MulDivExpression(0), v)
	y, _ := accept[ast.Expr](ctx.MulDivExpression(1), v)

	return &ast.BinaryExpr{
		X:     x,
		OpPos: token.Pos(ctx.GetAddSubOp().GetStart()),
		Op:    token.Token(ctx.GetAddSubOp().GetTokenType()),
		Y:     y,
	}
}

// VisitMulDivExpression implements the Visitor interface for MulDivExpression
func (v *Visitor) VisitMulDivExpression(ctx *generated.MulDivExpressionContext) interface{} {
	if ctx == nil {
		return nil
	}

	x, _ := accept[ast.Expr](ctx.IncDecExpression(0), v)
	y, _ := accept[ast.Expr](ctx.IncDecExpression(1), v)

	return &ast.BinaryExpr{
		X:     x,
		OpPos: token.Pos(ctx.GetMulDivOp().GetStart()),
		Op:    token.Token(ctx.GetMulDivOp().GetTokenType()),
		Y:     y,
	}
}

// VisitReturnStatement implements the Visitor interface for ReturnStatement
func (v *Visitor) VisitReturnStatement(ctx *generated.ReturnStatementContext) interface{} {
	if ctx == nil {
		return nil
	}

	var x ast.Expr
	if ctx.ExpressionList() != nil {
		exprList, _ := accept[[]ast.Expr](ctx.ExpressionList(), v)
		if len(exprList) > 0 {
			x = exprList[0]
		}
	}

	return &ast.ReturnStmt{
		Return: token.Pos(ctx.RETURN().GetSymbol().GetStart()),
		X:      x,
	}
}

// VisitBreakStatement implements the Visitor interface for BreakStatement
func (v *Visitor) VisitBreakStatement(ctx *generated.BreakStatementContext) interface{} {
	if ctx == nil {
		return nil
	}

	return &ast.BreakStmt{
		Break: token.Pos(ctx.BREAK().GetSymbol().GetStart()),
	}
}

// VisitContinueStatement implements the Visitor interface for ContinueStatement
func (v *Visitor) VisitContinueStatement(ctx *generated.ContinueStatementContext) interface{} {
	if ctx == nil {
		return nil
	}

	return &ast.ContinueStmt{
		Continue: token.Pos(ctx.CONTINUE().GetSymbol().GetStart()),
	}
}

// VisitIfStatement implements the Visitor interface for IfStatement
func (v *Visitor) VisitIfStatement(ctx *generated.IfStatementContext) interface{} {
	if ctx == nil {
		return nil
	}

	cond, _ := accept[ast.Expr](ctx.ConditionalExpression(), v)
	body, _ := accept[ast.Stmt](ctx.Block(0), v)

	var else_ ast.Stmt
	if len(ctx.AllBlock()) > 1 {
		else_, _ = accept[ast.Stmt](ctx.Block(1), v)
	}

	return &ast.IfStmt{
		If:   token.Pos(ctx.IF().GetSymbol().GetStart()),
		Cond: cond,
		Body: body.(*ast.BlockStmt),
		Else: else_,
	}
}

// VisitFile implements the Visitor interface for File
func (v *Visitor) VisitFile(ctx *generated.FileContext) interface{} {
	if ctx == nil {
		return nil
	}

	file := &ast.File{
		Imports: make(map[string]*ast.Import),
		Stmts:   make([]ast.Stmt, 0),
		Decls:   make([]ast.Decl, 0),
	}

	// Visit all statements
	for _, stmt := range ctx.AllStatement() {
		if s, ok := accept[ast.Stmt](stmt, v); ok {
			file.Stmts = append(file.Stmts, s)
		}
	}

	return file
}

// VisitStatement implements the Visitor interface for Statement
func (v *Visitor) VisitStatement(ctx *generated.StatementContext) interface{} {
	if ctx == nil {
		return nil
	}

	// Handle different types of statements
	if ctx.ExpressionStatement() != nil {
		return v.VisitExpressionStatement(ctx.ExpressionStatement().(*generated.ExpressionStatementContext))
	}
	if ctx.AssignStatement() != nil {
		return v.VisitAssignStatement(ctx.AssignStatement().(*generated.AssignStatementContext))
	}
	if ctx.ReturnStatement() != nil {
		return v.VisitReturnStatement(ctx.ReturnStatement().(*generated.ReturnStatementContext))
	}
	if ctx.BreakStatement() != nil {
		return v.VisitBreakStatement(ctx.BreakStatement().(*generated.BreakStatementContext))
	}
	if ctx.ContinueStatement() != nil {
		return v.VisitContinueStatement(ctx.ContinueStatement().(*generated.ContinueStatementContext))
	}
	if ctx.IfStatement() != nil {
		return v.VisitIfStatement(ctx.IfStatement().(*generated.IfStatementContext))
	}
	// TODO: Add more statement types

	return nil
}

// VisitExpressionStatement implements the Visitor interface for ExpressionStatement
func (v *Visitor) VisitExpressionStatement(ctx *generated.ExpressionStatementContext) interface{} {
	if ctx == nil {
		return nil
	}

	expr, _ := accept[ast.Expr](ctx.Expression(), v)
	return &ast.ExprStmt{
		X: expr,
	}
}

// VisitAssignStatement implements the Visitor interface for AssignStatement
func (v *Visitor) VisitAssignStatement(ctx *generated.AssignStatementContext) interface{} {
	if ctx == nil {
		return nil
	}

	lhs, _ := accept[ast.Expr](ctx.Expression(), v)
	rhs, _ := accept[ast.Expr](ctx.Expression(), v)

	return &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.ASSIGN, // Default to ASSIGN token
		Rhs: rhs,
	}
}
