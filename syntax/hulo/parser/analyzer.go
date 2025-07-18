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

		return &ast.ConditionalExpr{
			Cond:      cond,
			Quest:     token.Pos(ctx.QUEST().GetSymbol().GetStart()),
			WhenTrue:  then,
			Colon:     token.Pos(ctx.COLON().GetSymbol().GetStart()),
			WhneFalse: els,
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

		// 获取所有的 MulDivExpression
		mulDivExprs := ctx.AllMulDivExpression()
		if len(mulDivExprs) == 0 {
			return nil
		}

		// 如果只有一个表达式，直接返回
		if len(mulDivExprs) == 1 {
			ret, _ := accept[ast.Expr](mulDivExprs[0], a)
			return ret
		}

		// 处理多个表达式，构建二元表达式树
		var result ast.Expr
		for i, mulDivExpr := range mulDivExprs {
			expr, _ := accept[ast.Expr](mulDivExpr, a)
			if i == 0 {
				result = expr
			} else {
				// 获取对应的操作符
				opIndex := i - 1
				if opIndex < len(ctx.AllADD())+len(ctx.AllSUB()) {
					// 检查是 ADD 还是 SUB
					var op token.Token
					if opIndex < len(ctx.AllADD()) {
						op = token.PLUS
					} else {
						op = token.MINUS
					}
					result = &ast.BinaryExpr{
						X:     result,
						OpPos: token.Pos(ctx.GetAddSubOp().GetStart()),
						Op:    op,
						Y:     expr,
					}
				}
			}
		}

		return result
	})
}

