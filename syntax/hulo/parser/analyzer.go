// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package parser

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hulo-lang/hulo/syntax/hulo/ast"
	"github.com/hulo-lang/hulo/syntax/hulo/parser/generated"
	"github.com/hulo-lang/hulo/syntax/hulo/token"
)

type Analyzer struct {
	*generated.BasehuloParserVisitor
	tokens   antlr.TokenStream
	lexer    antlr.Lexer
	parser   antlr.Parser
	file     generated.IFileContext
	comments []*ast.Comment

	*Tracer
}

func NewAnalyzer(input antlr.CharStream, opts ...ParserOptions) (*Analyzer, error) {
	lexer := generated.NewhuloLexer(input)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := generated.NewhuloParser(tokens)

	analyzer := &Analyzer{
		file:   parser.File(),
		tokens: tokens,
		lexer:  lexer,
		parser: parser,
		Tracer: NewTracer(),
	}

	for _, opt := range opts {
		if err := opt(analyzer); err != nil {
			return nil, err
		}
	}

	return analyzer, nil
}

func (a *Analyzer) getPos(ctx antlr.Token) token.Pos {
	if ctx == nil {
		return token.NoPos
	}
	return token.Pos(ctx.GetStart())
}

func (a *Analyzer) getTerminalPos(node antlr.TerminalNode) token.Pos {
	if node == nil {
		return token.NoPos
	}
	return token.Pos(node.GetSymbol().GetStart())
}

func accept[T any](tree antlr.ParseTree, analyzer *Analyzer) (T, bool) {
	if tree == nil {
		return *new(T), false
	}

	result, ok := tree.Accept(analyzer).(T)
	if !ok {
		var expectedType string
		t := reflect.TypeOf((*T)(nil)).Elem()
		if t.Kind() == reflect.Interface && t.Name() == "" {
			expectedType = "interface{}"
		} else {
			expectedType = t.String()
		}

		analyzer.EmitError(fmt.Errorf(
			"type conversion failed: expected %s, got %T",
			expectedType,
			result,
		))
		return *new(T), false
	}
	return result, true
}

func (a *Analyzer) visitWrapper(name string, ctx antlr.ParserRuleContext, visit func() any) any {
	if ctx == nil {
		return nil
	}

	pos := Position{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}

	a.Enter(name, fmt.Sprintf("%T", ctx), pos)

	var result any
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in %s: %v", name, r)
		}
		a.Exit(result, err)
	}()

	result = visit()
	return result
}

