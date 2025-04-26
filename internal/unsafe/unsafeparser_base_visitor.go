// Code generated from UnsafeParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package unsafe // UnsafeParser
import "github.com/antlr4-go/antlr/v4"

type BaseUnsafeParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseUnsafeParserVisitor) VisitTemplate(ctx *TemplateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUnsafeParserVisitor) VisitContent(ctx *ContentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUnsafeParserVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUnsafeParserVisitor) VisitIfStatement(ctx *IfStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUnsafeParserVisitor) VisitLoopStatement(ctx *LoopStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUnsafeParserVisitor) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUnsafeParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUnsafeParserVisitor) VisitPipelineExpr(ctx *PipelineExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUnsafeParserVisitor) VisitPrimaryExpr(ctx *PrimaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseUnsafeParserVisitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}
