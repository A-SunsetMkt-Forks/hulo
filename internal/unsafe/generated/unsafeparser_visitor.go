// Code generated from unsafeParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // unsafeParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by unsafeParser.
type unsafeParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by unsafeParser#template.
	VisitTemplate(ctx *TemplateContext) interface{}

	// Visit a parse tree produced by unsafeParser#content.
	VisitContent(ctx *ContentContext) interface{}

	// Visit a parse tree produced by unsafeParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by unsafeParser#commentStatement.
	VisitCommentStatement(ctx *CommentStatementContext) interface{}

	// Visit a parse tree produced by unsafeParser#variableStatement.
	VisitVariableStatement(ctx *VariableStatementContext) interface{}

	// Visit a parse tree produced by unsafeParser#ifStatement.
	VisitIfStatement(ctx *IfStatementContext) interface{}

	// Visit a parse tree produced by unsafeParser#elseStatement.
	VisitElseStatement(ctx *ElseStatementContext) interface{}

	// Visit a parse tree produced by unsafeParser#loopStatement.
	VisitLoopStatement(ctx *LoopStatementContext) interface{}

	// Visit a parse tree produced by unsafeParser#expressionStatement.
	VisitExpressionStatement(ctx *ExpressionStatementContext) interface{}

	// Visit a parse tree produced by unsafeParser#macroStatement.
	VisitMacroStatement(ctx *MacroStatementContext) interface{}

	// Visit a parse tree produced by unsafeParser#pipelineExpression.
	VisitPipelineExpression(ctx *PipelineExpressionContext) interface{}

	// Visit a parse tree produced by unsafeParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by unsafeParser#logicalOrExpression.
	VisitLogicalOrExpression(ctx *LogicalOrExpressionContext) interface{}

	// Visit a parse tree produced by unsafeParser#logicalAndExpression.
	VisitLogicalAndExpression(ctx *LogicalAndExpressionContext) interface{}

	// Visit a parse tree produced by unsafeParser#equalityExpression.
	VisitEqualityExpression(ctx *EqualityExpressionContext) interface{}

	// Visit a parse tree produced by unsafeParser#comparisonExpression.
	VisitComparisonExpression(ctx *ComparisonExpressionContext) interface{}

	// Visit a parse tree produced by unsafeParser#primaryExpression.
	VisitPrimaryExpression(ctx *PrimaryExpressionContext) interface{}

	// Visit a parse tree produced by unsafeParser#varExpr.
	VisitVarExpr(ctx *VarExprContext) interface{}

	// Visit a parse tree produced by unsafeParser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by unsafeParser#logicalOrExpressionStmt.
	VisitLogicalOrExpressionStmt(ctx *LogicalOrExpressionStmtContext) interface{}

	// Visit a parse tree produced by unsafeParser#logicalAndExpressionStmt.
	VisitLogicalAndExpressionStmt(ctx *LogicalAndExpressionStmtContext) interface{}

	// Visit a parse tree produced by unsafeParser#equalityExpressionStmt.
	VisitEqualityExpressionStmt(ctx *EqualityExpressionStmtContext) interface{}

	// Visit a parse tree produced by unsafeParser#comparisonExpressionStmt.
	VisitComparisonExpressionStmt(ctx *ComparisonExpressionStmtContext) interface{}

	// Visit a parse tree produced by unsafeParser#primaryExpressionStmt.
	VisitPrimaryExpressionStmt(ctx *PrimaryExpressionStmtContext) interface{}
}
