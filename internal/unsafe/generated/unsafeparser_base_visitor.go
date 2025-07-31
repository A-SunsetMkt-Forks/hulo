// Code generated from unsafeParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // unsafeParser
import "github.com/antlr4-go/antlr/v4"

type BaseunsafeParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseunsafeParserVisitor) VisitTemplate(ctx *TemplateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitContent(ctx *ContentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitCommentStatement(ctx *CommentStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitVariableStatement(ctx *VariableStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitIfStatement(ctx *IfStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitElseStatement(ctx *ElseStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitLoopStatement(ctx *LoopStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitMacroStatement(ctx *MacroStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitPipelineExpression(ctx *PipelineExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitLogicalOrExpression(ctx *LogicalOrExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitLogicalAndExpression(ctx *LogicalAndExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitEqualityExpression(ctx *EqualityExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitComparisonExpression(ctx *ComparisonExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitPrimaryExpression(ctx *PrimaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitVarExpr(ctx *VarExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitLogicalOrExpressionStmt(ctx *LogicalOrExpressionStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitLogicalAndExpressionStmt(ctx *LogicalAndExpressionStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitEqualityExpressionStmt(ctx *EqualityExpressionStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitComparisonExpressionStmt(ctx *ComparisonExpressionStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseunsafeParserVisitor) VisitPrimaryExpressionStmt(ctx *PrimaryExpressionStmtContext) interface{} {
	return v.VisitChildren(ctx)
}
