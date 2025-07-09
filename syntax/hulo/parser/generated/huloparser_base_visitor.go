// Code generated from huloParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // huloParser
import "github.com/antlr4-go/antlr/v4"

type BasehuloParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BasehuloParserVisitor) VisitFile(ctx *FileContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitComment(ctx *CommentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitAssignStatement(ctx *AssignStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitLambdaAssignStatement(ctx *LambdaAssignStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitVariableNullableExpressions(ctx *VariableNullableExpressionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitVariableNullableExpression(ctx *VariableNullableExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitVariableNames(ctx *VariableNamesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitVariableNameList(ctx *VariableNameListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitConditionalExpression(ctx *ConditionalExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitConditionalBoolExpression(ctx *ConditionalBoolExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitLogicalExpression(ctx *LogicalExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitShiftExpression(ctx *ShiftExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitAddSubExpression(ctx *AddSubExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMulDivExpression(ctx *MulDivExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitIncDecExpression(ctx *IncDecExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitPreIncDecExpression(ctx *PreIncDecExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitPostIncDecExpression(ctx *PostIncDecExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitFactor(ctx *FactorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTripleExpression(ctx *TripleExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitListExpression(ctx *ListExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitNewDelExpression(ctx *NewDelExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMapExpression(ctx *MapExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitPair(ctx *PairContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitVariableExpression(ctx *VariableExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMethodExpression(ctx *MethodExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExpressionList(ctx *ExpressionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMemberAccess(ctx *MemberAccessContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitFileExpression(ctx *FileExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCallExpression(ctx *CallExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCallExpressionLinkedList(ctx *CallExpressionLinkedListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMemberAccessPoint(ctx *MemberAccessPointContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitReceiverArgumentList(ctx *ReceiverArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitNamedArgumentList(ctx *NamedArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitNamedArgument(ctx *NamedArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitVariableName(ctx *VariableNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitRangeExpression(ctx *RangeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitOption(ctx *OptionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitShortOption(ctx *ShortOptionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitLongOption(ctx *LongOptionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCommandExpression(ctx *CommandExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCommandJoin(ctx *CommandJoinContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCommandStream(ctx *CommandStreamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCommandMemberAccess(ctx *CommandMemberAccessContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCommandAccessPoint(ctx *CommandAccessPointContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitReceiverParameters(ctx *ReceiverParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitReceiverParameterList(ctx *ReceiverParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitReceiverParameter(ctx *ReceiverParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitNamedParameters(ctx *NamedParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitNamedParameterList(ctx *NamedParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitNamedParameter(ctx *NamedParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitReturnStatement(ctx *ReturnStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitFunctionDeclaration(ctx *FunctionDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitOperatorIdentifier(ctx *OperatorIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitStandardFunctionDeclaration(ctx *StandardFunctionDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitFunctionReturnValue(ctx *FunctionReturnValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitLambdaFunctionDeclaration(ctx *LambdaFunctionDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitLambdaExpression(ctx *LambdaExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitLambdaBody(ctx *LambdaBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitFunctionSignature(ctx *FunctionSignatureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitFunctionModifier(ctx *FunctionModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMacroStatement(ctx *MacroStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassDeclaration(ctx *ClassDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassModifier(ctx *ClassModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassSuper(ctx *ClassSuperContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassBody(ctx *ClassBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassMember(ctx *ClassMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassMemberModifier(ctx *ClassMemberModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassMethod(ctx *ClassMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassMethodModifier(ctx *ClassMethodModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassBuiltinMethod(ctx *ClassBuiltinMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassBuiltinMethodModifier(ctx *ClassBuiltinMethodModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassBuiltinParameters(ctx *ClassBuiltinParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassNamedParameters(ctx *ClassNamedParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassNamedParameterList(ctx *ClassNamedParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassNamedParameter(ctx *ClassNamedParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassNamedParameterAccessPoint(ctx *ClassNamedParameterAccessPointContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitClassInitializeExpression(ctx *ClassInitializeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumDeclaration(ctx *EnumDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumModifier(ctx *EnumModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumBodySimple(ctx *EnumBodySimpleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumBodyAssociated(ctx *EnumBodyAssociatedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumBodyADT(ctx *EnumBodyADTContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumValue(ctx *EnumValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumAssociatedFields(ctx *EnumAssociatedFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumAssociatedField(ctx *EnumAssociatedFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumFieldModifier(ctx *EnumFieldModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumField(ctx *EnumFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumAssociatedValues(ctx *EnumAssociatedValuesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumAssociatedValue(ctx *EnumAssociatedValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumAssociatedConstructor(ctx *EnumAssociatedConstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumConstructor(ctx *EnumConstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumConstructorName(ctx *EnumConstructorNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumConstructorParameters(ctx *EnumConstructorParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumConstructorDirectInit(ctx *EnumConstructorDirectInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumConstructorInit(ctx *EnumConstructorInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumAssociatedMethods(ctx *EnumAssociatedMethodsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumVariant(ctx *EnumVariantContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumMethods(ctx *EnumMethodsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumMethod(ctx *EnumMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumMethodModifier(ctx *EnumMethodModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumMember(ctx *EnumMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumMemberModifier(ctx *EnumMemberModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumBuiltinMethodModifier(ctx *EnumBuiltinMethodModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumInitialize(ctx *EnumInitializeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitEnumInitializeMember(ctx *EnumInitializeMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTraitDeclaration(ctx *TraitDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTraitModifier(ctx *TraitModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTraitBody(ctx *TraitBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTraitMember(ctx *TraitMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTraitMemberModifier(ctx *TraitMemberModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitImplDeclaration(ctx *ImplDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitImplDeclarationBinding(ctx *ImplDeclarationBindingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitImplDeclarationBody(ctx *ImplDeclarationBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExtendDeclaration(ctx *ExtendDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExtendEnum(ctx *ExtendEnumContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExtendClass(ctx *ExtendClassContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExtendTrait(ctx *ExtendTraitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExtendType(ctx *ExtendTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExtendMod(ctx *ExtendModContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitImportDeclaration(ctx *ImportDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitImportSingle(ctx *ImportSingleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitImportAll(ctx *ImportAllContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitImportMulti(ctx *ImportMultiContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitAsIdentifier(ctx *AsIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitIdentifierAsIdentifier(ctx *IdentifierAsIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitModuleDeclaration(ctx *ModuleDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitModuleStatement(ctx *ModuleStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitUseDeclaration(ctx *UseDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitUseSingle(ctx *UseSingleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitUseMulti(ctx *UseMultiContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitUseAll(ctx *UseAllContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTypeDeclaration(ctx *TypeDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitGenericArguments(ctx *GenericArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitGenericParameters(ctx *GenericParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitGenericParameterList(ctx *GenericParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitGenericParameter(ctx *GenericParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCompositeType(ctx *CompositeTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitType(ctx *TypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTypeLiteral(ctx *TypeLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitNullableType(ctx *NullableTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitUnionType(ctx *UnionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitIntersectionType(ctx *IntersectionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitArrayType(ctx *ArrayTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTypeAccessPoint(ctx *TypeAccessPointContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTypeList(ctx *TypeListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTypeofExpression(ctx *TypeofExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitAsExpression(ctx *AsExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitObjectType(ctx *ObjectTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitObjectTypeMember(ctx *ObjectTypeMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTupleType(ctx *TupleTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitFunctionType(ctx *FunctionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitTryStatement(ctx *TryStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCatchClause(ctx *CatchClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitCatchClauseReceiver(ctx *CatchClauseReceiverContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitFinallyClause(ctx *FinallyClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitThrowStatement(ctx *ThrowStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitBreakStatement(ctx *BreakStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitContinueStatement(ctx *ContinueStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitIfStatement(ctx *IfStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMatchStatement(ctx *MatchStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMatchCaseClause(ctx *MatchCaseClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMatchEnum(ctx *MatchEnumContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMatchEnumMember(ctx *MatchEnumMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMatchTriple(ctx *MatchTripleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMatchTripleValue(ctx *MatchTripleValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMatchCaseBody(ctx *MatchCaseBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitMatchDefaultClause(ctx *MatchDefaultClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitLoopStatement(ctx *LoopStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitLoopLabel(ctx *LoopLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitForeachStatement(ctx *ForeachStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitForeachClause(ctx *ForeachClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitForStatement(ctx *ForStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitForClause(ctx *ForClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitRangeStatement(ctx *RangeStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitRangeClause(ctx *RangeClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitDoWhileStatement(ctx *DoWhileStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitWhileStatement(ctx *WhileStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitDeferStatement(ctx *DeferStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitDeclareStatement(ctx *DeclareStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitChannelInputStatement(ctx *ChannelInputStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitChannelOutputExpression(ctx *ChannelOutputExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitUnsafeExpression(ctx *UnsafeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitComptimeExpression(ctx *ComptimeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExternDeclaration(ctx *ExternDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExternList(ctx *ExternListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasehuloParserVisitor) VisitExternItem(ctx *ExternItemContext) interface{} {
	return v.VisitChildren(ctx)
}