// VisitMulDivExpression implements the Visitor interface for MulDivExpression
func (a *Analyzer) VisitMulDivExpression(ctx *generated.MulDivExpressionContext) any {
	return a.visitWrapper("MulDivExpression", ctx, func() any {

		// 获取所有的 IncDecExpression
		incDecExprs := ctx.AllIncDecExpression()

		if len(incDecExprs) == 1 {
			// 只有一个表达式，直接返回
			ret, _ := accept[ast.Expr](incDecExprs[0], a)
			return ret
		} else if len(incDecExprs) == 2 {
			// 有两个表达式，说明有运算符
			x, _ := accept[ast.Expr](incDecExprs[0], a)
			y, _ := accept[ast.Expr](incDecExprs[1], a)

			// 获取运算符
			op := ctx.GetMulDivOp()
			if op == nil {
				// 如果没有运算符，返回第一个表达式
				return x
			}

			ret := &ast.BinaryExpr{
				X:     x,
				OpPos: token.Pos(op.GetStart()),
				Op:    a.convertTokenType(op.GetTokenType()),
				Y:     y,
			}
			return ret
		}

		// 默认情况，返回第一个表达式
		if len(incDecExprs) > 0 {
			ret, _ := accept[ast.Expr](incDecExprs[0], a)
			return ret
		}

		return nil
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
		if ctx.ELSE() != nil {
			// 检查 else 后面是什么
			if ctx.IfStatement() != nil {
				// else if 的情况
				else_, _ = accept[ast.Stmt](ctx.IfStatement(), a)
			} else if len(ctx.AllBlock()) > 1 {
				// else 的情况
				else_, _ = accept[ast.Stmt](ctx.Block(1), a)
			}
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
			Decls:   make([]ast.Stmt, 0),
		}

		// Visit all statements
		for _, stmt := range ctx.AllStatement() {
			result := stmt.Accept(a)
			if result != nil {
				// Check if it's a declaration
				if stmt, ok := result.(ast.Stmt); ok {
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
		if ctx.ClassDeclaration() != nil {
			return a.VisitClassDeclaration(ctx.ClassDeclaration().(*generated.ClassDeclarationContext))
		}
		if ctx.ImportDeclaration() != nil {
			return a.VisitImportDeclaration(ctx.ImportDeclaration().(*generated.ImportDeclarationContext))
		}
		if ctx.DeclareStatement() != nil {
			return a.VisitDeclareStatement(ctx.DeclareStatement().(*generated.DeclareStatementContext))
		}
		if ctx.TypeDeclaration() != nil {
			return a.VisitTypeDeclaration(ctx.TypeDeclaration().(*generated.TypeDeclarationContext))
		}
		if ctx.EnumDeclaration() != nil {
			return a.VisitEnumDeclaration(ctx.EnumDeclaration().(*generated.EnumDeclarationContext))
		}
		if ctx.MatchStatement() != nil {
			return a.VisitMatchStatement(ctx.MatchStatement().(*generated.MatchStatementContext))
		}
		if ctx.ExternDeclaration() != nil {
			return a.VisitExternDeclaration(ctx.ExternDeclaration().(*generated.ExternDeclarationContext))
		}
		// TODO: Add more statement types

		return nil
	})
}

func (a *Analyzer) collectComments() *ast.CommentGroup {
	ret := a.comments
	a.comments = nil
	return &ast.CommentGroup{
		List: ret,
	}
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

		var type_ ast.Expr
		if ctx.Type_() != nil {
			type_, _ = accept[ast.Expr](ctx.Type_(), a)
		}

		// Get right hand side expression
		var rhs ast.Expr
		if ctx.Expression() != nil {
			rhs, _ = accept[ast.Expr](ctx.Expression(), a)
		} else if ctx.MatchStatement() != nil {
			rhs, _ = accept[ast.Expr](ctx.MatchStatement(), a)
		}

		return &ast.AssignStmt{
			Docs:     a.collectComments(),
			Scope:    scope,
			ScopePos: scopePos,
			Lhs:      lhs,
			Type:     type_,
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
			return a.VisitUnsafeExpression(ctx.UnsafeExpression().(*generated.UnsafeExpressionContext))
		}
		if ctx.ComptimeExpression() != nil {
			return a.VisitComptimeExpression(ctx.ComptimeExpression().(*generated.ComptimeExpressionContext))
		}

		return nil
	})
}

// VisitExpressionList implements the Visitor interface for ExpressionList
func (a *Analyzer) VisitExpressionList(ctx *generated.ExpressionListContext) any {
	return a.visitWrapper("ExpressionList", ctx, func() any {
		var expressions []ast.Expr

		// Visit all expressions in the list
		for _, exprCtx := range ctx.AllExpression() {
			expr, _ := accept[ast.Expr](exprCtx, a)
			if expr != nil {
				expressions = append(expressions, expr)
			}
		}

		return expressions
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
		params, _ := accept[[]ast.Expr](ctx.ReceiverParameters(), a)

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

func (a *Analyzer) VisitOption(ctx *generated.OptionContext) any {
	return a.visitWrapper("Option", ctx, func() any {
		// Handle short option like -e
		if ctx.ShortOption() != nil {
			shortOpt := ctx.ShortOption()
			// Combine SUB (-) and Identifier (e) to form "-e"
			optionText := shortOpt.SUB().GetText() + shortOpt.Identifier().GetText()
			return &ast.Ident{
				NamePos: token.Pos(shortOpt.SUB().GetSymbol().GetStart()),
				Name:    optionText,
			}
		}

		// Handle long option like --verbose
		if ctx.LongOption() != nil {
			longOpt := ctx.LongOption()
			// Combine DEC (--) and Identifier (verbose) to form "--verbose"
			optionText := longOpt.DEC().GetText() + longOpt.Identifier().GetText()
			return &ast.Ident{
				NamePos: token.Pos(longOpt.DEC().GetSymbol().GetStart()),
				Name:    optionText,
			}
		}

		return nil
	})
}

func (a *Analyzer) VisitCommandArgument(ctx *generated.CommandArgumentContext) any {
	return a.visitWrapper("CommandArgument", ctx, func() any {
		// Check each possible child type and visit it if present
		if ctx.Option() != nil {
			// TODO: option 就提取字面量，貌似GetText就可以
			return &ast.Ident{
				NamePos: token.Pos(ctx.Option().GetStart().GetStart()),
				Name:    ctx.Option().GetText(),
			}
		}
		if ctx.ConditionalExpression() != nil {
			expr, _ := accept[ast.Expr](ctx.ConditionalExpression(), a)
			return expr
		}
		if ctx.MemberAccess() != nil {
			expr, _ := accept[ast.Expr](ctx.MemberAccess(), a)
			return expr
		}
		// If none of the expected children are present, return nil
		return nil
	})
}

func (a *Analyzer) VisitCommandExpression(ctx *generated.CommandExpressionContext) any {
	return a.visitWrapper("CommandExpression", ctx, func() any {
		var fun ast.Expr
		var recv []ast.Expr

		// Handle command string literal or identifier
		if ctx.CommandStringLiteral() != nil {
			fun = &ast.StringLiteral{
				Value:    ctx.CommandStringLiteral().GetText(),
				ValuePos: token.Pos(ctx.CommandStringLiteral().GetSymbol().GetStart()),
			}
		} else if ctx.MemberAccess() != nil {
			// First memberAccess is the command name
			if memberAccess, ok := accept[ast.Expr](ctx.MemberAccess(), a); ok {
				fun = memberAccess
			}
		}

		ce, isComptime := fun.(*ast.ComptimeExpr)
		if isComptime {
			fun = ce.X
		}

		// Collect arguments in order from commandArgument list
		for _, argCtx := range ctx.AllCommandArgument() {
			if arg, ok := accept[ast.Expr](argCtx, a); ok {
				recv = append(recv, arg)
			}
		}

		isAsync := false
		if ctx.BITAND() != nil {
			isAsync = true
		}

		builtinArgs := []ast.Expr{}
		if ctx.BuiltinCommandArgument() != nil {
			for _, argCtx := range ctx.BuiltinCommandArgument().AllCommandArgument() {
				arg, _ := accept[ast.Expr](argCtx, a)
				builtinArgs = append(builtinArgs, arg)
			}
		}

		var call ast.Expr
		// Create CallExpr
		call = &ast.CmdExpr{
			Cmd:         fun,
			Args:        recv,
			IsAsync:     isAsync,
			BuiltinArgs: builtinArgs,
		}

		// Handle command join or stream if present
		if ctx.CommandJoin() != nil {
			join, _ := accept[ast.Expr](ctx.CommandJoin(), a)
			if join != nil {
				tok := token.PIPE
				if ctx.CommandJoin().AND() != nil {
					tok = token.CONCAT
				}
				// Create a new CallExpr for the joined command
				call = &ast.BinaryExpr{
					X:  call,
					Op: tok,
					Y:  join,
				}
			}
		} else if ctx.CommandStream() != nil {
			stream, _ := accept[ast.Expr](ctx.CommandStream(), a)
			if stream != nil {
				var tok token.Token
				if ctx.CommandStream().LT() != nil {
					tok = token.LT
				} else if ctx.CommandStream().GT() != nil {
					tok = token.GT
				} else if ctx.CommandStream().SHL() != nil {
					tok = token.LT
				} else if ctx.CommandStream().SHR() != nil {
					tok = token.GT
				} else if ctx.CommandStream().BITOR() != nil {
					tok = token.PIPE
				}
				// Create a new CallExpr for the streamed command
				call = &ast.BinaryExpr{
					X:  call,
					Op: tok,
					Y:  stream,
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
			genericArgs := a.Visit(ctx.GenericArguments()).([]ast.Expr)
			ret := &ast.CallExpr{
				Fun: &ast.Ident{
					NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
					Name:    ident,
				},
				TypeParams: genericArgs,
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
			memberAccessPoint := ctx.MemberAccessPoint()
			if memberAccessPoint != nil {
				ret := a.visitMemberAccessPoint(base, memberAccessPoint.(*generated.MemberAccessPointContext))
				return a.wrapComptimeExprssion(ret, isComptime)
			}
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
			genericArgs := a.Visit(ctx.GenericArguments()).([]ast.Expr)
			selectExpr.Y = &ast.CallExpr{
				Fun:        selectExpr.Y,
				TypeParams: genericArgs,
			}
		}

		// Handle recursive member access point
		if ctx.MemberAccessPoint() != nil {
			return a.visitMemberAccessPoint(selectExpr, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return selectExpr
	}

	// Handle double colon access with identifier
	if ctx.DOUBLE_COLON() != nil && ctx.Identifier() != nil {
		ident := ctx.Identifier().GetText()
		modExpr := &ast.ModAccessExpr{
			X:        base,
			DblColon: token.Pos(ctx.DOUBLE_COLON().GetSymbol().GetStart()),
			Y: &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ident,
			},
		}

		// Handle recursive member access point
		if ctx.MemberAccessPoint() != nil {
			return a.visitMemberAccessPoint(modExpr, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return modExpr
	}

	// Handle double colon access with variable expression
	if ctx.DOUBLE_COLON() != nil && ctx.VariableExpression() != nil {
		varExpr := a.Visit(ctx.VariableExpression()).(ast.Expr)
		modExpr := &ast.ModAccessExpr{
			X:        base,
			DblColon: token.Pos(ctx.DOUBLE_COLON().GetSymbol().GetStart()),
			Y:        varExpr,
		}

		// Handle recursive member access point
		if ctx.MemberAccessPoint() != nil {
			return a.visitMemberAccessPoint(modExpr, ctx.MemberAccessPoint().(*generated.MemberAccessPointContext))
		}
		return modExpr
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
		if x == nil {
			return nil
		}

		// If there's only one logical expression, return it
		if len(ctx.AllLogicalExpression()) == 1 {
			return x
		}

		// Handle multiple logical expressions with operators
		for i := 1; i < len(ctx.AllLogicalExpression()); i++ {
			y, _ := accept[ast.Expr](ctx.LogicalExpression(i), a)
			if y == nil {
				continue
			}

			// Get the operator for this pair
			var op token.Token
			var opPos token.Pos

			// Find the corresponding operator
			opIndex := i - 1
			if opIndex < len(ctx.AllBITAND()) {
				op = token.CONCAT
				opPos = token.Pos(ctx.BITAND(opIndex).GetSymbol().GetStart())
			} else if opIndex < len(ctx.AllBITOR()) {
				op = token.OR
				opPos = token.Pos(ctx.BITOR(opIndex).GetSymbol().GetStart())
			} else if opIndex < len(ctx.AllAND()) {
				op = token.AND
				opPos = token.Pos(ctx.AND(opIndex).GetSymbol().GetStart())
			} else if opIndex < len(ctx.AllOR()) {
				op = token.OR
				opPos = token.Pos(ctx.OR(opIndex).GetSymbol().GetStart())
			}

			x = &ast.BinaryExpr{
				X:     x,
				OpPos: opPos,
				Op:    op,
				Y:     y,
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

		// 如果没有 INC 或 DEC 操作符，直接返回表达式
		if ctx.INC() == nil && ctx.DEC() == nil {
			return expr.(ast.Expr)
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

		// Handle list expressions (array literals)
		if ctx.ListExpression() != nil {
			return a.VisitListExpression(ctx.ListExpression().(*generated.ListExpressionContext))
		}

		// Handle map expressions (object literals)
		if ctx.MapExpression() != nil {
			return a.VisitMapExpression(ctx.MapExpression().(*generated.MapExpressionContext))
		}

		// Handle identifiers
		if ctx.VariableExpression() != nil {
			return a.VisitVariableExpression(ctx.VariableExpression().(*generated.VariableExpressionContext))
		}

		// Handle method expressions
		if ctx.MethodExpression() != nil {
			return a.VisitMethodExpression(ctx.MethodExpression().(*generated.MethodExpressionContext))
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

		// Handle comptime expressions
		if stmt, ok := ret.(*ast.ComptimeExpr); ok {
			callExpr := &ast.CallExpr{Fun: stmt.X.(ast.Expr)}
			if ctx.ReceiverArgumentList() != nil && ctx.ReceiverArgumentList().ExpressionList() != nil {
				for _, expr := range ctx.ReceiverArgumentList().ExpressionList().AllExpression() {
					e := a.VisitExpression(expr.(*generated.ExpressionContext))
					callExpr.Recv = append(callExpr.Recv, e.(ast.Expr))
				}
			}
			stmt.X = callExpr

			// Handle chained calls
			if ctx.CallExpressionLinkedList() != nil {
				stmt.X = a.processCallExpressionLinkedList(ctx.CallExpressionLinkedList().(*generated.CallExpressionLinkedListContext), stmt.X.(ast.Expr))
			}

			return ret
		}

		// Handle regular function calls
		if fun, ok := ret.(ast.Expr); ok {
			callExpr := &ast.CallExpr{Fun: fun}
			if ctx.ReceiverArgumentList() != nil && ctx.ReceiverArgumentList().ExpressionList() != nil {
				for _, expr := range ctx.ReceiverArgumentList().ExpressionList().AllExpression() {
					e := a.VisitExpression(expr.(*generated.ExpressionContext))
					callExpr.Recv = append(callExpr.Recv, e.(ast.Expr))
				}
			}

			// Handle chained calls
			if ctx.CallExpressionLinkedList() != nil {
				return a.processCallExpressionLinkedList(ctx.CallExpressionLinkedList().(*generated.CallExpressionLinkedListContext), callExpr)
			}

			return callExpr
		}

		return ret
	})
}

func (a *Analyzer) VisitCallExpressionLinkedList(ctx *generated.CallExpressionLinkedListContext) any {
	return a.visitWrapper("CallExpressionLinkedList", ctx, func() any {
		// This method should not be called directly by the visitor pattern
		// It should only be called from VisitCallExpression with a base expression
		return nil
	})
}

func (a *Analyzer) processCallExpressionLinkedList(ctx *generated.CallExpressionLinkedListContext, base ast.Expr) ast.Expr {
	// Handle dot access (e.g., .to_str())
	if ctx.DOT() != nil {
		// Create a SelectExpr for the method access
		selectExpr := &ast.SelectExpr{
			X:   base,
			Dot: token.Pos(ctx.DOT().GetSymbol().GetStart()),
			Y: &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ctx.Identifier().GetText(),
			},
		}

		// Create a CallExpr for the method call
		callExpr := &ast.CallExpr{
			Fun:    selectExpr,
			Lparen: token.Pos(ctx.LPAREN().GetSymbol().GetStart()),
			Rparen: token.Pos(ctx.RPAREN().GetSymbol().GetStart()),
		}

		// Handle arguments if present
		if ctx.ReceiverArgumentList() != nil && ctx.ReceiverArgumentList().ExpressionList() != nil {
			for _, expr := range ctx.ReceiverArgumentList().ExpressionList().AllExpression() {
				e := a.VisitExpression(expr.(*generated.ExpressionContext))
				if e != nil {
					callExpr.Recv = append(callExpr.Recv, e.(ast.Expr))
				}
			}
		}

		// Handle recursive chained calls
		if ctx.CallExpressionLinkedList() != nil {
			return a.processCallExpressionLinkedList(ctx.CallExpressionLinkedList().(*generated.CallExpressionLinkedListContext), callExpr)
		}

		return callExpr
	}

	// Handle double dot access (e.g., ..to_str()) - cascade operator
	if ctx.DOUBLE_DOT() != nil {
		// Create a CascadeExpr for the cascade operation
		cascadeExpr := &ast.CascadeExpr{
			X:      base,
			DblDot: token.Pos(ctx.DOUBLE_DOT().GetSymbol().GetStart()),
			Y: &ast.CallExpr{
				Fun: &ast.Ident{
					NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
					Name:    ctx.Identifier().GetText(),
				},
				Lparen: token.Pos(ctx.LPAREN().GetSymbol().GetStart()),
				Rparen: token.Pos(ctx.RPAREN().GetSymbol().GetStart()),
			},
		}

		// Handle arguments if present
		if ctx.ReceiverArgumentList() != nil && ctx.ReceiverArgumentList().ExpressionList() != nil {
			for _, expr := range ctx.ReceiverArgumentList().ExpressionList().AllExpression() {
				e := a.VisitExpression(expr.(*generated.ExpressionContext))
				if e != nil {
					cascadeExpr.Y.(*ast.CallExpr).Recv = append(cascadeExpr.Y.(*ast.CallExpr).Recv, e.(ast.Expr))
				}
			}
		}

		// Handle recursive chained calls
		if ctx.CallExpressionLinkedList() != nil {
			return a.processCallExpressionLinkedList(ctx.CallExpressionLinkedList().(*generated.CallExpressionLinkedListContext), cascadeExpr)
		}

		return cascadeExpr
	}

	return base
}

func (a *Analyzer) VisitVariableExpression(ctx *generated.VariableExpressionContext) any {
	return a.visitWrapper("VariableExpression", ctx, func() any {

		ret := a.VisitMemberAccess(ctx.MemberAccess().(*generated.MemberAccessContext))
		if ret == nil {
			return nil
		}

		return &ast.RefExpr{
			X: ret.(ast.Expr),
		}
	})
}

func (a *Analyzer) VisitMethodExpression(ctx *generated.MethodExpressionContext) any {
	return a.visitWrapper("MethodExpression", ctx, func() any {
		// Get the member access (e.g., p.greet)
		memberAccess := a.VisitMemberAccess(ctx.MemberAccess().(*generated.MemberAccessContext))
		if memberAccess == nil {
			return nil
		}

		// Wrap the member access in a RefExpr to represent the $ prefix
		refExpr := &ast.RefExpr{
			X: memberAccess.(ast.Expr),
		}

		// Create a call expression
		callExpr := &ast.CallExpr{
			Fun:    refExpr,
			Lparen: token.Pos(ctx.LPAREN().GetSymbol().GetStart()),
			Rparen: token.Pos(ctx.RPAREN().GetSymbol().GetStart()),
		}

		// Handle arguments if present - use the same approach as VisitCallExpression
		if ctx.ReceiverArgumentList() != nil && ctx.ReceiverArgumentList().ExpressionList() != nil {
			for _, expr := range ctx.ReceiverArgumentList().ExpressionList().AllExpression() {
				e := a.VisitExpression(expr.(*generated.ExpressionContext))
				if e != nil {
					callExpr.Recv = append(callExpr.Recv, e.(ast.Expr))
				}
			}
		}

		return callExpr
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
			Loop:  token.Pos(ctx.LOOP().GetSymbol().GetStart()),
			Index: index,
			In:    token.Pos(ctx.IN().GetSymbol().GetStart()),
			RangeExpr: ast.RangeExpr{
				Start: rangeExpr.(*ast.BinaryExpr).X,
				End_:  rangeExpr.(*ast.BinaryExpr).Y,
			},
			Body: body.(*ast.BlockStmt),
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

		// Get the foreach clause
		foreachClause := ctx.ForeachClause()

		var index ast.Expr
		var value ast.Expr
		var lparen, rparen token.Pos

		// Check if we have parentheses (multiple variables)
		if foreachClause.LPAREN() != nil {
			// Multiple variables: ($item, $index)
			lparen = token.Pos(foreachClause.LPAREN().GetSymbol().GetStart())
			rparen = token.Pos(foreachClause.RPAREN().GetSymbol().GetStart())

			// Get the first variable (index)
			if len(foreachClause.AllForeachVariableName()) > 0 {
				index, _ = accept[ast.Expr](foreachClause.ForeachVariableName(0), a)
			}

			// Get the second variable (value) if it exists
			if len(foreachClause.AllForeachVariableName()) > 1 {
				value, _ = accept[ast.Expr](foreachClause.ForeachVariableName(1), a)
			}
		} else {
			// Single variable: $item
			if len(foreachClause.AllForeachVariableName()) > 0 {
				index, _ = accept[ast.Expr](foreachClause.ForeachVariableName(0), a)
			}
		}

		expr, _ := accept[ast.Expr](ctx.Expression(), a)
		body, _ := accept[ast.Stmt](ctx.Block(), a)

		var tok token.Token
		if ctx.IN() != nil {
			tok = token.IN
		} else if ctx.OF() != nil {
			tok = token.OF
		}

		return &ast.ForeachStmt{
			Loop:   token.Pos(ctx.LOOP().GetSymbol().GetStart()),
			Lparen: lparen,
			Index:  index,
			Value:  value,
			Rparen: rparen,
			Tok:    tok,
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
		} else if ctx.FunctionSignature() != nil {
			return a.VisitFunctionSignature(ctx.FunctionSignature().(*generated.FunctionSignatureContext))
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
		params, _ := accept[[]ast.Expr](ctx.ReceiverParameters(), a)

		// Get function body
		body, _ := accept[*ast.BlockStmt](ctx.Block(), a)

		var isComptime bool
		// Create function modifier
		var funcMod []ast.Modifier
		for _, modCtx := range ctx.AllFunctionModifier() {
			if funcMod == nil {
				funcMod = []ast.Modifier{}
			}
			switch {
			case modCtx.PUB() != nil:
				funcMod = append(funcMod, &ast.PubModifier{Pub: token.Pos(modCtx.PUB().GetSymbol().GetStart())})
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
		params, _ := accept[[]ast.Expr](ctx.LambdaExpression().(*generated.LambdaExpressionContext).ReceiverParameters(), a)

		// Get lambda body
		body, _ := accept[ast.Stmt](ctx.LambdaExpression().(*generated.LambdaExpressionContext).LambdaBody(), a)

		// Create function modifier
		var mods []ast.Modifier
		for _, modCtx := range ctx.AllFunctionModifier() {
			if mods == nil {
				mods = []ast.Modifier{}
			}
			if modCtx.PUB() != nil {
				mods = append(mods, &ast.PubModifier{Pub: token.Pos(modCtx.PUB().GetSymbol().GetStart())})
			} else if modCtx.COMPTIME() != nil {
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

// VisitFunctionSignature implements the Visitor interface for FunctionSignature
func (a *Analyzer) VisitFunctionSignature(ctx *generated.FunctionSignatureContext) any {
	return a.visitWrapper("FunctionSignature", ctx, func() any {
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

		// Get parameters
		params, _ := accept[[]ast.Expr](ctx.ReceiverParameters(), a)

		// Get return type
		var returnType ast.Expr
		if ctx.FunctionReturnValue() != nil {
			returnType, _ = accept[ast.Expr](ctx.FunctionReturnValue(), a)
		}

		// Create function modifier
		var funcMod []ast.Modifier
		for _, modCtx := range ctx.AllFunctionModifier() {
			if funcMod == nil {
				funcMod = []ast.Modifier{}
			}
			switch {
			case modCtx.PUB() != nil:
				funcMod = append(funcMod, &ast.PubModifier{Pub: token.Pos(modCtx.PUB().GetSymbol().GetStart())})
			case modCtx.COMPTIME() != nil:
				// Handle comptime modifier if needed
			}
		}

		// Create FuncDecl with empty body for signature
		return &ast.FuncDecl{
			Modifiers: funcMod,
			Fn:        token.Pos(ctx.FN().GetSymbol().GetStart()),
			Name:      funcName,
			Recv:      params,
			Type:      returnType,
			Body:      nil, // Function signature has no body
		}
	})
}

// VisitReceiverParameters implements the Visitor interface for ReceiverParameters
func (a *Analyzer) VisitReceiverParameters(ctx *generated.ReceiverParametersContext) any {
	return a.visitWrapper("ReceiverParameters", ctx, func() any {
		var fields []ast.Expr

		// Handle receiver parameter list
		if ctx.ReceiverParameterList() != nil {
			for _, paramCtx := range ctx.ReceiverParameterList().AllReceiverParameter() {
				field := &ast.Parameter{}

				if paramCtx.ELLIPSIS() != nil {
					field.Modifier = &ast.EllipsisModifier{
						Ellipsis: token.Pos(paramCtx.ELLIPSIS().GetSymbol().GetStart()),
					}
				}

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

					// Check if parameter is optional (has QUEST) - make type nullable
					if paramCtx.QUEST() != nil {
						field.Type = &ast.NullableType{
							X: field.Type,
						}
					}
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

		return fields
	})
}

// VisitType implements the Visitor interface for Type
func (a *Analyzer) VisitType(ctx *generated.TypeContext) any {
	return a.visitWrapper("Type", ctx, func() any {
		var baseType ast.Expr

		// Handle basic types
		if ctx.STR() != nil {
			baseType = &ast.TypeReference{
				Name: &ast.Ident{
					NamePos: token.Pos(ctx.STR().GetSymbol().GetStart()),
					Name:    "str",
				},
			}
		} else if ctx.NUM() != nil {
			baseType = &ast.TypeReference{
				Name: &ast.Ident{
					NamePos: token.Pos(ctx.NUM().GetSymbol().GetStart()),
					Name:    "num",
				},
			}
		} else if ctx.BOOL() != nil {
			baseType = &ast.TypeReference{
				Name: &ast.Ident{
					NamePos: token.Pos(ctx.BOOL().GetSymbol().GetStart()),
					Name:    "bool",
				},
			}
		} else if ctx.ANY() != nil {
			baseType = &ast.TypeReference{
				Name: &ast.Ident{
					NamePos: token.Pos(ctx.ANY().GetSymbol().GetStart()),
					Name:    "any",
				},
			}
		} else if ctx.Identifier() != nil {
			baseType = &ast.TypeReference{
				Name: &ast.Ident{
					NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
					Name:    ctx.Identifier().GetText(),
				},
			}
		} else if ctx.StringLiteral() != nil {
			// Remove the quotes from the string literal
			raw := ctx.StringLiteral().GetText()
			value := raw[1 : len(raw)-1] // Remove first and last character (quotes)
			baseType = &ast.StringLiteral{
				Value:    value,
				ValuePos: token.Pos(ctx.StringLiteral().GetSymbol().GetStart()),
			}
		} else if ctx.MemberAccess() != nil {
			// Handle member access types
			memberAccess := a.VisitMemberAccess(ctx.MemberAccess().(*generated.MemberAccessContext))
			baseType = memberAccess.(ast.Expr)
		} else if ctx.FunctionType() != nil {
			// Handle function types
			funcType, _ := accept[ast.Expr](ctx.FunctionType(), a)
			baseType = funcType

		} else if ctx.ObjectType() != nil {
			// Handle object types
			objectType, _ := accept[ast.Expr](ctx.ObjectType(), a)
			baseType = objectType
		} else if ctx.TupleType() != nil {
			// Handle tuple types
			tupleType, _ := accept[ast.Expr](ctx.TupleType(), a)
			baseType = tupleType
		}

		// Handle type access points (array types)
		if baseType != nil && ctx.TypeAccessPoint() != nil {
			baseType = a.visitTypeAccessPoint(baseType, ctx.TypeAccessPoint().(*generated.TypeAccessPointContext))
		}

		// Handle nullable types (T?)
		if baseType != nil && ctx.QUEST() != nil {
			baseType = &ast.NullableType{
				X: baseType,
			}
		}

		// Handle composite types (union and intersection)
		if baseType != nil && ctx.CompositeType() != nil {
			compositeType, _ := accept[ast.Expr](ctx.CompositeType(), a)
			if compositeType != nil {
				// Merge the base type with the composite type
				if unionType, ok := compositeType.(*ast.UnionType); ok {
					unionType.Types = append([]ast.Expr{baseType}, unionType.Types...)
					baseType = unionType
				} else if intersectionType, ok := compositeType.(*ast.IntersectionType); ok {
					intersectionType.Types = append([]ast.Expr{baseType}, intersectionType.Types...)
					baseType = intersectionType
				} else {
					baseType = compositeType
				}
			}
		}

		return baseType
	})
}

// VisitClassDeclaration implements the Visitor interface for ClassDeclaration
func (a *Analyzer) VisitClassDeclaration(ctx *generated.ClassDeclarationContext) any {
	return a.visitWrapper("ClassDeclaration", ctx, func() any {
		// Get class name
		className := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}

		// Create class modifiers
		var pubPos token.Pos
		var modifiers []ast.Modifier
		for _, modCtx := range ctx.AllClassModifier() {
			if modCtx.PUB() != nil {
				pubPos = token.Pos(modCtx.PUB().GetSymbol().GetStart())
				modifiers = append(modifiers, &ast.PubModifier{
					Pub: pubPos,
				})
			}
		}

		// Extract fields and methods directly from class body
		var fields []*ast.Field
		var methods []*ast.FuncDecl
		var constructors []*ast.ConstructorDecl

		// Process all class members
		for _, member := range ctx.ClassBody().AllClassMember() {
			if member != nil {
				field, _ := accept[*ast.Field](member, a)
				if field != nil {
					fields = append(fields, field)
				}
			}
		}

		// Process all class methods
		for _, method := range ctx.ClassBody().AllClassMethod() {
			if method != nil {
				funcDecl, _ := accept[*ast.FuncDecl](method, a)
				if funcDecl != nil {
					methods = append(methods, funcDecl)
				}
			}
		}

		// Process all class builtin methods
		for _, builtinMethod := range ctx.ClassBody().AllClassBuiltinMethod() {
			if builtinMethod != nil {
				constructor, _ := accept[*ast.ConstructorDecl](builtinMethod, a)
				if constructor != nil {
					constructors = append(constructors, constructor)
				}
			}
		}

		return &ast.ClassDecl{
			Modifiers: modifiers,
			Pub:       pubPos,
			Class:     token.Pos(ctx.CLASS().GetSymbol().GetStart()),
			Name:      className,
			Lbrace:    token.Pos(ctx.ClassBody().LBRACE().GetSymbol().GetStart()),
			Fields: &ast.FieldList{
				Opening: token.Pos(ctx.ClassBody().LBRACE().GetSymbol().GetStart()),
				List:    fields,
				Closing: token.Pos(ctx.ClassBody().RBRACE().GetSymbol().GetStart()),
			},
			Methods: methods,
			Ctors:   constructors,
			Rbrace:  token.Pos(ctx.ClassBody().RBRACE().GetSymbol().GetStart()),
		}
	})
}

// VisitClassBody implements the Visitor interface for ClassBody
func (a *Analyzer) VisitClassBody(ctx *generated.ClassBodyContext) any {
	return a.visitWrapper("ClassBody", ctx, func() any {
		// Create a field list to hold all class members
		var fields []*ast.Field
		var methods []*ast.FuncDecl
		var constructors []*ast.ConstructorDecl

		// Process all class members
		for _, member := range ctx.AllClassMember() {
			if member != nil {
				field, _ := accept[*ast.Field](member, a)
				if field != nil {
					fields = append(fields, field)
				}
			}
		}

		for _, method := range ctx.AllClassMethod() {
			if method != nil {
				funcDecl, _ := accept[*ast.FuncDecl](method, a)
				if funcDecl != nil {
					methods = append(methods, funcDecl)
				}
			}
		}

		for _, builtinMethod := range ctx.AllClassBuiltinMethod() {
			if builtinMethod != nil {
				constructor, _ := accept[*ast.ConstructorDecl](builtinMethod, a)
				if constructor != nil {
					constructors = append(constructors, constructor)
				}
			}
		}

		// Create block statement with all members
		var stmts []ast.Stmt
		for _, field := range fields {
			stmts = append(stmts, &ast.AssignStmt{
				Lhs: field.Name,
				Tok: token.ASSIGN,
				Rhs: field.Value,
			})
		}

		for _, method := range methods {
			stmts = append(stmts, method)
		}

		for _, constructor := range constructors {
			stmts = append(stmts, constructor)
		}

		return &ast.BlockStmt{
			Lbrace: token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
			List:   stmts,
			Rbrace: token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
		}
	})
}

// VisitClassMember implements the Visitor interface for ClassMember
func (a *Analyzer) VisitClassMember(ctx *generated.ClassMemberContext) any {
	return a.visitWrapper("ClassMember", ctx, func() any {
		// Get field name
		fieldName := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}

		// Get field type
		var fieldType ast.Expr
		if ctx.Type_() != nil {
			fieldType, _ = accept[ast.Expr](ctx.Type_(), a)
		}

		// Get default value
		var defaultValue ast.Expr
		if ctx.Expression() != nil {
			defaultValue, _ = accept[ast.Expr](ctx.Expression(), a)
		}

		// Get field modifiers
		var modifiers []ast.Modifier
		for _, modCtx := range ctx.AllClassMemberModifier() {
			if modCtx.PUB() != nil {
				modifiers = append(modifiers, &ast.PubModifier{
					Pub: token.Pos(modCtx.PUB().GetSymbol().GetStart()),
				})
			} else if modCtx.STATIC() != nil {
				modifiers = append(modifiers, &ast.StaticModifier{
					Static: token.Pos(modCtx.STATIC().GetSymbol().GetStart()),
				})
			} else if modCtx.FINAL() != nil {
				modifiers = append(modifiers, &ast.FinalModifier{
					Final: token.Pos(modCtx.FINAL().GetSymbol().GetStart()),
				})
			} else if modCtx.CONST() != nil {
				modifiers = append(modifiers, &ast.ConstModifier{
					Const: token.Pos(modCtx.CONST().GetSymbol().GetStart()),
				})
			}
		}

		return &ast.Field{
			Modifiers: modifiers,
			Name:      fieldName,
			Type:      fieldType,
			Value:     defaultValue,
		}
	})
}

// VisitClassMethod implements the Visitor interface for ClassMethod
func (a *Analyzer) VisitClassMethod(ctx *generated.ClassMethodContext) any {
	return a.visitWrapper("ClassMethod", ctx, func() any {
		// Get class method modifiers
		var classMods []ast.Modifier
		for _, modCtx := range ctx.AllClassMethodModifier() {
			if modCtx.PUB() != nil {
				classMods = append(classMods, &ast.PubModifier{
					Pub: token.Pos(modCtx.PUB().GetSymbol().GetStart()),
				})
			} else if modCtx.STATIC() != nil {
				classMods = append(classMods, &ast.StaticModifier{
					Static: token.Pos(modCtx.STATIC().GetSymbol().GetStart()),
				})
			}
		}

		// Handle standard function declaration
		if ctx.StandardFunctionDeclaration() != nil {
			funcDecl := a.VisitStandardFunctionDeclaration(ctx.StandardFunctionDeclaration().(*generated.StandardFunctionDeclarationContext))
			if funcDecl != nil {
				// Merge class method modifiers with function modifiers
				if fdecl, ok := funcDecl.(*ast.FuncDecl); ok {
					if classMods != nil {
						fdecl.Modifiers = append(classMods, fdecl.Modifiers...)
					}
				}
			}
			return funcDecl
		}

		// Handle lambda function declaration
		if ctx.LambdaFunctionDeclaration() != nil {
			funcDecl := a.VisitLambdaFunctionDeclaration(ctx.LambdaFunctionDeclaration().(*generated.LambdaFunctionDeclarationContext))
			if funcDecl != nil {
				// Merge class method modifiers with function modifiers
				if fdecl, ok := funcDecl.(*ast.FuncDecl); ok {
					if classMods != nil {
						fdecl.Modifiers = append(classMods, fdecl.Modifiers...)
					}
				}
			}
			return funcDecl
		}

		return nil
	})
}

// VisitClassBuiltinMethod implements the Visitor interface for ClassBuiltinMethod
func (a *Analyzer) VisitClassBuiltinMethod(ctx *generated.ClassBuiltinMethodContext) any {
	return a.visitWrapper("ClassBuiltinMethod", ctx, func() any {
		// Get method name
		methodName := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}

		// Get parameters
		params, _ := accept[[]ast.Expr](ctx.ClassBuiltinParameters(), a)

		// Get method body
		var body *ast.BlockStmt
		if ctx.Block() != nil {
			body, _ = accept[*ast.BlockStmt](ctx.Block(), a)
		}

		// Create constructor declaration
		return &ast.ConstructorDecl{
			ClsName: methodName, // Use method name as class name for constructor
			Name:    methodName,
			Recv:    params,
			Body:    body,
		}
	})
}

// VisitImportDeclaration implements the Visitor interface for ImportDeclaration
func (a *Analyzer) VisitImportDeclaration(ctx *generated.ImportDeclarationContext) any {
	return a.visitWrapper("ImportDeclaration", ctx, func() any {
		// Handle different types of imports
		if ctx.ImportSingle() != nil {
			return a.VisitImportSingle(ctx.ImportSingle().(*generated.ImportSingleContext))
		} else if ctx.ImportAll() != nil {
			return a.VisitImportAll(ctx.ImportAll().(*generated.ImportAllContext))
		} else if ctx.ImportMulti() != nil {
			return a.VisitImportMulti(ctx.ImportMulti().(*generated.ImportMultiContext))
		}

		return nil
	})
}

// VisitImportSingle implements the Visitor interface for ImportSingle
func (a *Analyzer) VisitImportSingle(ctx *generated.ImportSingleContext) any {
	return a.visitWrapper("ImportSingle", ctx, func() any {
		// Get the path from string literal
		path := ctx.StringLiteral().GetText()
		// Remove quotes
		path = path[1 : len(path)-1]

		var asPos token.Pos
		var alias string

		// Check if there's an alias
		if ctx.AsIdentifier() != nil {
			asPos = token.Pos(ctx.AsIdentifier().AS().GetSymbol().GetStart())
			alias = ctx.AsIdentifier().Identifier().GetText()
		}

		return &ast.Import{
			ImportPos: token.Pos(ctx.GetStart().GetStart()),
			ImportSingle: &ast.ImportSingle{
				Path:  path,
				As:    asPos,
				Alias: alias,
			},
		}
	})
}

// VisitImportAll implements the Visitor interface for ImportAll
func (a *Analyzer) VisitImportAll(ctx *generated.ImportAllContext) any {
	return a.visitWrapper("ImportAll", ctx, func() any {
		mulPos := token.Pos(ctx.MUL().GetSymbol().GetStart())
		fromPos := token.Pos(ctx.FROM().GetSymbol().GetStart())
		path := ctx.StringLiteral().GetText()
		// Remove quotes
		path = path[1 : len(path)-1]

		var asPos token.Pos
		var alias string

		// Check if there's an alias
		if ctx.AsIdentifier() != nil {
			asPos = token.Pos(ctx.AsIdentifier().AS().GetSymbol().GetStart())
			alias = ctx.AsIdentifier().Identifier().GetText()
		}

		return &ast.Import{
			ImportPos: token.Pos(ctx.GetStart().GetStart()),
			ImportAll: &ast.ImportAll{
				Mul:   mulPos,
				As:    asPos,
				Alias: alias,
				From:  fromPos,
				Path:  path,
			},
		}
	})
}

// VisitImportMulti implements the Visitor interface for ImportMulti
func (a *Analyzer) VisitImportMulti(ctx *generated.ImportMultiContext) any {
	return a.visitWrapper("ImportMulti", ctx, func() any {
		lbracePos := token.Pos(ctx.LBRACE().GetSymbol().GetStart())
		rbracePos := token.Pos(ctx.RBRACE().GetSymbol().GetStart())
		fromPos := token.Pos(ctx.FROM().GetSymbol().GetStart())
		path := ctx.StringLiteral().GetText()
		// Remove quotes
		path = path[1 : len(path)-1]

		var importFields []*ast.ImportField

		// Process all identifier-as-identifier pairs
		for _, idAsId := range ctx.AllIdentifierAsIdentifier() {
			field, _ := accept[*ast.ImportField](idAsId, a)
			if field != nil {
				importFields = append(importFields, field)
			}
		}

		return &ast.Import{
			ImportPos: token.Pos(ctx.GetStart().GetStart()),
			ImportMulti: &ast.ImportMulti{
				Lbrace: lbracePos,
				List:   importFields,
				Rbrace: rbracePos,
				From:   fromPos,
				Path:   path,
			},
		}
	})
}

// VisitIdentifierAsIdentifier implements the Visitor interface for IdentifierAsIdentifier
func (a *Analyzer) VisitIdentifierAsIdentifier(ctx *generated.IdentifierAsIdentifierContext) any {
	return a.visitWrapper("IdentifierAsIdentifier", ctx, func() any {
		field := ctx.Identifier().GetText()

		var asPos token.Pos
		var alias string

		// Check if there's an alias
		if ctx.AsIdentifier() != nil {
			asPos = token.Pos(ctx.AsIdentifier().AS().GetSymbol().GetStart())
			alias = ctx.AsIdentifier().Identifier().GetText()
		}

		return &ast.ImportField{
			Field: field,
			As:    asPos,
			Alias: alias,
		}
	})
}

// VisitDeclareStatement implements the Visitor interface for DeclareStatement
func (a *Analyzer) VisitDeclareStatement(ctx *generated.DeclareStatementContext) any {
	return a.visitWrapper("DeclareStatement", ctx, func() any {

		// Handle different types of declarations that can be inside declare

		decl := &ast.DeclareDecl{
			Docs:    a.collectComments(),
			Declare: token.Pos(ctx.DECLARE().GetSymbol().GetStart()),
		}

		var x ast.Node
		if ctx.Block() != nil {
			// declare { ... }
			block, _ := accept[*ast.BlockStmt](ctx.Block(), a)
			x = block
		} else if ctx.ClassDeclaration() != nil {
			// declare class ...
			classDecl, _ := accept[*ast.ClassDecl](ctx.ClassDeclaration(), a)
			x = classDecl
		} else if ctx.EnumDeclaration() != nil {
			// declare enum ...
			enumDecl, _ := accept[*ast.EnumDecl](ctx.EnumDeclaration(), a)
			x = enumDecl
		} else if ctx.TraitDeclaration() != nil {
			// declare trait ...
			traitDecl, _ := accept[*ast.TraitDecl](ctx.TraitDeclaration(), a)
			x = traitDecl
		} else if ctx.FunctionDeclaration() != nil {
			// declare function ...
			funcDecl, _ := accept[*ast.FuncDecl](ctx.FunctionDeclaration(), a)
			x = funcDecl
		} else if ctx.ModuleDeclaration() != nil {
			// declare module ...
			modDecl, _ := accept[*ast.ModDecl](ctx.ModuleDeclaration(), a)
			x = modDecl
		} else if ctx.TypeDeclaration() != nil {
			// declare type ...
			typeDecl, _ := accept[*ast.TypeDecl](ctx.TypeDeclaration(), a)
			x = typeDecl
		}

		decl.X = x

		return decl
	})
}

// VisitTypeDeclaration implements the Visitor interface for TypeDeclaration
func (a *Analyzer) VisitTypeDeclaration(ctx *generated.TypeDeclarationContext) any {
	return a.visitWrapper("TypeDeclaration", ctx, func() any {
		// Get the type name
		typeName := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}

		// Get the type value
		typeValue, _ := accept[ast.Expr](ctx.Type_(), a)

		return &ast.TypeDecl{
			Type:   token.Pos(ctx.TYPE().GetSymbol().GetStart()),
			Name:   typeName,
			Assign: token.Pos(ctx.ASSIGN().GetSymbol().GetStart()),
			Value:  typeValue,
		}
	})
}

// visitTypeAccessPoint handles type access points like array types
func (a *Analyzer) visitTypeAccessPoint(baseType ast.Expr, ctx *generated.TypeAccessPointContext) ast.Expr {
	// Handle array types like T[5] or T[]
	if ctx.LBRACK() != nil && ctx.RBRACK() != nil {
		// For now, we'll create a simple array type
		// In the future, this could be enhanced to support more complex array types
		return &ast.ArrayType{
			Name: baseType,
		}
	}

	// Handle recursive type access points
	if ctx.TypeAccessPoint() != nil {
		return a.visitTypeAccessPoint(baseType, ctx.TypeAccessPoint().(*generated.TypeAccessPointContext))
	}

	return baseType
}

// VisitCompositeType implements the Visitor interface for CompositeType
func (a *Analyzer) VisitCompositeType(ctx *generated.CompositeTypeContext) any {
	return a.visitWrapper("CompositeType", ctx, func() any {
		// Get the right-hand type
		rightType, _ := accept[ast.Expr](ctx.Type_(), a)

		// Handle recursive composite types first
		if ctx.CompositeType() != nil {
			compositeType, _ := accept[ast.Expr](ctx.CompositeType(), a)
			if compositeType != nil {
				// Merge the types
				if unionType, ok := compositeType.(*ast.UnionType); ok {
					unionType.Types = append(unionType.Types, rightType)
					return unionType
				} else if intersectionType, ok := compositeType.(*ast.IntersectionType); ok {
					intersectionType.Types = append(intersectionType.Types, rightType)
					return intersectionType
				}
			}
		}

		// Handle union types (T | U)
		if ctx.BITOR() != nil {
			// For the base case, we need to get the left type from the parent context
			// This is tricky because we need to look up the call stack
			// For now, let's create a union type and let the parent handle it
			return &ast.UnionType{
				Types: []ast.Expr{rightType},
			}
		}

		// Handle intersection types (T & U)
		if ctx.BITAND() != nil {
			// For the base case, we need to get the left type from the parent context
			// This is tricky because we need to look up the call stack
			// For now, let's create an intersection type and let the parent handle it
			return &ast.IntersectionType{
				Types: []ast.Expr{rightType},
			}
		}

		return rightType
	})
}

// VisitFunctionType implements the Visitor interface for FunctionType
func (a *Analyzer) VisitFunctionType(ctx *generated.FunctionTypeContext) any {
	return a.visitWrapper("FunctionType", ctx, func() any {
		// Get parameters
		params, _ := accept[[]ast.Expr](ctx.ReceiverParameters(), a)

		// Get return type
		var returnType ast.Expr
		if ctx.FunctionReturnValue() != nil {
			returnType, _ = accept[ast.Expr](ctx.FunctionReturnValue(), a)
		}

		// Get positions safely - use default positions for now
		var lparen, rparent token.Pos
		if ctx.GetStart() != nil {
			lparen = token.Pos(ctx.GetStart().GetStart())
		}
		if ctx.GetStop() != nil {
			rparent = token.Pos(ctx.GetStop().GetStop())
		}

		return &ast.FunctionType{
			Lparen:  lparen,
			Recv:    params,
			Rparent: rparent,
			RetVal:  returnType,
		}
	})
}

// VisitObjectType implements the Visitor interface for ObjectType
func (a *Analyzer) VisitObjectType(ctx *generated.ObjectTypeContext) any {
	return a.visitWrapper("ObjectType", ctx, func() any {
		var members []ast.Expr

		// Get the first member
		if ctx.ObjectTypeMember(0) != nil {
			member, _ := accept[ast.Expr](ctx.ObjectTypeMember(0), a)
			if member != nil {
				members = append(members, member)
			}
		}

		// Get additional members
		for i := 1; i < len(ctx.AllObjectTypeMember()); i++ {
			member, _ := accept[ast.Expr](ctx.ObjectTypeMember(i), a)
			if member != nil {
				members = append(members, member)
			}
		}

		return &ast.TypeLiteral{
			Members: members,
		}
	})
}

// VisitObjectTypeMember implements the Visitor interface for ObjectTypeMember
func (a *Analyzer) VisitObjectTypeMember(ctx *generated.ObjectTypeMemberContext) any {
	return a.visitWrapper("ObjectTypeMember", ctx, func() any {
		// Get member name
		memberName := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}

		// Get member type
		memberType, _ := accept[ast.Expr](ctx.Type_(), a)

		return &ast.KeyValueExpr{
			Key:   memberName,
			Colon: token.Pos(ctx.COLON().GetSymbol().GetStart()),
			Value: memberType,
		}
	})
}

// VisitTupleType implements the Visitor interface for TupleType
func (a *Analyzer) VisitTupleType(ctx *generated.TupleTypeContext) any {
	return a.visitWrapper("TupleType", ctx, func() any {
		typeList, _ := accept[[]ast.Expr](ctx.TypeList(), a)
		return &ast.TupleType{
			Types: typeList,
		}
	})
}

// VisitTypeList implements the Visitor interface for TypeList
func (a *Analyzer) VisitTypeList(ctx *generated.TypeListContext) any {
	return a.visitWrapper("TypeList", ctx, func() any {
		var types []ast.Expr

		// Get the first type
		if ctx.Type_(0) != nil {
			firstType, _ := accept[ast.Expr](ctx.Type_(0), a)
			if firstType != nil {
				types = append(types, firstType)
			}
		}

		// Get additional types
		for i := 1; i < len(ctx.AllType_()); i++ {
			typ, _ := accept[ast.Expr](ctx.Type_(i), a)
			if typ != nil {
				types = append(types, typ)
			}
		}

		return types
	})
}

// VisitFunctionReturnValue implements the Visitor interface for FunctionReturnValue
func (a *Analyzer) VisitFunctionReturnValue(ctx *generated.FunctionReturnValueContext) any {
	return a.visitWrapper("FunctionReturnValue", ctx, func() any {
		// Handle single type
		if ctx.Type_() != nil {
			typ, _ := accept[ast.Expr](ctx.Type_(), a)
			return typ
		}

		// Handle tuple type (LPAREN typeList RPAREN)
		if ctx.LPAREN() != nil && ctx.TypeList() != nil && ctx.RPAREN() != nil {
			typeList, _ := accept[[]ast.Expr](ctx.TypeList(), a)
			return &ast.TupleType{
				Types: typeList,
			}
		}

		return nil
	})
}

// VisitGenericParameters implements the Visitor interface for GenericParameters
func (a *Analyzer) VisitGenericParameters(ctx *generated.GenericParametersContext) any {
	return a.visitWrapper("GenericParameters", ctx, func() any {
		// Get parameters from GenericParameterList
		if ctx.GenericParameterList() != nil {
			params, _ := accept[[]ast.Expr](ctx.GenericParameterList(), a)
			return params
		}

		return []ast.Expr{}
	})
}

// VisitGenericParameter implements the Visitor interface for GenericParameter
func (a *Analyzer) VisitGenericParameter(ctx *generated.GenericParameterContext) any {
	return a.visitWrapper("GenericParameter", ctx, func() any {
		// Get parameter name
		paramName := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}

		// Get constraint type if present
		var constraint ast.Expr
		if ctx.Type_() != nil {
			constraint, _ = accept[ast.Expr](ctx.Type_(), a)
		}

		return &ast.TypeParameter{
			Name:        paramName,
			Constraints: []ast.Expr{constraint},
		}
	})
}

// VisitGenericParameterList implements the Visitor interface for GenericParameterList
func (a *Analyzer) VisitGenericParameterList(ctx *generated.GenericParameterListContext) any {
	return a.visitWrapper("GenericParameterList", ctx, func() any {
		var params []ast.Expr

		// Get all generic parameters
		for i := 0; i < len(ctx.AllGenericParameter()); i++ {
			param, _ := accept[ast.Expr](ctx.GenericParameter(i), a)
			if param != nil {
				params = append(params, param)
			}
		}

		return params
	})
}

// VisitEnumDeclaration implements the Visitor interface for EnumDeclaration
func (a *Analyzer) VisitEnumDeclaration(ctx *generated.EnumDeclarationContext) any {
	return a.visitWrapper("EnumDeclaration", ctx, func() any {
		name := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}
		var typeParams []ast.Expr
		if ctx.GenericParameters() != nil {
			typeParams, _ = accept[[]ast.Expr](ctx.GenericParameters(), a)
		}
		var body ast.EnumBody
		switch {
		case ctx.EnumBodySimple() != nil:
			body, _ = accept[*ast.BasicEnumBody](ctx.EnumBodySimple(), a)
		case ctx.EnumBodyAssociated() != nil:
			body, _ = accept[*ast.AssociatedEnumBody](ctx.EnumBodyAssociated(), a)
		case ctx.EnumBodyADT() != nil:
			body, _ = accept[*ast.ADTEnumBody](ctx.EnumBodyADT(), a)
		}
		return &ast.EnumDecl{
			Enum:       token.Pos(ctx.ENUM().GetSymbol().GetStart()),
			Name:       name,
			TypeParams: typeParams,
			Body:       body,
		}
	})
}

func (a *Analyzer) VisitEnumBodySimple(ctx *generated.EnumBodySimpleContext) any {
	return a.visitWrapper("EnumBodySimple", ctx, func() any {
		var values []*ast.EnumValue
		for _, v := range ctx.AllEnumValue() {
			val, _ := accept[*ast.EnumValue](v, a)
			values = append(values, val)
		}
		return &ast.BasicEnumBody{
			Lbrace: token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
			Values: values,
			Rbrace: token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
		}
	})
}

func (a *Analyzer) VisitEnumValue(ctx *generated.EnumValueContext) any {
	return a.visitWrapper("EnumValue", ctx, func() any {
		name := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}
		var value ast.Expr
		var assign token.Pos
		if ctx.ASSIGN() != nil {
			assign = token.Pos(ctx.ASSIGN().GetSymbol().GetStart())
			value, _ = accept[ast.Expr](ctx.Expression(), a)
		}
		return &ast.EnumValue{
			Name:   name,
			Assign: assign,
			Value:  value,
		}
	})
}

func (a *Analyzer) VisitEnumBodyAssociated(ctx *generated.EnumBodyAssociatedContext) any {
	return a.visitWrapper("EnumBodyAssociated", ctx, func() any {
		var fields *ast.FieldList
		if ctx.EnumAssociatedFields() != nil {
			fields, _ = accept[*ast.FieldList](ctx.EnumAssociatedFields(), a)
		}
		var values []*ast.EnumValue
		if ctx.EnumAssociatedValues() != nil {
			for _, v := range ctx.EnumAssociatedValues().AllEnumAssociatedValue() {
				val, _ := accept[*ast.EnumValue](v, a)
				values = append(values, val)
			}
		}
		var methods []ast.Stmt
		if ctx.EnumAssociatedMethods() != nil {
			methods, _ = accept[[]ast.Stmt](ctx.EnumAssociatedMethods(), a)
		}
		// 处理构造函数
		if ctx.EnumAssociatedConstructor() != nil {
			constructors, _ := accept[[]ast.Stmt](ctx.EnumAssociatedConstructor(), a)
			if constructors != nil {
				methods = append(methods, constructors...)
			}
		}
		return &ast.AssociatedEnumBody{
			Lbrace:  token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
			Fields:  fields,
			Values:  values,
			Methods: methods,
			Rbrace:  token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
		}
	})
}

func (a *Analyzer) VisitEnumAssociatedFields(ctx *generated.EnumAssociatedFieldsContext) any {
	return a.visitWrapper("EnumAssociatedFields", ctx, func() any {
		var fields []*ast.Field
		for _, f := range ctx.AllEnumAssociatedField() {
			field, _ := accept[*ast.Field](f, a)
			fields = append(fields, field)
		}
		return &ast.FieldList{
			Opening: token.Pos(ctx.GetStart().GetStart()),
			List:    fields,
			Closing: token.Pos(ctx.GetStop().GetStop()),
		}
	})
}

func (a *Analyzer) VisitEnumField(ctx *generated.EnumFieldContext) any {
	return a.visitWrapper("EnumField", ctx, func() any {
		name := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}
		typ, _ := accept[ast.Expr](ctx.Type_(), a)
		return &ast.Field{
			Name: name,
			Type: typ,
		}
	})
}

func (a *Analyzer) VisitEnumAssociatedValue(ctx *generated.EnumAssociatedValueContext) any {
	return a.visitWrapper("EnumAssociatedValue", ctx, func() any {
		name := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}
		var data []ast.Expr
		if ctx.ExpressionList() != nil {
			data, _ = accept[[]ast.Expr](ctx.ExpressionList(), a)
		}
		return &ast.EnumValue{
			Name:   name,
			Lparen: token.Pos(ctx.LPAREN().GetSymbol().GetStart()),
			Data:   data,
			Rparen: token.Pos(ctx.RPAREN().GetSymbol().GetStart()),
		}
	})
}

func (a *Analyzer) VisitEnumBodyADT(ctx *generated.EnumBodyADTContext) any {
	return a.visitWrapper("EnumBodyADT", ctx, func() any {
		var variants []*ast.EnumVariant
		for _, v := range ctx.AllEnumVariant() {
			variant, _ := accept[*ast.EnumVariant](v, a)
			variants = append(variants, variant)
		}
		// 方法暂时略
		return &ast.ADTEnumBody{
			Lbrace:   token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
			Variants: variants,
			Rbrace:   token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
		}
	})
}

func (a *Analyzer) VisitEnumVariant(ctx *generated.EnumVariantContext) any {
	return a.visitWrapper("EnumVariant", ctx, func() any {
		name := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}
		var fields *ast.FieldList
		if len(ctx.AllEnumField()) > 0 {
			var fieldList []*ast.Field
			for _, f := range ctx.AllEnumField() {
				field, _ := accept[*ast.Field](f, a)
				fieldList = append(fieldList, field)
			}
			fields = &ast.FieldList{
				Opening: token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
				List:    fieldList,
				Closing: token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
			}
		}
		return &ast.EnumVariant{
			Name:   name,
			Lbrace: token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
			Fields: fields,
			Rbrace: token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
		}
	})
}

func (a *Analyzer) VisitEnumAssociatedField(ctx *generated.EnumAssociatedFieldContext) any {
	return a.visitWrapper("EnumAssociatedField", ctx, func() any {
		var modifiers []ast.Modifier
		for _, m := range ctx.AllEnumFieldModifier() {
			mod, _ := accept[ast.Modifier](m, a)
			if mod != nil {
				modifiers = append(modifiers, mod)
			}
		}
		name := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}
		typeExpr, _ := accept[ast.Expr](ctx.Type_(), a)
		return &ast.Field{
			Modifiers: modifiers,
			Name:      name,
			Colon:     token.Pos(ctx.COLON().GetSymbol().GetStart()),
			Type:      typeExpr,
		}
	})
}

func (a *Analyzer) VisitEnumFieldModifier(ctx *generated.EnumFieldModifierContext) any {
	return a.visitWrapper("EnumFieldModifier", ctx, func() any {
		if ctx.FINAL() != nil {
			return &ast.FinalModifier{
				Final: token.Pos(ctx.FINAL().GetSymbol().GetStart()),
			}
		}
		if ctx.CONST() != nil {
			return &ast.ConstModifier{
				Const: token.Pos(ctx.CONST().GetSymbol().GetStart()),
			}
		}
		return nil
	})
}

func (a *Analyzer) VisitEnumAssociatedConstructor(ctx *generated.EnumAssociatedConstructorContext) any {
	return a.visitWrapper("EnumAssociatedConstructor", ctx, func() any {
		var constructors []ast.Stmt
		for _, c := range ctx.AllEnumConstructor() {
			constructor, _ := accept[ast.Stmt](c, a)
			if constructor != nil {
				constructors = append(constructors, constructor)
			}
		}
		return constructors
	})
}

func (a *Analyzer) VisitEnumConstructor(ctx *generated.EnumConstructorContext) any {
	return a.visitWrapper("EnumConstructor", ctx, func() any {
		name, _ := accept[*ast.Ident](ctx.EnumConstructorName(), a)
		var recv []ast.Expr
		if ctx.EnumConstructorParameters() != nil {
			recv, _ = accept[[]ast.Expr](ctx.EnumConstructorParameters(), a)
		}

		// 处理构造函数体
		var body *ast.BlockStmt
		if ctx.Block() != nil {
			body, _ = accept[*ast.BlockStmt](ctx.Block(), a)
		}

		// 处理冒号后的直接初始化表达式
		var initFields []ast.Expr
		if ctx.EnumConstructorInit() != nil {
			init, _ := accept[ast.Expr](ctx.EnumConstructorInit(), a)
			if init != nil {
				initFields = append(initFields, init)
			}
		}

		return &ast.ConstructorDecl{
			Name:       name,
			Recv:       recv,
			InitFields: initFields,
			Body:       body,
		}
	})
}

func (a *Analyzer) VisitEnumConstructorName(ctx *generated.EnumConstructorNameContext) any {
	return a.visitWrapper("EnumConstructorName", ctx, func() any {
		mainName := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier(0).GetSymbol().GetStart()),
			Name:    ctx.Identifier(0).GetText(),
		}
		if len(ctx.AllIdentifier()) > 1 {
			// 处理 Protocol.One 这种情况
			subName := &ast.Ident{
				NamePos: token.Pos(ctx.Identifier(1).GetSymbol().GetStart()),
				Name:    ctx.Identifier(1).GetText(),
			}
			// 这里可以创建一个复合标识符或者使用其他方式表示
			return &ast.Ident{
				NamePos: mainName.NamePos,
				Name:    mainName.Name + "." + subName.Name,
			}
		}
		return mainName
	})
}

func (a *Analyzer) VisitEnumConstructorParameters(ctx *generated.EnumConstructorParametersContext) any {
	return a.visitWrapper("EnumConstructorParameters", ctx, func() any {
		var params []ast.Expr

		// 处理直接初始化: $this.port = -1
		if ctx.EnumConstructorDirectInit(0) != nil {
			init, _ := accept[ast.Expr](ctx.EnumConstructorDirectInit(0), a)
			if init != nil {
				params = append(params, init)
			}
		}

		// 处理普通参数: (v: num)
		if ctx.ReceiverParameterList() != nil {
			recvParams, _ := accept[[]ast.Expr](ctx.ReceiverParameterList(), a)
			params = append(params, recvParams...)
		}

		return params
	})
}

func (a *Analyzer) VisitEnumConstructorDirectInit(ctx *generated.EnumConstructorDirectInitContext) any {
	return a.visitWrapper("EnumConstructorDirectInit", ctx, func() any {
		// 处理直接赋值：$this.port = -1
		varExpr, _ := accept[ast.Expr](ctx.VariableExpression(), a)
		assignPos := token.Pos(ctx.ASSIGN().GetSymbol().GetStart())
		value, _ := accept[ast.Expr](ctx.Expression(), a)

		return &ast.BinaryExpr{
			X:     varExpr,
			OpPos: assignPos,
			Op:    token.ASSIGN,
			Y:     value,
		}
	})
}

func (a *Analyzer) VisitEnumConstructorInit(ctx *generated.EnumConstructorInitContext) any {
	return a.visitWrapper("EnumConstructorInit", ctx, func() any {
		// 处理冒号后的直接赋值：: $this.port = 1
		if ctx.EnumConstructorDirectInit() != nil {
			return a.VisitEnumConstructorDirectInit(ctx.EnumConstructorDirectInit().(*generated.EnumConstructorDirectInitContext))
		}
		return nil
	})
}

func (a *Analyzer) VisitEnumAssociatedMethods(ctx *generated.EnumAssociatedMethodsContext) any {
	return a.visitWrapper("EnumAssociatedMethods", ctx, func() any {
		var methods []ast.Stmt
		for _, m := range ctx.AllEnumMethod() {
			method, _ := accept[ast.Stmt](m, a)
			if method != nil {
				methods = append(methods, method)
			}
		}
		return methods
	})
}

func (a *Analyzer) VisitEnumMethod(ctx *generated.EnumMethodContext) any {
	return a.visitWrapper("EnumMethod", ctx, func() any {
		switch {
		case ctx.StandardFunctionDeclaration() != nil:
			return a.VisitStandardFunctionDeclaration(ctx.StandardFunctionDeclaration().(*generated.StandardFunctionDeclarationContext))
		case ctx.LambdaFunctionDeclaration() != nil:
			return a.VisitLambdaFunctionDeclaration(ctx.LambdaFunctionDeclaration().(*generated.LambdaFunctionDeclarationContext))
		}
		return nil
	})
}

// VisitMatchStatement implements the Visitor interface for MatchStatement
func (a *Analyzer) VisitMatchStatement(ctx *generated.MatchStatementContext) any {
	return a.visitWrapper("MatchStatement", ctx, func() any {
		// Get the match expression
		expr, _ := accept[ast.Expr](ctx.Expression(), a)

		// Get all case clauses
		var cases []*ast.CaseClause
		for _, caseCtx := range ctx.AllMatchCaseClause() {
			caseClause, _ := accept[*ast.CaseClause](caseCtx, a)
			if caseClause != nil {
				cases = append(cases, caseClause)
			}
		}

		// Get default clause if exists
		var defaultClause *ast.CaseClause
		if ctx.MatchDefaultClause() != nil {
			defaultClause, _ = accept[*ast.CaseClause](ctx.MatchDefaultClause(), a)
		}

		return &ast.MatchStmt{
			Match:   token.Pos(ctx.MATCH().GetSymbol().GetStart()),
			Expr:    expr,
			Lbrace:  token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
			Cases:   cases,
			Default: defaultClause,
			Rbrace:  token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
		}
	})
}

// VisitMatchCaseClause implements the Visitor interface for MatchCaseClause
func (a *Analyzer) VisitMatchCaseClause(ctx *generated.MatchCaseClauseContext) any {
	return a.visitWrapper("MatchCaseClause", ctx, func() any {
		// Get the condition expression
		var cond ast.Expr
		if ctx.Type_() != nil {
			cond, _ = accept[ast.Expr](ctx.Type_(), a)
		} else if ctx.MatchEnum() != nil {
			cond, _ = accept[ast.Expr](ctx.MatchEnum(), a)
		} else if ctx.MemberAccess() != nil {
			cond, _ = accept[ast.Expr](ctx.MemberAccess(), a)
		} else if ctx.MatchTriple() != nil {
			cond, _ = accept[ast.Expr](ctx.MatchTriple(), a)
		} else if ctx.Expression() != nil {
			cond, _ = accept[ast.Expr](ctx.Expression(), a)
		} else if ctx.RangeExpression() != nil {
			// RangeExpression returns *ast.RangeExpr which implements ast.Expr
			rangeExpr := a.VisitRangeExpression(ctx.RangeExpression().(*generated.RangeExpressionContext))
			if rangeExpr != nil {
				cond = rangeExpr.(ast.Expr)
			}
		}

		// Get the case body
		body, _ := accept[*ast.BlockStmt](ctx.MatchCaseBody(), a)

		return &ast.CaseClause{
			Cond:   cond,
			DArrow: token.Pos(ctx.DOUBLE_ARROW().GetSymbol().GetStart()),
			Body:   body,
		}
	})
}

// VisitMatchDefaultClause implements the Visitor interface for MatchDefaultClause
func (a *Analyzer) VisitMatchDefaultClause(ctx *generated.MatchDefaultClauseContext) any {
	return a.visitWrapper("MatchDefaultClause", ctx, func() any {
		// Get the case body
		body, _ := accept[*ast.BlockStmt](ctx.MatchCaseBody(), a)

		return &ast.CaseClause{
			Cond:   &ast.Ident{Name: "_"}, // Wildcard for default case
			DArrow: token.Pos(ctx.DOUBLE_ARROW().GetSymbol().GetStart()),
			Body:   body,
		}
	})
}

// VisitMatchEnum implements the Visitor interface for MatchEnum
func (a *Analyzer) VisitMatchEnum(ctx *generated.MatchEnumContext) any {
	return a.visitWrapper("MatchEnum", ctx, func() any {
		// For now, just return the member access
		return a.VisitMemberAccess(ctx.MemberAccess().(*generated.MemberAccessContext))
	})
}

// VisitMatchTriple implements the Visitor interface for MatchTriple
func (a *Analyzer) VisitMatchTriple(ctx *generated.MatchTripleContext) any {
	return a.visitWrapper("MatchTriple", ctx, func() any {
		// For now, just return a simple expression
		// This could be enhanced to handle tuple patterns properly
		return &ast.Ident{Name: "tuple"}
	})
}

// VisitMatchCaseBody implements the Visitor interface for MatchCaseBody
func (a *Analyzer) VisitMatchCaseBody(ctx *generated.MatchCaseBodyContext) any {
	return a.visitWrapper("MatchCaseBody", ctx, func() any {
		if ctx.Expression() != nil {
			expr, _ := accept[ast.Expr](ctx.Expression(), a)
			if expr != nil {
				// Wrap single expression in a block
				return &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{X: expr},
					},
				}
			}
		} else if ctx.ReturnStatement() != nil {
			retStmt, _ := accept[ast.Stmt](ctx.ReturnStatement(), a)
			if retStmt != nil {
				return &ast.BlockStmt{
					List: []ast.Stmt{retStmt},
				}
			}
		} else if ctx.Block() != nil {
			block, _ := accept[*ast.BlockStmt](ctx.Block(), a)
			if block != nil {
				return block
			}
		}

		// If none of the above conditions are met, return an empty block
		// This should not happen in valid syntax, but prevents nil returns
		return &ast.BlockStmt{
			List: []ast.Stmt{},
		}
	})
}

// VisitVariableName implements the Visitor interface for VariableName
func (a *Analyzer) VisitVariableName(ctx *generated.VariableNameContext) any {
	return a.visitWrapper("VariableName", ctx, func() any {
		// Handle wildcard
		if ctx.WILDCARD() != nil {
			return &ast.Ident{
				NamePos: token.Pos(ctx.WILDCARD().GetSymbol().GetStart()),
				Name:    "_",
			}
		}

		// Handle identifier with optional type
		if ctx.Identifier() != nil {
			ident := &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ctx.Identifier().GetText(),
			}

			// If there's a type annotation, create a typed parameter
			if ctx.Type_() != nil {
				typ, _ := accept[ast.Expr](ctx.Type_(), a)
				return &ast.Parameter{
					Name:  ident,
					Colon: token.Pos(ctx.COLON().GetSymbol().GetStart()),
					Type:  typ,
				}
			}

			return ident
		}

		return nil
	})
}

// VisitListExpression implements the Visitor interface for ListExpression
func (a *Analyzer) VisitListExpression(ctx *generated.ListExpressionContext) any {
	return a.visitWrapper("ListExpression", ctx, func() any {
		var elements []ast.Expr

		// Get all expressions in the list
		for _, exprCtx := range ctx.AllExpression() {
			expr, _ := accept[ast.Expr](exprCtx, a)
			if expr != nil {
				elements = append(elements, expr)
			}
		}

		return &ast.ArrayLiteralExpr{
			Lbrack: token.Pos(ctx.LBRACK().GetSymbol().GetStart()),
			Elems:  elements,
			Rbrack: token.Pos(ctx.RBRACK().GetSymbol().GetStart()),
		}
	})
}

// VisitMapExpression implements the Visitor interface for MapExpression
func (a *Analyzer) VisitMapExpression(ctx *generated.MapExpressionContext) any {
	return a.visitWrapper("MapExpression", ctx, func() any {
		var properties []ast.Expr

		// Get all pairs in the map
		for _, pairCtx := range ctx.AllPair() {
			pair, _ := accept[ast.Expr](pairCtx, a)
			if pair != nil {
				properties = append(properties, pair)
			}
		}

		return &ast.ObjectLiteralExpr{
			Lbrace: token.Pos(ctx.LBRACE().GetSymbol().GetStart()),
			Props:  properties,
			Rbrace: token.Pos(ctx.RBRACE().GetSymbol().GetStart()),
		}
	})
}

// VisitPair implements the Visitor interface for Pair
func (a *Analyzer) VisitPair(ctx *generated.PairContext) any {
	return a.visitWrapper("Pair", ctx, func() any {
		var key ast.Expr

		// Handle different types of keys
		if ctx.NumberLiteral() != nil {
			key = &ast.NumericLiteral{
				Value:    ctx.NumberLiteral().GetText(),
				ValuePos: token.Pos(ctx.NumberLiteral().GetSymbol().GetStart()),
			}
		} else if ctx.BoolLiteral() != nil {
			value := ctx.BoolLiteral().GetText()
			if value == "true" {
				key = &ast.TrueLiteral{
					ValuePos: token.Pos(ctx.BoolLiteral().GetSymbol().GetStart()),
				}
			} else {
				key = &ast.FalseLiteral{
					ValuePos: token.Pos(ctx.BoolLiteral().GetSymbol().GetStart()),
				}
			}
		} else if ctx.StringLiteral() != nil {
			raw := ctx.StringLiteral().GetText()
			value := raw[1 : len(raw)-1] // Remove quotes
			key = &ast.StringLiteral{
				Value:    value,
				ValuePos: token.Pos(ctx.StringLiteral().GetSymbol().GetStart()),
			}
		} else if ctx.Identifier() != nil {
			key = &ast.Ident{
				NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
				Name:    ctx.Identifier().GetText(),
			}
		}

		// Get the value expression
		value, _ := accept[ast.Expr](ctx.Expression(), a)

		return &ast.KeyValueExpr{
			Key:   key,
			Colon: token.Pos(ctx.COLON().GetSymbol().GetStart()),
			Value: value,
		}
	})
}

// VisitForeachVariableName implements the Visitor interface for ForeachVariableName
func (a *Analyzer) VisitForeachVariableName(ctx *generated.ForeachVariableNameContext) any {
	return a.visitWrapper("ForeachVariableName", ctx, func() any {
		// Handle variable name
		if ctx.VariableName() != nil {
			return a.VisitVariableName(ctx.VariableName().(*generated.VariableNameContext))
		}

		// Handle variable expression
		if ctx.VariableExpression() != nil {
			return a.VisitVariableExpression(ctx.VariableExpression().(*generated.VariableExpressionContext))
		}

		return nil
	})
}

// VisitRangeExpression implements the Visitor interface for RangeExpression
func (a *Analyzer) VisitRangeExpression(ctx *generated.RangeExpressionContext) any {
	return a.visitWrapper("RangeExpression", ctx, func() any {
		// Parse NumberLiteral as BasicLit
		startLit := &ast.BasicLit{
			Kind:     token.NUM,
			Value:    ctx.NumberLiteral(0).GetText(),
			ValuePos: token.Pos(ctx.NumberLiteral(0).GetSymbol().GetStart()),
		}

		endLit := &ast.BasicLit{
			Kind:     token.NUM,
			Value:    ctx.NumberLiteral(1).GetText(),
			ValuePos: token.Pos(ctx.NumberLiteral(1).GetSymbol().GetStart()),
		}

		return &ast.RangeExpr{
			Start:     startLit,
			DblColon1: token.Pos(ctx.DOUBLE_DOT().GetSymbol().GetStart()),
			End_:      endLit,
		}
	})
}

// VisitExternDeclaration implements the Visitor interface for ExternDeclaration
func (a *Analyzer) VisitExternDeclaration(ctx *generated.ExternDeclarationContext) any {
	return a.visitWrapper("ExternDeclaration", ctx, func() any {
		// Get the extern list
		externList, _ := accept[[]ast.Expr](ctx.ExternList(), a)

		return &ast.ExternDecl{
			Extern: token.Pos(ctx.EXTERN().GetSymbol().GetStart()),
			List:   externList,
		}
	})
}

// VisitExternList implements the Visitor interface for ExternList
func (a *Analyzer) VisitExternList(ctx *generated.ExternListContext) any {
	return a.visitWrapper("ExternList", ctx, func() any {
		var items []ast.Expr

		// Get all extern items
		for _, itemCtx := range ctx.AllExternItem() {
			item, _ := accept[ast.Expr](itemCtx, a)
			if item != nil {
				items = append(items, item)
			}
		}

		return items
	})
}

// VisitExternItem implements the Visitor interface for ExternItem
func (a *Analyzer) VisitExternItem(ctx *generated.ExternItemContext) any {
	return a.visitWrapper("ExternItem", ctx, func() any {
		// Get the identifier
		ident := &ast.Ident{
			NamePos: token.Pos(ctx.Identifier().GetSymbol().GetStart()),
			Name:    ctx.Identifier().GetText(),
		}

		var typ ast.Expr
		var colonPos token.Pos

		if ctx.Type_() != nil {
			typ, _ = accept[ast.Expr](ctx.Type_(), a)
			if ctx.COLON() != nil {
				colonPos = token.Pos(ctx.COLON().GetSymbol().GetStart())
			}
		}

		// Create a parameter to represent the extern item
		return &ast.Parameter{
			Name:  ident,
			Colon: colonPos,
			Type:  typ,
		}
	})
}

// VisitUnsafeExpression implements the Visitor interface for UnsafeExpression
// func (a *Analyzer) VisitUnsafeExpression(ctx *generated.UnsafeExpressionContext) any {
// 	return a.visitWrapper("UnsafeExpression", ctx, func() any {
// 		// Handle unsafe block: unsafe { ... }
// 		if ctx.UnsafeBlock() != nil {
// 			unsafeBlock := ctx.UnsafeBlock()
// 			// Get the text from the UnsafeBlock token
// 			blockText := unsafeBlock.GetText()

// 			// Remove the "unsafe {" prefix and "}" suffix
// 			if len(blockText) > 2 { // "unsafe {" is 8 characters
// 				blockText = blockText[2 : len(blockText)-1]
// 			}

// 			return &ast.UnsafeStmt{
// 				Unsafe: token.Pos(unsafeBlock.GetSymbol().GetStart()),
// 				Start:  token.Pos(unsafeBlock.GetSymbol().GetStart()),
// 				Text:   blockText, // Get the text inside the block
// 				EndPos: token.Pos(unsafeBlock.GetSymbol().GetStop()),
// 			}
// 		}

// 		// Handle unsafe literal: [[ ... ]]
// 		if ctx.UnsafeLiteral() != nil {
// 			text := ctx.UnsafeLiteral().GetText()
// 			// Remove the [[ and ]] delimiters
// 			if len(text) >= 4 {
// 				text = text[2 : len(text)-2]
// 			}
// 			return &ast.UnsafeStmt{
// 				Unsafe: token.Pos(ctx.UnsafeLiteral().GetSymbol().GetStart()),
// 				Text:   text,
// 				EndPos: token.Pos(ctx.UnsafeLiteral().GetSymbol().GetStop()),
// 			}
// 		}

// 		return nil
// 	})
// }

// convertTokenType 将 ANTLR token 类型转换为 Hulo token 类型
func (a *Analyzer) convertTokenType(antlrTokenType int) token.Token {
	switch antlrTokenType {
	case 57: // huloLexerMUL
		return token.ASTERISK // 44
	case 59: // huloLexerDIV
		return token.SLASH // 45
	case 58: // huloLexerMOD
		return token.MOD // 46
	case 56: // huloLexerEXP
		return token.POWER // 47
	case 60: // huloLexerADD
		return token.PLUS // 38
	case 61: // huloLexerSUB
		return token.MINUS // 39
	case 81: // huloLexerEQ
		return token.EQ // 48
	case 80: // huloLexerNEQ
		return token.NEQ // 49
	case 76: // huloLexerLT
		return token.LT // 50
	case 77: // huloLexerGT
		return token.GT // 51
	case 78: // huloLexerLE
		return token.LE // 52
	case 79: // huloLexerGE
		return token.GE // 53
	case 92: // huloLexerAND
		return token.AND // 54
	case 93: // huloLexerOR
		return token.OR // 55
	case 97: // huloLexerBITAND
		return token.CONCAT // 56
	case 68: // huloLexerASSIGN
		return token.ASSIGN // 37
	case 71: // huloLexerMUL_ASSIGN
		return token.ASTERISK_ASSIGN // 47
	case 73: // huloLexerDIV_ASSIGN
		return token.SLASH_ASSIGN // 48
	case 74: // huloLexerMOD_ASSIGN
		return token.MOD_ASSIGN // 49
	case 69: // huloLexerADD_ASSIGN
		return token.PLUS_ASSIGN // 40
	case 70: // huloLexerSUB_ASSIGN
		return token.MINUS_ASSIGN // 41
	case 95: // huloLexerINC
		return token.INC // 42
	case 96: // huloLexerDEC
		return token.DEC // 43
	default:
		// 对于未知的 token 类型，返回 ILLEGAL
		return token.ILLEGAL
	}
}
