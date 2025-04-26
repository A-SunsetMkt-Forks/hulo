// Code generated from UnsafeParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package unsafe // UnsafeParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by UnsafeParser.
type UnsafeParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by UnsafeParser#template.
	VisitTemplate(ctx *TemplateContext) interface{}

	// Visit a parse tree produced by UnsafeParser#content.
	VisitContent(ctx *ContentContext) interface{}

	// Visit a parse tree produced by UnsafeParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by UnsafeParser#ifStatement.
	VisitIfStatement(ctx *IfStatementContext) interface{}

	// Visit a parse tree produced by UnsafeParser#loopStatement.
	VisitLoopStatement(ctx *LoopStatementContext) interface{}

	// Visit a parse tree produced by UnsafeParser#expressionStatement.
	VisitExpressionStatement(ctx *ExpressionStatementContext) interface{}

	// Visit a parse tree produced by UnsafeParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by UnsafeParser#pipelineExpr.
	VisitPipelineExpr(ctx *PipelineExprContext) interface{}

	// Visit a parse tree produced by UnsafeParser#primaryExpr.
	VisitPrimaryExpr(ctx *PrimaryExprContext) interface{}

	// Visit a parse tree produced by UnsafeParser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}
}