// VisitIdentifier implements the Visitor interface for Identifier
func (a *Analyzer) VisitIdentifier(node antlr.TerminalNode) any {
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
func (a *Analyzer) VisitLiteral(ctx *generated.LiteralContext) any {
	return a.visitWrapper("Literal", ctx, func() any {

		var (
			value    string
			valuePos token.Pos
		)

		switch {
		case ctx.NumberLiteral() != nil:
			value = ctx.NumberLiteral().GetText()
			valuePos = a.getTerminalPos(ctx.NumberLiteral())
			return &ast.NumericLiteral{
				Value:    value,
				ValuePos: valuePos,
			}
		case ctx.BoolLiteral() != nil:
			value := ctx.BoolLiteral().GetText()
			if value == "true" {
				return &ast.TrueLiteral{
					ValuePos: a.getTerminalPos(ctx.BoolLiteral()),
				}
			} else {
				return &ast.FalseLiteral{
					ValuePos: a.getTerminalPos(ctx.BoolLiteral()),
				}
			}
		case ctx.StringLiteral() != nil:
			raw := ctx.StringLiteral().GetText()
			value = raw[1 : len(raw)-1]
			valuePos = a.getTerminalPos(ctx.StringLiteral())
			return &ast.StringLiteral{
				Value:    value,
				ValuePos: valuePos,
			}
		case ctx.NULL() != nil:
			return &ast.NullLiteral{
				ValuePos: a.getTerminalPos(ctx.NULL()),
			}
		default:
			return nil
		}
	})
}

// VisitConditionalExpression implements the Visitor interface for ConditionalExpression
func (a *Analyzer) VisitConditionalExpression(ctx *generated.ConditionalExpressionContext) any {
	return a.visitWrapper("ConditionalExpression", ctx, func() any {

		// Get the condition expression
		cond, _ := accept[ast.Expr](ctx.ConditionalBoolExpression(), a)

		// If there's no question mark, just return the condition
		if ctx.QUEST() == nil {
			return cond
		}

		// If there's a question mark, we have a ternary expression
		then, _ := accept[ast.Expr](ctx.ConditionalExpression(0), a)
		els, _ := accept[ast.Expr](ctx.ConditionalExpression(1), a)

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
	})
}

// VisitLogicalExpression implements the Visitor interface for LogicalExpression
func (a *Analyzer) VisitLogicalExpression(ctx *generated.LogicalExpressionContext) any {
	return a.visitWrapper("LogicalExpression", ctx, func() any {

		var ret ast.Expr
		switch len(ctx.AllShiftExpression()) {
		case 1:
			ret, _ = accept[ast.Expr](ctx.ShiftExpression(0), a)
		case 2:
			x, _ := accept[ast.Expr](ctx.ShiftExpression(0), a)
			y, _ := accept[ast.Expr](ctx.ShiftExpression(1), a)

			ret = &ast.BinaryExpr{
				X:     x,
				OpPos: token.Pos(ctx.GetLogicalOp().GetStart()),
				Op:    logicalOpMap[ctx.GetLogicalOp().GetText()],
				Y:     y,
			}
		}

		return ret
	})
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
func (a *Analyzer) VisitShiftExpression(ctx *generated.ShiftExpressionContext) any {
	return a.visitWrapper("ShiftExpression", ctx, func() any {
		var ret ast.Expr

		switch ctx.GetChildCount() {
		case 1:
			ret, _ = accept[ast.Expr](ctx.AddSubExpression(0), a)
		case 2:
			x, _ := accept[ast.Expr](ctx.AddSubExpression(0), a)
			y, _ := accept[ast.Expr](ctx.AddSubExpression(1), a)
			ret = &ast.BinaryExpr{
				X:     x,
				OpPos: token.Pos(ctx.GetShiftOp().GetStart()),
				Op:    token.Token(ctx.GetShiftOp().GetTokenType()),
				Y:     y,
			}
		}

		return ret
	})
}

// VisitAddSubExpression implements the Visitor interface for AddSubExpression
func (a *Analyzer) VisitAddSubExpression(ctx *generated.AddSubExpressionContext) any {
	return a.visitWrapper("AddSubExpression", ctx, func() any {

		var ret ast.Expr
		switch ctx.GetChildCount() {
		case 1:
			ret, _ = accept[ast.Expr](ctx.MulDivExpression(0), a)
		case 2:
			x, _ := accept[ast.Expr](ctx.MulDivExpression(0), a)
			y, _ := accept[ast.Expr](ctx.MulDivExpression(1), a)
			ret = &ast.BinaryExpr{
				X:     x,
				OpPos: token.Pos(ctx.GetAddSubOp().GetStart()),
				Op:    token.Token(ctx.GetAddSubOp().GetTokenType()),
				Y:     y,
			}
		}

		return ret
	})
}

// VisitMulDivExpression implements the Visitor interface for MulDivExpression
func (a *Analyzer) VisitMulDivExpression(ctx *generated.MulDivExpressionContext) any {
	return a.visitWrapper("MulDivExpression", ctx, func() any {

		var ret ast.Expr
		switch ctx.GetChildCount() {
		case 1:
			ret, _ = accept[ast.Expr](ctx.IncDecExpression(0), a)
		case 2:
			x, _ := accept[ast.Expr](ctx.IncDecExpression(0), a)
			y, _ := accept[ast.Expr](ctx.IncDecExpression(1), a)
			ret = &ast.BinaryExpr{
				X:     x,
				OpPos: token.Pos(ctx.GetMulDivOp().GetStart()),
				Op:    token.Token(ctx.GetMulDivOp().GetTokenType()),
				Y:     y,
			}
		}

		return ret
	})
}

// VisitReturnStatement implements the Visitor interface for ReturnStatement
func (a *Analyzer) VisitReturnStatement(ctx *generated.ReturnStatementContext) any {
	return a.visitWrapper("ReturnStatement", ctx, func() any {

		var x ast.Expr
		if ctx.ExpressionList() != nil {
			exprList, _ := accept[[]ast.Expr](ctx.ExpressionList(), a)
			if len(exprList) > 0 {
				x = exprList[0]
			}
		}

		return &ast.ReturnStmt{
			Return: token.Pos(ctx.RETURN().GetSymbol().GetStart()),
			X:      x,
		}
	})
}

// VisitBreakStatement implements the Visitor interface for BreakStatement
func (a *Analyzer) VisitBreakStatement(ctx *generated.BreakStatementContext) any {
	return a.visitWrapper("BreakStatement", ctx, func() any {

		return &ast.BreakStmt{
			Break: token.Pos(ctx.BREAK().GetSymbol().GetStart()),
		}
	})
}

// VisitContinueStatement implements the Visitor interface for ContinueStatement
func (a *Analyzer) VisitContinueStatement(ctx *generated.ContinueStatementContext) any {
	return a.visitWrapper("ContinueStatement", ctx, func() any {

		return &ast.ContinueStmt{
			Continue: token.Pos(ctx.CONTINUE().GetSymbol().GetStart()),
		}
	})
}

// VisitIfStatement implements the Visitor interface for IfStatement
func (a *Analyzer) VisitIfStatement(ctx *generated.IfStatementContext) any {
	return a.visitWrapper("IfStatement", ctx, func() any {

		cond, _ := accept[ast.Expr](ctx.ConditionalExpression(), a)
		body, _ := accept[ast.Stmt](ctx.Block(0), a)

		var else_ ast.Stmt
		if len(ctx.AllBlock()) > 1 {
			else_, _ = accept[ast.Stmt](ctx.Block(1), a)
		}

		return &ast.IfStmt{
			If:   token.Pos(ctx.IF().GetSymbol().GetStart()),
			Cond: cond,
			Body: body.(*ast.BlockStmt),
			Else: else_,
		}
	})
}

// VisitFile implements the Visitor interface for File
func (a *Analyzer) VisitFile(ctx *generated.FileContext) any {
	return a.visitWrapper("File", ctx, func() any {
		file := &ast.File{
			Imports: make(map[string]*ast.Import),
			Stmts:   make([]ast.Stmt, 0),
			Decls:   make([]ast.Decl, 0),
		}

		// Visit all statements
		for _, stmt := range ctx.AllStatement() {
			result := stmt.Accept(a)
			if result != nil {
				// Check if it's a declaration
				if decl, ok := result.(ast.Decl); ok {
					file.Decls = append(file.Decls, decl)
				} else if stmt, ok := result.(ast.Stmt); ok {
					file.Stmts = append(file.Stmts, stmt)
				}
			}
		}
		if len(a.comments) != 0 {
			file.Docs = append(file.Docs, &ast.CommentGroup{
				List: a.comments,
			})
			a.comments = nil
		}

		return file
	})
}

// VisitStatement implements the Visitor interface for Statement
func (a *Analyzer) VisitStatement(ctx *generated.StatementContext) any {
	return a.visitWrapper("Statement", ctx, func() any {

		if ctx.Comment() != nil {
			if ctx.Comment().LineComment() != nil {
				cmt := a.fmtLineComment(ctx.Comment().LineComment().GetText())
				cmt.Slash = token.Pos(ctx.Comment().LineComment().GetSymbol().GetStart())
				a.comments = append(a.comments, cmt)
			} else {
				cmts := a.fmtBlockComment(ctx.Comment().BlockComment().GetText())
				for _, cmt := range cmts {
					cmt.Slash = token.Pos(ctx.Comment().BlockComment().GetSymbol().GetStart())
					a.comments = append(a.comments, cmt)
				}
			}
		}
		// Handle different types of statements
		if ctx.FunctionDeclaration() != nil {
			return a.VisitFunctionDeclaration(ctx.FunctionDeclaration().(*generated.FunctionDeclarationContext))
		}
		if ctx.ExpressionStatement() != nil {
			return a.VisitExpressionStatement(ctx.ExpressionStatement().(*generated.ExpressionStatementContext))
		}
		if ctx.LambdaAssignStatement() != nil {
			return a.VisitLambdaAssignStatement(ctx.LambdaAssignStatement().(*generated.LambdaAssignStatementContext))
		}
		if ctx.AssignStatement() != nil {
			return a.VisitAssignStatement(ctx.AssignStatement().(*generated.AssignStatementContext))
		}
		if ctx.ReturnStatement() != nil {
			return a.VisitReturnStatement(ctx.ReturnStatement().(*generated.ReturnStatementContext))
		}
		if ctx.BreakStatement() != nil {
			return a.VisitBreakStatement(ctx.BreakStatement().(*generated.BreakStatementContext))
		}
		if ctx.ContinueStatement() != nil {
			return a.VisitContinueStatement(ctx.ContinueStatement().(*generated.ContinueStatementContext))
		}
		if ctx.IfStatement() != nil {
			return a.VisitIfStatement(ctx.IfStatement().(*generated.IfStatementContext))
		}
		if ctx.LoopStatement() != nil {
			return a.VisitLoopStatement(ctx.LoopStatement().(*generated.LoopStatementContext))
		}
		// TODO: Add more statement types

		return nil
	})
}

// VisitExpressionStatement implements the Visitor interface for ExpressionStatement
func (a *Analyzer) VisitExpressionStatement(ctx *generated.ExpressionStatementContext) any {
	return a.visitWrapper("ExpressionStatement", ctx, func() any {
		ret := a.VisitExpression(ctx.Expression().(*generated.ExpressionContext))
		if expr, ok := ret.(ast.Expr); ok {
			return &ast.ExprStmt{
				X: expr,
			}
		}
		return ret
	})
}

// VisitLambdaAssignStatement implements the Visitor interface for LambdaAssignStatement
func (a *Analyzer) VisitLambdaAssignStatement(ctx *generated.LambdaAssignStatementContext) any {
	return a.visitWrapper("ReturnStatement", ctx, func() any {

		lhs, _ := accept[ast.Expr](ctx.VariableExpression(), a)
		rhs, _ := accept[ast.Expr](ctx.Expression(), a)

		return &ast.AssignStmt{
			Lhs: lhs,
			Tok: token.COLON_ASSIGN, // Use COLON_ASSIGN token for lambda assignment
			Rhs: rhs,
		}
	})
}

// VisitAssignStatement implements the Visitor interface for AssignStatement
func (a *Analyzer) VisitAssignStatement(ctx *generated.AssignStatementContext) any {
	return a.visitWrapper("AssignStatement", ctx, func() any {

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
			lhs, _ = accept[ast.Expr](ctx.VariableNames(), a)
		} else if ctx.VariableExpression() != nil {
			lhs, _ = accept[ast.Expr](ctx.VariableExpression(), a)
		} else if ctx.VariableNullableExpressions() != nil {
			lhs, _ = accept[ast.Expr](ctx.VariableNullableExpressions(), a)
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
			rhs, _ = accept[ast.Expr](ctx.Expression(), a)
		} else if ctx.MatchStatement() != nil {
			rhs, _ = accept[ast.Expr](ctx.MatchStatement(), a)
		}

		return &ast.AssignStmt{
			Scope:    scope,
			ScopePos: scopePos,
			Lhs:      lhs,
			Tok:      tok,
			Rhs:      rhs,
		}
	})
}

// VisitExpression implements the Visitor interface for Expression
func (a *Analyzer) VisitExpression(ctx *generated.ExpressionContext) any {
	return a.visitWrapper("Expression", ctx, func() any {

		// Handle different types of expressions
		if ctx.LambdaExpression() != nil {
			return a.VisitLambdaExpression(ctx.LambdaExpression().(*generated.LambdaExpressionContext))
		}
		if ctx.ConditionalExpression() != nil {
			node := a.VisitConditionalExpression(ctx.ConditionalExpression().(*generated.ConditionalExpressionContext))
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
			return a.VisitCommandExpression(ctx.CommandExpression().(*generated.CommandExpressionContext))
		}
		if ctx.UnsafeExpression() != nil {
			// return v.VisitUnsafeExpression(ctx.UnsafeExpression().(*generated.UnsafeExpressionContext))
		}
		if ctx.ComptimeExpression() != nil {
			return a.VisitComptimeExpression(ctx.ComptimeExpression().(*generated.ComptimeExpressionContext))
		}

		return nil
	})
}

// VisitComptimeExpression implements the Visitor interface for ComptimeExpression
func (a *Analyzer) VisitComptimeExpression(ctx *generated.ComptimeExpressionContext) any {
	return a.visitWrapper("ComptimeExpression", ctx, func() any {
		block, _ := accept[*ast.BlockStmt](ctx.Block(), a)
		return &ast.ComptimeStmt{
			X: block,
		}
	})
}

// VisitLambdaExpression implements the Visitor interface for LambdaExpression
func (a *Analyzer) VisitLambdaExpression(ctx *generated.LambdaExpressionContext) any {
	return a.visitWrapper("LambdaExpression", ctx, func() any {
		// Get parameters
		params, _ := accept[*ast.FieldList](ctx.ReceiverParameters(), a)

		// Get lambda body
		body, _ := accept[ast.Stmt](ctx.LambdaBody(), a)

		// Create anonymous function name
		funcName := &ast.Ident{
			NamePos: token.Pos(ctx.GetStart().GetStart()),
			Name:    "anonymous", // Generate unique name if needed
		}

		// Convert to FuncDecl
		return &ast.FuncDecl{
			Fn:   token.Pos(ctx.GetStart().GetStart()),
			Name: funcName,
			Recv: params,
			Body: &ast.BlockStmt{
				Lbrace: token.Pos(ctx.GetStart().GetStart()),
				List:   []ast.Stmt{body},
				Rbrace: token.Pos(ctx.GetStop().GetStop()),
			},
		}
	})
}

// VisitLambdaBody implements the Visitor interface for LambdaBody
func (a *Analyzer) VisitLambdaBody(ctx *generated.LambdaBodyContext) any {
	return a.visitWrapper("LambdaBody", ctx, func() any {
		if ctx.Expression() != nil {
			// Single expression - wrap in return statement
			expr, _ := accept[ast.Expr](ctx.Expression(), a)
			return &ast.ReturnStmt{
				Return: token.Pos(ctx.GetStart().GetStart()),
				X:      expr,
			}
		} else if ctx.ExpressionList() != nil {
			// Multiple expressions - wrap in return statement
			exprList, _ := accept[[]ast.Expr](ctx.ExpressionList(), a)
			if len(exprList) > 0 {
				return &ast.ReturnStmt{
					Return: token.Pos(ctx.GetStart().GetStart()),
					X:      exprList[0], // Return first expression
				}
			}
		} else if ctx.Block() != nil {
			// Block - return as is
			block, _ := accept[*ast.BlockStmt](ctx.Block(), a)
			return block
		}

		return nil
	})
}

// VisitCommandExpression implements the Visitor interface for CommandExpression
func (a *Analyzer) VisitCommandExpression(ctx *generated.CommandExpressionContext) any {
	return a.visitWrapper("CommandExpression", ctx, func() any {
		var fun ast.Expr
		var recv []ast.Expr

		// Handle command string literal or member access
		if ctx.CommandStringLiteral() != nil {
			fun = &ast.StringLiteral{
				Value:    ctx.CommandStringLiteral().GetText(),
				ValuePos: token.Pos(ctx.CommandStringLiteral().GetSymbol().GetStart()),
			}
		} else if ctx.MemberAccess() != nil {
			fun, _ = accept[ast.Expr](ctx.MemberAccess(), a)
		}

		ce, isComptime := fun.(*ast.ComptimeExpr)
		if isComptime {
			fun = ce.X
		}

		// Handle options and arguments
		for _, opt := range ctx.AllOption() {
			optExpr, _ := accept[ast.Expr](opt, a)
			if optExpr != nil {
				recv = append(recv, optExpr)
			}
		}

		for _, expr := range ctx.AllConditionalExpression() {
			arg, _ := accept[ast.Expr](expr, a)
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
			join, _ := accept[ast.Expr](ctx.CommandJoin(), a)
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
			stream, _ := accept[ast.Expr](ctx.CommandStream(), a)
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
		return a.wrapComptimeExprssion(call, isComptime)
	})
}

func (a *Analyzer) wrapComptimeExprssion(ret ast.Expr, shouldWith bool) ast.Node {
	if shouldWith {
		return &ast.ComptimeExpr{
			X: ret,
		}
	}
	return ret
}

func (a *Analyzer) VisitMemberAccess(ctx *generated.MemberAccessContext) any {
	return a.visitWrapper("MemberAccess", ctx, func() any {
		var isComptime bool
		if ctx.NOT() != nil {
			isComptime = true
		}

		// Handle identifier with generic arguments
		if ctx.Identifier() != nil && ctx.GenericArguments() != nil {
			ident := ctx.Identifier().GetText()
			genericArgs := a.Visit(ctx.GenericArguments()).(*ast.GenericExpr)
			ret := &ast.CallExpr{
				Fun: &ast.Ident{
					NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
					Name:    ident,
				},
				Generic: genericArgs,
			}
			return a.wrapComptimeExprssion(ret, isComptime)
		}

		// Handle identifier with member access point
		if ctx.Identifier() != nil && ctx.MemberAccessPoint() != nil {
			ident := ctx.Identifier().GetText()
			base := &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ident,
			}
			ret := a.visitMemberAccessPoint(base, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
			return a.wrapComptimeExprssion(ret, isComptime)
		}

		// Handle STR, NUM, BOOL with member access point
		if ctx.STR() != nil || ctx.NUM() != nil || ctx.BOOL() != nil {
			var lhs ast.Expr

			if ctx.STR() != nil {
				lhs = &ast.StringLiteral{
					Value:    ctx.STR().GetText(),
					ValuePos: token.Pos(ctx.STR().GetSymbol().GetStart()),
				}
			} else if ctx.NUM() != nil {
				lhs = &ast.NumericLiteral{
					Value:    ctx.NUM().GetText(),
					ValuePos: token.Pos(ctx.NUM().GetSymbol().GetStart()),
				}
			} else {
				// For BOOL, we need to check if it's true or false
				if ctx.BOOL().GetText() == "true" {
					lhs = &ast.TrueLiteral{
						ValuePos: token.Pos(ctx.BOOL().GetSymbol().GetStart()),
					}
				} else {
					lhs = &ast.FalseLiteral{
						ValuePos: token.Pos(ctx.BOOL().GetSymbol().GetStart()),
					}
				}
			}

			ret := a.visitMemberAccessPoint(lhs, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
			return a.wrapComptimeExprssion(ret, isComptime)
		}

		// Handle literal with member access point
		if ctx.Literal() != nil {
			base := a.Visit(ctx.Literal()).(ast.Expr)
			ret := a.visitMemberAccessPoint(base, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
			return a.wrapComptimeExprssion(ret, isComptime)
		}

		// Handle THIS with optional member access point
		if ctx.THIS() != nil {
			base := &ast.Ident{
				NamePos: token.Pos(ctx.THIS().GetSymbol().GetStart()),
				Name:    "this",
			}
			if ctx.MemberAccessPoint() != nil {
				ret := a.visitMemberAccessPoint(base, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
				return a.wrapComptimeExprssion(ret, isComptime)
			}
			return a.wrapComptimeExprssion(base, isComptime)
		}

		// Handle SUPER with optional member access point
		if ctx.SUPER() != nil {
			base := &ast.Ident{
				NamePos: token.Pos(ctx.SUPER().GetSymbol().GetStart()),
				Name:    "super",
			}
			if ctx.MemberAccessPoint() != nil {
				ret := a.visitMemberAccessPoint(base, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
				return a.wrapComptimeExprssion(ret, isComptime)
			}
			return a.wrapComptimeExprssion(base, isComptime)
		}

		if ctx.Identifier() != nil {
			ret := &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ctx.Identifier().GetText(),
			}
			return a.wrapComptimeExprssion(ret, isComptime)
		}

		return nil
	})
}

func (a *Analyzer) visitMemberAccessPoint(base ast.Expr, ctx *generated.MemberAccessPointContext) ast.Expr {
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
			genericArgs := a.Visit(ctx.GenericArguments()).(*ast.GenericExpr)
			selectExpr.Y = &ast.CallExpr{
				Fun:     selectExpr.Y,
				Generic: genericArgs,
			}
		}

		// Handle recursive member access point
		if ctx.MemberAccessPoint() != nil {
			return a.visitMemberAccessPoint(selectExpr, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
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
			return a.visitMemberAccessPoint(selectExpr, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return selectExpr
	}

	// Handle index access
	if ctx.LBRACK() != nil && ctx.Expression() != nil && ctx.RBRACK() != nil {
		indexExpr := &ast.IndexExpr{
			X:      base,
			Lbrack: token.Pos(ctx.LBRACK().GetSymbol().GetStart()),
			Index:  a.Visit(ctx.Expression()).(ast.Expr),
			Rbrack: token.Pos(ctx.RBRACK().GetSymbol().GetStart()),
		}

		// Handle recursive member access point
		if ctx.MemberAccessPoint() != nil {
			return a.visitMemberAccessPoint(indexExpr, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return indexExpr
	}

	return base
}

// VisitConditionalBoolExpression implements the Visitor interface for ConditionalBoolExpression
func (a *Analyzer) VisitConditionalBoolExpression(ctx *generated.ConditionalBoolExpressionContext) any {
	return a.visitWrapper("ConditionalBoolExpression", ctx, func() any {

		// Get the first logical expression
		x, _ := accept[ast.Expr](ctx.LogicalExpression(0), a)

		if len(ctx.AllLogicalExpression()) > 1 {
			// Handle multiple logical expressions with operators
			for i := 1; i < len(ctx.AllLogicalExpression()); i++ {
				y, _ := accept[ast.Expr](ctx.LogicalExpression(i), a)
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

		return x
	})
}

// VisitIncDecExpression implements the Visitor interface for IncDecExpression
func (a *Analyzer) VisitIncDecExpression(ctx *generated.IncDecExpressionContext) any {
	return a.visitWrapper("IncDecExpression", ctx, func() any {

		var ret ast.Expr

		if ctx.PreIncDecExpression() != nil {
			ret, _ = accept[ast.Expr](ctx.PreIncDecExpression(), a)
		}

		if ctx.PostIncDecExpression() != nil {
			ret, _ = accept[ast.Expr](ctx.PostIncDecExpression(), a)
		}

		return ret
	})
}

// VisitPreIncDecExpression implements the Visitor interface for PreIncDecExpression
func (a *Analyzer) VisitPreIncDecExpression(ctx *generated.PreIncDecExpressionContext) any {
	return a.visitWrapper("PreIncDecExpression", ctx, func() any {

		expr := a.VisitFactor(ctx.Factor().(*generated.FactorContext))
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

		return ret
	})
}

// VisitPostIncDecExpression implements the Visitor interface for PostIncDecExpression
func (a *Analyzer) VisitPostIncDecExpression(ctx *generated.PostIncDecExpressionContext) any {
	return a.visitWrapper("PostIncDecExpression", ctx, func() any {

		expr := a.VisitFactor(ctx.Factor().(*generated.FactorContext))

		if expr == nil {
			return nil
		}

		var tok token.Token
		if ctx.INC() != nil {
			tok = token.INC
		} else if ctx.DEC() != nil {
			tok = token.DEC
		}

		return &ast.IncDecExpr{
			Pre:    false,
			X:      expr.(ast.Expr),
			Tok:    tok,
			TokPos: token.Pos(ctx.GetStart().GetStart()),
		}
	})
}

// VisitFactor implements the Visitor interface for Factor
func (a *Analyzer) VisitFactor(ctx *generated.FactorContext) any {
	return a.visitWrapper("Factor", ctx, func() any {

		// Handle unary expressions
		if ctx.SUB() != nil {
			expr := a.VisitFactor(ctx.Factor().(*generated.FactorContext))
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
			return a.VisitLiteral(ctx.Literal().(*generated.LiteralContext))
		}

		// Handle identifiers
		if ctx.VariableExpression() != nil {
			return a.VisitVariableExpression(ctx.VariableExpression().(*generated.VariableExpressionContext))
		}

		// Handle parenthesized expressions
		if ctx.LPAREN() != nil {
			return a.VisitFactor(ctx.Factor().(*generated.FactorContext))
		}

		if ctx.CallExpression() != nil {
			return a.VisitCallExpression(ctx.CallExpression().(*generated.CallExpressionContext))
		}

		return nil
	})
}

func (a *Analyzer) VisitCallExpression(ctx *generated.CallExpressionContext) any {
	return a.visitWrapper("CallExpression", ctx, func() any {
		ret := a.VisitMemberAccess(ctx.MemberAccess().(*generated.MemberAccessContext))
		stmt, ok := ret.(*ast.ComptimeExpr)
		if !ok {
			return ret
		}
		callExpr := &ast.CallExpr{Fun: stmt.X.(ast.Expr)}
		if ctx.ReceiverArgumentList() != nil && ctx.ReceiverArgumentList().ExpressionList() != nil {
			for _, expr := range ctx.ReceiverArgumentList().ExpressionList().AllExpression() {
				e := a.VisitExpression(expr.(*generated.ExpressionContext))
				callExpr.Recv = append(callExpr.Recv, e.(ast.Expr))
			}
		}
		stmt.X = callExpr
		return ret
	})
}

func (a *Analyzer) VisitVariableExpression(ctx *generated.VariableExpressionContext) any {
	return a.visitWrapper("VariableExpression", ctx, func() any {

		ret := a.VisitMemberAccess(ctx.MemberAccess().(*generated.MemberAccessContext))

		return &ast.RefExpr{
			X: ret.(ast.Expr),
		}
	})
}

func (a *Analyzer) fmtLineComment(cmt string) *ast.Comment {
	// remove the first two characters
	return &ast.Comment{Text: cmt[2:]}
}

var removeStar = regexp.MustCompile(`^\s*\*`)

// fmtBlockComment formats a block comment.
// Because most compiled languages do not support multi-line comments,
// so we need to convert them into single-line comments
func (a *Analyzer) fmtBlockComment(cmt string) (ret []*ast.Comment) {
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
func (a *Analyzer) VisitBlock(ctx *generated.BlockContext) any {
	return a.visitWrapper("Block", ctx, func() any {

		// Create a new block statement
		block := &ast.BlockStmt{
			Lbrace: token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
			Rbrace: token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
		}

		// Process all statements in the block
		for _, stmtCtx := range ctx.AllStatement() {
			stmt, _ := accept[ast.Stmt](stmtCtx, a)
			if stmt != nil {
				block.List = append(block.List, stmt)
			}
		}

		return block
	})
}

// VisitLoopStatement implements the Visitor interface for LoopStatement
func (a *Analyzer) VisitLoopStatement(ctx *generated.LoopStatementContext) any {
	return a.visitWrapper("LoopStatement", ctx, func() any {

		// Handle different types of loops
		if ctx.WhileStatement() != nil {
			return a.VisitWhileStatement(ctx.WhileStatement().(*generated.WhileStatementContext))
		} else if ctx.DoWhileStatement() != nil {
			return a.VisitDoWhileStatement(ctx.DoWhileStatement().(*generated.DoWhileStatementContext))
		} else if ctx.RangeStatement() != nil {
			return a.VisitRangeStatement(ctx.RangeStatement().(*generated.RangeStatementContext))
		} else if ctx.ForStatement() != nil {
			return a.VisitForStatement(ctx.ForStatement().(*generated.ForStatementContext))
		} else if ctx.ForeachStatement() != nil {
			return a.VisitForeachStatement(ctx.ForeachStatement().(*generated.ForeachStatementContext))
		}

		return nil
	})
}

// VisitWhileStatement implements the Visitor interface for WhileStatement
func (a *Analyzer) VisitWhileStatement(ctx *generated.WhileStatementContext) any {
	return a.visitWrapper("WhileStatement", ctx, func() any {
		body, _ := accept[ast.Stmt](ctx.Block(), a)

		var cond ast.Expr
		if ctx.Expression() != nil {
			cond, _ = accept[ast.Expr](ctx.Expression(), a)
		}

		return &ast.WhileStmt{
			Loop: token.Pos(ctx.LOOP().GetSymbol().GetStart()),
			Cond: cond,
			Body: body.(*ast.BlockStmt),
		}
	})
}

// VisitDoWhileStatement implements the Visitor interface for DoWhileStatement
func (a *Analyzer) VisitDoWhileStatement(ctx *generated.DoWhileStatementContext) any {
	return a.visitWrapper("DoWhileStatement", ctx, func() any {

		body, _ := accept[ast.Stmt](ctx.Block(), a)
		cond, _ := accept[ast.Expr](ctx.Expression(), a)

		return &ast.DoWhileStmt{
			Do:     token.Pos(ctx.DO().GetSymbol().GetStart()),
			Body:   body.(*ast.BlockStmt),
			Loop:   token.Pos(ctx.LOOP().GetSymbol().GetStart()),
			Lparen: token.Pos(ctx.LPAREN().GetSymbol().GetStart()),
			Cond:   cond,
			Rparen: token.Pos(ctx.RPAREN().GetSymbol().GetStart()),
		}
	})
}

// VisitRangeStatement implements the Visitor interface for RangeStatement
func (a *Analyzer) VisitRangeStatement(ctx *generated.RangeStatementContext) any {
	return a.visitWrapper("RangeStatement", ctx, func() any {

		index := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}

		rangeExpr, _ := accept[ast.Expr](ctx.RangeClause(), a)
		body, _ := accept[ast.Stmt](ctx.Block(), a)

		return &ast.ForInStmt{
			Loop:   token.Pos(ctx.LOOP().GetSymbol().GetStart()),
			Index:  index,
			In:     token.Pos(ctx.IN().GetSymbol().GetStart()),
			Range:  token.Pos(ctx.RangeClause().RANGE().GetSymbol().GetStart()),
			Lparen: token.Pos(ctx.RangeClause().LPAREN().GetSymbol().GetStart()),
			RangeExpr: ast.RangeExpr{
				Start: rangeExpr.(*ast.BinaryExpr).X,
				End_:  rangeExpr.(*ast.BinaryExpr).Y,
			},
			Rparen: token.Pos(ctx.RangeClause().RPAREN().GetSymbol().GetStart()),
			Body:   body.(*ast.BlockStmt),
		}
	})
}

// VisitForStatement implements the Visitor interface for ForStatement
func (a *Analyzer) VisitForStatement(ctx *generated.ForStatementContext) any {
	return a.visitWrapper("ForStatement", ctx, func() any {

		var (
			init           ast.Stmt
			cond           ast.Expr
			post           ast.Expr
			comma1, comma2 token.Pos
		)

		if ctx.ForClause() != nil {
			if ctx.ForClause().Statement() != nil {
				init, _ = accept[ast.Stmt](ctx.ForClause().Statement(), a)
			}
			if ctx.ForClause().Expression(0) != nil {
				cond, _ = accept[ast.Expr](ctx.ForClause().Expression(0), a)
			}
			if ctx.ForClause().Expression(1) != nil {
				post, _ = accept[ast.Expr](ctx.ForClause().Expression(1), a)
			}
			if len(ctx.ForClause().AllSEMI()) > 0 {
				comma1 = token.Pos(ctx.ForClause().SEMI(0).GetSymbol().GetStart())
			}
			if len(ctx.ForClause().AllSEMI()) > 1 {
				comma2 = token.Pos(ctx.ForClause().SEMI(1).GetSymbol().GetStart())
			}
		}

		body, _ := accept[ast.Stmt](ctx.Block(), a)

		var lparen, rparen token.Pos
		if ctx.LPAREN() != nil {
			lparen = token.Pos(ctx.LPAREN().GetSymbol().GetStart())
		}
		if ctx.RPAREN() != nil {
			rparen = token.Pos(ctx.RPAREN().GetSymbol().GetStart())
		}

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
	})
}

// VisitForeachStatement implements the Visitor interface for ForeachStatement
func (a *Analyzer) VisitForeachStatement(ctx *generated.ForeachStatementContext) any {
	return a.visitWrapper("ForeachStatement", ctx, func() any {

		index, _ := accept[ast.Expr](ctx.ForeachClause().VariableName(0), a)
		var value ast.Expr
		if len(ctx.ForeachClause().AllVariableName()) > 1 {
			value, _ = accept[ast.Expr](ctx.ForeachClause().VariableName(1), a)
		}

		expr, _ := accept[ast.Expr](ctx.Expression(), a)
		body, _ := accept[ast.Stmt](ctx.Block(), a)

		return &ast.ForeachStmt{
			Loop:   token.Pos(ctx.LOOP().GetSymbol().GetStart()),
			Lparen: token.Pos(ctx.LPAREN().GetSymbol().GetStart()),
			Index:  index,
			Value:  value,
			Rparen: token.Pos(ctx.RPAREN().GetSymbol().GetStart()),
			Tok:    token.Pos(ctx.IN().GetSymbol().GetStart()),
			Var:    expr,
			Body:   body.(*ast.BlockStmt),
		}
	})
}

// VisitFunctionDeclaration implements the Visitor interface for FunctionDeclaration
func (a *Analyzer) VisitFunctionDeclaration(ctx *generated.FunctionDeclarationContext) any {
	return a.visitWrapper("FunctionDeclaration", ctx, func() any {
		if ctx.StandardFunctionDeclaration() != nil {
			return a.VisitStandardFunctionDeclaration(ctx.StandardFunctionDeclaration().(*generated.StandardFunctionDeclarationContext))
		} else if ctx.LambdaFunctionDeclaration() != nil {
			return a.VisitLambdaFunctionDeclaration(ctx.LambdaFunctionDeclaration().(*generated.LambdaFunctionDeclarationContext))
		}
		return nil
	})
}

// VisitStandardFunctionDeclaration implements the Visitor interface for StandardFunctionDeclaration
func (a *Analyzer) VisitStandardFunctionDeclaration(ctx *generated.StandardFunctionDeclarationContext) any {
	return a.visitWrapper("StandardFunctionDeclaration", ctx, func() any {
		// Get function name
		var funcName *ast.Ident
		if ctx.Identifier() != nil {
			funcName = &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ctx.Identifier().GetText(),
			}
		} else if ctx.OperatorIdentifier() != nil {
			// Handle operator functions
			opText := ctx.OperatorIdentifier().GetText()
			funcName = &ast.Ident{
				NamePos: token.Pos(ctx.OperatorIdentifier().GetStart().GetStart()),
				Name:    opText,
			}
		}

		// Get parameters
		params, _ := accept[*ast.FieldList](ctx.ReceiverParameters(), a)

		// Get function body
		body, _ := accept[*ast.BlockStmt](ctx.Block(), a)

		var isComptime bool
		// Create function modifier
		var funcMod *ast.FuncModifier
		for _, modCtx := range ctx.AllFunctionModifier() {
			if funcMod == nil {
				funcMod = &ast.FuncModifier{}
			}
			switch {
			case modCtx.PUB() != nil:
				funcMod.Pub = token.Pos(modCtx.PUB().GetSymbol().GetStart())
			case modCtx.COMPTIME() != nil:
				isComptime = true
			}
		}

		decl := &ast.FuncDecl{
			Modifiers: funcMod,
			Fn:        token.Pos(ctx.FN().GetSymbol().GetStart()),
			Name:      funcName,
			Recv:      params,
			Body:      body,
		}

		if isComptime {
			return &ast.ComptimeStmt{
				X: decl,
			}
		}

		// Create FuncDecl
		return decl
	})
}

// VisitLambdaFunctionDeclaration implements the Visitor interface for LambdaFunctionDeclaration
func (a *Analyzer) VisitLambdaFunctionDeclaration(ctx *generated.LambdaFunctionDeclarationContext) any {
	return a.visitWrapper("LambdaFunctionDeclaration", ctx, func() any {
		// Get function name
		var funcName *ast.Ident
		if ctx.Identifier() != nil {
			funcName = &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ctx.Identifier().GetText(),
			}
		} else if ctx.OperatorIdentifier() != nil {
			opText := ctx.OperatorIdentifier().GetText()
			funcName = &ast.Ident{
				NamePos: token.Pos(ctx.OperatorIdentifier().GetStart().GetStart()),
				Name:    opText,
			}
		}

		// Get parameters from lambda expression
		params, _ := accept[*ast.FieldList](ctx.LambdaExpression().(*generated.LambdaExpressionContext).ReceiverParameters(), a)

		// Get lambda body
		body, _ := accept[ast.Stmt](ctx.LambdaExpression().(*generated.LambdaExpressionContext).LambdaBody(), a)

		// Create function modifier
		var mods []ast.Modifier
		for _, modCtx := range ctx.AllFunctionModifier() {
			if mods == nil {
				mods = []ast.Modifier{}
			}
			if modCtx.GetText() == "pub" {
				mods = append(mods, &ast.PubModifier{Pub: token.Pos(modCtx.GetStart().GetStart())})
			} else if modCtx.GetText() == "comptime" {
				// funcMod.Static = token.Pos(modCtx.GetStart().GetStart())
			}
		}

		return &ast.FuncDecl{
			Modifiers: mods,
			Fn:        token.Pos(ctx.FN().GetSymbol().GetStart()),
			Name:      funcName,
			Recv:      params,
			Body: &ast.BlockStmt{
				Lbrace: token.Pos(ctx.GetStart().GetStart()),
				List:   []ast.Stmt{body},
				Rbrace: token.Pos(ctx.GetStop().GetStop()),
			},
		}
	})
}

// VisitReceiverParameters implements the Visitor interface for ReceiverParameters
func (a *Analyzer) VisitReceiverParameters(ctx *generated.ReceiverParametersContext) any {
	return a.visitWrapper("ReceiverParameters", ctx, func() any {
		var fields []*ast.Field

		// Handle receiver parameter list
		if ctx.ReceiverParameterList() != nil {
			for _, paramCtx := range ctx.ReceiverParameterList().AllReceiverParameter() {
				field := &ast.Field{}

				// Get parameter name
				if paramCtx.Identifier() != nil {
					field.Name = &ast.Ident{
						NamePos: token.Pos(paramCtx.Identifier().GetSymbol().GetStart()),
						Name:    paramCtx.Identifier().GetText(),
					}
				}

				// Get parameter type
				if paramCtx.Type_() != nil {
					field.Type, _ = accept[ast.Expr](paramCtx.Type_(), a)
				}

				// Get default value
				if paramCtx.Expression() != nil {
					field.Value, _ = accept[ast.Expr](paramCtx.Expression(), a)
					field.Assign = token.Pos(paramCtx.ASSIGN().GetSymbol().GetStart())
				}

				fields = append(fields, field)
			}
		}

		// Handle named parameters
		if ctx.NamedParameters() != nil {
			// TODO: Implement named parameters handling
		}

		var opening, closing token.Pos
		if ctx.LPAREN() != nil {
			opening = token.Pos(ctx.LPAREN().GetSymbol().GetStart())
		}
		if ctx.RPAREN() != nil {
			closing = token.Pos(ctx.RPAREN().GetSymbol().GetStart())
		}

		return &ast.FieldList{
			Opening: opening,
			List:    fields,
			Closing: closing,
		}
	})
}

// VisitType implements the Visitor interface for Type
func (a *Analyzer) VisitType(ctx *generated.TypeContext) any {
	return a.visitWrapper("Type", ctx, func() any {
		// Handle basic types
		if ctx.STR() != nil {
			return &ast.Ident{
				NamePos: token.Pos(ctx.STR().GetSymbol().GetStart()),
				Name:    "str",
			}
		}
		if ctx.NUM() != nil {
			return &ast.Ident{
				NamePos: token.Pos(ctx.NUM().GetSymbol().GetStart()),
				Name:    "num",
			}
		}
		if ctx.BOOL() != nil {
			return &ast.Ident{
				NamePos: token.Pos(ctx.BOOL().GetSymbol().GetStart()),
				Name:    "bool",
			}
		}
		if ctx.ANY() != nil {
			return &ast.Ident{
				NamePos: token.Pos(ctx.ANY().GetSymbol().GetStart()),
				Name:    "any",
			}
		}
		if ctx.Identifier() != nil {
			return &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ctx.Identifier().GetText(),
			}
		}
		if ctx.StringLiteral() != nil {
			return &ast.StringLiteral{
				Value:    ctx.StringLiteral().GetText(),
				ValuePos: token.Pos(ctx.StringLiteral().GetSymbol().GetStart()),
			}
		}

		// Handle member access types
		if ctx.MemberAccess() != nil {
			return a.VisitMemberAccess(ctx.MemberAccess().(*generated.MemberAccessContext))
		}

		// Handle function types
		if ctx.FunctionType() != nil {
			// TODO: Implement function type handling
			return &ast.Ident{
				NamePos: token.Pos(ctx.GetStart().GetStart()),
				Name:    "function",
			}
		}

		// Handle ellipsis types
		if ctx.EllipsisType() != nil {
			// TODO: Implement ellipsis type handling
			return &ast.Ident{
				NamePos: token.Pos(ctx.GetStart().GetStart()),
				Name:    "...",
			}
		}

		return nil
	})
}
