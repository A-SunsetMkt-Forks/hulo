package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/caarlos0/log"
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
		return &ast.BasicLit{
			Kind:     token.IDENT, // Using IDENT for boolean literals
			Value:    ctx.BoolLiteral().GetText(),
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

	defer log.Info("exit ConditionalExpression")
	defer log.DecreasePadding()
	
	// Get the condition expression
	cond, _ := accept[ast.Expr](ctx.ConditionalBoolExpression(), v)

	// If there's no question mark, just return the condition
	if ctx.QUEST() == nil {
		return cond
	}

	// If there's a question mark, we have a ternary expression
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
func (v *Visitor) VisitLogicalExpression(ctx *generated.LogicalExpressionContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter LogicalExpression")
	log.IncreasePadding()

	var ret ast.Expr
	switch ctx.GetChildCount() {
	case 1:
		ret, _ = accept[ast.Expr](ctx.ShiftExpression(0), v)
	case 2:
		x, _ := accept[ast.Expr](ctx.ShiftExpression(0), v)
		y, _ := accept[ast.Expr](ctx.ShiftExpression(1), v)
		ret = &ast.BinaryExpr{
			X:     x,
			OpPos: token.Pos(ctx.GetLogicalOp().GetStart()),
			Op:    token.Token(ctx.GetLogicalOp().GetTokenType()),
			Y:     y,
		}
	}

	log.DecreasePadding()
	log.Info("exit LogicalExpression")

	return ret
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
func (v *Visitor) VisitExpressionStatement(ctx *generated.ExpressionStatementContext) any {
	if ctx == nil {
		return nil
	}
	log.Info("enter ExpressionStatement")
	log.IncreasePadding()
	defer log.Info("exit ExpressionStatement")
	defer log.DecreasePadding()

	expr, _ := accept[ast.Expr](ctx.Expression(), v)

	return &ast.ExprStmt{
		X: expr,
	}
}

// VisitAssignStatement implements the Visitor interface for AssignStatement
func (v *Visitor) VisitAssignStatement(ctx *generated.AssignStatementContext) any {
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
		return v.VisitConditionalExpression(ctx.ConditionalExpression().(*generated.ConditionalExpressionContext))
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
	// Handle pre-increment/decrement
	if ctx.PreIncDecExpression() != nil {
		ret, _ = accept[ast.Expr](ctx.PreIncDecExpression(), v)
	}

	// Handle post-increment/decrement
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
	defer log.Info("exit PostIncDecExpression")
	defer log.DecreasePadding()

	expr := v.Visit(ctx.Factor())
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
	defer log.Info("exit VariableExpression")
	defer log.DecreasePadding()

	return v.VisitMemberAccess(ctx.MemberAccess().(*generated.MemberAccessContext))
}
