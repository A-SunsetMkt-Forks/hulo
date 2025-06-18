// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import (
	"regexp"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/caarlos0/log"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser/generated"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

type Visitor struct {
	*generated.BasehuloParserVisitor
	comments []*ast.Comment
}

func accept[T any](tree antlr.ParseTree, visitor antlr.ParseTreeVisitor) (T, bool) {
	if tree == nil {
		return *new(T), false
	}

	t, ok := tree.Accept(visitor).(T)
	return t, ok
}

// VisitIdentifier implements the Visitor interface for Identifier
func (v *Visitor) VisitIdentifier(node antlr.TerminalNode) any {
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
func (v *Visitor) VisitLiteral(ctx *generated.LiteralContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter Literal")
	log.IncreasePadding()
	defer log.Info("exit Literal")
	defer log.DecreasePadding()

	if ctx.NumberLiteral() != nil {
		return &ast.BasicLit{
			Kind:     token.NUM,
			Value:    ctx.NumberLiteral().GetText(),
			ValuePos: token.Pos(ctx.NumberLiteral().GetSymbol().GetStart()),
		}
	}
	if ctx.BoolLiteral() != nil {
		if ctx.BoolLiteral().GetText() == "true" {
			return &ast.BasicLit{
				Kind:     token.TRUE,
				Value:    "true",
				ValuePos: token.Pos(ctx.BoolLiteral().GetSymbol().GetStart()),
			}
		}
		return &ast.BasicLit{
			Kind:     token.FALSE,
			Value:    "false",
			ValuePos: token.Pos(ctx.BoolLiteral().GetSymbol().GetStart()),
		}
	}
	if ctx.StringLiteral() != nil {
		raw := ctx.StringLiteral().GetText()
		return &ast.BasicLit{
			Kind:     token.STR,
			Value:    raw[1 : len(raw)-1],
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
func (v *Visitor) VisitConditionalExpression(ctx *generated.ConditionalExpressionContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter ConditionalExpression")
	log.IncreasePadding()

	// Get the condition expression
	cond, _ := accept[ast.Expr](ctx.ConditionalBoolExpression(), v)

	// If there's no question mark, just return the condition
	if ctx.QUEST() == nil {
		log.DecreasePadding()
		log.Info("exit ConditionalExpression")
		return cond
	}

	// If there's a question mark, we have a ternary expression
	then, _ := accept[ast.Expr](ctx.ConditionalExpression(0), v)
	els, _ := accept[ast.Expr](ctx.ConditionalExpression(1), v)

	log.DecreasePadding()
	log.Info("exit ConditionalExpression")

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
func (v *Visitor) VisitLogicalExpression(ctx *generated.LogicalExpressionContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter LogicalExpression")
	log.IncreasePadding()

	var ret ast.Expr
	switch len(ctx.AllShiftExpression()) {
	case 1:
		ret, _ = accept[ast.Expr](ctx.ShiftExpression(0), v)
	case 2:
		x, _ := accept[ast.Expr](ctx.ShiftExpression(0), v)
		y, _ := accept[ast.Expr](ctx.ShiftExpression(1), v)

		ret = &ast.BinaryExpr{
			X:     x,
			OpPos: token.Pos(ctx.GetLogicalOp().GetStart()),
			Op:    logicalOpMap[ctx.GetLogicalOp().GetText()],
			Y:     y,
		}
	}

	log.DecreasePadding()
	log.Info("exit LogicalExpression")
	return ret
}

var logicalOpMap = map[string]token.Token{
	"&&": token.AND,
	"||": token.OR,
	"==": token.EQ,
	"!=": token.NEQ,
	"<":  token.LT,
	">":  token.GT,
	"<=": token.LE,
	">=": token.GE,
}

// VisitShiftExpression implements the Visitor interface for ShiftExpression
func (v *Visitor) VisitShiftExpression(ctx *generated.ShiftExpressionContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter ShiftExpression")
	log.IncreasePadding()

	var ret ast.Expr

	switch ctx.GetChildCount() {
	case 1:
		ret, _ = accept[ast.Expr](ctx.AddSubExpression(0), v)
	case 2:
		x, _ := accept[ast.Expr](ctx.AddSubExpression(0), v)
		y, _ := accept[ast.Expr](ctx.AddSubExpression(1), v)
		ret = &ast.BinaryExpr{
			X:     x,
			OpPos: token.Pos(ctx.GetShiftOp().GetStart()),
			Op:    token.Token(ctx.GetShiftOp().GetTokenType()),
			Y:     y,
		}
	}

	log.DecreasePadding()
	log.Info("exit ShiftExpression")
	return ret
}

// VisitAddSubExpression implements the Visitor interface for AddSubExpression
func (v *Visitor) VisitAddSubExpression(ctx *generated.AddSubExpressionContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter AddSubExpression")
	log.IncreasePadding()

	var ret ast.Expr
	switch ctx.GetChildCount() {
	case 1:
		ret, _ = accept[ast.Expr](ctx.MulDivExpression(0), v)
	case 2:
		x, _ := accept[ast.Expr](ctx.MulDivExpression(0), v)
		y, _ := accept[ast.Expr](ctx.MulDivExpression(1), v)
		ret = &ast.BinaryExpr{
			X:     x,
			OpPos: token.Pos(ctx.GetAddSubOp().GetStart()),
			Op:    token.Token(ctx.GetAddSubOp().GetTokenType()),
			Y:     y,
		}
	}

	log.DecreasePadding()
	log.Info("exit AddSubExpression")
	return ret
}

// VisitMulDivExpression implements the Visitor interface for MulDivExpression
func (v *Visitor) VisitMulDivExpression(ctx *generated.MulDivExpressionContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter MulDivExpression")
	log.IncreasePadding()

	var ret ast.Expr
	switch ctx.GetChildCount() {
	case 1:
		ret, _ = accept[ast.Expr](ctx.IncDecExpression(0), v)
	case 2:
		x, _ := accept[ast.Expr](ctx.IncDecExpression(0), v)
		y, _ := accept[ast.Expr](ctx.IncDecExpression(1), v)
		ret = &ast.BinaryExpr{
			X:     x,
			OpPos: token.Pos(ctx.GetMulDivOp().GetStart()),
			Op:    token.Token(ctx.GetMulDivOp().GetTokenType()),
			Y:     y,
		}
	}

	log.DecreasePadding()
	log.Info("exit MulDivExpression")
	return ret
}

// VisitReturnStatement implements the Visitor interface for ReturnStatement
func (v *Visitor) VisitReturnStatement(ctx *generated.ReturnStatementContext) any {
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
func (v *Visitor) VisitBreakStatement(ctx *generated.BreakStatementContext) any {
	if ctx == nil {
		return nil
	}

	return &ast.BreakStmt{
		Break: token.Pos(ctx.BREAK().GetSymbol().GetStart()),
	}
}

// VisitContinueStatement implements the Visitor interface for ContinueStatement
func (v *Visitor) VisitContinueStatement(ctx *generated.ContinueStatementContext) any {
	if ctx == nil {
		return nil
	}

	return &ast.ContinueStmt{
		Continue: token.Pos(ctx.CONTINUE().GetSymbol().GetStart()),
	}
}

// VisitIfStatement implements the Visitor interface for IfStatement
func (v *Visitor) VisitIfStatement(ctx *generated.IfStatementContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter IfStatement")
	log.IncreasePadding()

	cond, _ := accept[ast.Expr](ctx.ConditionalExpression(), v)
	body, _ := accept[ast.Stmt](ctx.Block(0), v)

	var else_ ast.Stmt
	if len(ctx.AllBlock()) > 1 {
		else_, _ = accept[ast.Stmt](ctx.Block(1), v)
	}

	log.DecreasePadding()
	log.Info("exit IfStatement")
	return &ast.IfStmt{
		If:   token.Pos(ctx.IF().GetSymbol().GetStart()),
		Cond: cond,
		Body: body.(*ast.BlockStmt),
		Else: else_,
	}
}

// VisitFile implements the Visitor interface for File
func (v *Visitor) VisitFile(ctx *generated.FileContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter file")
	log.IncreasePadding()
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
	if len(v.comments) != 0 {
		file.Docs = append(file.Docs, &ast.CommentGroup{
			List: v.comments,
		})
		v.comments = nil
	}

	log.DecreasePadding()
	log.Info("exit file")
	return file
}

// VisitStatement implements the Visitor interface for Statement
func (v *Visitor) VisitStatement(ctx *generated.StatementContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter statement")
	log.IncreasePadding()
	defer log.Info("exit statement")
	defer log.DecreasePadding()

	if ctx.Comment() != nil {
		if ctx.Comment().LineComment() != nil {
			cmt := v.fmtLineComment(ctx.Comment().LineComment().GetText())
			cmt.Slash = token.Pos(ctx.Comment().LineComment().GetSymbol().GetStart())
			v.comments = append(v.comments, cmt)
		} else {
			cmts := v.fmtBlockComment(ctx.Comment().BlockComment().GetText())
			for _, cmt := range cmts {
				cmt.Slash = token.Pos(ctx.Comment().BlockComment().GetSymbol().GetStart())
				v.comments = append(v.comments, cmt)
			}
		}
	}
	// Handle different types of statements
	if ctx.ExpressionStatement() != nil {
		return v.VisitExpressionStatement(ctx.ExpressionStatement().(*generated.ExpressionStatementContext))
	}
	if ctx.LambdaAssignStatement() != nil {
		return v.VisitLambdaAssignStatement(ctx.LambdaAssignStatement().(*generated.LambdaAssignStatementContext))
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
	if ctx.LoopStatement() != nil {
		return v.VisitLoopStatement(ctx.LoopStatement().(*generated.LoopStatementContext))
	}
	// TODO: Add more statement types

	return nil
}

// VisitExpressionStatement implements the Visitor interface for ExpressionStatement
func (v *Visitor) VisitExpressionStatement(ctx *generated.ExpressionStatementContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter ExpressionStatement")
	log.IncreasePadding()

	expr, _ := accept[ast.Expr](ctx.Expression(), v)

	log.DecreasePadding()
	log.Info("exit ExpressionStatement")
	return &ast.ExprStmt{
		X: expr,
	}
}

// VisitLambdaAssignStatement implements the Visitor interface for LambdaAssignStatement
func (v *Visitor) VisitLambdaAssignStatement(ctx *generated.LambdaAssignStatementContext) any {
	if ctx == nil {
		return nil
	}

	lhs, _ := accept[ast.Expr](ctx.VariableExpression(), v)
	rhs, _ := accept[ast.Expr](ctx.Expression(), v)

	return &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.COLON_ASSIGN, // Use COLON_ASSIGN token for lambda assignment
		Rhs: rhs,
	}
}

// VisitAssignStatement implements the Visitor interface for AssignStatement
func (v *Visitor) VisitAssignStatement(ctx *generated.AssignStatementContext) any {
	if ctx == nil {
		return nil
	}

	var scope token.Token
	var scopePos token.Pos

	// Handle scope modifiers (LET, CONST, VAR)
	if ctx.LET() != nil {
		scope = token.LET
		scopePos = token.Pos(ctx.LET().GetSymbol().GetStart())
	} else if ctx.CONST() != nil {
		scope = token.CONST
		scopePos = token.Pos(ctx.CONST().GetSymbol().GetStart())
	} else if ctx.VAR() != nil {
		scope = token.VAR
		scopePos = token.Pos(ctx.VAR().GetSymbol().GetStart())
	}

	// Get left hand side expression
	var lhs ast.Expr
	if ctx.Identifier() != nil {
		lhs = &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}
	} else if ctx.VariableNames() != nil {
		lhs, _ = accept[ast.Expr](ctx.VariableNames(), v)
	} else if ctx.VariableExpression() != nil {
		lhs, _ = accept[ast.Expr](ctx.VariableExpression(), v)
	} else if ctx.VariableNullableExpressions() != nil {
		lhs, _ = accept[ast.Expr](ctx.VariableNullableExpressions(), v)
	}

	// Get assignment operator
	var tok token.Token
	if ctx.ASSIGN() != nil {
		tok = token.ASSIGN
	} else if ctx.ADD_ASSIGN() != nil {
		tok = token.PLUS_ASSIGN
	} else if ctx.SUB_ASSIGN() != nil {
		tok = token.MINUS_ASSIGN
	} else if ctx.MUL_ASSIGN() != nil {
		tok = token.ASTERISK_ASSIGN
	} else if ctx.DIV_ASSIGN() != nil {
		tok = token.SLASH_ASSIGN
	} else if ctx.MOD_ASSIGN() != nil {
		tok = token.MOD_ASSIGN
	} else if ctx.AND_ASSIGN() != nil {
		tok = token.AND_ASSIGN
	} else if ctx.EXP_ASSIGN() != nil {
		tok = token.POWER_ASSIGN
	}

	// Get right hand side expression
	var rhs ast.Expr
	if ctx.Expression() != nil {
		rhs, _ = accept[ast.Expr](ctx.Expression(), v)
	} else if ctx.MatchStatement() != nil {
		rhs, _ = accept[ast.Expr](ctx.MatchStatement(), v)
	}

	return &ast.AssignStmt{
		Scope:    scope,
		ScopePos: scopePos,
		Lhs:      lhs,
		Tok:      tok,
		Rhs:      rhs,
	}
}

// VisitExpression implements the Visitor interface for Expression
func (v *Visitor) VisitExpression(ctx *generated.ExpressionContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter Expression")
	log.IncreasePadding()
	defer log.Info("exit Expression")
	defer log.DecreasePadding()

	// Handle different types of expressions
	if ctx.LambdaExpression() != nil {
		// return v.VisitLambdaExpression(ctx.LambdaExpression().(*generated.LambdaExpressionContext))
	}
	if ctx.ConditionalExpression() != nil {
		node := v.VisitConditionalExpression(ctx.ConditionalExpression().(*generated.ConditionalExpressionContext))
		return node
	}
	if ctx.NewDelExpression() != nil {
		// return v.VisitNewDelExpression(ctx.NewDelExpression().(*generated.NewDelExpressionContext))
	}
	if ctx.ClassInitializeExpression() != nil {
		// return v.VisitClassInitializeExpression(ctx.ClassInitializeExpression().(*generated.ClassInitializeExpressionContext))
	}
	if ctx.TypeofExpression() != nil {
		// return v.VisitTypeofExpression(ctx.TypeofExpression().(*generated.TypeofExpressionContext))
	}
	if ctx.ChannelOutputExpression() != nil {
		// return v.VisitChannelOutputExpression(ctx.ChannelOutputExpression().(*generated.ChannelOutputExpressionContext))
	}
	if ctx.CommandExpression() != nil {
		return v.VisitCommandExpression(ctx.CommandExpression().(*generated.CommandExpressionContext))
	}
	if ctx.UnsafeExpression() != nil {
		// return v.VisitUnsafeExpression(ctx.UnsafeExpression().(*generated.UnsafeExpressionContext))
	}
	if ctx.ComptimeExpression() != nil {
		// return v.VisitComptimeExpression(ctx.ComptimeExpression().(*generated.ComptimeExpressionContext))
	}

	return nil
}

// VisitCommandExpression implements the Visitor interface for CommandExpression
func (v *Visitor) VisitCommandExpression(ctx *generated.CommandExpressionContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter CommandExpression")
	log.IncreasePadding()

	var fun ast.Expr
	var recv []ast.Expr

	// Handle command string literal or member access
	if ctx.CommandStringLiteral() != nil {
		fun = &ast.BasicLit{
			Kind:     token.STR,
			Value:    ctx.CommandStringLiteral().GetText(),
			ValuePos: token.Pos(ctx.CommandStringLiteral().GetSymbol().GetStart()),
		}
	} else if ctx.MemberAccess() != nil {
		fun, _ = accept[ast.Expr](ctx.MemberAccess(), v)
	}

	// Handle options and arguments
	for _, opt := range ctx.AllOption() {
		optExpr, _ := accept[ast.Expr](opt, v)
		if optExpr != nil {
			recv = append(recv, optExpr)
		}
	}

	for _, expr := range ctx.AllConditionalExpression() {
		arg, _ := accept[ast.Expr](expr, v)
		if arg != nil {
			recv = append(recv, arg)
		}
	}

	// Create CallExpr
	call := &ast.CallExpr{
		Fun:    fun,
		Lparen: token.Pos(ctx.GetStart().GetStart()),
		Recv:   recv,
		Rparen: token.Pos(ctx.GetStop().GetStop()),
	}

	// Handle command join or stream if present
	if ctx.CommandJoin() != nil {
		join, _ := accept[ast.Expr](ctx.CommandJoin(), v)
		if join != nil {
			// Create a new CallExpr for the joined command
			call = &ast.CallExpr{
				Fun:    call,
				Lparen: token.Pos(ctx.GetStart().GetStart()),
				Recv:   []ast.Expr{join},
				Rparen: token.Pos(ctx.GetStop().GetStop()),
			}
		}
	} else if ctx.CommandStream() != nil {
		stream, _ := accept[ast.Expr](ctx.CommandStream(), v)
		if stream != nil {
			// Create a new CallExpr for the streamed command
			call = &ast.CallExpr{
				Fun:    call,
				Lparen: token.Pos(ctx.GetStart().GetStart()),
				Recv:   []ast.Expr{stream},
				Rparen: token.Pos(ctx.GetStop().GetStop()),
			}
		}
	}

	log.DecreasePadding()
	log.Info("exit CommandExpression")
	return call
}

func (v *Visitor) VisitMemberAccess(ctx *generated.MemberAccessContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter member access")
	log.IncreasePadding()
	defer log.Info("exit member access")
	defer log.DecreasePadding()

	// Handle identifier with generic arguments
	if ctx.Identifier() != nil && ctx.GenericArguments() != nil {
		ident := ctx.Identifier().GetText()
		genericArgs := v.Visit(ctx.GenericArguments()).(*ast.GenericExpr)
		return &ast.CallExpr{
			Fun: &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ident,
			},
			Generic: genericArgs,
		}
	}

	// Handle identifier with member access point
	if ctx.Identifier() != nil && ctx.MemberAccessPoint() != nil {
		ident := ctx.Identifier().GetText()
		base := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ident,
		}
		return v.visitMemberAccessPoint(base, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
	}

	// Handle STR, NUM, BOOL with member access point
	if ctx.STR() != nil || ctx.NUM() != nil || ctx.BOOL() != nil {
		var kind token.Token
		var value string
		var pos token.Pos

		if ctx.STR() != nil {
			kind = token.STR
			value = ctx.STR().GetText()
			pos = token.Pos(ctx.STR().GetSymbol().GetStart())
		} else if ctx.NUM() != nil {
			kind = token.NUM
			value = ctx.NUM().GetText()
			pos = token.Pos(ctx.NUM().GetSymbol().GetStart())
		} else {
			// For BOOL, we need to check if it's true or false
			if ctx.BOOL().GetText() == "true" {
				kind = token.TRUE
				value = "true"
				pos = token.Pos(ctx.BOOL().GetSymbol().GetStart())
			} else {
				kind = token.FALSE
				value = "false"
				pos = token.Pos(ctx.BOOL().GetSymbol().GetStart())
			}
		}

		base := &ast.BasicLit{
			Kind:     kind,
			Value:    value,
			ValuePos: pos,
		}
		return v.visitMemberAccessPoint(base, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
	}

	// Handle literal with member access point
	if ctx.Literal() != nil {
		base := v.Visit(ctx.Literal()).(ast.Expr)
		return v.visitMemberAccessPoint(base, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
	}

	// Handle THIS with optional member access point
	if ctx.THIS() != nil {
		base := &ast.Ident{
			NamePos: token.Pos(ctx.THIS().GetSymbol().GetStart()),
			Name:    "this",
		}
		if ctx.MemberAccessPoint() != nil {
			return v.visitMemberAccessPoint(base, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return base
	}

	// Handle SUPER with optional member access point
	if ctx.SUPER() != nil {
		base := &ast.Ident{
			NamePos: token.Pos(ctx.SUPER().GetSymbol().GetStart()),
			Name:    "super",
		}
		if ctx.MemberAccessPoint() != nil {
			return v.visitMemberAccessPoint(base, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return base
	}

	if ctx.Identifier() != nil {
		return &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}
	}

	return nil
}

func (v *Visitor) visitMemberAccessPoint(base ast.Expr, ctx *generated.MemberAccessPointContext) ast.Expr {
	if ctx == nil {
		return base
	}
	log.Info("enter member access point")
	log.IncreasePadding()
	defer log.Info("exit member access point")
	defer log.DecreasePadding()

	// Handle dot access
	if ctx.DOT() != nil && ctx.Identifier() != nil {
		ident := ctx.Identifier().GetText()
		selectExpr := &ast.SelectExpr{
			X:   base,
			Dot: token.Pos(ctx.DOT().GetSymbol().GetStart()),
			Y: &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ident,
			},
		}

		// Handle generic arguments if present
		if ctx.GenericArguments() != nil {
			genericArgs := v.Visit(ctx.GenericArguments()).(*ast.GenericExpr)
			selectExpr.Y = &ast.CallExpr{
				Fun:     selectExpr.Y,
				Generic: genericArgs,
			}
		}

		// Handle recursive member access point
		if ctx.MemberAccessPoint() != nil {
			return v.visitMemberAccessPoint(selectExpr, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return selectExpr
	}

	// Handle double colon access
	if ctx.DOUBLE_COLON() != nil && ctx.Identifier() != nil {
		ident := ctx.Identifier().GetText()
		selectExpr := &ast.SelectExpr{
			X:   base,
			Dot: token.Pos(ctx.DOUBLE_COLON().GetSymbol().GetStart()),
			Y: &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ident,
			},
		}

		// Handle recursive member access point
		if ctx.MemberAccessPoint() != nil {
			return v.visitMemberAccessPoint(selectExpr, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return selectExpr
	}

	// Handle index access
	if ctx.LBRACK() != nil && ctx.Expression() != nil && ctx.RBRACK() != nil {
		indexExpr := &ast.IndexExpr{
			X:      base,
			Lbrack: token.Pos(ctx.LBRACK().GetSymbol().GetStart()),
			Index:  v.Visit(ctx.Expression()).(ast.Expr),
			Rbrack: token.Pos(ctx.RBRACK().GetSymbol().GetStart()),
		}

		// Handle recursive member access point
		if ctx.MemberAccessPoint() != nil {
			return v.visitMemberAccessPoint(indexExpr, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return indexExpr
	}

	return base
}

// VisitConditionalBoolExpression implements the Visitor interface for ConditionalBoolExpression
func (v *Visitor) VisitConditionalBoolExpression(ctx *generated.ConditionalBoolExpressionContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter ConditionalBoolExpression")
	log.IncreasePadding()

	// Get the first logical expression
	x, _ := accept[ast.Expr](ctx.LogicalExpression(0), v)

	if len(ctx.AllLogicalExpression()) > 1 {
		// Handle multiple logical expressions with operators
		for i := 1; i < len(ctx.AllLogicalExpression()); i++ {
			y, _ := accept[ast.Expr](ctx.LogicalExpression(i), v)
			op := token.Token(ctx.GetConditionalOp().GetTokenType())
			opPos := token.Pos(ctx.GetConditionalOp().GetStart())

			x = &ast.BinaryExpr{
				X:     x,
				OpPos: opPos,
				Op:    op,
				Y:     y,
			}
		}
	}

	log.DecreasePadding()
	log.Info("exit ConditionalBoolExpression")

	return x
}

// VisitIncDecExpression implements the Visitor interface for IncDecExpression
func (v *Visitor) VisitIncDecExpression(ctx *generated.IncDecExpressionContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter IncDecExpression")
	log.IncreasePadding()

	var ret ast.Expr

	if ctx.PreIncDecExpression() != nil {
		ret, _ = accept[ast.Expr](ctx.PreIncDecExpression(), v)
	}

	if ctx.PostIncDecExpression() != nil {
		ret, _ = accept[ast.Expr](ctx.PostIncDecExpression(), v)
	}

	log.DecreasePadding()
	log.Info("exit IncDecExpression")

	return ret
}

// VisitPreIncDecExpression implements the Visitor interface for PreIncDecExpression
func (v *Visitor) VisitPreIncDecExpression(ctx *generated.PreIncDecExpressionContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter PreIncDecExpression")
	log.IncreasePadding()

	expr := v.VisitFactor(ctx.Factor().(*generated.FactorContext))
	if expr == nil {
		return nil
	}

	var ret ast.Expr
	if ctx.INC() != nil {
		ret = &ast.IncDecExpr{
			Pre:    true,
			X:      expr.(ast.Expr),
			Tok:    token.INC,
			TokPos: token.Pos(ctx.GetStart().GetStart()),
		}
	} else if ctx.DEC() != nil {
		ret = &ast.IncDecExpr{
			Pre:    true,
			X:      expr.(ast.Expr),
			Tok:    token.DEC,
			TokPos: token.Pos(ctx.GetStart().GetStart()),
		}
	} else {
		ret = expr.(ast.Expr)
	}

	log.DecreasePadding()
	log.Info("exit PreIncDecExpression")

	return ret
}

// VisitPostIncDecExpression implements the Visitor interface for PostIncDecExpression
func (v *Visitor) VisitPostIncDecExpression(ctx *generated.PostIncDecExpressionContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter PostIncDecExpression")
	log.IncreasePadding()

	expr := v.VisitFactor(ctx.Factor().(*generated.FactorContext))

	if expr == nil {
		log.DecreasePadding()
		log.Info("exit PostIncDecExpression")
		return nil
	}

	var tok token.Token
	if ctx.INC() != nil {
		tok = token.INC
	} else if ctx.DEC() != nil {
		tok = token.DEC
	}

	log.DecreasePadding()
	log.Info("exit PostIncDecExpression")
	return &ast.IncDecExpr{
		Pre:    false,
		X:      expr.(ast.Expr),
		Tok:    tok,
		TokPos: token.Pos(ctx.GetStart().GetStart()),
	}
}

// VisitFactor implements the Visitor interface for Factor
func (v *Visitor) VisitFactor(ctx *generated.FactorContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter Factor")
	log.IncreasePadding()
	defer log.Info("exit Factor")
	defer log.DecreasePadding()

	// Handle unary expressions
	if ctx.SUB() != nil {
		expr := v.VisitFactor(ctx.Factor().(*generated.FactorContext))
		if expr == nil {
			return nil
		}
		return &ast.UnaryExpr{
			OpPos: token.Pos(ctx.SUB().GetSymbol().GetStart()),
			Op:    token.MINUS,
			X:     expr.(ast.Expr),
		}
	}

	// Handle literals
	if ctx.Literal() != nil {
		return v.VisitLiteral(ctx.Literal().(*generated.LiteralContext))
	}

	// Handle identifiers
	if ctx.VariableExpression() != nil {
		return v.VisitVariableExpression(ctx.VariableExpression().(*generated.VariableExpressionContext))
	}

	// Handle parenthesized expressions
	if ctx.LPAREN() != nil {
		return v.VisitFactor(ctx.Factor().(*generated.FactorContext))
	}

	return nil
}

func (v *Visitor) VisitVariableExpression(ctx *generated.VariableExpressionContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter VariableExpression")
	log.IncreasePadding()

	ret := v.VisitMemberAccess(ctx.MemberAccess().(*generated.MemberAccessContext))

	log.DecreasePadding()
	log.Info("exit VariableExpression")
	return &ast.RefExpr{
		X: ret.(ast.Expr),
	}
}

func (v *Visitor) fmtLineComment(cmt string) *ast.Comment {
	// remove the first two characters
	return &ast.Comment{Text: cmt[2:]}
}

var removeStar = regexp.MustCompile(`^\s*\*`)

// fmtBlockComment formats a block comment.
// Because most compiled languages do not support multi-line comments,
// so we need to convert them into single-line comments
func (v *Visitor) fmtBlockComment(cmt string) (ret []*ast.Comment) {
	ret = make([]*ast.Comment, 0)
	if strings.HasPrefix(cmt, "/**") {
		cmt = cmt[3 : len(cmt)-2]
		lines := strings.SplitSeq(cmt, "\n")
		for line := range lines {
			// remove the first * and the space before it
			line = removeStar.ReplaceAllString(line, "")
			ret = append(ret, &ast.Comment{Text: line})
		}
		return ret
	}
	lines := strings.SplitSeq(cmt[2:len(cmt)-2], "\n")
	for line := range lines {
		ret = append(ret, &ast.Comment{Text: line})
	}

	return ret
}

// VisitBlock implements the Visitor interface for Block
func (v *Visitor) VisitBlock(ctx *generated.BlockContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter Block")
	log.IncreasePadding()

	// Create a new block statement
	block := &ast.BlockStmt{
		Lbrace: token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
		Rbrace: token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
	}

	// Process all statements in the block
	for _, stmtCtx := range ctx.AllStatement() {
		stmt, _ := accept[ast.Stmt](stmtCtx, v)
		if stmt != nil {
			block.List = append(block.List, stmt)
		}
	}

	log.DecreasePadding()
	log.Info("exit Block")

	return block
}

// VisitLoopStatement implements the Visitor interface for LoopStatement
func (v *Visitor) VisitLoopStatement(ctx *generated.LoopStatementContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter LoopStatement")
	log.IncreasePadding()
	defer log.Info("exit LoopStatement")
	defer log.DecreasePadding()

	// Handle different types of loops
	if ctx.WhileStatement() != nil {
		return v.VisitWhileStatement(ctx.WhileStatement().(*generated.WhileStatementContext))
	} else if ctx.DoWhileStatement() != nil {
		return v.VisitDoWhileStatement(ctx.DoWhileStatement().(*generated.DoWhileStatementContext))
	} else if ctx.RangeStatement() != nil {
		return v.VisitRangeStatement(ctx.RangeStatement().(*generated.RangeStatementContext))
	} else if ctx.ForStatement() != nil {
		return v.VisitForStatement(ctx.ForStatement().(*generated.ForStatementContext))
	} else if ctx.ForeachStatement() != nil {
		return v.VisitForeachStatement(ctx.ForeachStatement().(*generated.ForeachStatementContext))
	}

	return nil
}

// VisitWhileStatement implements the Visitor interface for WhileStatement
func (v *Visitor) VisitWhileStatement(ctx *generated.WhileStatementContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter WhileStatement")
	log.IncreasePadding()

	body, _ := accept[ast.Stmt](ctx.Block(), v)

	var cond ast.Expr
	if ctx.Expression() != nil {
		cond, _ = accept[ast.Expr](ctx.Expression(), v)
	}

	log.DecreasePadding()
	log.Info("exit WhileStatement")

	return &ast.WhileStmt{
		Loop: token.Pos(ctx.LOOP().GetSymbol().GetStart()),
		Cond: cond,
		Body: body.(*ast.BlockStmt),
	}
}

// VisitDoWhileStatement implements the Visitor interface for DoWhileStatement
func (v *Visitor) VisitDoWhileStatement(ctx *generated.DoWhileStatementContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter DoWhileStatement")
	log.IncreasePadding()

	body, _ := accept[ast.Stmt](ctx.Block(), v)
	cond, _ := accept[ast.Expr](ctx.Expression(), v)

	log.DecreasePadding()
	log.Info("exit DoWhileStatement")

	return &ast.DoWhileStmt{
		Do:     token.Pos(ctx.DO().GetSymbol().GetStart()),
		Body:   body.(*ast.BlockStmt),
		Loop:   token.Pos(ctx.LOOP().GetSymbol().GetStart()),
		Lparen: token.Pos(ctx.LPAREN().GetSymbol().GetStart()),
		Cond:   cond,
		Rparen: token.Pos(ctx.RPAREN().GetSymbol().GetStart()),
	}
}

// VisitRangeStatement implements the Visitor interface for RangeStatement
func (v *Visitor) VisitRangeStatement(ctx *generated.RangeStatementContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter RangeStatement")
	log.IncreasePadding()

	index := &ast.Ident{
		NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
		Name:    ctx.Identifier().GetText(),
	}

	rangeExpr, _ := accept[ast.Expr](ctx.RangeClause(), v)
	body, _ := accept[ast.Stmt](ctx.Block(), v)

	log.DecreasePadding()
	log.Info("exit RangeStatement")

	return &ast.RangeStmt{
		Loop:   token.Pos(ctx.LOOP().GetSymbol().GetStart()),
		Index:  index,
		In:     token.Pos(ctx.IN().GetSymbol().GetStart()),
		Range:  token.Pos(ctx.RangeClause().RANGE().GetSymbol().GetStart()),
		Lparen: token.Pos(ctx.RangeClause().LPAREN().GetSymbol().GetStart()),
		RangeClauseExpr: ast.RangeClauseExpr{
			Start: rangeExpr.(*ast.BinaryExpr).X,
			End:   rangeExpr.(*ast.BinaryExpr).Y,
		},
		Rparen: token.Pos(ctx.RangeClause().RPAREN().GetSymbol().GetStart()),
		Body:   body.(*ast.BlockStmt),
	}
}

// VisitForStatement implements the Visitor interface for ForStatement
func (v *Visitor) VisitForStatement(ctx *generated.ForStatementContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter ForStatement")
	log.IncreasePadding()

	var (
		init           ast.Stmt
		cond           ast.Expr
		post           ast.Expr
		comma1, comma2 token.Pos
	)

	if ctx.ForClause() != nil {
		if ctx.ForClause().Statement() != nil {
			init, _ = accept[ast.Stmt](ctx.ForClause().Statement(), v)
		}
		if ctx.ForClause().Expression(0) != nil {
			cond, _ = accept[ast.Expr](ctx.ForClause().Expression(0), v)
		}
		if ctx.ForClause().Expression(1) != nil {
			post, _ = accept[ast.Expr](ctx.ForClause().Expression(1), v)
		}
		if len(ctx.ForClause().AllSEMI()) > 0 {
			comma1 = token.Pos(ctx.ForClause().SEMI(0).GetSymbol().GetStart())
		}
		if len(ctx.ForClause().AllSEMI()) > 1 {
			comma2 = token.Pos(ctx.ForClause().SEMI(1).GetSymbol().GetStart())
		}
	}

	body, _ := accept[ast.Stmt](ctx.Block(), v)

	var lparen, rparen token.Pos
	if ctx.LPAREN() != nil {
		lparen = token.Pos(ctx.LPAREN().GetSymbol().GetStart())
	}
	if ctx.RPAREN() != nil {
		rparen = token.Pos(ctx.RPAREN().GetSymbol().GetStart())
	}

	log.DecreasePadding()
	log.Info("exit ForStatement")

	return &ast.ForStmt{
		Loop:   token.Pos(ctx.LOOP().GetSymbol().GetStart()),
		Lparen: lparen,
		Init:   init,
		Comma1: comma1,
		Cond:   cond,
		Comma2: comma2,
		Post:   post,
		Rparen: rparen,
		Body:   body.(*ast.BlockStmt),
	}
}

// VisitForeachStatement implements the Visitor interface for ForeachStatement
func (v *Visitor) VisitForeachStatement(ctx *generated.ForeachStatementContext) any {
	if ctx == nil {
		return nil
	}

	log.Info("enter ForeachStatement")
	log.IncreasePadding()

	index, _ := accept[ast.Expr](ctx.ForeachClause().VariableName(0), v)
	var value ast.Expr
	if len(ctx.ForeachClause().AllVariableName()) > 1 {
		value, _ = accept[ast.Expr](ctx.ForeachClause().VariableName(1), v)
	}

	expr, _ := accept[ast.Expr](ctx.Expression(), v)
	body, _ := accept[ast.Stmt](ctx.Block(), v)

	log.DecreasePadding()
	log.Info("exit ForeachStatement")

	return &ast.ForeachStmt{
		Loop:   token.Pos(ctx.LOOP().GetSymbol().GetStart()),
		Lparen: token.Pos(ctx.LPAREN().GetSymbol().GetStart()),
		Index:  index,
		Value:  value,
		Rparen: token.Pos(ctx.RPAREN().GetSymbol().GetStart()),
		In:     token.Pos(ctx.IN().GetSymbol().GetStart()),
		Var:    expr,
		Body:   body.(*ast.BlockStmt),
	}
}
